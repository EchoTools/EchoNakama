import {
  Account,
  ClientProfile,
  ServerProfile,
  AuthSecrets,
  RemoteConfig,
  Document,
  AccessControlList,
  ChannelInfo,
  LoginSettings
} from './types';

const errBadInput: nkruntime.Error = {
  message: 'input contained invalid data',
  code: nkruntime.Codes.INVALID_ARGUMENT
};

const errMissingId: nkruntime.Error = {
  message: 'no "id" provided.',
  code: nkruntime.Codes.NOT_FOUND
}

const errInternal: nkruntime.Error = {
  message: 'Internal Error.',
  code: nkruntime.Codes.INTERNAL
}

/**
 * Converts Nakama account-related data into an Echo Relay Account.
 *
 * @param ctx - The context for the execution.
 * @param logger - The server logger.
 * @param nk - The Nakama server APIs.
 * @param userId - The user ID for which to retrieve and convert account data.
 *
 * @throws {nkruntime.Error} If there is an issue with reading data from storage.
 *
 * @returns An Echo Relay Account object populated with data from Nakama storage.
 */
let accountAsEchoRelayAccount = function (ctx: nkruntime.Context, logger:
  nkruntime.Logger, nk:
    nkruntime.Nakama, userId:
    string) {
  // get the account objects

  // The Echo Relay Account skeleton
  let account: Account = {
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

/**
 * Retrieves and returns the Echo Relay account information for the current user.
 *
 * @param ctx - The context for the execution.
 * @param logger - The server logger.
 * @param nk - The Nakama server APIs.
 * @param payload - Unused payload parameter (included for RPC signature compatibility).
 *
 * @throws {nkruntime.Error} If there is an issue with retrieving the account data.
 *
 * @returns A JSON string representing the Echo Relay account information for the current user.
 */
const echoRelayGetAccountRpc =
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

    let userId = ctx.userId;
    let account = accountAsEchoRelayAccount(ctx, logger, nk, userId);

    return JSON.stringify(account);
  }


/**
 * Updates the Echo Relay account information for a user.
 *
 * @param ctx - The context for the execution.
 * @param logger - The server logger.
 * @param nk - The Nakama server APIs.
 * @param payload - The input data containing Echo Relay account information.
 *
 * @throws {nkruntime.Error} If there is an issue with parsing the payload, missing ID, or storage write operation failure.
 *
 * @returns A JSON string indicating the success of the operation.
 */
let echoRelaySetAccountRpc: nkruntime.RpcFunction =
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

    const success = JSON.stringify({ success: true });
    let userId = ctx.userId;

    if (!payload || payload === "") {
      throw errMissingId;
    }

    try {
      var account = JSON.parse(payload) as Account;
    } catch {
      throw errBadInput;
    }

    // Syncronize the displayname given by the client
    account.profile.server.displayname = account.profile.client.displayname;
    nk.accountUpdateId(userId, null, account.profile.client.displayname, null, null, null, null, null);

    // Extract auth secrets to go into another object
    let authSecrets = {
      AccountLockHash: account.account_lock_hash,
      AccountLockSalt: account.account_lock_salt,
    };

    // Write objects with appopriate permissions
    let newObjects: nkruntime.StorageWriteRequest[] = [
      { collection: 'profile', key: 'client', userId, value: account.profile.client, permissionRead: 2, permissionWrite: 1 },
      { collection: 'profile', key: 'server', userId, value: account.profile.server, permissionRead: 2, permissionWrite: 1 },
      { collection: 'relayConfig', key: 'authSecrets', userId, value: authSecrets, permissionRead: 1, permissionWrite: 1 },
    ];

    // Write the updated inventory to storage.
    const storageWriteAck = nk.storageWrite(newObjects);

    // Return a failure if the write does not succeed.
    if (!storageWriteAck || storageWriteAck.length == 0) {
      throw errInternal;
    }

    return success;
  }


/**
 * Generates an RPC function for retrieving data from a Nakama storage collection.
 *
 * @param dataType - The type of data being retrieved.
 * @param collection - The name of the storage collection.
 * @param keyFunc - A function that generates the storage key from the provided data.
 *
 * @returns An RPC function that retrieves data from the specified storage collection.
 *          Returns the retrieved data or throws an error if the operation fails.
 *
 * @template T - The type of data being retrieved.
 *
 * @throws {nkruntime.Error} If there is an issue parsing the payload or the payload does not match the expected data type.
 *   - `code: nkruntime.Codes.INTERNAL` if the payload is invalid or corrupt.
 *   - `code: nkruntime.Codes.INVALID_ARGUMENT` if the storage key cannot be generated from the provided data.
 *   - `code: nkruntime.Codes.NOT_FOUND` if the requested data is not found in the storage collection.
 */
let generateRpcGetFunction = <T>(dataType: T, collection: string, keyFunc: Function): nkruntime.RpcFunction => {
  return function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {
    let userId = ctx.userId;
    try {
      var data = JSON.parse(payload) as typeof dataType;

    } catch (error) {
      throw {
        message: `Invalid/corrupt ${collection} data: ${error}`,
        code: nkruntime.Codes.INTERNAL
      } as nkruntime.Error;
    }

    try {
      var key = keyFunc(data);

    } catch (error) {
      throw {
        message: `${collection} key not found in data: ${keyFunc} ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    try {
      let objects: nkruntime.StorageObject[] = nk.storageRead([{ collection, key, userId }]);

      if (objects.length == 0) throw {
        message: `${collection}/${key} not found.`,
        code: nkruntime.Codes.NOT_FOUND
      } as nkruntime.Error;

      return JSON.stringify(objects[0].value);

    } catch (error) {
      logger.error('storageRead error: %s', error.message);
      throw error;
    }
  }
}
/**
 * Generates an RPC function for storing data in a Nakama storage collection.
 *
 * @param dataType - The type of data being stored.
 * @param collection - The name of the storage collection.
 * @param keyFunc - A function that generates the storage key from the provided data.
 *
 * @returns An RPC function that stores data in the specified storage collection.
 *          Returns a success response or throws an error if the operation fails.
 *
 * @template T - The type of data being stored.
 *
 * @param ctx - The context for the execution.
 * @param logger - The server logger.
 * @param nk - The Nakama server APIs.
 * @param payload - The input data to the function call. This is usually an escaped JSON object.
 *
 * @throws {nkruntime.Error} If there is an issue parsing the payload or the payload does not match the expected data type.
 *   - `code: nkruntime.Codes.INVALID_ARGUMENT` if the payload is invalid or does not match the expected data type.
 *   - `code: nkruntime.Codes.INTERNAL` if there is an issue during the storage write operation.
 *
 * @throws {nkruntime.Error} If the storage write operation fails or if the response from storage write is falsy or empty.
 *   - `code: nkruntime.Codes.INTERNAL` if the storage write operation fails.
 *
 * @throws {nkruntime.Error} If the generated storage key is not found in the data.
 *   - `code: nkruntime.Codes.INVALID_ARGUMENT` if the storage key cannot be generated from the provided data.
 */
let generateRpcSetFunction = <T>(dataType: T, collection: string, keyFunc: Function): nkruntime.RpcFunction => {
  return function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {
    let userId = ctx.userId;
    try {
      var data = JSON.parse(payload) as typeof dataType;

    } catch (error) {
      throw {
        message: `Invalid ${collection} data: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }
    try {
      let storageWriteAck = nk.storageWrite([{ collection: collection, key: keyFunc(data), value: data, userId: userId, permissionRead: 2, permissionWrite: 1 }]);
      if (!storageWriteAck || storageWriteAck.length == 0) {
        throw errInternal;
      }
      return JSON.stringify({ success: true });

    } catch (error) {
      throw {
        message: `Invalid ${collection} data: ${error}`,
        code: nkruntime.Codes.INTERNAL
      } as nkruntime.Error;
    }
  }
}


// Due to Nakama's mapping method, all Rpc functions must be globally assigned
let setRemoteConfigRpc = generateRpcSetFunction<RemoteConfig>({} as RemoteConfig, 'ClientConfig', (e: any) => e.id);
let getRemoteConfigRpc = generateRpcGetFunction<RemoteConfig>({} as RemoteConfig, 'ClientConfig', (e: any) => e.id);
let setDocumentRpc = generateRpcSetFunction<Document>({} as Document, 'Info', (e: any) => `${e.type}_${e.lang}`);
let getDocumenRpc = generateRpcGetFunction<Document>({} as Document, 'Info', (e: any) => `${e.type}_${e.lang}`);
let setAccessControlListRpc = generateRpcSetFunction({} as AccessControlList, 'RelayConfig', (e: any) => 'AccessControlList');
let getAccessControlListRpc = generateRpcGetFunction({} as AccessControlList, 'RelayConfig', (e: any) => 'AccessControlList');
let setChannelInfoRpc = generateRpcSetFunction({} as ChannelInfo, 'Info', (e: any) => 'Channels');
let getChannelInfoRpc = generateRpcGetFunction({} as ChannelInfo, 'Info', (e: any) => 'Channels');
let setLoginSettingsRpc = generateRpcSetFunction({} as LoginSettings, 'ClientConfig', (e: any) => 'LoginSettings');
let getLoginSettingsRpc = generateRpcGetFunction({} as LoginSettings, 'ClientConfig', (e: any) => 'LoginSettings');


// Initialize the server module
let InitModule: nkruntime.InitModule =
  function (ctx: nkruntime.Context,
    logger: nkruntime.Logger,
    nk: nkruntime.Nakama,
    initializer: nkruntime.Initializer) {

    initializer.registerRpc('echorelay/setRemoteConfig', setRemoteConfigRpc);
    initializer.registerRpc('echorelay/getRemoteConfig', getRemoteConfigRpc);
    initializer.registerRpc('echorelay/setDocument', setDocumentRpc);
    initializer.registerRpc('echorelay/getDocument', getDocumenRpc);
    initializer.registerRpc('echorelay/setAccessControlList', setAccessControlListRpc);
    initializer.registerRpc('echorelay/getAccessControlList', getAccessControlListRpc);
    initializer.registerRpc('echorelay/setChannelInfo', setChannelInfoRpc);
    initializer.registerRpc('echorelay/getChannelInfo', getChannelInfoRpc);
    initializer.registerRpc('echorelay/setLoginSettings', setLoginSettingsRpc);
    initializer.registerRpc('echorelay/getLoginSettings', getLoginSettingsRpc);
  }


// Reference InitModule to avoid it getting removed on build
!InitModule && InitModule.bind(null);
