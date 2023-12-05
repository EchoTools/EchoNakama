import _ from 'lodash';
import { LinkCode } from './types';
import { getStorageObject } from './utils';
import { errInternal } from './errors';

const systemUserId = "00000000-0000-0000-0000-000000000000";
const LINKCODE_COLLECTION = "LinkCode";
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
        collection: LINKCODE_COLLECTION,
        key: linkData.code,
        value: linkData,
        userId: systemUserId,
        version: '*',
        permissionRead: 1,
        permissionWrite: 1
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
  var { oauthCode, oauthRedirectUrl, linkObject } = _parseAndValidateRequest(logger, payload, nk);

  // Exchange the oauthCode for the user's access token, and retrieve account metadata
  let { discordUser, accessToken } = _retrieveDiscordAccount(ctx, oauthCode, oauthRedirectUrl, logger, nk);
  
  // Construct the username from the discord user data
  let username = discordUser.id;
  var accountId = null;

  // Link the device to the account, or create one if it doesn't exist
  accountId = _LinkOrCreateAccount(nk, username, accountId, linkObject, logger, accessToken);

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
 * Deletes a link code from the storage.
 * @param nk - The Nakama instance.
 * @param linkObject - The link code object to delete.
 * @param logger - The logger instance.
 */
function _deleteLinkCode(nk: nkruntime.Nakama, linkObject: LinkCode, logger: nkruntime.Logger) {
  try {
    nk.storageDelete([{ collection: LINKCODE_COLLECTION, key: linkObject.code, userId: systemUserId }]);
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
function _LinkOrCreateAccount(nk: nkruntime.Nakama, username: any, accountId: any, linkObject: LinkCode, logger: nkruntime.Logger, accessToken: any) {
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
  return accountId;
}

/**
 * Retrieves the Discord account information by exchanging the OAuth code for an access token
 * and then using the access token to fetch the user's Discord data.
 * @param ctx - The Nakama runtime context.
 * @param oauthCode - The OAuth code obtained from the client.
 * @param oauthRedirectUrl - The redirect URL used during the OAuth process.
 * @param logger - The logger instance for logging debug and error messages.
 * @param nk - The Nakama instance for making HTTP requests.
 * @returns An object containing the user's Discord data and access token.
 * @throws {Error} If the code exchange or user lookup fails.
 */
function _retrieveDiscordAccount(ctx: nkruntime.Context, oauthCode: any, oauthRedirectUrl: any, logger: nkruntime.Logger, nk: nkruntime.Nakama) {
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
  return { discordUser, accessToken };
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
function _parseAndValidateRequest(logger: nkruntime.Logger, payload: string, nk: nkruntime.Nakama) {
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
    linkObject = getStorageObject(nk, logger, LINKCODE_COLLECTION, deviceLinkCode, systemUserId) as LinkCode;
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

export {
  getDeviceLinkCodeRpc,
  discordLinkDeviceRpc,
  generateLinkCode,
}