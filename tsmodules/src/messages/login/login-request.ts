import { XPlatformId } from "../../game/x-platform-id";

class LoginRequest {
    // The unique 64-bit symbol denoting the type of message.
public messageTypeSymbol: bigint = BigInt('-4777159589668118518');

// The user's session token.
public session: string;

// The user identifier.
public userId: XPlatformId;

// The client account information supplied for the sign-in request.
public accountInfo: LoginAccountInfo;

    // Initializes a new LoginRequest message.
    constructor(session?: string, userId?: XPlatformId, accountData?: LoginAccountInfo) {
        this.accountInfo = new LoginAccountInfo();
        this.session = session || '';
        this.userId = userId || new XPlatformId();
        this.accountInfo = accountData || new LoginAccountInfo();
    }

}
export interface LoginRequest {
    accountid:                   number;
    displayname:                 null;
    bypassauth:                  boolean;
    access_token:                string;
    nonce:                       string;
    buildversion:                number;
    lobbyversion:                number;
    appid:                       number;
    publisher_lock:              string;
    hmdserialnumber:             string;
    desiredclientprofileversion: number;
    game_settings:               GameSettings;
    system_info:                 SystemInfo;
    graphics_settings:           GraphicsSettings;
}

export interface GameSettings {
    experimentalthrowing: number;
    smoothrotationspeed:  number;
    HUD:                  boolean;
    voipmode:             number;
    MatchTagDisplay:      boolean;
    EnableNetStatusHUD:   boolean;
    EnableGhostAll:       boolean;
    EnablePitch:          boolean;
    EnablePersonalBubble: boolean;
    releasedistance:      number;
    EnableYaw:            boolean;
    EnableNetStatusPause: boolean;
    dynamicmusicmode:     number;
    EnableRoll:           boolean;
    EnableMuteAll:        boolean;
    EnableMaxLoudness:    boolean;
    EnableSmoothRotation: boolean;
    EnableAPIAccess:      boolean;
    EnableMuteEnemyTeam:  boolean;
    EnableVoipLoudness:   boolean;
    voiploudnesslevel:    number;
    voipmodeffect:        number;
    personalbubblemode:   number;
    announcer:            number;
    ghostallmode:         number;
    music:                number;
    personalbubbleradius: number;
    sfx:                  number;
    voip:                 number;
    grabdeadzone:         number;
    wristangleoffset:     number;
    muteallmode:          number;
}

export interface GraphicsSettings {
    temporalaa:                 boolean;
    fullscreen:                 boolean;
    display:                    number;
    adaptiverestargetframerate: number;
    adaptiveresmaxscale:        number;
    adaptiveresolution:         boolean;
    adaptiveresminscale:        number;
    resolutionscale:            number;
    qualitylevel:               number;
    adaptiveresheadroom:        number;
    quality:                    Quality;
    msaa:                       number;
    sharpening:                 number;
    multires:                   boolean;
    gamma:                      number;
    capturefov:                 number;
}

export interface Quality {
    shadowresolution:   number;
    fx:                 number;
    bloom:              boolean;
    cascaderesolution:  number;
    cascadedistance:    number;
    textures:           number;
    shadowmsaa:         number;
    meshes:             number;
    shadowfilterscale:  number;
    staggerfarcascades: boolean;
    volumetrics:        boolean;
    lights:             number;
    shadows:            number;
    anims:              number;
}

export interface SystemInfo {
    headset_type:         string;
    driver_version:       string;
    network_type:         string;
    video_card:           string;
    cpu:                  string;
    num_physical_cores:   number;
    num_logical_cores:    number;
    memory_total:         number;
    memory_used:          number;
    dedicated_gpu_memory: number;
}
class LoginAccountInfo {
    // The account identifier.

    public accountId: number = 0;

    // The user's display name.

    public displayName?: string;

    // TODO: Unknown
    @JsonProperty("bypassauth")
    public bypassAuth?: boolean = false;

    // Oculus-related access token for authentication.
    @JsonProperty("access_token")
    public accessToken?: string;

    // Authentication-related nonce.
    @JsonProperty("nonce")
    public nonce?: string;

    // Client build version.
    @JsonProperty("buildversion")
    public buildVersion: number = 0;

    // The lobby build timestamp.
    @JsonProperty("lobbyversion")
    public lobbyVersion?: number;

    // The identifier for the application.
    @JsonProperty("appid")
    public appId?: number;

    // An environment lock for different sandboxes.
    @JsonProperty("publisher_lock")
    public publisherLock?: string = 'rad15_live';

    // Headset serial number
    @JsonProperty("hmdserialnumber")
    public hmdSerialNumber?: string;

    // Requested version for clients to receive their profile in.
    @JsonProperty("desiredclientprofileversion")
    public desiredClientProfileVersion?: number;

    // Additional fields which are not caught explicitly are retained here.
    @JsonExtensionData
    public additionalData: Record<string, any> = {};
}

export { LoginRequest, LoginAccountInfo };