package main

import (
	"context"
	"database/sql"
	"echonakama/discordbot"
	"echonakama/server"
	"echonakama/server/services/login"

	"github.com/heroiclabs/nakama-common/runtime"
	_ "google.golang.org/protobuf/proto"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	vars, _ := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)

	if err := initializer.RegisterRpc("relay/loginrequest", server.LoginRequestRpc); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("signin/discord", server.DiscordSignInRpc); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("link/device", server.LinkDeviceRpc); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	login.RegisterIndexes(initializer)
	initializer.RegisterBeforeAuthenticateCustom(login.BeforeAuthenticateCustom)

	initializer.RegisterAfterAuthenticateCustom(login.AfterAuthenticateCustom)

	logger.Info("Initialized module.")

	// Start the bot
	discordbot.Bot(ctx, logger, nk, vars["DISCORD_BOT_TOKEN"])
	return nil
}
