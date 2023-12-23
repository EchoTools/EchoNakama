import _ from 'lodash';
import { LinkTicket, DiscordAccessToken } from './types';
import { getStorageObject } from './utils';
import { errBadInput, errInternal } from './errors';
import { discordExchangeCode, discordGetCurrentUser, discordRefreshAccessToken } from './discord';
import { StoragePermissions } from './utils';
import { systemUserId } from './utils';
import { CollectionMap } from './utils';
import { JwtDecode } from 'jwt-decode'

let discordLinkDeviceRpc: nkruntime.RpcFunction = function (ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, payload: string) {

  // Parse and validate the payload data
  var { sessionToken, oauthRedirectUrl, linkObject } = _validateLinkRequest(logger, payload, nk);
  let secretKey = ctx.env["SESSION_ENCRYPTION_KEY"];
  let nkUserId = null;
  try {
    let decoded = jwt.verify(sessionToken, secretKey, { ignoreExpiration: true, complete: true })
    nkUserId = decoded.payload['uid'];

  } catch (error) {
    logger.debug('Token is invalid:', error.message);
    throw errInternal(`Invalid token: ${error}`);
  }

  nk.linkDevice(nkUserId, linkObject.nk_device_auth_token);
  _deleteLinkTicket(nk, linkObject, logger);
  
  let account = {} as nkruntime.Account;
 
  return JSON.stringify({ "success": true });
};


/**
 * Deletes a link code from the storage.
 * @param nk - The Nakama instance.
 * @param linkTicket - The link code object to delete.
 * @param logger - The logger instance.
 */
function _deleteLinkTicket(nk: nkruntime.Nakama, linkTicket: LinkTicket, logger: nkruntime.Logger) {
  try {
    nk.storageDelete([{ collection: CollectionMap.linkTicket, key: linkTicket.link_code, userId: systemUserId }]);
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
 * @param linkTicket - The link code object.
 * @param logger - The logger instance.
 * @param accessToken - The access token.
 * @returns The ID of the linked or created account.
 */
function _linkOrCreateAccount(ctx: nkruntime.Context, nk: nkruntime.Nakama, logger: nkruntime.Logger, username: string, linkTicket: LinkTicket,  accessToken: DiscordAccessToken) {
  let accountId = null;
  let deviceId = linkTicket.nk_device_auth_token;

  let users = nk.usersGetUsername([username]);
  if (users.length == 1) {
    accountId = users[0].userId;
    // Link the device to the account
    nk.linkDevice(accountId, deviceId);
  }

  try {
    let authResult = nk.authenticateDevice(deviceId, null, false);
    logger.debug("Auth result: %s", authResult);
    accountId = authResult.userId;

  } catch (error) {
    logger.debug("%s", error);
  }

  // Authenticate with the device ID, creating the account if it doesn't exist
  try {
    let result = nk.authenticateDevice(deviceId, username, true);
    accountId = result.userId;

  } catch (error) {
    logger.error('Failed to authenticate device (%s) to user %s: %s', deviceId, username, error);
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
  var sessionToken = null;
  var oauthRedirectUrl = null;
  try {
    logger.debug("Payload:", payload);
    let data = JSON.parse(payload);

    deviceLinkCode = data.linkCode;
    sessionToken = data.sessionToken;
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
  let linkTicket = {} as LinkTicket;
  try {
    linkTicket = getStorageObject(nk, logger, CollectionMap.linkTicket, deviceLinkCode, systemUserId) as LinkTicket;
  } catch (error) {
    throw {
      message: `Link code not found: ${error.message}`,
      code: nkruntime.Codes.NOT_FOUND
    } as nkruntime.Error;
  }

  // ensure the payload contains the discord code
  if (!sessionToken) {
    throw {
      message: `session token code is missing from payload: ${payload}`,
      code: nkruntime.Codes.INVALID_ARGUMENT
    } as nkruntime.Error;
  }
  return { sessionToken: sessionToken, oauthRedirectUrl, linkObject: linkTicket };
}

export function refreshDiscordLink(ctx: nkruntime.Context, nk: nkruntime.Nakama, logger: nkruntime.Logger, userId: string, accessToken?: DiscordAccessToken, noRefresh = false) {
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
  /*
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
  */
  var customIdToken = `${newAccessToken.access_token}:${newAccessToken.refresh_token}`;
  // relink custom with new access token
  try {
    nk.unlinkCustom(userId);
  } catch (noterror) {
  }
  try {
    nk.linkCustom(userId, customIdToken);
  } catch (error) {
    logger.error("Failed to link discord access token: %s", error);
    throw errInternal(`Failed to link discord access token: ${error}`);
  }
  return;
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

  discordLinkDeviceRpc,

}
