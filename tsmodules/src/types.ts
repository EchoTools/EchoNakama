export interface AccessControlList {
    allow: string[];
    deny: string[];
  }

  export interface Account {
    profile:           Profile;
    is_moderator:      boolean;
    account_lock_hash: string;
    account_lock_salt: string;
    banned_until?: string;
}

export interface Profile {
    client: Client;
    server: Server;
}

export interface Client {
    displayname:   string;
    xplatformid:   string;
    weapon:        string;
    grenade:       string;
    weaponarm:     number;
    ability:       string;
    legal:         Legal;
    npe:           Npe;
    customization: Customization;
    social:        Social;
    newunlocks:    any[];
}

export interface Customization {
    battlepass_season_poi_version: number;
    new_unlocks_poi_version:       number;
    store_entry_poi_version:       number;
    clear_new_unlocks_version:     number;
}

export interface Legal {
    points_policy_version: number;
    eula_version:          number;
    game_admin_version:    number;
    splash_screen_version: number;
}

export interface Npe {
    lobby:                Arenabasics;
    firstmatch:           Arenabasics;
    movement:             Arenabasics;
    arenabasics:          Arenabasics;
    social_tab_seen:      BlueTintTabSeen;
    pointer:              BlueTintTabSeen;
    blue_tint_tab_seen:   BlueTintTabSeen;
    heraldry_tab_seen:    BlueTintTabSeen;
    orange_tint_tab_seen: BlueTintTabSeen;
}

export interface Arenabasics {
    completed: boolean;
}

export interface BlueTintTabSeen {
    version: number;
}

export interface Social {
    community_values_version: number;
    setup_version:            number;
    group:                    string;
}

export interface Server {
    displayname:     string;
    xplatformid:     string;
    _version:        number;
    publisher_lock:  string;
    purchasedcombat: number;
    lobbyversion:    number;
    modifytime:      number;
    logintime:       number;
    updatetime:      number;
    createtime:      number;
    stats:           Stats;
    unlocks:         Unlocks;
    loadout:         Loadout;
    dev: DeveloperSettings;
}

export interface DeveloperSettings {
    disable_afk_timeout: boolean;
    xplatformid: string;
}
export interface Loadout {
    instances: Instances;
    number:    number;
}

export interface Instances {
    unified: Unified;
}

export interface Unified {
    slots: Slots;
}

export interface Slots {
    decal:            string;
    decal_body:       string;
    emote:            string;
    secondemote:      string;
    tint:             string;
    tint_body:        string;
    tint_alignment_a: string;
    tint_alignment_b: string;
    pattern:          string;
    pattern_body:     string;
    pip:              string;
    chassis:          string;
    bracer:           string;
    booster:          string;
    title:            string;
    tag:              string;
    banner:           string;
    medal:            string;
    goal_fx:          string;
    emissive:         string;
}

export interface Stats {
    arena:  Arena;
    combat: StatsCombat;
}

export interface Arena {
    Level:                        Level;
    Goals:                        ArenaLosses;
    TopSpeedsTotal:               ArenaLosses;
    HighestArenaWinStreak:        ArenaLosses;
    ArenaWinPercentage:           ArenaLosses;
    ArenaWins:                    ArenaLosses;
    GoalsPerGame:                 ArenaLosses;
    Points:                       ArenaLosses;
    Interceptions:                ArenaLosses;
    ThreePointGoals:              ArenaLosses;
    Clears:                       ArenaLosses;
    BounceGoals:                  ArenaLosses;
    PossessionTime:               ArenaLosses;
    HatTricks:                    ArenaLosses;
    ShotsOnGoal:                  ArenaLosses;
    HighestPoints:                ArenaLosses;
    GoalScorePercentage:          ArenaLosses;
    AveragePossessionTimePerGame: ArenaLosses;
    AverageTopSpeedPerGame:       ArenaLosses;
    AveragePointsPerGame:         ArenaLosses;
    ArenaMVPPercentage:           ArenaLosses;
    ArenaMVPs:                    ArenaLosses;
    CurrentArenaWinStreak:        ArenaLosses;
    CurrentArenaMVPStreak:        ArenaLosses;
    HighestArenaMVPStreak:        ArenaLosses;
    XP:                           ArenaLosses;
    ShotsOnGoalAgainst:           ArenaLosses;
    ArenaLosses:                  ArenaLosses;
    Catches:                      ArenaLosses;
    StunsPerGame:                 ArenaLosses;
    HighestStuns:                 ArenaLosses;
    Steals:                       ArenaLosses;
    Stuns:                        ArenaLosses;
    PunchesReceived:              ArenaLosses;
    Passes:                       ArenaLosses;
    Blocks:                       ArenaLosses;
    BlockPercentage:              ArenaLosses;
}

export interface ArenaLosses {
    op:  Op;
    val: number;
}

export enum Op {
    Add = "add",
    Max = "max",
    Rep = "rep",
}

export interface Level {
    cnt: number;
    op:  Op;
    val: number;
}

export interface StatsCombat {
    Level: Level;
}

export interface Unlocks {
    arena:  { [key: string]: boolean };
    combat: UnlocksCombat;
}

export interface UnlocksCombat {
    rwd_booster_s10:            boolean;
    rwd_chassis_body_s10_a:     boolean;
    rwd_medal_s1_combat_gold:   boolean;
    rwd_title_title_b:          boolean;
    rwd_medal_s1_combat_bronze: boolean;
    rwd_medal_s1_combat_silver: boolean;
}
export interface ChannelInfo {
    group: Group[];
  }
  
  export interface Group {
    channeluuid:  string;
    name:         string;
    description:  string;
    rules:        string;
    rulesVersion: number;
    link:         string;
    priority:     number;
    rad:          boolean;
  }
  export interface Document {
    type: string;
    lang: string;
    version: number;
    versionGa: number;
    text: string;
    textGa: string;
    markAsReadProfileKey: string;
    markAsReadProfileKeyGa: string;
    linkCc: string;
    linkPp: string;
    linkVR: string;
    linkCp: string;
    linkEc: string;
    linkEa: string;
    linkGa: string;
    linkTc: string;
  }
  export interface Config {
    type: string
    id: string
    _ts: number
    news: NewsItem[]
    splash: NewsItem[]
    splash_version: number
    help_link: string
    news_link: string
    discord_link: string
  }
  
  export interface NewsItem {
    texture: string
    link: string
  }
  export interface LoginSettings {
    iapUnlocked:           boolean;
    remoteLogSocial:       boolean;
    remoteLogWarnings:     boolean;
    remoteLogErrors:       boolean;
    remoteLogRichPresence: boolean;
    remoteLogMetrics:      boolean;
    env:                   string;
    matchmakerQueueMode:   string;
    configData:            ConfigData;
  }
  export interface ConfigData {
}

export interface LinkCode {
    deviceId: string;
    code: string;
}

export interface DiscordAccessToken {
    access_token:  string;
    token_type:    string;
    expires_in:    number;
    refresh_token: string;
    scope:         string;
}

export interface DiscordUser {
    accent_color:           number;
    avatar:                 string;
    avatar_decoration_data: null;
    banner:                 null;
    banner_color:           string;
    discriminator:          string;
    flags:                  number;
    global_name:            null;
    id:                     string;
    locale:                 string;
    mfa_enabled:            boolean;
    premium_type:           number;
    public_flags:           number;
    username:               string;
}

// RemoteLogSetv3 log
export interface GamePauseSettings {
    announcer:            number;
    music:                number;
    sfx:                  number;
    voip:                 number;
    wristangleoffset:     number;
    smoothrotationspeed:  number;
    personalbubbleradius: number;
    grabdeadzone:         number;
    releasedistance:      number;
    personalbubblemode:   number;
    personalspacemode:    number;
    voipmode:             number;
    voipmodeffect:        number;
    voiploudnesslevel:    number;
    dynamicmusicmode:     number;
    HUD:                  boolean;
    EnableYaw:            boolean;
    EnablePitch:          boolean;
    EnableRoll:           boolean;
    EnableSmoothRotation: boolean;
    EnablePersonalBubble: boolean;
    EnablePersonalSpace:  boolean;
    EnableNetStatusHUD:   boolean;
    EnableNetStatusPause: boolean;
    EnableAPIAccess:      boolean;
    EnableGhostAll:       boolean;
    EnableMuteAll:        boolean;
    EnableMuteEnemyTeam:  boolean;
    MatchTagDisplay:      boolean;
    EnableVoipLoudness:   boolean;
    EnableMaxLoudness:    boolean;
    EnableStreamerMode:   boolean;
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
