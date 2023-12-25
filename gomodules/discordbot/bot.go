package discordbot

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/heroiclabs/nakama-common/runtime"
)

var s *discordgo.Session

func Bot(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, BotToken string) (*discordgo.Session, error) {
	var err error

	logger.Info("Starting bot")
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		return nil, err
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

	s.StateEnabled = true

	s.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		logger.Info("Bot is up")
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	return s, nil
}
