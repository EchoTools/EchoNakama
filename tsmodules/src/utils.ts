import _ from 'lodash';
import { LinkCode } from './types';

/**
 * Generates a random 4-character link code using a specified set of characters.
 * The characters include uppercase letters (excluding I and O) and digits.
 * 
 * @returns {string} A random 4-character link code.
 */
let generateLinkCode = (): string =>  {

  // Define the set of valid characters for the link code
  const characters = 'ABCDEFGHJKLMNPQRSTUVWXYZ';
  
  // Use lodash.range to create an array with 4 elements
  const indices = _.range(4);
  
  // Use lodash.sample to randomly select an index from the array
  const code = indices.map(() => _.sample(characters)).join('');

  return code;
}

let parsePayload = function (payload : string): any {
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
let getStorageObject = function (nk: nkruntime.Nakama, logger: nkruntime.Logger, collection: string, key: string, userId: string): LinkCode {
  try {
    logger.info(`looking up ${collection}/${key}/${userId}`);
    let objects: nkruntime.StorageObject[] = nk.storageRead([{ collection, key, userId }]);
    logger.info("%s", objects);
    if (objects.length == 0) throw {
      message: `'${collection}/${key}' not found.`,
      code: nkruntime.Codes.NOT_FOUND
    } as nkruntime.Error;

    return objects[0].value as LinkCode;

  } catch (error) {
    logger.error('getStorageObject error: %s', error.message);
    throw error;
  }
}

export {
  getStorageObject,
  generateLinkCode,
  parsePayload,
}