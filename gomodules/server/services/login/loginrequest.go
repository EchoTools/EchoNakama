// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    loginAccountInfo, err := UnmarshalLoginAccountInfo(bytes)
//    bytes, err = loginAccountInfo.Marshal()

package login

import (
	"encoding/json"

	"echonakama/game"

	"github.com/google/uuid"
)

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
