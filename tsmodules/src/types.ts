export interface AccessControlList {
  allow: string[];
  deny: string[];
}export interface AccountMetadata {
  suspended_until: number;
}

export interface Account {
  account_lock_hash: string | null,
  account_lock_salt: string | null,
  banned_until: number,
  is_moderator: boolean,
  profile: {
    client: ClientProfile,
    server: ServerProfile,
  }
}

export interface ClientProfile {
  customization: {
    battlepass_season_poi_version: number;
    clear_new_unlocks_version: number;
    new_unlocks_poi_version: number;
    store_entry_poi_version: number;
  };
  displayname: string;
  grenade: string;
  newunlocks: any[]; // Replace 'any' with the actual type if you have more specific information
  npe: {
    arenabasics: {
      completed: boolean;
    };
    firstmatch: {
      completed: boolean;
    };
    lobby: {
      completed: boolean;
    };
    movement: {
      completed: boolean;
    };
  };
  social: {
    community_values_version: number;
    setup_version: number;
  };
  weapon: string;
  weaponarm: number;
  xplatformid: string;
}

export interface ServerProfile {
  _version: number;
  createtime: number;
  displayname: string;
  loadout: {
    instances: {
      unified: {
        slots: {
          banner: string;
          booster: string;
          bracer: string;
          chassis: string;
          decal: string;
          decal_body: string;
          emissive: string;
          emote: string;
          goal_fx: string;
          medal: string;
          pattern: string;
          pattern_body: string;
          pip: string;
          secondemote: string;
          tag: string;
          tint: string;
          tint_alignment_a: string;
          tint_alignment_b: string;
          tint_body: string;
          title: string;
        };
      };
    };
    number: number;
  };
  lobbyversion: number;
  logintime: number;
  modifytime: number;
  publisher_lock: string;
  purchasedcombat: number;
  stats: {
    arena: {
      Level: {
        cnt: number;
        op: "add";
        val: number;
      };
    };
    combat: {
      Level: {
        cnt: number;
        op: "add";
        val: number;
      };
    };
  };
  unlocks: {
    arena: {
      [key: string]: boolean;
    };
    combat: {
      [key: string]: boolean;
    };
  };
  updatetime: number;
  xplatformid: string;
}

export interface AuthSecrets {
  AccountLockHash: string | null;
  AccountLockSalt: string | null;
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
