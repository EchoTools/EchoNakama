import _ from 'lodash';
import { LinkCode, DiscordAccessToken } from './types';
import { getStorageObject } from './utils';
import { errInternal } from './errors';
import { discordExchangeCode, discordGetCurrentUser, discordRefreshAccessToken } from './discord';
import { StoragePermissions } from './utils';
import { systemUserId } from './utils';
import { CollectionMap } from './utils';


/**
 * Generates a random 4-character link code using a specified set of characters.
 * The characters include uppercase letters (excluding I and O) and digits.
 * 
 * @returns {string} A random 4-character link code.
 */
const generateLinkCode = (): string => {
  // Define the set of valid characters for the link code
  const characters = 'ABCDEFGHJKLMNPQRSTUVWXYZ';

  // Use lodash.range to create an array with 4 elements
  const indices = _.range(4);

  // Use lodash.sample to randomly select an index from the array
  const code = indices.map(() => _.sample(characters)).join('');

  return code;
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
const getDeviceLinkCodeRpc: nkruntime.RpcFunction = function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {
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
      // Check if this link code exists
      let linkObj = getStorageObject(nk, logger, CollectionMap.linkCode, linkData.code, systemUserId)
    } catch (error) {
      // The link code doesn't exist, so we can use it
      linkData = newLink;
    }
  }

  // Try to create a new link code in storage.
  try {
    nk.storageWrite([
      {
        collection: CollectionMap.linkCode,
        key: linkData.code,
        value: linkData,
        userId: systemUserId,
        version: '*',
        permissionRead: StoragePermissions.OWNER_READ,
        permissionWrite: StoragePermissions.OWNER_WRITE
      }
    ]);
  } catch (error) {
    logger.error("Failed to create link code: %s", error);
    throw errInternal(`Failed to create link code: ${error}`);
  }

  logger.info("Generated link code: %s for %s", linkData.code, linkData.deviceId);
  // Retrieve and return the link code for the device.
  return JSON.stringify(linkData);
};



let discordLinkDeviceRpc: nkruntime.RpcFunction = function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

  // Parse and validate the payload data
  var { oauthCode, oauthRedirectUrl, linkObject } = _validateLinkRequest(logger, payload, nk);

  // Exchange the oauthCode for the user's access token, and retrieve user data
  let accessToken = discordExchangeCode(ctx, nk, logger, oauthCode, oauthRedirectUrl);
  let discordUser = discordGetCurrentUser(ctx, nk, logger, accessToken);
  // Construct the username from the discord user data
  let username = discordUser.id;

  // Link the device to the account, or create one if it doesn't exist
  //  
  let accountId = _LinkOrCreateAccount(ctx, nk, logger, username, linkObject,  accessToken);

  _deleteLinkCode(nk, linkObject, logger);

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

  // Retrieve and return the link code for the device.
  return JSON.stringify({ "success": true });
};


/**
 * Deletes a link code from the storage.
 * @param nk - The Nakama instance.
 * @param linkObject - The link code object to delete.
 * @param logger - The logger instance.
 */
function _deleteLinkCode(nk: nkruntime.Nakama, linkObject: LinkCode, logger: nkruntime.Logger) {
  try {
    nk.storageDelete([{ collection: CollectionMap.linkCode, key: linkObject.code, userId: systemUserId }]);
  } catch (error) {
    logger.error("Failed to delete link code: %s", error);
    throw errInternal(`Failed to delete link code: ${error}`);
  }
}


/**
 * Links or creates an account based on the provided parameters.
 * 
 * @param nk - The Nakama instance.
 * @param username - The username of the account.
 * @param accountId - The ID of the account.
 * @param linkObject - The link code object.
 * @param logger - The logger instance.
 * @param accessToken - The access token.
 * @returns The ID of the linked or created account.
 */
function _LinkOrCreateAccount(ctx: nkruntime.Context, nk: nkruntime.Nakama, logger: nkruntime.Logger, username: string, linkObject: LinkCode,  accessToken: DiscordAccessToken) {
  let accountId = null;

  let users = nk.usersGetUsername([username]);
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
  refreshDiscordLink(ctx, nk, logger, accountId, accessToken, true)

  return accountId;
}


/**
 * Parses and validates the request payload.
 * 
 * @param logger - The logger object for logging debug and error messages.
 * @param payload - The payload string to be parsed and validated.
 * @param nk - The Nakama object for accessing Nakama runtime functionality.
 * @returns An object containing the parsed and validated oauthCode, oauthRedirectUrl, and linkObject.
 * @throws {nkruntime.Error} If the payload is invalid or missing required data.
 */
function _validateLinkRequest(logger: nkruntime.Logger, payload: string, nk: nkruntime.Nakama) {
  var deviceLinkCode = null;
  var oauthCode = null;
  var oauthRedirectUrl = null;
  try {
    logger.debug("Payload:", payload);
    let data = JSON.parse(payload);

    deviceLinkCode = data.linkCode;
    oauthCode = data.oauthCode;
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

  // Validate the deviceId and throw an error if missing.
  if (!deviceLinkCode || deviceLinkCode.length != 4) {
    throw {
      message: `Link code is invalid/missing from payload: ${payload}`,
      code: nkruntime.Codes.INVALID_ARGUMENT
    } as nkruntime.Error;
  }

  // Sanitize the linkCode and make sure it's only capital letters
  deviceLinkCode = deviceLinkCode.toUpperCase().replace(/[^A-Z]/g, '').slice(0, 4);

  // Retrieve the linkCode and deviceId from storage.
  let linkObject = {} as LinkCode;
  try {
    linkObject = getStorageObject(nk, logger, CollectionMap.linkCode, deviceLinkCode, systemUserId) as LinkCode;
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
  return { oauthCode, oauthRedirectUrl, linkObject };
}

function refreshDiscordLink(ctx: nkruntime.Context, nk: nkruntime.Nakama, logger: nkruntime.Logger, userId: string, accessToken?: DiscordAccessToken, noRefresh = false) {
  // retrieve the discord access_token storage object
  let collection = CollectionMap.discord;
  let key = CollectionMap.discordAccessToken;

  if (!accessToken) {
    try {
      accessToken = getStorageObject(nk, logger, collection, key, userId) as DiscordAccessToken;
    } catch (error) {
      logger.error("Failed to retrieve discord/accessToken: %s", error);
      throw errInternal(`Failed to retrieve discord/accessToken: ${error}`);
    }
  }

  // refresh the access token
  let newAccessToken = {} as DiscordAccessToken;
  if (noRefresh) {
    newAccessToken = accessToken;
  } else {
    try {
      newAccessToken = discordRefreshAccessToken(ctx, nk, logger, accessToken);
    } catch (error) {
      logger.error("Failed to refresh discord access token: %s", error.message);
      throw errInternal(`Failed to refresh discord access token: ${error}`);
    }
  }

  // update the storage object
  try {
    nk.storageWrite([
      {
        collection,
        key,
        value: newAccessToken,
        userId,
        permissionRead: StoragePermissions.NO_READ,
        permissionWrite: StoragePermissions.NO_WRITE,
      }
    ]);
  } catch (error) {
    logger.error("Failed to update discord/accessToken: %s", error);
    throw errInternal(`Failed to update discord/accessToken: ${error}`);
  }

  // relink custom with new access token
  try {
    nk.unlinkCustom(userId);
  } catch (noterror) {
  }
  try {
    nk.linkCustom(userId, newAccessToken.access_token);
  } catch (error) {
    logger.error("Failed to link discord access token: %s", error);
    throw errInternal(`Failed to link discord access token: ${error}`);
  }
  // Retrieve the full User object
  let user = discordGetCurrentUser(ctx, nk, logger, accessToken);

  // Update the user's info
  let displayName = user.global_name ?? user.username;
  nk.accountUpdateId(userId, user.id, displayName, null, null, null, null);
  // Create a storage object with the discord user data
  try {
    nk.storageWrite([
      {
        collection: CollectionMap.discord,
        key: CollectionMap.discordUser,
        value: user,
        userId,
        permissionRead: StoragePermissions.OWNER_READ,
        permissionWrite: StoragePermissions.NO_WRITE,
      }
    ]);
  } catch (error) {
    logger.error("Failed to update discord/user: %s", error);
    throw errInternal(`Failed to update discord/user: ${error}`);

  }
}


export {
  getDeviceLinkCodeRpc,
  discordLinkDeviceRpc,
  generateLinkCode,
}