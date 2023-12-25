package login

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/heroiclabs/nakama-common/api"
)

// Select the display name for the user by prioritizing
// the guild nickname, then the discord global displayname,
// then fallback to the discord username
func DetermineDisplayName(nakamaAccount *api.Account, discordUser *discordgo.User, guildMember *discordgo.Member) string {
	// Try the guild nickname, then the discord global displayname, then the discord username
	if guildMember != nil && guildMember.Nick != "" {
		displayName := FilterDisplayName(guildMember.Nick)
		if displayName != "" {
			return displayName
		}
	}

	if discordUser != nil && discordUser.GlobalName != "" {
		displayName := FilterDisplayName(discordUser.GlobalName)
		if displayName != "" {
			return displayName
		}
	}

	if discordUser != nil && discordUser.Username != "" {
		displayName := FilterDisplayName(discordUser.Username)
		if displayName != "" {
			return displayName
		}
	}

	if nakamaAccount != nil && nakamaAccount.User != nil && nakamaAccount.User.Username != "" {
		displayName := FilterDisplayName(nakamaAccount.User.Username)
		if displayName != "" {
			return displayName
		}
	}

	return FilterDisplayName(nakamaAccount.User.Id)
}

func FilterDisplayName(displayName string) string {
	// Use a regular expression to match allowed characters
	re := regexp.MustCompile("[^-_[]A-Za-z]")

	// Filter the string using the regular expression

	filteredUsername := re.ReplaceAllString(displayName, "")

	// two characters minimum
	if len(filteredUsername) < 2 {
		filteredUsername = ""
	}

	// twenty characters maximum
	if len(filteredUsername) > 20 {
		filteredUsername = filteredUsername[:20]
	}

	return filteredUsername
}
