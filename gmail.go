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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func Gmail(ctx context.Context) *GmailClient {
	g := &GmailClient{ctx: ctx}
	g.setup()
	return g
}

type GmailClient struct {
	ctx     context.Context
	client  *http.Client
	service *gmail.Service
}

func (g *GmailClient) SendEmail(msg string, subs ...Subscription) error {
	if len(subs) == 0 {
		return nil
	}

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
	return err
}

type setupState struct {
	file string
	cont []byte
	conf *oauth2.Config
	err  error
}

func (g *GmailClient) setup() error {
	state := setupState{"credentials.json", []byte{}, nil, nil}
	state.loadCredentials()
	state.parseConfig()
	state.getClient(g)
	state.getService(g)
	return state.err
}

func (s *setupState) loadCredentials() {
	if s.err != nil {
		return
	}

	s.cont, s.err = ioutil.ReadFile("credentials.json")
}

func (s *setupState) parseConfig() {
	if s.err != nil {
		return
	}

	s.conf, s.err = google.ConfigFromJSON(s.cont, gmail.GmailReadonlyScope)
}

func (s *setupState) getClient(g *GmailClient) {
	if s.err != nil {
		return
	}

	g.client, s.err = g.getClient(g.ctx, s.conf)
}

func (s *setupState) getService(g *GmailClient) {
	if s.err != nil {
		return
	}

	g.service, s.err = gmail.NewService(g.ctx, option.WithHTTPClient(g.client))
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
