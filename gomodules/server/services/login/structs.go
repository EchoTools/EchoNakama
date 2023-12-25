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

func UnmarshalLoginAccountInfo(data []byte) (LoginMetadata, error) {
	var r LoginMetadata
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LoginMetadata) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// This is sent by the relay to login a connected game client/broadcaster
type LoginRequest struct {
	Metadata                LoginMetadata   `json:"metadata"`                   // the metadata sent by the game client
	SessionGuid             uuid.UUID       `json:"echo_session_guid"`          // the session guid sent by the game client
	EchoUserId              game.EchoUserId `json:"echo_user_id"`               // the game user id sent by the game client
	UserPassword            string          `json:"user_password"`              // the password query param set in the config.json
	HmdSerialNumberOverride string          `json:"hmd_serial_number_override"` // the hmd serial number override query param set in the config.json
	DisplayNameOverride     string          `json:"display_name_override"`      // the display name override query param set in the config.json
	ClientIpAddress         string          `json:"client_ip_address"`          // the client ip address
}

// Extract the identifying information used for Device Authentication
// WARNING: If this is changed, then device "links" will be invalidated
func (l *LoginRequest) DeviceId() DeviceId {
	return DeviceId{
		AppId:           l.Metadata.AppId,
		UserIdToken:     l.EchoUserId.String(),
		HmdSerialNumber: l.Metadata.HmdSerialNumber,
	}
}

// This is the payload sent by the game client to the relay

type LoginMetadata struct {
	// WARNING: EchoVR dictates this schema.
	AccountId                   int64            `json:"accountid"`
	DisplayName                 string           `json:"displayname"`
	BypassAuth                  bool             `json:"bypassauth"`
	AccessToken                 string           `json:"access_token"`
	Nonce                       string           `json:"nonce"`
	BuildVersion                int64            `json:"buildversion"`
	LobbyVersion                int64            `json:"lobbyversion"`
	AppId                       int64            `json:"appid"`
	PublisherLock               string           `json:"publisher_lock"`
	HmdSerialNumber             string           `json:"hmdserialnumber"`
	DesiredClientProfileVersion int64            `json:"desiredclientprofileversion"`
	GameSettings                GameSettings     `json:"game_settings"`
	SystemInfo                  SystemInfo       `json:"system_info"`
	GraphicsSettings            GraphicsSettings `json:"graphics_settings"`
}

type GameSettings struct {
	// WARNING: EchoVR dictates this schema.
	Experimentalthrowing int64   `json:"experimentalthrowing"`
	Smoothrotationspeed  float32 `json:"smoothrotationspeed"`
	HUD                  bool    `json:"HUD"`
	VoIPMode             int64   `json:"voipmode"`
	MatchTagDisplay      bool    `json:"MatchTagDisplay"`
	EnableNetStatusHUD   bool    `json:"EnableNetStatusHUD"`
	EnableGhostAll       bool    `json:"EnableGhostAll"`
	EnablePitch          bool    `json:"EnablePitch"`
	EnablePersonalBubble bool    `json:"EnablePersonalBubble"`
	ReleaseDistance      float32 `json:"releasedistance"`
	EnableYaw            bool    `json:"EnableYaw"`
	EnableNetStatusPause bool    `json:"EnableNetStatusPause"`
	DynamicMusicMode     int64   `json:"dynamicmusicmode"`
	EnableRoll           bool    `json:"EnableRoll"`
	EnableMuteAll        bool    `json:"EnableMuteAll"`
	EnableMaxLoudness    bool    `json:"EnableMaxLoudness"`
	EnableSmoothRotation bool    `json:"EnableSmoothRotation"`
	EnableAPIAccess      bool    `json:"EnableAPIAccess"`
	EnableMuteEnemyTeam  bool    `json:"EnableMuteEnemyTeam"`
	EnableVoIPLoudness   bool    `json:"EnableVoipLoudness"`
	VoIPLoudnessLevel    int64   `json:"voiploudnesslevel"`
	VoIPModEffect        int64   `json:"voipmodeffect"`
	PersonalBubbleMode   float32 `json:"personalbubblemode"`
	Announcer            int64   `json:"announcer"`
	GhostAllMode         int64   `json:"ghostallmode"`
	Music                int64   `json:"music"`
	PersonalBubbleRadius float32 `json:"personalbubbleradius"`
	SFX                  int64   `json:"sfx"`
	VoIP                 int64   `json:"voip"`
	GrabDeadZone         float32 `json:"grabdeadzone"`
	WristAngleOffset     float32 `json:"wristangleoffset"`
	MuteAllMode          int64   `json:"muteallmode"`
}

type GraphicsSettings struct {
	// WARNING: EchoVR dictates this schema.
	TemporalAA                        bool    `json:"temporalaa"`
	Fullscreen                        bool    `json:"fullscreen"`
	Display                           int64   `json:"display"`
	ResolutionScale                   float32 `json:"resolutionscale"`
	AdaptiveResolutionTargetFramerate int64   `json:"adaptiverestargetframerate"`
	AdaptiveResolutionMaxScale        float32 `json:"adaptiveresmaxscale"`
	AdaptiveResolution                bool    `json:"adaptiveresolution"`
	AdaptiveResolutionMinScale        float32 `json:"adaptiveresminscale"`
	AdaptiveResolutionHeadroom        float32 `json:"adaptiveresheadroom"`
	QualityLevel                      int64   `json:"qualitylevel"`
	Quality                           Quality `json:"quality"`
	MSAA                              int64   `json:"msaa"`
	Sharpening                        float32 `json:"sharpening"`
	MultiResolution                   bool    `json:"multires"`
	Gamma                             float32 `json:"gamma"`
	CaptureFOV                        float32 `json:"capturefov"`
}

type Quality struct {
	// WARNING: EchoVR dictates this schema.
	ShadowResolution   int64   `json:"shadowresolution"`
	FX                 int64   `json:"fx"`
	Bloom              bool    `json:"bloom"`
	CascadeResolution  int64   `json:"cascaderesolution"`
	CascadeDistance    float32 `json:"cascadedistance"`
	Textures           int64   `json:"textures"`
	ShadowMSAA         int64   `json:"shadowmsaa"`
	Meshes             int64   `json:"meshes"`
	ShadowFilterScale  float32 `json:"shadowfilterscale"`
	StaggerFarCascades bool    `json:"staggerfarcascades"`
	Volumetrics        bool    `json:"volumetrics"`
	Lights             int64   `json:"lights"`
	Shadows            int64   `json:"shadows"`
	Anims              int64   `json:"anims"`
}

type SystemInfo struct {
	// WARNING: EchoVR dictates this schema.
	HeadsetType        string `json:"headset_type"`
	DriverVersion      string `json:"driver_version"`
	NetworkType        string `json:"network_type"`
	VideoCard          string `json:"video_card"`
	CPU                string `json:"cpu"`
	NumPhysicalCores   int64  `json:"num_physical_cores"`
	NumLogicalCores    int64  `json:"num_logical_cores"`
	MemoryTotal        int64  `json:"memory_total"`
	MemoryUsed         int64  `json:"memory_used"`
	DedicatedGPUMemory int64  `json:"dedicated_gpu_memory"`
}
