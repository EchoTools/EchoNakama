package login

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"echonakama/game"
	"echonakama/server/services"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	SystemUserId = "00000000-0000-0000-0000-000000000000"

	PasswordURLParam          = "password"
	HMDSerialOverrideURLParam = "hmdserial"

	LinkTicketCollection                = "Login:linkTicket"
	LinkTicketIndex                     = "Index_" + LinkTicketCollection
	DiscordAccessTokenCollection        = "Login:discordAccessToken"
	DiscordAccessTokenKey               = "accessToken"
	GameClientSettingsStorageCollection = "Login:gameSettings"
	GameClientSettingsStorageKey        = "gameSettings"
	GameProfileStorageCollection        = "Profile"
	ServerGameProfileStorageKey         = "server"
	ClientGameProfileStorageKey         = "client"
	XPlatformIdStorageCollection        = "XPlatformId"
	IpAddressIndex                      = "Index_" + XPlatformIdStorageCollection

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

// ProcessLoginRequest processes a login request and returns the login success response or an error.
// It returns a string representing the login success response and a *runtime.Error object if there is an error.
func ProcessLoginRequest(serviceContext *services.ServiceContext, request *LoginRequest) (*LoginSuccessResponse, *runtime.Error) {
	ctx := serviceContext.Ctx
	logger := serviceContext.Logger
	nk := serviceContext.NakamaModule

	// immediately strip the password out of the requests to avoid
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

	logger.WithField("relayUserName", relayUserName).Debug("Processing login request for user %s on relay %s", request.EchoUserId, relayNkUserID)

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
	gameProfiles.Server.LobbyVersion = request.Metadata.LobbyVersion
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
			PermissionRead:  1,
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

	logger.Debug("Logged %s in successfully.", gameProfiles.Server.DisplayName)
	return &loginSuccess, nil
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

// authenticateClient is a function that authenticates a client using the provided service context and login request.
// It returns the Nakama user ID, username, and any error that occurred during authentication.
func authenticateAccountDevice(serviceContext *services.ServiceContext, loginRequest *LoginRequest, authPassword string) (*api.Account, *runtime.Error) {
	ctx := serviceContext.Ctx
	nk := serviceContext.NakamaModule
	logger := serviceContext.Logger
	discordBot := serviceContext.DiscordBot
	// Get the linking page url from the environment
	vars, _ := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)
	linkingPageUrl := vars["LINK_PAGE_URL"]
	placeholderEmailDomain := vars["PLACEHOLDER_EMAIL_DOMAIN"]
	var nkUserId string
	var err error
	UserIdToken := loginRequest.EchoUserId.String()

	// Validate the user identifier
	if !loginRequest.EchoUserId.Valid() {
		return nil, runtime.NewError(fmt.Sprintf("invalid Game User ID: %q", UserIdToken), StatusInvalidArgument)
	}

	// Check if the account is linked
	nkUserId, _, _, err = nk.AuthenticateDevice(ctx, loginRequest.DeviceId().Token(), "", false)
	if err != nil { // account is missing, this is okay.
		logger.WithField("err", err).WithField("device_auth_token", loginRequest.DeviceId().Token()).Debug("Device not linked.")
		err = nil
	}

	// If the account is not linked, create a link ticket and return an error
	if nkUserId == "" {
		// No Account. Create link ticket and return error
		linkTicket, err := loginRequest.LinkTicket(serviceContext, SystemUserId)
		if err != nil {
			logger.WithField("err", err).Error("unable to generate link ticket.")
			return nil, runtime.NewError(fmt.Sprintf("unable to generate link ticket: %q", UserIdToken), StatusPermissionDenied)
		}

		logger.WithField("linkTicket", linkTicket).Debug("Link ticket found/generated.")
		// Return the link ticket to the client
		return nil, runtime.NewError(fmt.Sprintf("visit %s and enter code: %s", linkingPageUrl, linkTicket.Code), StatusInvalidArgument)

	}

	// Authorize the authenticated account
	account, err := nk.AccountGetId(ctx, nkUserId)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("unable to get account for Id: %q", UserIdToken), StatusInternalError)
	}

	// Check if the account is disabled/banned
	if account.GetDisableTime() != nil {
		return nil, runtime.NewError(fmt.Sprintf("account Permanently Banned: %q", UserIdToken), StatusPermissionDenied)
	}

	// Authenticate if the accounts has a password set (i.e. an email is set)
	if account.Email != "" {
		_, _, _, err = nk.AuthenticateEmail(ctx, account.Email, authPassword, "", false)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("invalid password for account: %q", UserIdToken), StatusUnauthenticated)
		}

	} else if authPassword != "" {
		// if the login contains a password, but there is no password set. set the password.
		err = nk.LinkEmail(ctx, nkUserId, account.User.Id+"@"+placeholderEmailDomain, authPassword)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("unable to set password for account: %q", UserIdToken), StatusInternalError)
		}
	}

	if account.CustomId == "" {
		// if the account does not have a customId, the account needs to be linked to discord.
		// return nothing, and let the client know that they need to link their account
		return nil, runtime.NewError(fmt.Sprintf("Re-link %s at %s", UserIdToken, linkingPageUrl), StatusInternalError)
	}

	/*
		// The tokens expire too quickly to use this method
		// get the discord access token from storage
		accessToken, err := ReadAccessTokenFromStorage(ctx, logger, nk, account.User.Id, vars["DISCORD_CLIENT_ID"], vars["DISCORD_CLIENT_SECRET"])
		if err != nil {
			logger.Warn("error reading discord access token from storage: %v", err)
			return nil, runtime.NewError("error reading discord access token from storage", StatusInternalError)
		}
		if accessToken == nil {
			return nil, runtime.NewError(fmt.Sprintf("Re-link Discord at %s", linkingPageUrl), StatusUnauthenticated)
		}

		// Refresh the access token
		if err := accessToken.Refresh(vars["DISCORD_CLIENT_ID"], vars["DISCORD_CLIENT_SECRET"]); err != nil {
			logger.Warn("error refreshing DiscordAccessToken: %v", err)
			nk.UnlinkCustom(ctx, account.User.Id, account.CustomId)
			return nil, runtime.NewError("error refreshing DiscordAccessToken", StatusUnauthenticated)
		}

		// Write the refreshed token to storage
		if err := WriteAccessTokenToStorage(ctx, logger, nk, account.User.Id, accessToken); err != nil {
			logger.Warn("error writing DiscordAccessToken to storage: %v", err)
			return nil, runtime.NewError("error writing DiscordAccessToken to storage", StatusInternalError)
		}
	*/

	// Use the discordbot to get the guild members ID
	// Get the Discord guildMember

	botGuildId := vars["DISCORD_BOT_GUILD"]
	guildMember, err := discordBot.GuildMember(botGuildId, account.User.Username)
	if err != nil {
		logger.Warn("error getting guild member: %v", err)
		return nil, runtime.NewError("error getting guild member", StatusInternalError)
	}

	// if the nakama custom id isn't composed of only numbers, then update the customId to be the discord ID
	// this is to support legacy accounts that were created before the discord ID was used as the customId
	// use a regular expression to check if the customId is only numbers
	// if it is, then it is a discord ID, and we don't need to update it
	// if it is not, then it is a legacy customId, and we need to update it to the discord ID
	re := regexp.MustCompile("^[0-9]+$")
	if !re.Match([]byte(account.CustomId)) {
		logger.Warn("Migrating legacy account to discord ID")
		nk.UnlinkCustom(ctx, account.User.Id, account.CustomId)
		nk.LinkCustom(ctx, account.User.Id, guildMember.User.ID)
	}

	displayName := DetermineDisplayName(account, guildMember.User, guildMember)
	// Check if the displayName matches an existing nakama username
	// TODO: message the user that their trying to use a name that is in use.
	// This is to prevent duplicate displayNames
	/*
		users, err := nk.UsersGetUsername(ctx, []string{displayName})
		if err != nil {
			logger.WithField("err", err).Error("Users get username error.")
		} else {
			if len(users) > 0 {
				logger.Warn("displayName: %s already exists as a username. Setting displayName to discord username: %s", displayName, guildMember.User.Username)
				// Add 3 random digits to teh end of the users name to make it unique
				// check if it is in use
				// it is, loop and try again
				// it is not, set the name
				displayName = guildMember.User.Username + strconv.Itoa(rand.Intn(999))
				}

				displayName = guildMember.User.Username
			}
		}
	*/

	// Update the Nakama user
	if err := nk.AccountUpdateId(ctx, account.User.Id, "", nil, displayName, "", "", "", guildMember.AvatarURL("")); err != nil {
		logger.Warn("error updating nakama user: %v", err)
		return nil, runtime.NewError(fmt.Sprintf("%v", err), StatusInternalError)
	}

	return account, nil
}
