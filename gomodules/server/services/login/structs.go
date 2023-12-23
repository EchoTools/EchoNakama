package login

import (
	"echonakama/game"
	"echonakama/server/services"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
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
