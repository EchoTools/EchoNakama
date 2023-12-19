import { errInternal } from './errors';
import { DiscordAccessToken, DiscordUser } from './types';

export function discordExchangeCode(ctx: nkruntime.Context, nk: nkruntime.Nakama, logger: nkruntime.Logger, oauthCode: string, oauthRedirectUrl: string): DiscordAccessToken {
  var params = `client_id=${ctx.env["DISCORD_CLIENT_ID"]}&` +
    `client_secret=${ctx.env["DISCORD_CLIENT_SECRET"]}&` +
    `code=${oauthCode}&` +
    `grant_type=authorization_code&` +
    `redirect_uri=${oauthRedirectUrl}&` +
    `scope=identify`;
  
  // exchange the oauthCode for the user's access token
  let response = nk.httpRequest("https://discord.com/api/v10/oauth2/token", "post",
    {
      'Accept': 'application/json',
      "Content-Type": "application/x-www-form-urlencoded",
      "Authorization": `${ctx.env["DISCORD_CLIENT_ID"]}:${ctx.env["DISCORD_CLIENT_SECRET"]}`,
    }, params);

  if (response.code != 200) {
    if ("error" in response) {
      logger.error("Discord code exchange failed: %s", response.body);
      throw {
        message: `Discord code exchange failed: ${response.body}`,
        code: nkruntime.Codes.UNAUTHENTICATED
      } as nkruntime.Error;
    }
  }

  try {
    return JSON.parse(response.body) as DiscordAccessToken;

  } catch (error) {
    logger.error("Could not decode discord response body: %s", response.body);
    throw errInternal(`Could not decode discord response body: ${response.body}`);
  }
}


export function discordRefreshAccessToken(ctx: nkruntime.Context, nk: nkruntime.Nakama, logger: nkruntime.Logger, accessToken: DiscordAccessToken): DiscordAccessToken {
  var params = `grant_type=refresh_token&refresh_token=${accessToken.refresh_token}`;


  let response = nk.httpRequest("https://discord.com/api/v10/oauth2/token", "post",
    {
      'Accept': 'application/json',
      "Content-Type": "application/x-www-form-urlencoded",
      "Authorization": `${ctx.env["DISCORD_CLIENT_ID"]}:${ctx.env["DISCORD_CLIENT_SECRET"]}`,
    }, params);

  let responseJson = null;
  try {
    responseJson = JSON.parse(response.body);
  } catch (error) {
    logger.error("Could not decode discord response body: %s", response.body);
    throw errInternal(`Could not decode discord response body: ${response.body}`);
  }

  if (response.code != 200) {
    if ("error" in responseJson) {
      logger.error("Discord code exchange failed: %s", response.body);
      throw {
        message: `Discord code exchange failed: ${response.body}`,
        code: nkruntime.Codes.UNAUTHENTICATED
      } as nkruntime.Error;
    }
  }
  return responseJson as DiscordAccessToken;
}


export function discordGetCurrentUser(ctx: nkruntime.Context, nk: nkruntime.Nakama, logger: nkruntime.Logger, accessToken: DiscordAccessToken): DiscordUser {

  let response = nk.httpRequest("https://discord.com/api/v10/users/@me", "get",
    {
      "Content-Type": "application/x-www-form-urlencoded",
      "Authorization": `Bearer ${accessToken.access_token}`,
      'Accept': 'application/json'
    }, null);

  
  if (response.code != 200) {
    logger.error("Discord user lookup failed: %s", response.body);
    throw errInternal(`Discord user lookup failed: ${response.body}`);
  }

  return JSON.parse(response.body) as DiscordUser;
}

