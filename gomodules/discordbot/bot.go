package discordbot

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/heroiclabs/nakama-common/runtime"
)

var bot *discordgo.Session

func Bot(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, BotToken string) (*discordgo.Session, error) {
	var err error

	logger.Info("Starting bot")
	bot, err = discordgo.New("Bot " + BotToken)

	if err != nil {
		return nil, err
	}
	bot.Identify.Intents |= discordgo.IntentAutoModerationExecution
	bot.Identify.Intents |= discordgo.IntentMessageContent
	bot.Identify.Intents |= discordgo.IntentGuilds
	bot.Identify.Intents |= discordgo.IntentGuildMembers
	bot.Identify.Intents |= discordgo.IntentGuildBans
	bot.Identify.Intents |= discordgo.IntentGuildEmojis
	bot.Identify.Intents |= discordgo.IntentGuildWebhooks
	bot.Identify.Intents |= discordgo.IntentGuildInvites
	bot.Identify.Intents |= discordgo.IntentGuildPresences
	bot.Identify.Intents |= discordgo.IntentGuildMessages
	bot.Identify.Intents |= discordgo.IntentGuildMessageReactions
	bot.Identify.Intents |= discordgo.IntentDirectMessages
	bot.Identify.Intents |= discordgo.IntentDirectMessageReactions
	bot.Identify.Intents |= discordgo.IntentMessageContent
	bot.Identify.Intents |= discordgo.IntentAutoModerationConfiguration
	bot.Identify.Intents |= discordgo.IntentAutoModerationExecution

	// slashcommands

	RegisterSlashCommands(bot, logger, nk)

	// list the guilds the bot is in
	bot.StateEnabled = true

	bot.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		logger.Info("Bot is up")
	})

	err = bot.Open()
	if err != nil {
		return nil, err
	}

	return bot, nil
}
