package discorddata

import (
	"time"

	"emperror.dev/errors"
	"github.com/jonas747/discordgo"
	"github.com/karlseguin/ccache"
	"golang.org/x/oauth2"
)

var applicationCache = ccache.New(ccache.Configure().MaxSize(10000).ItemsToPrune(100))

func keySession(raw string) string {
	return "discord_session:" + raw
}

func keyUserInfo(token string) string {
	return "user_info_token:" + token
}

func GetSession(raw string, tokenDecoder func(string) (*oauth2.Token, error)) (*discordgo.Session, error) {
	result, err := applicationCache.Fetch(keySession(raw), time.Minute*10, func() (interface{}, error) {
		decoded, err := tokenDecoder(raw)
		if err != nil {
			return nil, errors.WithStackIf(err)
		}

		session, err := discordgo.New(decoded.Type() + " " + decoded.AccessToken)
		if err != nil {
			return nil, errors.WithStackIf(err)
		}

		return session, nil
	})
	if err != nil {
		return nil, err
	}

	return result.Value().(*discordgo.Session), nil
}

func GetUserInfo(token string, session *discordgo.Session) (*discordgo.User, error) {
	result, err := applicationCache.Fetch(keyUserInfo(token), time.Minute*10, func() (interface{}, error) {
		user, err := session.UserMe()
		if err != nil {
			return nil, errors.WithStackIf(err)
		}

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return result.Value().(*discordgo.User), nil
}