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