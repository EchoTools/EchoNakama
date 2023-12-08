package login

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"echo-nakama/game"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

const SYSTEM_USER_ID = "00000000-0000-0000-0000-000000000000"

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

// ILoginService defines the interface for the login service.
type ILoginService interface {
	CheckUserSessionValid(session uuid.UUID, userId game.XPlatformID) bool
}

// LoginService is the implementation of the login service.
type LoginService struct {
	userSessions map[uuid.UUID]game.XPlatformID
}

// NewLoginService creates a new instance of LoginService.
func NewLoginService() *LoginService {
	return &LoginService{
		userSessions: make(map[uuid.UUID]game.XPlatformID),
	}
}

// CheckUserSessionValid checks if a provided user session token is valid.
func (ls *LoginService) CheckUserSessionValid(session uuid.UUID, userId game.XPlatformID) bool {
	storedUserId, ok := ls.userSessions[session]
	return ok && userId == storedUserId
}

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
	XPlatformId   game.XPlatformID `json:"xplatform_id"`
	Session       string           `json:"session_guid"`
	Token         string           `json:"session_token"`
	LoginSettings LoginSettings    `json:"login_settings"`
	PlayerData    game.PlayerData  `json:"account_data"`
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

// GenerateLinkCode generates a random link code consisting of 4 characters from the set of valid characters.
// The set of valid characters is defined as "ABCDEFGHJKLMNPQRSTUVWXYZ".
// The random number generator is seeded with the current time to ensure randomness.
// Returns the generated link code as a string.
func GenerateLinkCode() string {
	// Define the set of valid characters for the link code
	characters := "ABCDEFGHJKLMNPQRSTUVWXYZ"

	// Create a new local random generator with a known seed value
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create an array with 4 elements
	var indices []int
	for i := 0; i < 4; i++ {
		indices = append(indices, i)
	}

	// Randomly select an index from the array and generate the code
	var code string
	for _, _ = range indices {
		code += string(characters[rng.Intn(len(characters))])
	}

	return code
}

type LinkTicket struct {
	LinkCode        string `json:"link_code"`
	XplatformIDStr  string `json:"xplatform_id_str"`
	HmdSerialNumber string `json:"hmd_serial_number"`
}

// call GenerateLinkCode and attempt write it to storage
// if it fails, call GenerateLinkCode again
// if it succeeds, return the link code
func GenerateLinkTicket(
	ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, userId string,
	xplatformId game.XPlatformID, hmdSerialNumber string) (*LinkTicket, *runtime.Error) {

	// loop until we have a unique link code
	for {
		linkTicket := LinkTicket{
			LinkCode:        GenerateLinkCode(),
			XplatformIDStr:  xplatformId.String(),
			HmdSerialNumber: hmdSerialNumber,
		}

		linkTicketJson, err := json.Marshal(linkTicket)
		if err != nil {
			return nil, runtime.NewError(fmt.Sprintf("error marshaling client profile playerData: %v", err), INTERNAL)
		}

		// Write the link code to storage
		objectIDs := []*runtime.StorageWrite{
			{
				Collection:      "Login:linkTicket",
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
		return &linkTicket, nil
	}
}

// ProcessLoginRequest processes a LoginRequest.
func ProcessLoginRequest(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, request *LoginRequest) (string, *runtime.Error) {
	// If we have existing session data for this peer's connection, invalidate it.
	// Note: The client may have multiple connections, represented as different peers.
	// This only invalidates the current connection prior to accepting a new login.
	//InvalidatePeerUserSession(sender)

	relayUserId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		logger.Error("must authenticate before processing request")
		return "", runtime.NewError("relay did not authenticate first", UNAUTHENTICATED)
	}

	relayUserName, ok := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)
	if !ok {
		logger.Error("must authenticate before processing request")
		return "", runtime.NewError("relay did not authenticate first", UNAUTHENTICATED)
	}
	logger.WithField("relayUserName", relayUserName).WithField("request", request).Debug("Processing login request.")

	if !ok {
		logger.Error("must authenticate before processing request")
		return "", runtime.NewError("relay did not authenticate first", UNAUTHENTICATED)
	}

	currentTimestamp := time.Now().UTC().Unix()
	hmdSerialNumber := request.AccountInfo.Hmdserialnumber
	xplatformIdStr := request.XPlatformId.String()
	sessionGuid := uuid.New()
	var nkUserId string
	var account *api.Account

	var username string
	var err error
	var isOvrClient bool = true

	// Validate the user identifier
	if !request.XPlatformId.Valid() {
		return "", runtime.NewError(fmt.Sprintf("Invalid User Identifier\n'%v' is invalid.", xplatformIdStr), INVALID_ARGUMENT)
	}

	if hmdSerialNumber == "N/A" {
		isOvrClient = false
	} else if request.AccountInfo.AccountId > 0 {
		isOvrClient = false
	} else if !regexp.MustCompile("^[A-Z0-9]{8,12}$").MatchString(hmdSerialNumber) {
		isOvrClient = false
	}

	if request.XPlatformId.AccountId == 3963667097037078 {
		isOvrClient = false
	}
	if isOvrClient {
		logger.WithField("hmdSerialNumber", hmdSerialNumber).Debug("OVR Client detected, authenticating with HMD Serial Number")
		nkUserId, username, _, err = nk.AuthenticateDevice(ctx, hmdSerialNumber, "", false)

	} else {
		logger.WithField("xplatform_id", xplatformIdStr).Debug("Non-OVR Client detected, authenticating with XPlatformId")
		nkUserId, username, _, err = nk.AuthenticateDevice(ctx, xplatformIdStr, "", false)

	}
	if err != nil {
		// No Account. Create link ticket and return error
		logger.
			WithField("err", err).
			WithField("isOvrClient", isOvrClient).
			WithField("xplatformIdStr", xplatformIdStr).
			WithField("hmdSerialNumber", hmdSerialNumber).
			Error("Unable to authenticate device.")
		linkTicket, err := GenerateLinkTicket(ctx, logger, db, nk, SYSTEM_USER_ID, request.XPlatformId, hmdSerialNumber)
		if err != nil {
			logger.WithField("err", err).Error("Unable to generate link ticket.")
			return "", runtime.NewError(fmt.Sprintf("Unable to generate link ticket: %s", xplatformIdStr), INTERNAL)
		}
		logger.WithField("linkTicket", linkTicket).Debug("Generated link ticket.")

		// TODO: allow custom linking urls
		return "", runtime.NewError(fmt.Sprintf("Visit https://echovrce.com/link and enter code: %s", linkTicket.LinkCode), NOT_FOUND)

	}

	// Generate a session token with the Guid
	token, _, err := nk.AuthenticateTokenGenerate(nkUserId, username, 0, map[string]string{"sessionGuid": sessionGuid.String()})
	if err != nil {
		logger.WithField("err", err).Error("Authenticate token generate error.")
		return "", runtime.NewError("Authenticate token generate error.", INTERNAL)
	}

	account, err = nk.AccountGetId(ctx, nkUserId)
	if err != nil {
		logger.WithField("err", err).Error("Unable to get account: %s", xplatformIdStr)
		return "", runtime.NewError(fmt.Sprintf("Unable to get account: %s", xplatformIdStr), INTERNAL)
	}

	if account.GetDisableTime() != nil {
		return "", runtime.NewError(fmt.Sprintf("Account Permanently Banned: %s", xplatformIdStr), PERMISSION_DENIED)
	}

	// generate a blank playerData object
	playerData := game.DefaultPlayerData(request.XPlatformId, request.AccountInfo.Displayname)

	// read the client profile from the storage layer

	objectIds := []*runtime.StorageRead{{
		Collection: "Profile",
		Key:        "client",
		UserID:     nkUserId,
	}, {
		Collection: "Profile",
		Key:        "server",
		UserID:     nkUserId,
	}, {
		Collection: "Login:login_settings",
		Key:        "login_settings",
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
			if record.Key == "client" {
				err = json.Unmarshal([]byte(record.Value), &playerData.Profile.Client)
				if err != nil {
					return "", runtime.NewError(fmt.Sprintf("error unmarshaling client playerData: %v", err), INTERNAL)
				}
			} else if record.Key == "server" {
				err = json.Unmarshal([]byte(record.Value), &playerData.Profile.Server)
				if err != nil {
					return "", runtime.NewError(fmt.Sprintf("error unmarshaling server playerData: %v", err), INTERNAL)
				}
			} else if record.Key == "login_settings" {
				err = json.Unmarshal([]byte(record.Value), &loginSettings)
				if err != nil {
					return "", runtime.NewError(fmt.Sprintf("error unmarshaling server playerData: %v", err), INTERNAL)
				}
			}
		}
	}

	// Update the server profile's logintime and updatetime.
	playerData.Profile.Server.LobbyVersion = int(request.AccountInfo.LobbyVersion)
	playerData.Profile.Server.LoginTime = currentTimestamp
	playerData.Profile.Server.ModifyTime = account.User.UpdateTime.Seconds
	playerData.Profile.Server.UpdateTime = account.User.UpdateTime.Seconds
	playerData.Profile.Server.DisplayName = account.User.DisplayName
	playerData.Profile.Client.DisplayName = account.User.DisplayName

	// Write the profile data to storage
	jsonAccountInfo, err := request.AccountInfo.Marshal()
	if err != nil {
		logger.WithField("err", err).Error("Invalid account info.")
	}
	clientProfileJson, err := json.Marshal(playerData.Profile.Client)
	if err != nil {
		return "", runtime.NewError(fmt.Sprintf("error marshaling client profile playerData: %v", err), INTERNAL)
	}
	serverProfileJson, err := json.Marshal(playerData.Profile.Server)
	if err != nil {
		return "", runtime.NewError(fmt.Sprintf("error marshaling server profile playerData: %v", err), INTERNAL)
	}

	// Write the account info to storage using the XplatformID
	objectIDs := []*runtime.StorageWrite{
		{
			Collection:      "Xplatformid",
			Key:             request.XPlatformId.String(),
			UserID:          nkUserId,
			Value:           string(jsonAccountInfo),
			PermissionRead:  2,
			PermissionWrite: 1,
		},
		{
			Collection:      "Profile",
			Key:             "client",
			UserID:          nkUserId,
			Value:           string(clientProfileJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		},
		{
			Collection:      "Profile",
			Key:             "server",
			UserID:          nkUserId,
			Value:           string(serverProfileJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		}}

	var loginSettingsJson []byte
	// if loginSettings is empty, create a default loginSettings object, and write it storage

	if loginSettings.Env == "" {
		loginSettings = DefaultLoginSettings()
		loginSettingsJson, err = json.Marshal(loginSettings)
		if err != nil {
			return "", runtime.NewError(fmt.Sprintf("error marshalling LoginSettings: %v", err), INTERNAL)
		}

		objectIDs = append(objectIDs, &runtime.StorageWrite{
			Collection:      "Login:login_settings",
			Key:             "login_settings",
			UserID:          relayUserId,
			Value:           string(loginSettingsJson),
			PermissionRead:  2,
			PermissionWrite: 0,
		})
	}

	_, err = nk.StorageWrite(ctx, objectIDs)
	if err != nil {
		logger.WithField("err", err).Error("Storage write error.")
		return "", runtime.NewError(fmt.Sprintf("error writing profile data: %v", err), INTERNAL)
	}

	loginSuccess := LoginSuccess{
		XPlatformId:   request.XPlatformId,
		Session:       sessionGuid.String(),
		Token:         token,
		LoginSettings: loginSettings,
		PlayerData:    playerData,
	}

	logger.WithField("loginSuccess", loginSuccess).Debug("Login Success.")

	loginSuccessJson, err := json.Marshal(loginSuccess)
	if err != nil {
		return "", runtime.NewError(fmt.Sprintf("error marshalling LoginSuccess response: %v", err), INTERNAL)
	}

	return string(loginSuccessJson), nil
}
