package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordSecrets struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// DiscordAccessToken represents the Discord access token structure.
type DiscordAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Secrets      DiscordSecrets
}

func (t *DiscordAccessToken) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

func UnmarshalDiscordAccessToken(data []byte, s *DiscordSecrets) (DiscordAccessToken, error) {
	r := DiscordAccessToken{
		Secrets: *s,
	}

	err := json.Unmarshal(data, &r)
	return r, err
}

func (t *DiscordAccessToken) Refresh() (*DiscordAccessToken, error) {
	url := "https://discord.com/api/v10/oauth2/token"
	params := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", t.RefreshToken)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(params))
	if err != nil {
		return t, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", t.Secrets.ClientID, t.Secrets.ClientSecret))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return t, err
	}
	if resp.StatusCode != http.StatusOK {
		return t, fmt.Errorf("discord refresh failed: %s", resp.Status)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}

	return t, nil
}

// DiscordUser represents the Discord user structure.
type DiscordUser struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	Discriminator      string `json:"discriminator"`
	Avatar             string `json:"avatar"`
	DiscordAccessToken *DiscordAccessToken
}

// discordGetCurrentUser gets the current Discord user.
func (u *DiscordUser) User() (*DiscordUser, error) {
	url := "https://discord.com/api/v10/users/@me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return u, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+u.DiscordAccessToken.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return u, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return u, err
	}

	if resp.StatusCode != http.StatusOK {
		return u, fmt.Errorf("discord user lookup failed: %s", resp.Status)
	}

	return u, nil
}
