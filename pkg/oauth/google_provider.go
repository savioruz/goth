package oauth

import (
	"golang.org/x/oauth2"
)

// GoogleProviderIface defines the methods that a GoogleProvider must implement
type GoogleProviderIface interface {
	GetAuthURL() string
	Exchange(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*GoogleUserInfo, error)
}
