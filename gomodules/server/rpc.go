package server

import (
	"context"
	"database/sql"
	"echonakama/server/services"
	"echonakama/server/services/login"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgrijalva/jwt-go"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	// Websocket Error Codes
	StatusOK                 = 0  // StatusOK indicates a successful operation.
	StatusCanceled           = 1  // StatusCanceled indicates the operation was canceled.
	StatusUnknown            = 2  // StatusUnknown indicates an unknown error occurred.
	StatusInvalidArgument    = 3  // StatusInvalidArgument indicates an invalid argument was provided.
	StatusDeadlineExceeded   = 4  // StatusDeadlineExceeded indicates the operation exceeded the deadline.
	StatusNotFound           = 5  // StatusNotFound indicates the requested resource was not found.
	StatusAlreadyExists      = 6  // StatusAlreadyExists indicates the resource already exists.
	StatusPermissionDenied   = 7  // StatusPermissionDenied indicates the operation was denied due to insufficient permissions.
	StatusResourceExhausted  = 8  // StatusResourceExhausted indicates the resource has been exhausted.
	StatusFailedPrecondition = 9  // StatusFailedPrecondition indicates a precondition for the operation was not met.
	StatusAborted            = 10 // StatusAborted indicates the operation was aborted.
	StatusOutOfRange         = 11 // StatusOutOfRange indicates a value is out of range.
	StatusUnimplemented      = 12 // StatusUnimplemented indicates the operation is not implemented.
	StatusInternalError      = 13 // StatusInternal indicates an internal server error occurred.
	StatusUnavailable        = 14 // StatusUnavailable indicates the service is currently unavailable.
	StatusDataLoss           = 15 // StatusDataLoss indicates a loss of data occurred.
	StatusUnauthenticated    = 16 // StatusUnauthenticated indicates the request lacks valid authentication credentials.
)

// Handles the user login request from Echo Relay
// LoginRequestRpc is a function that handles a login request RPC.
// It takes a context, logger, database connection, Nakama module, and payload as input.
// It returns a string and an error.
// The payload is expected to be in JSON format and will be parsed into a LoginRequest object.
// The function creates a ServiceContext object and passes it to the login service for processing.
// If the login request is successful, it marshals the LoginSuccess object into JSON and returns it as a string.
// If there is an error during the process, it returns an error with an appropriate message.
func LoginRequestRpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string, discordBot *discordgo.Session) (string, error) {

	// Parse the payload into a LoginRequest object
	var request login.LoginRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		logger.WithField("err", err).Error("Unable to unmarshal payload")
		return "", runtime.NewError("Unable to unmarshal payload", StatusInvalidArgument)
	}
	// Create a ServiceContext object to pass to the login service
	serviceContext := &services.ServiceContext{
		Ctx:          ctx,
		Logger:       logger,
		DbConnection: db,
		NakamaModule: nk,
		DiscordBot:   discordBot,
	}

	// Process the login request
	success, nkerr := login.ProcessLoginRequest(serviceContext, &request)
	if nkerr != nil {
		logger.WithField("err", nkerr).Error("login failed")

		return "", runtime.NewError(nkerr.Message, StatusInternalError)
	}

	// Marshal the LoginSuccess object into JSON
	loginSuccessJson, err := json.Marshal(success)
	if err != nil {
		return "", runtime.NewError(fmt.Sprintf("error marshalling LoginSuccess response: %v", err), StatusInternalError)
	}

	return string(loginSuccessJson), nil
}

// DiscordSignInRpc is a function that handles the Discord sign-in RPC.
// It takes in the context, logger, database connection, Nakama module, and payload as parameters.
// The function exchanges the provided code for an access token, creates a Discord client,
// retrieves the Discord user, checks if a user exists with the Discord ID as a Nakama username,
// creates a user if necessary, gets the account data, relinks the custom ID if necessary,
// writes the access token to storage, updates the account information, generates a session token,
// stores the JWT in the user's metadata, and returns the session token and Discord username as a JSON response.
// If any error occurs during the process, an error message is returned.
func DiscordSignInRpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.WithField("payload", payload).Info("DiscordSignInRpc")

	vars, _ := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)
	clientId := vars["DISCORD_CLIENT_ID"]
	clientSecret := vars["DISCORD_CLIENT_SECRET"]
	nkUserId := ""

	type DiscordSignInRequest struct {
		Code             string `json:"code"`
		OAuthRedirectUrl string `json:"oauth_redirect_url"`
	}

	// Parse the payload into a LoginRequest object
	var request DiscordSignInRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		logger.WithField("err", err).WithField("payload", payload).Error("Unable to unmarshal payload")
		return "", runtime.NewError("Unable to unmarshal payload", StatusInvalidArgument)
	}
	if request.Code == "" {
		logger.Error("DiscordSignInRpc: Code is empty")
		return "", runtime.NewError("Code is empty", StatusInvalidArgument)
	}
	if request.OAuthRedirectUrl == "" {
		logger.Error("DiscordSignInRpc: OAuthRedirectUrl is empty")
		return "", runtime.NewError("OAuthRedirectUrl is empty", StatusInvalidArgument)
	}

	// Exchange the code for an access token
	accessToken, err := login.ExchangeCodeForAccessToken(logger, request.Code, clientId, clientSecret, request.OAuthRedirectUrl)
	if err != nil {
		logger.WithField("err", err).Error("Unable to exchange code for access token")
		return "", runtime.NewError("Unable to exchange code for access token", StatusInternalError)
	}

	// Create a Discord client
	discord, err := discordgo.New("Bearer " + accessToken.AccessToken)
	if err != nil {
		logger.WithField("err", err).Error("Unable to create Discord client")
		return "", runtime.NewError("Unable to create Discord client", StatusInternalError)
	}

	// Get the Discord user
	user, err := discord.User("@me")
	if err != nil {
		logger.WithField("err", err).Error("Unable to get Discord user")
		return "", runtime.NewError("Unable to get Discord user", StatusInternalError)
	}

	// check if a user exists with the Discord ID as a Nk username
	results, err := nk.UsersGetUsername(ctx, []string{user.ID})
	if err != nil {
		return "", runtime.NewError("Unable to get user", StatusInternalError)
	}
	if len(results) == 0 {
		// create the user
		nkUserId, _, _, err = nk.AuthenticateCustom(ctx, user.ID, user.ID, true)
		if err != nil {
			return "", runtime.NewError("Unable to create user", StatusInternalError)
		}
	} else {
		nkUserId = results[0].Id
	}

	// get the account data
	account, err := nk.AccountGetId(ctx, nkUserId)
	if err != nil {
		logger.WithField("err", err).Error("Unable to get account")
		return "", runtime.NewError("Unable to get account", StatusInternalError)
	}

	// If the customId doesn't match the discord token, relink it
	if account.CustomId != accessToken.AccessToken {
		_ = nk.UnlinkCustom(ctx, nkUserId, account.CustomId)
		err = nk.LinkCustom(ctx, nkUserId, accessToken.AccessToken)
		if err != nil {
			logger.WithField("err", err).Error("Unable to link custom")
			return "", runtime.NewError("Unable to link custom", StatusInternalError)
		}
	}
	// Write the access token to storage
	login.WriteAccessTokenToStorage(ctx, logger, nk, nkUserId, accessToken)
	if err != nil {
		logger.WithField("err", err).Error("Unable to write access token to storage")
		return "", runtime.NewError("Unable to write access token to storage", StatusInternalError)
	}
	logger.WithField("user.Username", user.Username).Info("DiscordSignInRpc: Wrote access token to storage")

	if err := nk.AccountUpdateId(ctx, nkUserId, "", nil, user.Username, "", "", "", ""); err != nil {
		return "", runtime.NewError("Unable to update account", StatusInternalError)
	}

	// Generate a session token for the user to use to authenticate for device linking
	sessionToken, _, err := nk.AuthenticateTokenGenerate(nkUserId, user.ID, time.Now().UTC().Unix()+3600, nil)
	if err != nil {
		logger.WithField("err", err).Error("Unable to generate session token")
		return "", runtime.NewError("Unable to generate session token", StatusInternalError)
	}

	// store the jwt in the user's metadata so we can verify it later

	type DiscordSignInResponse struct {
		SessionToken    string `json:"sessionToken"`
		DiscordUsername string `json:"discordUsername"`
	}
	response := DiscordSignInResponse{
		SessionToken:    sessionToken,
		DiscordUsername: user.Username,
	}

	// Marshal the Discord Token object into JSON
	responseJson, err := json.Marshal(response)
	if err != nil {
		return "", runtime.NewError(fmt.Sprintf("error marshalling LoginSuccess response: %v", err), StatusInternalError)
	}

	return string(responseJson), nil
}

// LinkDeviceRpc is a function that handles the linking of a device to a user account.
// It takes in the context, logger, database connection, Nakama module, and payload as parameters.
// The payload should be a JSON string containing the session token and link code.
// It returns an empty string and an error.
// The function performs the following steps:
// 1. Unmarshalls the payload to extract the session token and link code.
// 2. Validates the session token and retrieves the UID from it.
// 3. Retrieves the link ticket from storage using the link code.
// 4. Verifies the session token using the link ticket's device auth token.
// 5. Retrieves the user account using the UID.
// 6. Links the device to the user account.
// 7. Deletes the link ticket from storage.
func LinkDeviceRpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	vars, _ := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)
	// unmarshall the payload
	type LinkDeviceRequest struct {
		SessionToken string `json:"sessionToken"`
		LinkCode     string `json:"linkCode"`
	}
	var request LinkDeviceRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		logger.WithField("err", err).WithField("payload", payload).Error("Unable to unmarshal payload")
		return "", runtime.NewError("Unable to unmarshal payload", StatusInvalidArgument)
	}
	if request.SessionToken == "" {
		logger.Error("linkDeviceRpc: SessionToken is empty")
		return "", runtime.NewError("SessionToken is empty", StatusInvalidArgument)
	}
	if request.LinkCode == "" {
		logger.Error("linkDeviceRpc: LinkCode is empty")
		return "", runtime.NewError("LinkCode is empty", StatusInvalidArgument)
	}

	// verify the sessionToken. It's a JWT signed by the server.
	// pull the uid out of it
	logger.WithField("sessionToken", request.SessionToken).Info("Verifying session token")
	token, err := verifySignedJwt(request.SessionToken, []byte(vars["SESSION_ENCRYPTION_KEY"]))
	if err != nil {
		logger.WithField("err", err).Error("Unable to verify session token")
		return "", runtime.NewError("Unable to verify session token", StatusInternalError)
	}
	uid := token.Claims.(jwt.MapClaims)["uid"].(string)

	if err := LinkAccountDevice(ctx, nk, logger, request.LinkCode, uid); err != nil {
		return "", err
	}

	return "", nil
}

func LinkDiscordDevice(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, linkCode string, discordId string) error {
	results, err := nk.UsersGetUsername(ctx, []string{discordId})
	if err != nil {
		return errors.New("unable to get user: %v")
	}
	if len(results) == 0 {
		return errors.New("user not found")
	}
	LinkAccountDevice(ctx, nk, logger, linkCode, results[0].Id)
	return nil
}

func LinkAccountDevice(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, linkCode string, uid string) error {
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: login.LinkTicketCollection,
			Key:        linkCode,
			UserID:     login.SystemUserId,
		},
	})
	if err != nil {
		logger.WithField("err", err).Error("Unable to read link ticket from storage")
		return runtime.NewError("Unable to read link ticket from storage", StatusInternalError)
	}
	if len(objects) == 0 {
		logger.WithField("linkCode", linkCode).Error("Unable to find link ticket")
		return runtime.NewError("Unable to find link ticket", StatusNotFound)
	}
	var linkTicket login.LinkTicket
	if err := json.Unmarshal([]byte(objects[0].Value), &linkTicket); err != nil {
		logger.WithField("err", err).Error("Unable to unmarshal link ticket")
		return runtime.NewError("Unable to unmarshal link ticket", StatusInternalError)
	}

	account, err := nk.AccountGetId(ctx, uid)
	if err != nil {
		logger.WithField("err", err).Error("Unable to get account")
		return runtime.NewError("Unable to get account", StatusInternalError)
	}

	if err := nk.LinkDevice(ctx, account.GetUser().GetId(), linkTicket.DeviceAuthToken); err != nil {
		logger.WithField("err", err).Error("Unable to link device")
		return runtime.NewError("Unable to link device", StatusInternalError)
	}

	if err := nk.StorageDelete(ctx, []*runtime.StorageDelete{
		{
			Collection: login.LinkTicketCollection,
			Key:        linkCode,
			UserID:     login.SystemUserId,
		},
	}); err != nil {
		logger.WithField("err", err).Error("Unable to delete link ticket")
		return runtime.NewError("Unable to delete link ticket", StatusInternalError)
	}
	return nil
}

// verifyJWT parses and verifies a JWT token using the provided key function.
// It returns the parsed token if it is valid, otherwise it returns an error.
// Nakama JWT's are signed by the `session.session_encryption_key` in the Nakama config.
func verifySignedJwt(tokenString string, hmacSampleSecret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
