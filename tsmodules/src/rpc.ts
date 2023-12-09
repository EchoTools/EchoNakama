import {
  CollectionMap,
  StoragePermissions,
  getStorageObject,
  parsePayload,
} from './utils';

import * as Errors from './errors';

import {
  Account,
  Server as AccountProfile,
  Client as ClientProfile,
  Server as ServerProfile,
  Config,
  Document,
  AccessControlList,
  ChannelInfo,
  LoginSettings
} from './types';

import {
  getDeviceLinkCodeRpc,
  discordLinkDeviceRpc,
} from './linking';


/**
 * Retrieves the Echo Relay account information for a user.
 *
 * @param ctx - The context for the execution.
 * @param logger - The server logger.
 * @param nk - The Nakama server APIs.
 * @param userId - The user ID for which to retrieve the account data.
 *
 * @throws {nkruntime.Error} If there is an issue with reading data from storage.
 *
 * @returns The Echo Relay account object populated with data from Nakama storage.
 */
const accountAsEchoRelayAccount = function (
  ctx: nkruntime.Context,
  logger: nkruntime.Logger,
  nk: nkruntime.Nakama,
  userId: string
): Account {
  const account = nk.accountGetId(userId);
  
  // The Echo Relay Account skeleton
  const echoAccount: Account = {
    is_moderator: false,
    banned_until: null,
    account_lock_hash: null,
    account_lock_salt: null,
    profile: {
      client: {} as ClientProfile,
      server: {} as ServerProfile,
    },
  };

  const storageReadReqs: nkruntime.StorageReadRequest[] = [
    { collection: CollectionMap.echoProfile, key: CollectionMap.echoProfileClient, userId },
    { collection: CollectionMap.echoProfile, key: CollectionMap.echoProfileServer, userId },
  ];

  let objects: nkruntime.StorageObject[] = [];

  try {
    objects = nk.storageRead(storageReadReqs);

  } catch (error) {
    logger.error('storageRead error: %s', error.message);
    throw error;
  }

  // Populate the Echo Relay account from storage objects
  objects.forEach((object) => {
    switch (object.key) {
      case CollectionMap.echoProfileClient:
        echoAccount.profile.client = object.value as ClientProfile;
        break;
      case CollectionMap.echoProfileServer:
        echoAccount.profile.server = object.value as ServerProfile;
        break;
    }
  });

  // force the displayName to be the same on the profiles and account
  echoAccount.profile.server.displayname = echoAccount.profile.client.displayname = account.user.displayName;

  return echoAccount;
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
const getAccountRpc =
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
let setAccountRpc: nkruntime.RpcFunction =
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

    const success = JSON.stringify({ success: true });
    const userId = ctx.userId;
 
    var echoAccount = {} as Account; 
    try {
    echoAccount = parsePayload(payload);
    } catch (error) {
      logger.error('parsePayload error: %s', error.message);
      throw {
        message: `Invalid account data: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }
    
    // Override the display name with the Nakama account
    echoAccount.profile.client.displayname = echoAccount.profile.server.displayname = nk.accountGetId(userId).user.displayName;

    // Write objects with appopriate permissions
    let newObjects: nkruntime.StorageWriteRequest[] = [
      {
        collection: CollectionMap.echoProfile,
        key: CollectionMap.echoProfileClient,
        userId,
        value: echoAccount.profile.client,
        permissionRead: StoragePermissions.PUBLIC_READ,
        permissionWrite: StoragePermissions.OWNER_WRITE
      },
      {
        collection: CollectionMap.echoProfile,
        key: CollectionMap.echoProfileServer,
        userId,
        value: echoAccount.profile.server,
        permissionRead: StoragePermissions.OWNER_READ,
        permissionWrite: StoragePermissions.NO_WRITE
      },

    ];

    const storageWriteAck = nk.storageWrite(newObjects);

    // Return a failure if the write does not succeed.
    if (!storageWriteAck || storageWriteAck.length == 0) {
      throw Errors.errInternal('storageWrite failed.');
    }

    return success;
  }


/**
 * Generates an RPC function for retrieving data from a Nakama storage collection.
 *
 * @param dataType - The type of data being retrieved.
 * @param collection - The name of the storage collection.
 * @param keySelector - A function that generates the storage key from the provided data.
 *
 * @returns An RPC function that retrieves data from the specified storage collection.
 *          Returns the retrieved data or throws an error if the operation fails.
 *
 * @template T - The type of data being retrieved.
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
 */
let generateRpcGetFunction = <T>(dataType: T, collectionSelect: Function, keySelect: Function): nkruntime.RpcFunction => {
  return function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {
    let userId = ctx.userId;
    let key = null;
    let collection = null;
    let data = {};
    logger.debug(payload)

    try {
      data = parsePayload(payload);

    } catch (error) {
      logger.error('parsePayload error: %s', error.message);
      throw {
        message: `Invalid/corrupt data: ${error}`,
        code: nkruntime.Codes.INTERNAL
      } as nkruntime.Error;
    }

    try {
      collection = collectionSelect(data);

    } catch (error) {
      logger.error('collectionSelect error: %s', error.message);
      throw {
        message: `Collection function (${collectionSelect}) failed to obtain collection name: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    try {
      key = keySelect(data);

    } catch (error) {
      logger.error('keySelect error: %s', error.message);
      throw {
        message: `Key function (${keySelect}) failed to obtain key name: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    let responseData = {} as typeof dataType;
    try {
      let objects: nkruntime.StorageObject[] = nk.storageRead([{ collection, key, userId }]);


      if (objects.length == 0) {
        logger.error('storageRead error: %s', 'storageRead returned empty');
        throw {
          message: `${collection}/${key} not found.`,
          code: nkruntime.Codes.NOT_FOUND
        } as nkruntime.Error;
      }
      
      responseData = objects[0].value as typeof dataType;

    } catch (error) {
      logger.error('storageRead error: %s', error.message);
      throw error;
    }

    logger.debug(JSON.stringify(responseData));
    return JSON.stringify(responseData);

  };
};


/**
 * Generates an RPC function for storing data in a Nakama storage collection.
 *
 * @param dataType - The type of data being stored.
 * @param collection - The name of the storage collection.
 * @param keySelect - A function that generates the storage key from the provided data.
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
let generateRpcSetFunction = <T>(dataType: T, collectionSelect: Function, keySelect: Function): nkruntime.RpcFunction => {
  return function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {
    const userId = ctx.userId;
    logger.debug(payload)
    let data = {} as typeof dataType;
    let collection = null;
    let key = null;

    try {
      data = parsePayload(payload) as typeof dataType;
    } catch (error) {
      logger.error('parsePayload error: %s', error.message);
      throw {
        message: `Invalid ${collection} data: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }
    try {
      collection = collectionSelect(data);

    } catch (error) {
      logger.error('collectionSelect error: %s', error.message);
      throw {
        message: `Collection function (${collectionSelect}) failed to obtain collection name: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    try {
      key = keySelect(data);

    } catch (error) {
      logger.error('keySelect error: %s', error.message);
      throw {
        message: `Key function (${keySelect}) failed to obtain key name: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    logger.error("collection: %s, key: %s", collection, key);

    let storageWriteAck = [] as nkruntime.StorageWriteAck[];

    try {
      storageWriteAck = nk.storageWrite([
        {
          collection,
          key,
          value: data,
          userId,
          permissionRead: 2,
          permissionWrite: 1
        }
      ]);
    } catch (error) {
      logger.error('storageWrite error: %s', error.message);
      throw {
        message: `Invalid ${collection} data: ${error}`,
        code: nkruntime.Codes.INTERNAL
      } as nkruntime.Error;
    }

    if (!storageWriteAck || storageWriteAck.length === 0) {
      logger.error('storageWrite error: %s', 'storageWriteAck is empty');
      throw Errors.errInternal('storageWrite failed.');
    }
    logger.debug("success");
    return JSON.stringify({ success: true });
  };
};


// NOTE: Due to Nakama's mapping method, all Rpc functions must be globally assigned
let setConfigRpc = generateRpcSetFunction<Config>({} as Config, (e: any) => `Config:${e.type}`, (e: any) => e.id);
let getConfigRpc = generateRpcGetFunction<Config>({} as Config, (e: any) => `Config:${e.type}`, (e: any) => e.id);
let setDocumentRpc = generateRpcSetFunction<Document>({} as Document, (e: any) => `Document:${e.type}`, (e: any) => `${e.type}_${e.lang}`);
let getDocumentRpc = generateRpcGetFunction<Document>({} as Document, (e: any) => `Document:${e.type}`, (e: any) => e.id);
let setAccessControlListRpc = generateRpcSetFunction<AccessControlList>({} as AccessControlList, (e: any) => `Relay:acl`, (e: any) => 'allow_deny_list');
let getAccessControlListRpc = generateRpcGetFunction<AccessControlList>({} as AccessControlList, (e: any) => `Relay:acl`, (e: any) => 'allow_deny_list');
let setChannelInfoRpc = generateRpcSetFunction<ChannelInfo>({} as ChannelInfo, (e: any) => `Login:channel_info`, (e: any) => 'channel_info');
let getChannelInfoRpc = generateRpcGetFunction<ChannelInfo>({} as ChannelInfo, (e: any) => `Login:channel_info`, (e: any) => 'channel_info');
let setLoginSettingsRpc = generateRpcSetFunction<LoginSettings>({} as LoginSettings, (e: any) => `Login:login_settings`, (e: any) => 'login_settings');
let getLoginSettingsRpc = generateRpcGetFunction<LoginSettings>({} as LoginSettings, (e: any) => `Login:login_settings`, (e: any) => 'login_settings');


export {
  setAccountRpc,
  getAccountRpc,
  setConfigRpc,
  getConfigRpc,
  setDocumentRpc,
  getDocumentRpc,
  setAccessControlListRpc,
  getAccessControlListRpc,
  setChannelInfoRpc,
  getChannelInfoRpc,
  setLoginSettingsRpc,
  getLoginSettingsRpc,
  getDeviceLinkCodeRpc,
  discordLinkDeviceRpc,
}
