package login

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

// Register Indexes for the login service
func RegisterIndexes(initializer runtime.Initializer) error {
	// Register the LinkTicket Index that prevents multiple LinkTickets with the same device_id_str
	name := LinkTicketIndex
	collection := LinkTicketCollection
	key := ""                                                        // Set to empty string to match all keys instead
	fields := []string{"game_user_id_token", "nk_device_auth_token"} // index on these fields
	maxEntries := 1000000
	indexOnly := false

	if err := initializer.RegisterStorageIndex(name, collection, key, fields, maxEntries, indexOnly); err != nil {
		return err
	}

	// Register the IP Address index for looking up user's by IP Address
	name = IpAddressIndex
	collection = XPlatformIdStorageCollection
	key = ""                               // Set to empty string to match all keys instead
	fields = []string{"client_ip_address"} // index on these fields
	maxEntries = 1000000
	indexOnly = false

	err := initializer.RegisterStorageIndex(name, collection, key, fields, maxEntries, indexOnly)
	if err != nil {
		return err
	}
	return nil
}

func BeforeAuthenticateCustom(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, in *api.AuthenticateCustomRequest) (*api.AuthenticateCustomRequest, error) {
	// Get third-party API URL from the runtime context
	vars, _ := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)

	// Parse the auth token into a DiscordAccessToken
	discordAccessToken := &DiscordAccessToken{}

	// Refresh the access token
	if err := discordAccessToken.Refresh(vars["DISCORD_CLIENT_ID"], vars["DISCORD_CLIENT_SECRET"]); err != nil {
		logger.Warn("error refreshing DiscordAccessToken: %v", err)
		return in, runtime.NewError("error refreshing DiscordAccessToken", 13)
	}

	return in, nil
}

// Update the user's Nakama account after authenticating
func AfterAuthenticateCustom(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateCustomRequest) error {
	logger.Info("Updating Nakama user after authentication")
	vars, _ := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)

	// Get the Nakama user ID from the runtime context
	nakamaUserId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return runtime.NewError("error getting userId", 13)
	}
	// Get the Nakama Account
	nakamaAccount, err := nk.AccountGetId(ctx, nakamaUserId)
	if err != nil {
		logger.Warn("error getting nakama user: %v", err)
		return runtime.NewError("error getting nakama user", 13)
	}

	// Parse the auth token into a DiscordAccessToken
	discordAccessToken := &DiscordAccessToken{}

	// Get the Discord user
	discord, err := discordgo.New("Bearer " + discordAccessToken.AccessToken)
	if err != nil {
		logger.Warn("error creating discord session: %v", err)
		return runtime.NewError("error creating discord session", 13)
	}
	defer discord.Close()
	discordUser, err := discord.User("@me")
	if err != nil {
		logger.Warn("error getting discord user: %v", err)
		return runtime.NewError("error getting discord user", 13)
	}

	// Get the Discord guildMember
	guildMember, err := discord.GuildMember(vars["DISCORD_GUILD_ID"], discordUser.ID)
	if err != nil {
		logger.Warn("error getting discord member: %v", err)
		return runtime.NewError("error getting discord member", 13)
	}

	// Construct the user info from the discord user and guildMember
	avatarUrl := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordUser.ID, discordUser.Avatar)
	locale := discordUser.Locale

	displayName := DetermineDisplayName(nakamaAccount, discordUser, guildMember)
	// Check if the displayName matches an existing nakama username
	// If it does, throw a warning and set the displayName to the discord username
	// This is to prevent duplicate displayNames
	users, err := nk.UsersGetUsername(ctx, []string{displayName})
	if err != nil {
		logger.WithField("err", err).Error("Users get username error.")
	} else {
		if len(users) > 0 {
			logger.Warn("displayName: %s already exists as a username. Setting displayName to discord username: %s", displayName, discordUser.Username)
			displayName = discordUser.Username
		}
	}

	// Update the Nakama user
	logger.Info("Updating Nakama user: %v with displayName: %v", nakamaUserId, displayName)
	if err := nk.AccountUpdateId(ctx, nakamaUserId, "", nil, displayName, "", "", locale, avatarUrl); err != nil {
		logger.Warn("error updating nakama user: %v", err)
		return runtime.NewError("error updating nakama user", 13)
	}
	return nil
}
