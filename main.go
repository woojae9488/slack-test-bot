package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var configToken = "SLACK_CONFIG_TOKEN"
var refreshToken = "SLACK_REFRESH_TOKEN"
var signingSecret = "SLACK_SIGNING_SECRET"

var api = slack.New("SLACK_ACCESS_TOKEN")

var UserStatusChanged = slackevents.EventsAPIType("user_status_changed")

func rotateTokens() {
	freshTokens, err := api.RotateTokens(configToken, refreshToken)
	if err != nil {
		fmt.Printf("error rotating tokens: %v\n", err)
		return
	}

	fmt.Printf("new access token: %s\n", freshTokens.Token)
	fmt.Printf("new refresh token: %s\n", freshTokens.RefreshToken)
	fmt.Printf("new tokenset expires at: %d\n", freshTokens.ExpiresAt)
}

func main() {
	slackevents.EventsAPIInnerEventMapping[UserStatusChanged] = slackevents.UserProfileChangedEvent{}
	// rotateTokens()

	http.HandleFunc("/events-endpoint", handleEvents)
	http.HandleFunc("/oauth", handleOAuth)
	http.HandleFunc("/index", handleIndex)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8010", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func handleOAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	fmt.Printf("%s\n", code)

	res, err := slack.GetOAuthV2ResponseContext(context.Background(), &http.Client{}, "SLACK_CLIENT_ID", "SLACK_CLIENT_SECRET", code, "REDIRECT_URL")
	fmt.Printf("%v, %v\n", res.AccessToken, err)
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("[Event] %v(%v)\n", eventsAPIEvent.Type, eventsAPIEvent.InnerEvent.Type)

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		case *slackevents.UserProfileChangedEvent:
			messageText := func(statusText string) string {
				if statusText == "" {
					return fmt.Sprintf("User `%s` has cleared the status.", ev.User.Name)
				}
				return fmt.Sprintf("User `%s`'s status has changed to `%s`.", ev.User.Name, ev.User.Profile.StatusText)
			}(ev.User.Profile.StatusText)

			a, b, e := api.PostMessage("playground", slack.MsgOptionText(messageText, false))
			fmt.Printf("%v, %v, %v\n", a, b, e)
		}
	}
}
