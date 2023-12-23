package discordbot

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/heroiclabs/nakama-common/runtime"
)

func Bot(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, BotToken string) error {
	s, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		return err
	}
	s.Identify.Intents |= discordgo.IntentAutoModerationExecution
	s.Identify.Intents |= discordgo.IntentMessageContent
	s.Identify.Intents |= discordgo.IntentGuilds
	s.Identify.Intents |= discordgo.IntentGuildMembers
	s.Identify.Intents |= discordgo.IntentGuildBans
	s.Identify.Intents |= discordgo.IntentGuildEmojis
	s.Identify.Intents |= discordgo.IntentGuildWebhooks
	s.Identify.Intents |= discordgo.IntentGuildInvites
	s.Identify.Intents |= discordgo.IntentGuildPresences
	s.Identify.Intents |= discordgo.IntentGuildMessages
	s.Identify.Intents |= discordgo.IntentGuildMessageReactions
	s.Identify.Intents |= discordgo.IntentDirectMessages
	s.Identify.Intents |= discordgo.IntentDirectMessageReactions
	s.Identify.Intents |= discordgo.IntentMessageContent
	s.Identify.Intents |= discordgo.IntentAutoModerationConfiguration
	s.Identify.Intents |= discordgo.IntentAutoModerationExecution

	// list the guilds the bot is in
	s.StateEnabled = true

	s.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		logger.Info("Bot is up")
	})

	return nil
}
