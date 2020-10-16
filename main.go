package main

import (
	"fmt"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack"
)

func Virsh_CallCommands(command []string) ([]byte, error) {
	// var err error
	// var out []byte
	command[0] = "virsh"
	fmt.Println(command)
	switch len(command) {
	// case 0:
	// 	return
	case 1:
		out, err := exec.Command(command[0]).CombinedOutput()
		// fmt.Println(string(out))
		return out, err
	default:
		out, err := exec.Command(command[0], command[1:]...).CombinedOutput()
		// fmt.Println(string(out))
		return out, err
	}


	// if err != nil {
	// 	log.Println(err)
	// }
}

func UrlVerification(w http.ResponseWriter, body []byte) (error) {
		var res *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &res); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(res.Challenge)); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		return nil
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		if err := UrlVerification(w, body); err != nil {
			log .Println(err)
		}
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch event := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			message := strings.Split(event.Text, " ")
			if len(message) < 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			out, err := Virsh_CallCommands(message)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(out))

			// command := message[1]
			// switch command {
			// case "ping":
			if _, _, err := api.PostMessage(event.Channel, slack.MsgOptionText(string(out), false)); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// }
		}
	}
}

func main() {
	log.Println("test-slack")

	http.HandleFunc("/slack/events", StartHandler)

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatal(err)
	}
}

