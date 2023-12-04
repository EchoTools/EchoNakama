import _, { mapKeys, without } from 'lodash';

import {
  Account,
  Server as AccountProfile,
  Client as ClientProfile,
  Server as ServerProfile,
  Config,
  Document,
  AccessControlList,
  ChannelInfo,
  LoginSettings,
  LinkCode
} from './types';

import {
  generateLinkCode,
  getStorageObject,
  parsePayload,

} from './utils';

const systemUserId = "00000000-0000-0000-0000-000000000000";
const LINKCODE_COLLECTION = "LinkCode";

const errMissingPayload: nkruntime.Error = {
  message: 'no payload provided.',
  code: nkruntime.Codes.NOT_FOUND
}

const errBadInput: nkruntime.Error = {
  message: 'input contained invalid data.',
  code: nkruntime.Codes.INVALID_ARGUMENT
};

const errMissingId: nkruntime.Error = {
  message: 'no "id" provided.',
  code: nkruntime.Codes.NOT_FOUND
}

const errInternal = (message: string): nkruntime.Error => ({
  message: 'errInternal: ' + message,
  code: nkruntime.Codes.INTERNAL
} as nkruntime.Error);

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
let accountAsEchoRelayAccount = function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, userId: string) {
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
    { collection: 'RelayConfig', key: 'AuthSecrets', userId },
    { collection: 'Profile', key: 'Client', userId },
    { collection: 'Profile', key: 'Server', userId },
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
      case 'Client':
        account.profile.client = object.value as ClientProfile;
        break;
      case 'Server':
        account.profile.server = object.value as ServerProfile;
        break;
      case 'AuthSecrets':
        let authSecrets = object.value;
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
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string | void) {

    const success = JSON.stringify({ success: true });
    let userId = ctx.userId;

    if (!payload || payload === "") {
      logger.error('Request missing payload.');
      throw errMissingPayload;
    }

    try {
      var account = JSON.parse(payload);
    } catch (error) {
      throw {
        message: `Invalid/corrupt data: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
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
      { collection: 'Profile', key: 'Client', userId, value: account.profile.client, permissionRead: 2, permissionWrite: 1 },
      { collection: 'Profile', key: 'Server', userId, value: account.profile.server, permissionRead: 2, permissionWrite: 1 },
      { collection: 'RelayConfig', key: 'AuthSecrets', userId, value: authSecrets, permissionRead: 1, permissionWrite: 1 },
    ];

    const storageWriteAck = nk.storageWrite(newObjects);

    // Return a failure if the write does not succeed.
    if (!storageWriteAck || storageWriteAck.length == 0) {
      throw errInternal('storageWrite failed.');
    }

    return success;
  }

/**
 * RPC function to retrieve or generate a link code for a device.
 *
 * @param ctx - The Nakama runtime context.
 * @param logger - The Nakama logger instance for logging messages.
 * @param nk - The Nakama runtime instance.
 * @param payload - The payload containing information about the device (e.g., deviceId).
 * @returns A JSON string representing the link code for the device.
 * @throws If the deviceId is missing in the payload, an error with code `errMissingId` is thrown.
 *         If there's any other error during the process, an error with code `errInternal` is thrown.
 *         If a link code already exists for the device, it is retrieved and returned.
 */
let getDeviceLinkCodeRpc: nkruntime.RpcFunction =
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {
    // Parse the payload data.
    let deviceId = null;
    try {
      logger.info(payload);
      let data = JSON.parse(payload);
      deviceId = data.id;
    } catch (error) {
      throw {
        message: `Invalid/corrupt data: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    let linkData = null;
    while (linkData == null) {
      // Generate a new link code
      let newLink = { deviceId, "code": generateLinkCode() };
      try {
        // check if this link code exists
        let linkObj = getStorageObject(nk, logger, LINKCODE_COLLECTION, linkData.code, systemUserId)
      } catch (error) {
        // The link code doesn't exist, so we can use it
        linkData = newLink;
      }
    }

    // Try to create a new link code in storage.
    try {
      nk.storageWrite([
        {
          collection: LINKCODE_COLLECTION, key: linkData.code,
          value: linkData,
          userId: systemUserId, version: '*', permissionRead: 1, permissionWrite: 1
        }
      ]);
    } catch (error) {
      logger.error("Failed to create link code: %s", error);
      throw errInternal(`Failed to create link code: ${error}`);
    }

    logger.info("Generated link code: %s for %s", linkData.code, linkData.deviceId)
    // Retrieve and return the link code for the device.
    return JSON.stringify(linkData);
  };

let DiscordLinkDeviceRpc: nkruntime.RpcFunction =
  function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

    // Parse the payload data.

    var deviceLinkCode = null;
    var oauthCode = null;
    var oauthRedirectUrl = null;
    try {
      logger.debug("Payload:", payload);
      let data = JSON.parse(payload);

      deviceLinkCode = data.linkCode;
      oauthCode = data.oauthCode
      oauthRedirectUrl = data.oauthRedirectUrl;

    } catch (error) {
      throw {
        message: `Invalid/corrupt data: ${error}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    if (!oauthRedirectUrl) {
      logger.error("Didn't get oauth redirect in url");
      throw {
        message: `OAuth redirect URL is missing from payload: ${payload}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    if (!oauthCode) {
      logger.error("Didn't get oauth code");
      throw {
        message: `OAuth code is missing from payload: ${payload}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    // Validate the deviceId and throw an error if missing.
    if (!deviceLinkCode || deviceLinkCode.length != 4) {
      throw {
        message: `Link code is invalid/missing from payload: ${payload}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }
    // Sanitize the linkCode and make sure it's only capital letters
    deviceLinkCode = deviceLinkCode.toUpperCase().replace(/[^A-Z]/g, '').slice(0, 4)

    // Retrieve the linkCode and deviceId from storage.
    let linkObject = {} as LinkCode;
    try {
      linkObject = getStorageObject(nk, logger, LINKCODE_COLLECTION, deviceLinkCode, systemUserId);
    } catch (error) {
      throw {
        message: `Link code not found: ${error.message}`,
        code: nkruntime.Codes.NOT_FOUND
      } as nkruntime.Error;
    }

    // ensure the payload contains the discord code
    if (!oauthCode) {
      throw {
        message: `Discord code is missing from payload: ${payload}`,
        code: nkruntime.Codes.INVALID_ARGUMENT
      } as nkruntime.Error;
    }

    var params = `client_id=${ctx.env["DISCORD_CLIENT_ID"]}&` +
      `client_secret=${ctx.env["DISCORD_CLIENT_SECRET"]}&` +
      `code=${oauthCode}&` +
      `grant_type=authorization_code&` +
      `redirect_uri=${oauthRedirectUrl}&` +
      `scope=identify`;
    logger.debug("%s", params);

    // exchange the oauthCode for the user's access token
    let response = nk.httpRequest("https://discord.com/api/v10/oauth2/token", "post",
      {
        'Accept': 'application/json',
        "Content-Type": "application/x-www-form-urlencoded",
        "Authorization": `${ctx.env["DISCORD_CLIENT_ID"]}:${ctx.env["DISCORD_CLIENT_SECRET"]}`,

      }, params);
    
      let discordResponse = null;
    try {
      discordResponse = JSON.parse(response.body);
    } catch (error) {
      logger.error("Could not decode discord response body: %s", response.body);
      throw errInternal(`Could not decode discord response body: ${response.body}`);
    }
      
    if (response.code != 200) {
      if ("error" in discordResponse && discordResponse.error == "invalid_grant") {
        logger.error("Discord code exchange failed: %s", response.body);
        throw {
          message: `Discord code exchange failed: ${response.body}`,
          code: nkruntime.Codes.UNAUTHENTICATED
        } as nkruntime.Error;
      }
    }

    // get the user's discord data with the access token  
    let accessToken = JSON.parse(response.body);
    response = nk.httpRequest("https://discord.com/api/v10/users/@me", "get",
      {
        "Content-Type": "application/x-www-form-urlencoded",
        "Authorization": `Bearer ${accessToken.access_token}`,
        'Accept': 'application/json'
      }, null);

    logger.debug(response.body);

    if (response.code != 200) {
      logger.error("Discord user lookup failed: %s", response.body);
      throw errInternal(`Discord user lookup failed: ${response.body}`);
    }

    let discordUser = JSON.parse(response.body);
    // Construct the username from the discord user data
    let username = discordUser.id;
    var accountId = null;

    let users = nk.usersGetUsername([username])
    if (users.length == 1) {
      accountId = users[0].userId;
      // Link the device to the account
      nk.linkDevice(accountId, linkObject.deviceId);
    }

  
    try {
      let authResult = nk.authenticateDevice(linkObject.deviceId, null, false);
      logger.debug("Auth result: %s", authResult);
      accountId = authResult.userId;
    } catch (error) {
        logger.debug("%s", error);      
    }

    // Authenticate with the device ID, creating the account if it doesn't exist
    try {
      let result = nk.authenticateDevice(linkObject.deviceId, username, true);
      accountId = result.userId;
    } catch (error) {
      logger.error('Failed to authenticate device (%s) to user %s: %s', linkObject.deviceId, username, error);
      throw errInternal(`Failed to authenticate device: ${error}`);
    }
    // Link the discord token to the account
    try {
      nk.unlinkCustom(accountId);
    } catch (noterror) {
    }
    try {
      nk.linkCustom(accountId, accessToken.access_token);
    } catch (error) {
      logger.error('Could not link %s to %s', accessToken.access_token, accountId);
      throw errInternal(`Could not link ${accessToken.access_token} to ${accountId}`);
    }

    try {
      // remove the link code
      nk.storageDelete([{ collection: LINKCODE_COLLECTION, key: linkObject.code, userId: systemUserId }]);
    } catch (error) {
      logger.error("Failed to delete link code: %s", error);
      throw errInternal(`Failed to delete link code: ${error}`);
    }

    // Retrieve the full account
    let account = {} as nkruntime.Account;
    try {
      account = nk.accountGetId(accountId);
    } catch (error) {
      logger.error("Failed to get account: %s", error);
      throw errInternal(`Failed to get account: ${error}`);
    }

    // Update the user's info
    _.merge(account.user.metadata, { discord: { user: discordUser, oauth: accessToken } });

    try {
      let displayName = account.user.displayName ?? discordUser.global_name ?? discordUser.username;
      nk.accountUpdateId(accountId, username, displayName, null, null, null, null, account.user.metadata);
    } catch (error) {
      logger.error("Failed to update user: %s", error);
      throw errInternal(`Failed to update user: ${error}`);
    }

    // Retrieve and return the link code for the device.
    return JSON.stringify({ "success": true });
  };

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
        throw errInternal('storageWrite failed.');
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
let setConfigRpc = generateRpcSetFunction<Config>({} as Config, 'ClientConfig', (e: any) => e.id);
let getConfigRpc = generateRpcGetFunction<Config>({} as Config, 'ClientConfig', (e: any) => e.id);
let setDocumentRpc = generateRpcSetFunction<Document>({} as Document, 'Info', (e: any) => `${e.type}_${e.lang}`);
let getDocumentRpc = generateRpcGetFunction<Document>({} as Document, 'Info', (e: any) => `${e.type}_${e.lang}`);
let setAccessControlListRpc = generateRpcSetFunction({} as AccessControlList, 'RelayConfig', (e: any) => 'AccessControlList');
let getAccessControlListRpc = generateRpcGetFunction({} as AccessControlList, 'RelayConfig', (e: any) => 'AccessControlList');
let setChannelInfoRpc = generateRpcSetFunction({} as ChannelInfo, 'Info', (e: any) => 'Channels');
let getChannelInfoRpc = generateRpcGetFunction({} as ChannelInfo, 'Info', (e: any) => 'Channels');
let setLoginSettingsRpc = generateRpcSetFunction({} as LoginSettings, 'ClientConfig', (e: any) => 'LoginSettings');
let getLoginSettingsRpc = generateRpcGetFunction({} as LoginSettings, 'ClientConfig', (e: any) => 'LoginSettings');


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
  DiscordLinkDeviceRpc,
}
