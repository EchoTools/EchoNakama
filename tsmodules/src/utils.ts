import _ from 'lodash';
export const systemUserId = "00000000-0000-0000-0000-000000000000";

export const StoragePermissions = {
  PUBLIC_READ: 2 as nkruntime.ReadPermissionValues,
  OWNER_READ: 1 as nkruntime.ReadPermissionValues,
  NO_READ: 0 as nkruntime.ReadPermissionValues,
  OWNER_WRITE: 1 as nkruntime.WritePermissionValues,
  NO_WRITE: 0 as nkruntime.WritePermissionValues,
};

export const CollectionMap = {
  linkTicket: "Login:linkTicket",
  discord: "Discord",
  discordAccessToken: "accesstoken",
  discordUser: "user",
  echoProfile: 'Profile',
  echoProfileServer: 'server',
  echoProfileClient: 'client',
};


/**
 * Parses a JSON payload string into an object.
 * @param payload - The JSON payload string to parse.
 * @returns The parsed object.
 * @throws {nkruntime.Error} If the payload is invalid JSON.
 */
const parsePayload = function (payload: string): any {
  try {
    return JSON.parse(payload);
  } catch (error) {
    throw {
      message: `Invalid data: ${error}`,
      code: nkruntime.Codes.INVALID_ARGUMENT
    } as nkruntime.Error;
  }
}


/**
 * Retrieves a specific storage object from Nakama's storage.
 *
 * @param nk - The Nakama runtime instance.
 * @param logger - The Nakama logger instance for logging messages.
 * @param collection - The storage collection where the object is stored.
 * @param key - The key of the storage object to retrieve.
 * @param userId - The user ID associated with the storage object.
 * @returns The retrieved storage object value.
 * @throws If the specified storage object is not found, an error with code `nkruntime.Codes.NOT_FOUND` is thrown.
 *          If there's any other error during the storage read operation, it is logged, and the original error is rethrown.
 */
let getStorageObject = function (nk: nkruntime.Nakama, logger: nkruntime.Logger, collection: string, key: string, userId: string) {
  try {
    logger.info(`looking up ${collection}/${key}/${userId}`);
    let objects: nkruntime.StorageObject[] = nk.storageRead([{ collection, key, userId }]);
    logger.info("%s", objects);

    if (objects.length == 0) {
      throw {
        message: `'${collection}/${key}' not found.`,
        code: nkruntime.Codes.NOT_FOUND
      } as nkruntime.Error;
    }

    return objects[0].value;

  } catch (error) {
    logger.error('getStorageObject error: %s', error.message);
    throw error;
  }
}

export {
  getStorageObject,
  parsePayload,
}
