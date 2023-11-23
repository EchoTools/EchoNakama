const errBadInput: nkruntime.Error = {
  message: 'input contained invalid data',
  code: nkruntime.Codes.INVALID_ARGUMENT
};

const errAccountParseError: nkruntime.Error = {
  message: 'could not parse profile data on server.',
  code: nkruntime.Codes.INTERNAL
}

const errAccountNotFound: nkruntime.Error = {
  message: 'account not found.',
  code: nkruntime.Codes.NOT_FOUND
}

interface AccountMetadata {
  suspended_until: number;
}

interface EchoRelayAccount {
  account_lock_hash: string | null,
  account_lock_salt: string | null,
  banned_until: number,
  is_moderator: boolean,
  profile: {
    client: ClientProfile,
    server: ServerProfile,
  }
}
interface ClientProfile {
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

interface ServerProfile {
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

interface AuthSecrets {
  AccountLockHash: string | null;
  AccountLockSalt: string | null;
}


let accountAsEchoRelayAccount = function (ctx: nkruntime.Context, logger:
	nkruntime.Logger, nk:
	nkruntime.Nakama, userId:
	string) {
  // get the account objects

  // The Echo Relay Account skeleton
  let account: EchoRelayAccount = {
    is_moderator: false,
    banned_until: null,
    account_lock_hash: null,
    account_lock_salt: null,
    profile: {
      client: null,
      server: null,
    },
  };

  let storageReadReqs: nkruntime.StorageReadRequest[] = [
    { collection: 'relayConfig', key: 'authSecrets', userId },
    { collection: 'profile', key: 'client', userId },
    { collection: 'profile', key: 'server', userId },
  ];

  let objects: nkruntime.StorageObject[] = [];

  try {
    objects = nk.storageRead(storageReadReqs);

  } catch (error) {
    logger.error('storageRead error: %s', error.message);
    throw error;
  }

  // populate the Echo Relay account from storage objects
  objects.forEach((object) => {
    switch (object.key) {
      case 'client':
        account.profile.client = object.value as ClientProfile;
        break;
      case 'server':
        account.profile.server = object.value as ServerProfile;
        break;
      case 'authSecrets':
        let authSecrets = object.value as AuthSecrets;
        account.account_lock_hash = authSecrets.AccountLockHash;
        account.account_lock_salt = authSecrets.AccountLockSalt;
        break;
    }
  });

  return account;
}


let InitModule: nkruntime.InitModule =
	function (ctx: nkruntime.Context,
		logger: nkruntime.Logger,
		nk: nkruntime.Nakama,
		initializer: nkruntime.Initializer) {

		initializer.registerRpc("echorelay/getAccount", echoRelayGetAccountRpc);
		initializer.registerRpc("echorelay/setAccount", echoRelaySetAccountRpc);
	}


let echoRelayGetAccountRpc: nkruntime.RpcFunction =
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

    let userId = ctx.userId;
    let account = accountAsEchoRelayAccount(ctx, logger, nk, userId);

    return JSON.stringify(account);
  }


let echoRelaySetAccountRpc: nkruntime.RpcFunction =
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

    let userId = ctx.userId;

    //let accountExisting = EchoRelayGetAccount(ctx, logger, nk, userId);

    let account = JSON.parse(payload) as EchoRelayAccount;

    nk.accountUpdateId(userId, null, account.profile.client.displayname, null, null, null, null, null);

    let authSecrets = {
      AccountLockHash: account.account_lock_hash,
      AccountLockSalt: account.account_lock_salt,
        };

    let newObjects: nkruntime.StorageWriteRequest[] = [
      { collection: 'profile', key: 'client', userId, value: account.profile.client, permissionRead: 2, permissionWrite: 1 },
      { collection: 'profile', key: 'server', userId, value: account.profile.server, permissionRead: 2, permissionWrite: 1 },
      { collection: 'relayConfig', key: 'authSecrets', userId, value: authSecrets, permissionRead: 1, permissionWrite: 1 },
    ];

    try {
      nk.storageWrite(newObjects);
    } catch (error) {
      logger.error('storageWrite error: %s', error.message);
      throw error;
    }
    return JSON.stringify({ "success": true, })
  }


/*
let echoRelayDocumentRpc: nkruntime.RpcFunction =
function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {
  let userId = ctx.userId;
 
  try {
    var jsonData = JSON.parse(payload);
    documentId = payload['Id'];
  }
 
 
  // get the document
  let storageReadReqs: nkruntime.StorageReadRequest[] = [
    { collection: 'relayConfig', key: 'authSecrets', userId },
    { collection: 'profile', key: 'client', userId },
    { collection: 'profile', key: 'server', userId },
  ];
 
  let objects: nkruntime.StorageObject[] = [];
 
  try {
    objects = nk.storageRead(storageReadReqs);
 
  } catch (error) {
    logger.error('storageRead error: %s', error.message);
    throw error;
  }
*/

// Reference InitModule to avoid it getting removed on build
!InitModule && InitModule.bind(null);


