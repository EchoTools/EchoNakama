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
	SystemUserID = "00000000-0000-0000-0000-000000000000"

	PasswordURLParam          = "password"
	HMDSerialOverrideURLParam = "hmdsn"

	LinkingPageUrl = "https://echovrce.com/link"

	LinkTicketCollection           = "Login:linkTicket"
	LinkTicketIndex                = "Index_" + LinkTicketCollection
	DiscordAccessTokenCollection   = "Login:discordAccessToken"
	LoginSettingsStorageCollection = "Login:login_settings"
	LoginSettingsStorageKey        = "login_settings"
	GameProfileStorageCollection   = "Profile"
	ServerGameProfileStorageKey    = "server"
	ClientGameProfileStorageKey    = "client"
	XPlatformIDStorageCollection   = "XPlatformID"

	// The Application ID for Echo VR
	AppNoOVR = 0
	AppQuest = 2215004568539258
	AppPCVR  = 1369078409873402

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
	XPlatformID   game.XPlatformID  `json:"xplatform_id"`
	DeviceIDStr   string            `json:"device_id_str"`
	SessionGUID   string            `json:"session_guid"`
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
	XPlatformIDStr string        `json:"xplatform_id_str"`
	LoginRequest   *LoginRequest `json:"login_request"`
}

func (l *LinkTicket) StorageObject() (*runtime.StorageWrite, error) {
	linkTicketJson, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return &runtime.StorageWrite{
		Collection:      LinkTicketCollection,
		Key:             l.LinkCode,
		UserID:          SystemUserID,
		Value:           string(linkTicketJson),
		PermissionRead:  0,
		PermissionWrite: 0,
		Version:         "*",
	}, nil
}

type DeviceID struct {
	AppID           int64  `json:"app_id"`
	XPlatformIDStr  string `json:"xplatform_id_str"`
	HMDSerialNumber string `json:"hmd_serial_number"`
}

// Generate the string used for device authentication
// WARNING: If this is changed, then device "links" will be invalidated
func (d DeviceID) String() string {
	return fmt.Sprintf("%d:%s:%s", d.AppID, d.XPlatformIDStr, d.HMDSerialNumber)
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
func (request *LoginRequest) LinkTicket(serviceContext *services.ServiceContext, userID string) (*LinkTicket, *runtime.Error) {
	linkTicket := &LinkTicket{}
	ctx := serviceContext.Ctx
	nk := serviceContext.NakamaModule
	logger := serviceContext.Logger
	// Check if a link ticket already exists for the provided xplatformId and hmdSerialNumber
	objectIDs, err := nk.StorageIndexList(ctx, SystemUserID, LinkTicketIndex, fmt.Sprintf("+value.xplatform_id_str:%s", request.DeviceId().XPlatformIDStr), 10)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("error listing link tickets: `%q`  %v", "+value.xplatform_id_str:"+request.DeviceId().XPlatformIDStr, err), StatusInternalError)
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
			LinkCode:       GenerateLinkCode(),
			DeviceIdStr:    request.DeviceId().String(),
			XPlatformIDStr: request.DeviceId().XPlatformIDStr,
			LoginRequest:   request,
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
	var nkUserId string
	var err error
	xPlatformIDStr := loginRequest.XPlatformID.String()

	// Validate the user identifier
	if !loginRequest.XPlatformID.Valid() {
		return nil, runtime.NewError(fmt.Sprintf("Invalid XPlatformId: %q", xPlatformIDStr), StatusInvalidArgument)
	}

	nkUserId, _, _, err = nk.AuthenticateDevice(ctx, loginRequest.DeviceId().String(), "", false)
	if err != nil { // account is missing, this is okay.
		logger.WithField("err", err).WithField("device_id_str", loginRequest.DeviceId().String()).Debug("Device not linked.")
		err = nil
	}

	if nkUserId == "" {
		// No Account. Create link ticket and return error
		linkTicket, err := loginRequest.LinkTicket(serviceContext, SystemUserID)
		if err != nil {
			logger.WithField("err", err).Error("Unable to generate link ticket.")
			return nil, runtime.NewError(fmt.Sprintf("Unable to generate link ticket: %q", xPlatformIDStr), StatusInternalError)
		}
		logger.WithField("linkTicket", linkTicket).Debug("Link ticket found/generated.")
		return nil, runtime.NewError(fmt.Sprintf("Visit %s and enter code: %s", LinkingPageUrl, linkTicket.LinkCode), StatusNotFound)

	}

	// Authorize the authenticated account
	account, err := nk.AccountGetId(ctx, nkUserId)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("Unable to get account for Id: %q", xPlatformIDStr), StatusInternalError)
	}

	// Check if the account is disabled/banned
	if account.GetDisableTime() != nil {
		return nil, runtime.NewError(fmt.Sprintf("Account Permanently Banned: %q", xPlatformIDStr), StatusPermissionDenied)
	}

	// If the account has a password set, authenticate the 'auth=' query param as the password

	if account.Email != "" {

		_, _, _, err = nk.AuthenticateEmail(ctx, account.Email, authPassword, "", false)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("Invalid password for account: %q", xPlatformIDStr), StatusUnauthenticated)
		}
	} else if authPassword != "" {
		// if the 'auth=' query param is set, set the password to the 'auth=' query param
		nk.LinkEmail(ctx, nkUserId, account.User.Id+"@null.echovrce.com", authPassword)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("Unable to set password for account: %q", xPlatformIDStr), StatusInternalError)
		}
	}

	if account.CustomID == "" {
		// if the account does not have a customId, return nothing, and let the client know that they need to link their account
		return nil, runtime.NewError(fmt.Sprintf("Account (%q) not linked to Discord. visit %s", xPlatformIDStr, LinkingPageUrl), StatusInternalError)
	}
	// verify that the account has a valid customId by authenticating to it (this activates the validation/refresh hook)
	// if the account does not have a valid customId, this will return an error
	_, _, _, err = nk.AuthenticateCustom(ctx, account.CustomID, "", false)
	if err != nil {
		return nil, runtime.NewError(fmt.Sprintf("Discord link is invalid. visit %s", LinkingPageUrl), StatusInternalError)
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

	relayUserID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
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
		logger.WithField("nkerr", nkerr).Error("Authentication Errored.")
		return nil, nkerr
	}

	// Authorize the client to use the authenticated account
	nkUserID := account.User.Id
	currentTimestamp := time.Now().UTC().Unix()
	sessionGuid := uuid.New()

	// Generate a session token with the Guid
	token, _, err := nk.AuthenticateTokenGenerate(account.User.Id, account.User.Username, 0, map[string]string{"sessionGuid": sessionGuid.String()})
	if err != nil {
		logger.WithField("err", err).Error("Authenticate token generate error.")
		return nil, runtime.NewError("Authenticate token generation error.", StatusInternalError)
	}

	// generate a blank playerData object
	gameProfiles := game.DefaultGameProfiles(request.XPlatformID, request.LoginData.DisplayName)

	// read the client profile from the storage layer
	// TODO: Extact method
	objectIds := []*runtime.StorageRead{{
		Collection: GameProfileStorageCollection,
		Key:        ClientGameProfileStorageKey,
		UserID:     nkUserID,
	}, {
		Collection: GameProfileStorageCollection,
		Key:        ServerGameProfileStorageKey,
		UserID:     nkUserID,
	}, {
		Collection: LoginSettingsStorageCollection,
		Key:        LoginSettingsStorageKey,
		UserID:     relayUserID,
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
			} else if record.Key == LoginSettingsStorageKey {
				err = json.Unmarshal([]byte(record.Value), &loginSettings)
				if err != nil {
					return nil, runtime.NewError(fmt.Sprintf("error unmarshaling server playerData: %v", err), StatusInternalError)
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
			Collection:      XPlatformIDStorageCollection,
			Key:             request.DeviceId().XPlatformIDStr,
			UserID:          nkUserID,
			Value:           string(jsonRequest),
			PermissionRead:  0,
			PermissionWrite: 0,
		},
		{
			Collection:      GameProfileStorageCollection,
			Key:             ClientGameProfileStorageKey,
			UserID:          nkUserID,
			Value:           string(clientProfileJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		},
		{
			Collection:      GameProfileStorageCollection,
			Key:             ServerGameProfileStorageKey,
			UserID:          nkUserID,
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
		loginSettings = DefaultLoginSettings()
		loginSettingsJson, err = json.Marshal(loginSettings)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("error marshalling LoginSettings: %v", err), StatusInternalError)
		}

		objectIDs = append(objectIDs, &runtime.StorageWrite{
			Collection:      LoginSettingsStorageCollection,
			Key:             LoginSettingsStorageKey,
			UserID:          relayUserID,
			Value:           string(loginSettingsJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		})
	}

	_, err = nk.StorageWrite(ctx, objectIDs)
	if err != nil {
		logger.WithField("err", err).Error("Storage write error.")
		return nil, runtime.NewError(fmt.Sprintf("error writing profile data: %v", err), StatusInternalError)
	}

	loginSuccess := LoginSuccess{
		XPlatformID:   request.XPlatformID,
		DeviceIDStr:   request.DeviceId().String(),
		SessionGUID:   sessionGuid.String(),
		SessionToken:  token,
		LoginSettings: loginSettings,
		GameProfiles:  gameProfiles,
	}

	logger.WithField("loginSuccess", loginSuccess).Debug("Login Success.")
	return &loginSuccess, nil
}
