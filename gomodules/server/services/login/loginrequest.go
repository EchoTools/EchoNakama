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

func UnmarshalLoginAccountInfo(data []byte) (LoginAccountInfo, error) {
	var r LoginAccountInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LoginAccountInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LoginRequest struct {
	AccountInfo LoginAccountInfo `json:"account_info"`
	Session     uuid.UUID        `json:"session_guid"`
	XPlatformId game.XPlatformID `json:"xplatform_id"`
}

type LoginAccountInfo struct {
	AccountId                   int64            `json:"accountid"`
	Displayname                 string           `json:"displayname"`
	Bypassauth                  bool             `json:"bypassauth"`
	AccessToken                 string           `json:"access_token"`
	Nonce                       string           `json:"nonce"`
	Buildversion                int64            `json:"buildversion"`
	LobbyVersion                int64            `json:"lobbyversion"`
	Appid                       int64            `json:"appid"`
	PublisherLock               string           `json:"publisher_lock"`
	Hmdserialnumber             string           `json:"hmdserialnumber"`
	Desiredclientprofileversion int64            `json:"desiredclientprofileversion"`
	GameSettings                GameSettings     `json:"game_settings"`
	SystemInfo                  SystemInfo       `json:"system_info"`
	GraphicsSettings            GraphicsSettings `json:"graphics_settings"`
}

type GameSettings struct {
	Experimentalthrowing int64   `json:"experimentalthrowing"`
	Smoothrotationspeed  float32 `json:"smoothrotationspeed"`
	Hud                  bool    `json:"HUD"`
	Voipmode             int64   `json:"voipmode"`
	MatchTagDisplay      bool    `json:"MatchTagDisplay"`
	EnableNetStatusHUD   bool    `json:"EnableNetStatusHUD"`
	EnableGhostAll       bool    `json:"EnableGhostAll"`
	EnablePitch          bool    `json:"EnablePitch"`
	EnablePersonalBubble bool    `json:"EnablePersonalBubble"`
	Releasedistance      float32 `json:"releasedistance"`
	EnableYaw            bool    `json:"EnableYaw"`
	EnableNetStatusPause bool    `json:"EnableNetStatusPause"`
	Dynamicmusicmode     int64   `json:"dynamicmusicmode"`
	EnableRoll           bool    `json:"EnableRoll"`
	EnableMuteAll        bool    `json:"EnableMuteAll"`
	EnableMaxLoudness    bool    `json:"EnableMaxLoudness"`
	EnableSmoothRotation bool    `json:"EnableSmoothRotation"`
	EnableAPIAccess      bool    `json:"EnableAPIAccess"`
	EnableMuteEnemyTeam  bool    `json:"EnableMuteEnemyTeam"`
	EnableVoipLoudness   bool    `json:"EnableVoipLoudness"`
	Voiploudnesslevel    int64   `json:"voiploudnesslevel"`
	Voipmodeffect        int64   `json:"voipmodeffect"`
	Personalbubblemode   float32 `json:"personalbubblemode"`
	Announcer            int64   `json:"announcer"`
	Ghostallmode         int64   `json:"ghostallmode"`
	Music                int64   `json:"music"`
	Personalbubbleradius float32 `json:"personalbubbleradius"`
	Sfx                  int64   `json:"sfx"`
	Voip                 int64   `json:"voip"`
	Grabdeadzone         float32 `json:"grabdeadzone"`
	Wristangleoffset     float32 `json:"wristangleoffset"`
	Muteallmode          int64   `json:"muteallmode"`
}

type GraphicsSettings struct {
	Temporalaa                 bool    `json:"temporalaa"`
	Fullscreen                 bool    `json:"fullscreen"`
	Display                    int64   `json:"display"`
	Adaptiverestargetframerate int64   `json:"adaptiverestargetframerate"`
	Adaptiveresmaxscale        float32 `json:"adaptiveresmaxscale"`
	Adaptiveresolution         bool    `json:"adaptiveresolution"`
	Adaptiveresminscale        float32 `json:"adaptiveresminscale"`
	Resolutionscale            float32 `json:"resolutionscale"`
	Qualitylevel               int64   `json:"qualitylevel"`
	Adaptiveresheadroom        float32 `json:"adaptiveresheadroom"`
	Quality                    Quality `json:"quality"`
	Msaa                       int64   `json:"msaa"`
	Sharpening                 float32  `json:"sharpening"`
	Multires                   bool    `json:"multires"`
	Gamma                      float32 `json:"gamma"`
	Capturefov                 float32   `json:"capturefov"`
}

type Quality struct {
	Shadowresolution   int64 `json:"shadowresolution"`
	Fx                 int64 `json:"fx"`
	Bloom              bool  `json:"bloom"`
	Cascaderesolution  int64 `json:"cascaderesolution"`
	Cascadedistance    float32 `json:"cascadedistance"`
	Textures           int64 `json:"textures"`
	Shadowmsaa         int64 `json:"shadowmsaa"`
	Meshes             int64 `json:"meshes"`
	Shadowfilterscale  float32 `json:"shadowfilterscale"`
	Staggerfarcascades bool  `json:"staggerfarcascades"`
	Volumetrics        bool  `json:"volumetrics"`
	Lights             int64 `json:"lights"`
	Shadows            int64 `json:"shadows"`
	Anims              int64 `json:"anims"`
}

type SystemInfo struct {
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
