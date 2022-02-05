package bugle

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailClient struct {
	ctx     context.Context
	client  *http.Client
	service *gmail.Service
}

func NewGmailClient(ctx context.Context) *GmailClient {
	g := GmailClient{ctx: ctx}
	g.setup()
	return &g
}

func (g *GmailClient) SendEmail(msg string, emails ...string) error {
	messageParts := []string{
		"From: 'me'",
		"reply-to: help@bugle.email",
		"To: connormc7@gmail.com",
		"Subject: Testing Gmail API",
		"",
		msg,
	}
	messageString := strings.Join(messageParts, "\r\n")
	messageBytes := base64.StdEncoding.EncodeToString([]byte(messageString))
	_, err := g.service.Users.Messages.Send("", &gmail.Message{Raw: messageBytes}).Do()
	return errors.Wrap(err, "Error sending email")

}

func (g *GmailClient) setup() error {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return err
	}

	g.client, err = g.getClient(g.ctx, config)
	if err != nil {
		return err
	}

	g.service, err = gmail.NewService(g.ctx, option.WithHTTPClient(g.client))
	if err != nil {
		return err
	}

	return err
}

func (g *GmailClient) getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	tokFile := "token.json"
	tok, err := g.tokenFromFile(tokFile)
	if err != nil {
		tok, err = g.getTokenFromWeb(ctx, config)
		g.saveToken(tokFile, tok)
	}

	return config.Client(ctx, tok), err
}

func (g *GmailClient) getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Auth url(soon to deprecate): %s\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	return config.Exchange(ctx, authCode)
}

func (g *GmailClient) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (g *GmailClient) saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
