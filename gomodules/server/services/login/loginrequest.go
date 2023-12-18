// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    loginAccountInfo, err := UnmarshalLoginAccountInfo(bytes)
//    bytes, err = loginAccountInfo.Marshal()

package login

import (
	"encoding/json"

	"echo-nakama/game"

	"github.com/google/uuid"
)

func UnmarshalLoginAccountInfo(data []byte) (LoginData, error) {
	var r LoginData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LoginData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LoginRequest struct {
	LoginData       LoginData        `json:"login_data"`
	Session         uuid.UUID        `json:"session_guid"`
	XPlatformId     game.XPlatformID `json:"xplatform_id"`
	AuthPassword    string           `json:"auth_password"`
	HMDSerialNumber string           `json:"hmd_serial_number"`
}

func (l *LoginRequest) DeviceId() DeviceId {
	return DeviceId{
		AppId:           l.LoginData.AppId,
		XPlatformIdStr:  l.XPlatformId.String(),
		HmdSerialNumber: l.LoginData.HmdSerialNumber,
	}
}

type LoginData struct {
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
	Experimentalthrowing int64   `json:"experimentalthrowing"`
	Smoothrotationspeed  float32 `json:"smoothrotationspeed"`
	Hud                  bool    `json:"HUD"`
	VoipMode             int64   `json:"voipmode"`
	MatchTagDisplay      bool    `json:"MatchTagDisplay"`
	EnableNetStatusHud   bool    `json:"EnableNetStatusHUD"`
	EnableGhostAll       bool    `json:"EnableGhostAll"`
	EnablePitch          bool    `json:"EnablePitch"`
	EnablePersonalBubble bool    `json:"EnablePersonalBubble"`
	Releasedistance      float32 `json:"releasedistance"`
	EnableYaw            bool    `json:"EnableYaw"`
	EnableNetStatusPause bool    `json:"EnableNetStatusPause"`
	DynamicMusicMode     int64   `json:"dynamicmusicmode"`
	EnableRoll           bool    `json:"EnableRoll"`
	EnableMuteAll        bool    `json:"EnableMuteAll"`
	EnableMaxLoudness    bool    `json:"EnableMaxLoudness"`
	EnableSmoothRotation bool    `json:"EnableSmoothRotation"`
	EnableApiAccess      bool    `json:"EnableAPIAccess"`
	EnableMuteEnemyTeam  bool    `json:"EnableMuteEnemyTeam"`
	EnableVoipLoudness   bool    `json:"EnableVoipLoudness"`
	VoipLoudnessLevel    int64   `json:"voiploudnesslevel"`
	VoipModEffect        int64   `json:"voipmodeffect"`
	PersonalBubbleMode   float32 `json:"personalbubblemode"`
	Announcer            int64   `json:"announcer"`
	GhostAllMode         int64   `json:"ghostallmode"`
	Music                int64   `json:"music"`
	PersonalBubbleRadius float32 `json:"personalbubbleradius"`
	Sfx                  int64   `json:"sfx"`
	Voip                 int64   `json:"voip"`
	GrabDeadZone         float32 `json:"grabdeadzone"`
	WristAngleOffset     float32 `json:"wristangleoffset"`
	MuteAllMode          int64   `json:"muteallmode"`
}

type GraphicsSettings struct {
	TemporalAa                 bool    `json:"temporalaa"`
	Fullscreen                 bool    `json:"fullscreen"`
	Display                    int64   `json:"display"`
	AdaptiveResTargetFramerate int64   `json:"adaptiverestargetframerate"`
	AdaptiveResMaxScale        float32 `json:"adaptiveresmaxscale"`
	AdaptiveResolution         bool    `json:"adaptiveresolution"`
	AdaptiveResMinScale        float32 `json:"adaptiveresminscale"`
	ResolutionScale            float32 `json:"resolutionscale"`
	QualityLevel               int64   `json:"qualitylevel"`
	AdaptiveResHeadroom        float32 `json:"adaptiveresheadroom"`
	Quality                    Quality `json:"quality"`
	Msaa                       int64   `json:"msaa"`
	Sharpening                 float32 `json:"sharpening"`
	MultiRes                   bool    `json:"multires"`
	Gamma                      float32 `json:"gamma"`
	CaptureFov                 float32 `json:"capturefov"`
}

type Quality struct {
	ShadowResolution   int64   `json:"shadowresolution"`
	Fx                 int64   `json:"fx"`
	Bloom              bool    `json:"bloom"`
	CascadeResolution  int64   `json:"cascaderesolution"`
	CascadeDistance    float32 `json:"cascadedistance"`
	Textures           int64   `json:"textures"`
	ShadowMsaa         int64   `json:"shadowmsaa"`
	Meshes             int64   `json:"meshes"`
	ShadowFilterScale  float32 `json:"shadowfilterscale"`
	StaggerFarCascades bool    `json:"staggerfarcascades"`
	Volumetrics        bool    `json:"volumetrics"`
	Lights             int64   `json:"lights"`
	Shadows            int64   `json:"shadows"`
	Anims              int64   `json:"anims"`
}

type SystemInfo struct {
	HeadsetType        string `json:"headset_type"`
	DriverVersion      string `json:"driver_version"`
	NetworkType        string `json:"network_type"`
	VideoCard          string `json:"video_card"`
	Cpu                string `json:"cpu"`
	NumPhysicalCores   int64  `json:"num_physical_cores"`
	NumLogicalCores    int64  `json:"num_logical_cores"`
	MemoryTotal        int64  `json:"memory_total"`
	MemoryUsed         int64  `json:"memory_used"`
	DedicatedGpuMemory int64  `json:"dedicated_gpu_memory"`
}
