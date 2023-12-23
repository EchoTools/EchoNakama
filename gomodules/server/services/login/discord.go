package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/heroiclabs/nakama-common/runtime"
	"golang.org/x/oauth2"
)

// DiscordAccessToken represents the Discord access token structure.
type DiscordAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (t *DiscordAccessToken) Config() *oauth2.Config {
	return &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

func (t *DiscordAccessToken) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

func UnmarshalDiscordAccessToken(data []byte) (DiscordAccessToken, error) {
	r := DiscordAccessToken{}

	err := json.Unmarshal(data, &r)
	return r, err
}

func (t *DiscordAccessToken) Refresh(clientId string, clientSecret string) error {
	oauthUrl := "https://discord.com/api/v10/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", t.RefreshToken)
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", oauthUrl, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discord refresh failed: %s", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return err
	}

	return nil
}

func GetUser(discord *discordgo.Session) (*discordgo.User, error) {

	url := "https://discord.com/api/v10/users/@me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+discord.Identify.Token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	u := &discordgo.User{}

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return u, err
	}

	if resp.StatusCode != http.StatusOK {
		return u, fmt.Errorf("discord user lookup failed: %s", resp.Status)
	}

	return u, nil
}

func ExchangeCodeForAccessToken(logger runtime.Logger, code string, clientId string, clientSecret string, redirectUrl string) (*DiscordAccessToken, error) {

	conf := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		Scopes:       []string{"identify"},
		RedirectURL:  redirectUrl,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	accessToken := &DiscordAccessToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.Expiry.Second(),
	}

	return accessToken, nil
}
func WriteAccessTokenToStorage(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userId string, accessToken *DiscordAccessToken) error {
	// Create a StorageWrite object to write the access token to storage
	jsonData, err := accessToken.Marshal()
	if err != nil {
		logger.WithField("err", err).Error("error marshalling access token")
		return runtime.NewError(fmt.Sprintf("error marshalling access token: %v", err), StatusInternalError)
	}
	objectIDs := []*runtime.StorageWrite{
		{
			Collection:      DiscordAccessTokenCollection,
			Key:             DiscordAccessTokenKey,
			UserID:          userId,
			Value:           string(jsonData),
			PermissionRead:  0,
			PermissionWrite: 0,
		},
	}
	_, err = nk.StorageWrite(ctx, objectIDs)
	if err != nil {
		logger.WithField("err", err).Error("storage write error.")
	}

	return err
}

func ReadAccessTokenFromStorage(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userId string, clientId string, clientSecret string) (*DiscordAccessToken, error) {
	// Create a StorageRead object to read the access token from storage
	objectIds := []*runtime.StorageRead{{
		Collection: DiscordAccessTokenCollection,
		Key:        DiscordAccessTokenKey,
		UserID:     userId,
	},
	}

	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("storage read error.")
	}
	if len(records) == 0 {
		return nil, nil
	}

	accessToken, err := UnmarshalDiscordAccessToken([]byte(records[0].Value))
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
