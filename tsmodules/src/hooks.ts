/*
let userAddFriendLevelCheck: nkruntime.BeforeHookFunction<AddFriendsRequest> =
function(ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama, data: nkruntime.AddFriendsRequest): nkruntime.AddFriendsRequest {
	let userId = ctx.userId;

	let users: nkruntime.User[];
	try {
			users = nk.usersGetId([ userId ]);
	} catch (error) {
			logger.error('Failed to get user: %s', error.message);
			throw error;
	}

	// Let's assume we've stored a user's level in their metadata.
	if (users[0].metadata.level < 10) {
			throw Error('Must reach level 10 before you can add friends.');
	}

	// important!
	return data;
};
*/

import { errInternal } from "./errors";

import { refreshDiscordLink } from './linking';

import { DiscordAccessToken } from "./types";
import { getStorageObject } from "./utils";

import { CollectionMap } from "./utils";



// a nakama hook for the authenticatecustom function 
export function refreshDiscordLinkHook(ctx: nkruntime.Context, logger: nkruntime.Logger, nk: nkruntime.Nakama) {
  // read the discord access token from the storage object
  let collection = CollectionMap.discord;
  let key = CollectionMap.discordAccessToken;
  let userId = ctx.userId;
  let accessToken = null;
  try {
    accessToken = getStorageObject(nk, logger, collection, key, userId) as DiscordAccessToken;
  } catch (error) {
    logger.error("Failed to retrieve discord/accessToken: %s", error);
    throw errInternal(`Failed to retrieve discord/accessToken: ${error}`);
  }
    refreshDiscordLink(ctx, nk, logger, userId, accessToken);
}
