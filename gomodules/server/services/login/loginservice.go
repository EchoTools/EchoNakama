package login

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"

	"echo-nakama/game"
	"echo-nakama/server/services"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	SYSTEM_USER_ID                  = "00000000-0000-0000-0000-000000000000"
	LINKTICKET_COLLECTION           = "Login:linkTicket"
	LINKTICKET_INDEX                = "Index_" + LINKTICKET_COLLECTION
	DISCORD_ACCESS_TOKEN_COLLECTION = "Login:discordAccessToken"

	PASSWORD_QUERY_PARAM = "auth"
	LINK_URL             = "https://echovrce.com/link"

	LOGIN_SETTINGS_COLLECTION = "Login:login_settings"
	LOGIN_SETTINGS_KEY        = "login_settings"
	PROFILE_COLLECTION        = "Profile"
	PROFILE_SERVER_KEY        = "server"
	PROFILE_CLIENT_KEY        = "client"
	XPLATFORMID_COLLECTION    = "Xplatformid"
	LOGIN_REQUEST_COLLECTION  = "Login:login_request"
)

const (
	AUTH_TYPE_FAIL         = iota
	AUTH_TYPE_DEVICE_ID    = iota
	AUTH_TYPE_ACCESS_TOKEN = iota
)

const (
	APPID_NOOVR = 0
	APPID_QUEST = 2215004568539258
	APPID_PCVR  = 1369078409873402
)

// Constants representing different error codes.
const (
	OK                  = 0  // OK indicates a successful operation.
	CANCELED            = 1  // CANCELED indicates the operation was canceled.
	UNKNOWN             = 2  // UNKNOWN indicates an unknown error occurred.
	INVALID_ARGUMENT    = 3  // INVALID_ARGUMENT indicates an invalid argument was provided.
	DEADLINE_EXCEEDED   = 4  // DEADLINE_EXCEEDED indicates the operation exceeded the deadline.
	NOT_FOUND           = 5  // NOT_FOUND indicates the requested resource was not found.
	ALREADY_EXISTS      = 6  // ALREADY_EXISTS indicates the resource already exists.
	PERMISSION_DENIED   = 7  // PERMISSION_DENIED indicates the operation was denied due to insufficient permissions.
	RESOURCE_EXHAUSTED  = 8  // RESOURCE_EXHAUSTED indicates the resource has been exhausted.
	FAILED_PRECONDITION = 9  // FAILED_PRECONDITION indicates a precondition for the operation was not met.
	ABORTED             = 10 // ABORTED indicates the operation was aborted.
	OUT_OF_RANGE        = 11 // OUT_OF_RANGE indicates a value is out of range.
	UNIMPLEMENTED       = 12 // UNIMPLEMENTED indicates the operation is not implemented.
	INTERNAL            = 13 // INTERNAL indicates an internal server error occurred.
	UNAVAILABLE         = 14 // UNAVAILABLE indicates the service is currently unavailable.
	DATA_LOSS           = 15 // DATA_LOSS indicates a loss of data occurred.
	UNAUTHENTICATED     = 16 // UNAUTHENTICATED indicates the request lacks valid authentication credentials.
)

// Assuming you have a struct named RequestInfoType in Go that corresponds to AccountInfo in C#
type RequestInfoType struct {
	HMDSerialNumber string
	AccountID       uint64
}

// Assuming you have a struct named RequestType in Go
type RequestType struct {
	AccountInfo *RequestInfoType
}

type LoginFailure struct {
	UserID game.XPlatformID
	Status int
	Reason string
}

type SessionVars struct {
	SessionGuid uuid.UUID `json:"sessionGuid"`
}

type LoginSuccess struct {
	XPlatformId   game.XPlatformID  `json:"xplatform_id"`
	DeviceIdStr   string            `json:"device_id_str"`
	SessionGuid   string            `json:"session_guid"`
	SessionToken  string            `json:"session_token"`
	LoginSettings LoginSettings     `json:"login_settings"`
	GameProfiles  game.GameProfiles `json:"game_profiles"`
}

// Config represents the structure of the provided JSON.
type LoginSettings struct {
	ConfigData            map[string]interface{} `json:"config_data"`
	Env                   string                 `json:"env"`
	IAPUnlocked           bool                   `json:"iap_unlocked"`
	MatchmakerQueueMode   string                 `json:"matchmaker_queue_mode"`
	RemoteLogErrors       bool                   `json:"remote_log_errors"`
	RemoteLogMetrics      bool                   `json:"remote_log_metrics"`
	RemoteLogRichPresence bool                   `json:"remote_log_rich_presence"`
	RemoteLogSocial       bool                   `json:"remote_log_social"`
	RemoteLogWarnings     bool                   `json:"remote_log_warnings"`
}

// DefaultConfig returns a default Config object.
func DefaultLoginSettings() LoginSettings {
	return LoginSettings{
		ConfigData:            make(map[string]interface{}),
		Env:                   "live",
		IAPUnlocked:           false,
		MatchmakerQueueMode:   "disabled",
		RemoteLogErrors:       false,
		RemoteLogMetrics:      true,
		RemoteLogRichPresence: true,
		RemoteLogSocial:       true,
		RemoteLogWarnings:     false,
	}
}

type DiscordAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// LinkTicket represents a ticket used for linking accounts to Discord.
// It contains the link code, xplatform ID string, and HMD serial number.
type LinkTicket struct {
	LinkCode       string        `json:"link_code"`
	DeviceIdStr    string        `json:"device_id_str"`
	XplatformIdStr string        `json:"xplatform_id_str"`
	LoginRequest   *LoginRequest `json:"login_request"`
}

type DeviceId struct {
	AppId           int64  `json:"appid"`
	XPlatformIdStr  string `json:"xplatform_id_str"`
	HmdSerialNumber string `json:"hmd_serial_number"`
}

func (d DeviceId) String() string {
	return fmt.Sprintf("%d:%s:%s", d.AppId, d.XPlatformIdStr, d.HmdSerialNumber)
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

// call GenerateLinkCode and attempt write it to storage
// if it fails, call GenerateLinkCode again
// if it succeeds, return the link code
func (request *LoginRequest) LinkTicket(serviceContext *services.ServiceContext, userId string) (*LinkTicket, *runtime.Error) {
	linkTicket := &LinkTicket{}
	ctx := serviceContext.Ctx
	nk := serviceContext.NakamaModule
	logger := serviceContext.Logger
	// Check if a link ticket already exists for the provided xplatformId and hmdSerialNumber
	objectIds, err := nk.StorageIndexList(ctx, SYSTEM_USER_ID, LINKTICKET_INDEX, fmt.Sprintf("+value.xplatform_id_str:%s", request.DeviceId().XPlatformIdStr), 10)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("error listing link tickets: `%q`  %v", "+value.xplatform_id_str:"+request.DeviceId().XPlatformIdStr, err), INTERNAL)
	}
	logger.WithField("objectIds", objectIds).Debug("Link ticket found/generated.")
	// Link ticket was found. Return the link ticket.
	if objectIds != nil {
		for _, record := range objectIds.Objects {
			json.Unmarshal([]byte(record.Value), &linkTicket)

			return linkTicket, nil
		}
	}

	// Generate a link code and attempt to write it to storage
	for {
		// loop until we have a unique link code

		linkTicket = &LinkTicket{
			LinkCode:       GenerateLinkCode(),
			DeviceIdStr:    request.DeviceId().String(),
			XplatformIdStr: request.DeviceId().XPlatformIdStr,
			LoginRequest:   request,
		}

		linkTicketJson, err := json.Marshal(linkTicket)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("error marshaling client profile playerData: %v", err), INTERNAL)
		}

		// Write the link code to storage
		objectIDs := []*runtime.StorageWrite{
			{
				Collection:      LINKTICKET_COLLECTION,
				Key:             linkTicket.LinkCode,
				UserID:          SYSTEM_USER_ID,
				Value:           string(linkTicketJson),
				PermissionRead:  0,
				PermissionWrite: 0,
				Version:         "*",
			}}

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
	var nkUserId string
	var err error
	xplatformIdStr := loginRequest.XPlatformId.String()

	// Validate the user identifier
	if !loginRequest.XPlatformId.Valid() {
		return nil, runtime.NewError(fmt.Sprintf("Invalid XPlatformId: %q", xplatformIdStr), INVALID_ARGUMENT)
	}

	nkUserId, _, _, err = nk.AuthenticateDevice(ctx, loginRequest.DeviceId().String(), "", false)
	if err != nil { // account is missing, this is okay.
		logger.WithField("err", err).WithField("device_id_str", loginRequest.DeviceId().String()).Debug("Device not linked.")
		err = nil
	}

	if nkUserId == "" {
		// No Account. Create link ticket and return error
		linkTicket, err := loginRequest.LinkTicket(serviceContext, SYSTEM_USER_ID)
		if err != nil {
			logger.WithField("err", err).Error("Unable to generate link ticket.")
			return nil, runtime.NewError(fmt.Sprintf("Unable to generate link ticket: %q", xplatformIdStr), INTERNAL)
		}
		logger.WithField("linkTicket", linkTicket).Debug("Link ticket found/generated.")
		return nil, runtime.NewError(fmt.Sprintf("Visit %s and enter code: %s", LINK_URL, linkTicket.LinkCode), NOT_FOUND)

	}

	// Authorize the authenticated account
	account, err := nk.AccountGetId(ctx, nkUserId)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("Unable to get account for Id: %q", xplatformIdStr), INTERNAL)
	}

	// Check if the account is disabled/banned
	if account.GetDisableTime() != nil {
		return nil, runtime.NewError(fmt.Sprintf("Account Permanently Banned: %q", xplatformIdStr), PERMISSION_DENIED)
	}

	// If the account has a password set, authenticate the 'auth=' query param as the password

	if account.Email != "" {

		_, _, _, err = nk.AuthenticateEmail(ctx, account.Email, authPassword, "", false)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("Invalid password for account: %q", xplatformIdStr), UNAUTHENTICATED)
		}
	} else if authPassword != "" {
		// if the 'auth=' query param is set, set the password to the 'auth=' query param
		nk.LinkEmail(ctx, nkUserId, account.User.Id+"@null.echovrce.com", authPassword)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("Unable to set password for account: %q", xplatformIdStr), INTERNAL)
		}
	}

	if account.CustomId == "" {
		// if the account does not have a customId, return nothing, and let the client know that they need to link their account
		return nil, runtime.NewError(fmt.Sprintf("Account (%q) not linked to Discord. visit %s", xplatformIdStr, LINK_URL), INTERNAL)
	}
	// verify that the account has a valid customId by authenticating to it (this activates the validation/refresh hook)
	// if the account does not have a valid customId, this will return an error
	_, _, _, err = nk.AuthenticateCustom(ctx, account.CustomId, "", false)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("Discord link is invalid. visit %s", LINK_URL), INTERNAL)
	}

	return account, nil
}

// ProcessLoginRequest processes a login request and returns the login success response or an error.
// It returns a string representing the login success response and a *runtime.Error object if there is an error.
func ProcessLoginRequest(serviceContext *services.ServiceContext, request *LoginRequest) (*LoginSuccess, *runtime.Error) {
	ctx := serviceContext.Ctx
	logger := serviceContext.Logger
	nk := serviceContext.NakamaModule

	// immediately pop the password out of the requests to avoid
	// displaying it in the logs
	authPassword := request.AuthPassword
	request.AuthPassword = ""

	relayUserId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return nil, runtime.NewError("relay must authenticate", UNAUTHENTICATED)
	}

	relayUserName, ok := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)
	if !ok {
		return nil, runtime.NewError("relay must authenticate", UNAUTHENTICATED)
	}

	logger.WithField("relayUserName", relayUserName).WithField("request", request).Debug("Processing login request.")

	account, nkerr := authenticateAccountDevice(serviceContext, request, authPassword)
	if nkerr != nil {
		logger.WithField("nkerr", nkerr).Error("Authentication Errored.")
		return nil, nkerr
	}

	// Authorize the client to use the authenticated account
	nkUserId := account.User.Id
	currentTimestamp := time.Now().UTC().Unix()
	sessionGuid := uuid.New()

	// Generate a session token with the Guid
	token, _, err := nk.AuthenticateTokenGenerate(account.User.Id, account.User.Username, 0, map[string]string{"sessionGuid": sessionGuid.String()})
	if err != nil {
		logger.WithField("err", err).Error("Authenticate token generate error.")
		return nil, runtime.NewError("Authenticate token generation error.", INTERNAL)
	}

	// generate a blank playerData object
	gameProfiles := game.DefaultGameProfiles(request.XPlatformId, request.LoginData.DisplayName)

	// read the client profile from the storage layer
	// TODO: Extact method
	objectIds := []*runtime.StorageRead{{
		Collection: PROFILE_COLLECTION,
		Key:        PROFILE_CLIENT_KEY,
		UserID:     nkUserId,
	}, {
		Collection: PROFILE_COLLECTION,
		Key:        PROFILE_SERVER_KEY,
		UserID:     nkUserId,
	}, {
		Collection: LOGIN_SETTINGS_COLLECTION,
		Key:        LOGIN_SETTINGS_KEY,
		UserID:     relayUserId,
	}, /* {
		Collection: "Login:linking_settings",
		Key:        "linking_settings",
		UserID:     relayUserId,
	}, */
	}

	var loginSettings LoginSettings

	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
	} else {
		for _, record := range records {
			if record.Key == PROFILE_CLIENT_KEY {
				err = json.Unmarshal([]byte(record.Value), &gameProfiles.Client)
				if err != nil {
					return nil, runtime.NewError(fmt.Sprintf("error unmarshaling client playerData: %v", err), INTERNAL)
				}
			} else if record.Key == PROFILE_SERVER_KEY {
				err = json.Unmarshal([]byte(record.Value), &gameProfiles.Server)
				if err != nil {
					return nil, runtime.NewError(fmt.Sprintf("error unmarshaling server playerData: %v", err), INTERNAL)
				}
			} else if record.Key == LOGIN_SETTINGS_KEY {
				err = json.Unmarshal([]byte(record.Value), &loginSettings)
				if err != nil {
					return nil, runtime.NewError(fmt.Sprintf("error unmarshaling server playerData: %v", err), INTERNAL)
				}
			}
		}
	}

	// Update the server profile's logintime and updatetime.
	gameProfiles.Server.LobbyVersion = int(request.LoginData.LobbyVersion)
	gameProfiles.Server.LoginTime = currentTimestamp
	gameProfiles.Server.ModifyTime = account.User.UpdateTime.Seconds
	gameProfiles.Server.UpdateTime = account.User.UpdateTime.Seconds

	gameProfiles.Server.DisplayName = account.User.DisplayName
	gameProfiles.Client.DisplayName = account.User.DisplayName

	// Write the profile data to storage
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		logger.WithField("err", err).Error("Invalid account info.")
		return nil, runtime.NewError(fmt.Sprintf("error marshaling accountInfo: %v", err), INTERNAL)
	}
	clientProfileJson, err := json.Marshal(gameProfiles.Client)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("error marshaling client profile playerData: %v", err), INTERNAL)
	}
	serverProfileJson, err := json.Marshal(gameProfiles.Server)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("error marshaling server profile playerData: %v", err), INTERNAL)
	}

	// Write the latest profile data to storage
	objectIDs := []*runtime.StorageWrite{
		{
			Collection:      XPLATFORMID_COLLECTION,
			Key:             request.DeviceId().XPlatformIdStr,
			UserID:          nkUserId,
			Value:           string(jsonRequest),
			PermissionRead:  0,
			PermissionWrite: 0,
		},
		{
			Collection:      PROFILE_COLLECTION,
			Key:             PROFILE_CLIENT_KEY,
			UserID:          nkUserId,
			Value:           string(clientProfileJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		},
		{
			Collection:      PROFILE_COLLECTION,
			Key:             PROFILE_SERVER_KEY,
			UserID:          nkUserId,
			Value:           string(serverProfileJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		}}

	// Get the login settings from storage
	// TODO: extract method
	var loginSettingsJson []byte
	// if loginSettings is empty, create a default loginSettings object, and write it storage

	if loginSettings.Env == "" {
		loginSettings = DefaultLoginSettings()
		loginSettingsJson, err = json.Marshal(loginSettings)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("error marshalling LoginSettings: %v", err), INTERNAL)
		}

		objectIDs = append(objectIDs, &runtime.StorageWrite{
			Collection:      LOGIN_SETTINGS_COLLECTION,
			Key:             LOGIN_SETTINGS_KEY,
			UserID:          relayUserId,
			Value:           string(loginSettingsJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		})
	}

	_, err = nk.StorageWrite(ctx, objectIDs)
	if err != nil {
		logger.WithField("err", err).Error("Storage write error.")
		return nil, runtime.NewError(fmt.Sprintf("error writing profile data: %v", err), INTERNAL)
	}

	loginSuccess := LoginSuccess{
		XPlatformId:   request.XPlatformId,
		DeviceIdStr:   request.DeviceId().String(),
		SessionGuid:   sessionGuid.String(),
		SessionToken:  token,
		LoginSettings: loginSettings,
		GameProfiles:  gameProfiles,
	}

	logger.WithField("loginSuccess", loginSuccess).Debug("Login Success.")
	return &loginSuccess, nil
}
