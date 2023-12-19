package login

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"

	"echonakama/game"
	"echonakama/server/services"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	SystemUserId           = "00000000-0000-0000-0000-000000000000"
	PlaceholderEmailSuffix = "@null.echovrce.com"

	PasswordURLParam          = "password"
	HMDSerialOverrideURLParam = "hmdserial"

	LinkingPageUrl = "https://echovrce.com/link"

	LinkTicketCollection                = "Login:linkTicket"
	LinkTicketIndex                     = "Index_" + LinkTicketCollection
	DiscordAccessTokenCollection        = "Login:discordAccessToken"
	GameClientSettingsStorageCollection = "Login:gameSettings"
	GameClientSettingsStorageKey        = "gameSettings"
	GameProfileStorageCollection        = "Profile"
	ServerGameProfileStorageKey         = "server"
	ClientGameProfileStorageKey         = "client"
	XPlatformIdStorageCollection        = "XPlatformID"

	// The Application ID for Echo VR
	NoOvrAppId = 0
	QuestAppId = 2215004568539258
	PcvrAppId  = 1369078409873402

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

// The data sent to the relay when a client is sucessfully authenticated/authorized.
type LoginSuccessResponse struct {
	EchoUserId         game.EchoUserId         `json:"echo_user_id"`
	DeviceAuthToken    string                  `json:"nk_device_auth_token"`
	EchoSessionToken   string                  `json:"echo_session_token"`
	NkSessionToken     string                  `json:"nk_session_token"`
	EchoClientSettings game.EchoClientSettings `json:"client_settings"`
	GameProfiles       game.GameProfiles       `json:"game_profiles"`
}

// LinkTicket represents a ticket used for linking accounts to Discord.
// It contains the link code, xplatform ID string, and HMD serial number.
type LinkTicket struct {
	Code            string `json:"link_code"`            // the code the user will exchange to link the account
	DeviceAuthToken string `json:"nk_device_auth_token"` // the device ID token to be linked

	// NOTE: The UserIDToken has an index that is created in the InitModule function
	UserIDToken  string        `json:"game_user_id_token"` // the xplatform ID used by EchoVR as a UserID
	LoginRequest *LoginRequest `json:"game_login_request"` // the login request payload that generated this link ticket
}

func (l *LinkTicket) StorageObject() (*runtime.StorageWrite, error) {
	linkTicketJson, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	return &runtime.StorageWrite{
		Collection:      LinkTicketCollection,
		Key:             l.Code,
		UserID:          SystemUserId,
		Value:           string(linkTicketJson),
		PermissionRead:  0,
		PermissionWrite: 0,
		Version:         "*", // do not overwrite existing link tickets
	}, nil
}

// The data used to generate the Device ID authentication string.
type DeviceId struct {
	AppId           int64  `json:"game_app_id"`        // The application ID for the game
	UserIdToken     string `json:"game_user_id_token"` // The xplatform ID string
	HmdSerialNumber string `json:"hmd_serial_number"`  // The HMD serial number
}

// Generate the string used for device authentication
// WARNING: If this is changed, then device "links" will be invalidated
func (d DeviceId) Token() string {
	return fmt.Sprintf("%d:%s:%s", d.AppId, d.UserIdToken, d.HmdSerialNumber)
}

func filterDisplayName(displayName string) string {
	// Use a regular expression to match allowed characters
	re := regexp.MustCompile("[^a-zA-Z0-9_-]")

	// Find the index of the first non-ASCII character
	index := strings.IndexFunc(displayName, func(r rune) bool {
		return r > unicode.MaxASCII
	})

	// If non-ASCII character found, truncate the string up to that index
	if index != -1 {
		displayName = displayName[:index]
	}

	// Filter the string using the regular expression
	filteredUsername := re.ReplaceAllString(displayName, "")

	return filteredUsername
}

// GenerateLinkCode generates a 4 character random link code.
// The character set excludes homoglyphs.
// The random number generator is seeded with the current time to ensure randomness.
// Returns the generated link code as a string.
func GenerateLinkCode() string {
	// Define the set of valid characters for the link code
	characters := "ABCDEFGHJKLMNPRSTUVWXYZ"

	// Create a new local random generator with a known seed value
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create an array with 4 elements
	var indices []int
	for i := 0; i < 4; i++ {
		indices = append(indices, i)
	}

	// Randomly select an index from the array and generate the code
	var code string
	for range indices {
		code += string(characters[rng.Intn(len(characters))])
	}

	return code
}

// LinkTicket generates a link ticket for the provided xplatformId and hmdSerialNumber.
func (request *LoginRequest) LinkTicket(serviceContext *services.ServiceContext, userID string) (*LinkTicket, *runtime.Error) {
	linkTicket := &LinkTicket{}
	ctx := serviceContext.Ctx
	nk := serviceContext.NakamaModule
	logger := serviceContext.Logger
	// Check if a link ticket already exists for the provided xplatformId and hmdSerialNumber
	objectIDs, err := nk.StorageIndexList(ctx, SystemUserId, LinkTicketIndex, fmt.Sprintf("+value.game_user_id_token:%s", request.DeviceId().UserIdToken), 10)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("error listing link tickets: `%q`  %v", request.DeviceId().UserIdToken, err), StatusInternalError)
	}
	logger.WithField("objectIds", objectIDs).Debug("Link ticket found/generated.")
	// Link ticket was found. Return the link ticket.
	if objectIDs != nil {
		for _, record := range objectIDs.Objects {
			json.Unmarshal([]byte(record.Value), &linkTicket)

			return linkTicket, nil
		}
	}

	// Generate a link code and attempt to write it to storage
	for {

		// loop until we have a unique link code
		linkTicket = &LinkTicket{
			Code:            GenerateLinkCode(),
			DeviceAuthToken: request.DeviceId().Token(),
			UserIDToken:     request.DeviceId().UserIdToken,
			LoginRequest:    request,
		}

		linkTicketStorageObject, err := linkTicket.StorageObject()
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("error preparing link ticket storage object: %v", err), StatusInternalError)
		}

		// Write the link code to storage
		objectIDs := []*runtime.StorageWrite{
			linkTicketStorageObject,
		}

		_, err = nk.StorageWrite(ctx, objectIDs)

		if err != nil {
			continue
		}
		return linkTicket, nil
	}
}

// authenticateClient is a function that authenticates a client using the provided service context and login request.
// It returns the Nakama user ID, username, and any error that occurred during authentication.
func authenticateAccountDevice(serviceContext *services.ServiceContext, loginRequest *LoginRequest, authPassword string) (*api.Account, *runtime.Error) {
	ctx := serviceContext.Ctx
	nk := serviceContext.NakamaModule
	logger := serviceContext.Logger
	var nkUserID string
	var err error
	UserIdToken := loginRequest.EchoUserId.String()

	// Validate the user identifier
	if !loginRequest.EchoUserId.Valid() {
		return nil, runtime.NewError(fmt.Sprintf("invalid Game User ID: %q", UserIdToken), StatusInvalidArgument)
	}

	nkUserID, _, _, err = nk.AuthenticateDevice(ctx, loginRequest.DeviceId().Token(), "", false)
	if err != nil { // account is missing, this is okay.
		logger.WithField("err", err).WithField("device_auth_token", loginRequest.DeviceId().Token()).Debug("Device not linked.")
		err = nil
	}

	// If the account is not linked, create a link ticket and return an error
	if nkUserID == "" {
		// No Account. Create link ticket and return error
		linkTicket, err := loginRequest.LinkTicket(serviceContext, SystemUserId)
		if err != nil {
			logger.WithField("err", err).Error("unable to generate link ticket.")
			return nil, runtime.NewError(fmt.Sprintf("unable to generate link ticket: %q", UserIdToken), StatusInternalError)
		}
		logger.WithField("linkTicket", linkTicket).Debug("Link ticket found/generated.")
		return nil, runtime.NewError(fmt.Sprintf("visit %s and enter code: %s", LinkingPageUrl, linkTicket.Code), StatusNotFound)

	}

	// Authorize the authenticated account
	account, err := nk.AccountGetId(ctx, nkUserID)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("unable to get account for Id: %q", UserIdToken), StatusInternalError)
	}

	// Check if the account is disabled/banned
	if account.GetDisableTime() != nil {
		return nil, runtime.NewError(fmt.Sprintf("account Permanently Banned: %q", UserIdToken), StatusPermissionDenied)
	}

	// If the account has a password set, authenticate the 'auth=' query param as the password

	if account.Email != "" {

		_, _, _, err = nk.AuthenticateEmail(ctx, account.Email, authPassword, "", false)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("invalid password for account: %q", UserIdToken), StatusUnauthenticated)
		}
	} else if authPassword != "" {
		// if the 'auth=' query param is set, set the password to the 'auth=' query param
		nk.LinkEmail(ctx, nkUserID, account.User.Id+PlaceholderEmailSuffix, authPassword)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("unable to set password for account: %q", UserIdToken), StatusInternalError)
		}
	}

	// reauthenticate with custom device auth token
	_, _, _, err = nk.AuthenticateDevice(ctx, loginRequest.DeviceId().Token(), "", false)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("invalid device auth token for account: %q", UserIdToken), StatusUnauthenticated)
	}

	if account.CustomId == "" {
		// if the account does not have a customId, return nothing, and let the client know that they need to link their account
		return nil, runtime.NewError(fmt.Sprintf("account (%q) not linked to Discord. visit %s", UserIdToken, LinkingPageUrl), StatusInternalError)
	}
	// verify that the account has a valid customId by authenticating to it (this activates the validation/refresh hook)
	// if the account does not have a valid customId, this will return an error
	_, _, _, err = nk.AuthenticateCustom(ctx, account.CustomId, "", false)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("discord link is invalid. visit %s", LinkingPageUrl), StatusInternalError)
	}

	return account, nil
}

// ProcessLoginRequest processes a login request and returns the login success response or an error.
// It returns a string representing the login success response and a *runtime.Error object if there is an error.
func ProcessLoginRequest(serviceContext *services.ServiceContext, request *LoginRequest) (*LoginSuccessResponse, *runtime.Error) {
	ctx := serviceContext.Ctx
	logger := serviceContext.Logger
	nk := serviceContext.NakamaModule

	// immediately pop the password out of the requests to avoid
	// displaying it in the logs
	authPassword := request.UserPassword
	request.UserPassword = ""

	relayNkUserID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return nil, runtime.NewError("relay must authenticate", StatusUnauthenticated)
	}

	relayUserName, ok := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)
	if !ok {
		return nil, runtime.NewError("relay must authenticate", StatusUnauthenticated)
	}

	logger.WithField("relayUserName", relayUserName).WithField("request", request).Debug("Processing login request.")

	account, nkerr := authenticateAccountDevice(serviceContext, request, authPassword)
	if nkerr != nil {
		logger.WithField("nkerr", nkerr).Error("authentication errored.")
		return nil, nkerr
	}

	// Authorize the client to use the authenticated account
	playerNkUserID := account.User.Id
	currentTimestamp := time.Now().UTC().Unix()
	sessionGuid := uuid.New()

	// Generate a session token with the Guid
	token, _, err := nk.AuthenticateTokenGenerate(account.User.Id, account.User.Username, 0, map[string]string{"sessionGuid": sessionGuid.String()})
	if err != nil {
		logger.WithField("err", err).Error("authenticate token generate error.")
		return nil, runtime.NewError("authenticate token generation error.", StatusInternalError)
	}

	// generate a blank playerData object
	gameProfiles := game.DefaultGameProfiles(request.EchoUserId, request.Metadata.DisplayName)

	// read the client profile from the storage layer
	// TODO: Extact method
	objectIds := []*runtime.StorageRead{{
		Collection: GameProfileStorageCollection,
		Key:        ClientGameProfileStorageKey,
		UserID:     playerNkUserID,
	}, {
		Collection: GameProfileStorageCollection,
		Key:        ServerGameProfileStorageKey,
		UserID:     playerNkUserID,
	}, {
		Collection: GameClientSettingsStorageCollection,
		Key:        GameClientSettingsStorageKey,
		UserID:     relayNkUserID,
	}, /* {
		Collection: "Login:linking_settings",
		Key:        "linking_settings",
		UserID:     relayUserId,
	}, */
	}

	var loginSettings game.EchoClientSettings

	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("storage read error.")
	} else {
		for _, record := range records {
			if record.Key == ClientGameProfileStorageKey {
				err = json.Unmarshal([]byte(record.Value), &gameProfiles.Client)
				if err != nil {
					return nil, runtime.NewError(fmt.Sprintf("error unmarshaling client playerData: %v", err), StatusInternalError)
				}
			} else if record.Key == ServerGameProfileStorageKey {
				err = json.Unmarshal([]byte(record.Value), &gameProfiles.Server)
				if err != nil {
					return nil, runtime.NewError(fmt.Sprintf("error unmarshaling server playerData: %v", err), StatusInternalError)
				}
			} else if record.Key == GameClientSettingsStorageKey {
				err = json.Unmarshal([]byte(record.Value), &loginSettings)
				if err != nil {
					return nil, runtime.NewError(fmt.Sprintf("error unmarshaling server playerData: %v", err), StatusInternalError)
				}
			}
		}
	}

	// Update the server profile's logintime and updatetime.
	gameProfiles.Server.LobbyVersion = int64(request.Metadata.LobbyVersion)
	gameProfiles.Server.LoginTime = currentTimestamp
	gameProfiles.Server.ModifyTime = account.User.UpdateTime.Seconds
	gameProfiles.Server.UpdateTime = account.User.UpdateTime.Seconds

	gameProfiles.Server.DisplayName = account.User.DisplayName
	gameProfiles.Client.DisplayName = account.User.DisplayName

	// Write the profile data to storage
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		logger.WithField("err", err).Error("invalid account info.")
		return nil, runtime.NewError(fmt.Sprintf("error marshaling accountInfo: %v", err), StatusInternalError)
	}
	clientProfileJson, err := json.Marshal(gameProfiles.Client)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("error marshaling client profile playerData: %v", err), StatusInternalError)
	}
	serverProfileJson, err := json.Marshal(gameProfiles.Server)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("error marshaling server profile playerData: %v", err), StatusInternalError)
	}

	// Write the latest profile data to storage
	objectIDs := []*runtime.StorageWrite{
		{
			Collection:      XPlatformIdStorageCollection,
			Key:             request.DeviceId().UserIdToken,
			UserID:          playerNkUserID,
			Value:           string(jsonRequest),
			PermissionRead:  0,
			PermissionWrite: 0,
		},
		{
			Collection:      GameProfileStorageCollection,
			Key:             ClientGameProfileStorageKey,
			UserID:          playerNkUserID,
			Value:           string(clientProfileJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		},
		{
			Collection:      GameProfileStorageCollection,
			Key:             ServerGameProfileStorageKey,
			UserID:          playerNkUserID,
			Value:           string(serverProfileJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		}}

	// Get the login settings from storage
	// TODO: extract method
	var loginSettingsJson []byte

	// if loginSettings is empty, create a default loginSettings object, and write it storage
	// TODO move this to a hook when a relay is authenticated
	if loginSettings.Env == "" {
		loginSettings = game.DefaultEchoClientSettings()
		loginSettingsJson, err = json.Marshal(loginSettings)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("error marshalling LoginSettings: %v", err), StatusInternalError)
		}

		objectIDs = append(objectIDs, &runtime.StorageWrite{
			Collection:      GameClientSettingsStorageCollection,
			Key:             GameClientSettingsStorageKey,
			UserID:          relayNkUserID,
			Value:           string(loginSettingsJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		})
	}

	_, err = nk.StorageWrite(ctx, objectIDs)
	if err != nil {
		logger.WithField("err", err).Error("storage write error.")
		return nil, runtime.NewError(fmt.Sprintf("error writing profile data: %v", err), StatusInternalError)
	}

	loginSuccess := LoginSuccessResponse{
		EchoUserId:         request.EchoUserId,
		DeviceAuthToken:    request.DeviceId().Token(),
		EchoSessionToken:   sessionGuid.String(),
		NkSessionToken:     token,
		EchoClientSettings: loginSettings,
		GameProfiles:       gameProfiles,
	}

	logger.WithField("loginSuccess", loginSuccess).Debug("Login Success.")
	return &loginSuccess, nil
}
