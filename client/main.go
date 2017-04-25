package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"net/http"
	"net/url"

	"github.com/rs/xid"
)

func signin() error {
	signinID := xid.New().String()
	authUrl := "https://slack.com/oauth/authorize?scope=identity.basic,identity.email,identity.team,identity.avatar&client_id=158986125361.158956389232&state=" + url.QueryEscape(signinID) + "&redirect_uri=" + url.QueryEscape("https://slackme.pagekite.me/authenticate")
	if err := exec.Command("open", authUrl).Run(); err != nil {
		return err
	}

	authCompleteURL := fmt.Sprintf("https://slackme.pagekite.me/auth/%v", url.QueryEscape(signinID))

	for {
		response, err := http.Get(authCompleteURL)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusOK {
			body := make(map[string]interface{})
			decoder := json.NewDecoder(response.Body)
			if err := decoder.Decode(&body); err != nil {
				return err
			}
			fmt.Printf("%v", body)
			return nil
		} else {
			println("nope...")
			time.Sleep(time.Second)
			continue
		}
	}

	return nil
}

func main() {
	signin()
	// _, err := LoadConfig()
	//
	// if err != nil {
	// 	bootstrap()
	// 	return
	// }

	// args := os.Args
	// if len(os.Args) <= 1 {
	// 	log.Fatal("missing command")
	// }
	//
	// name := args[1]
	// arg := make([]string, 0)
	//
	// if len(args) > 2 {
	// 	arg = args[2:]
	// }
	//
	// cmd := exec.Command(name, arg...)
	// output, _ := cmd.CombinedOutput()
	//
	// os.Stdout.Write(output)
	//
	// api := slack.New("xoxp-158986125361-159701941572-159692993154-0ad0370934136efee11baf57f66bff62")
	// api.SetDebug(true)
	// params := slack.PostMessageParameters{}
	//
	// _, _, err := api.PostMessage("general", string(output), params)
	//
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
}

func bootstrap() error {
	bootstrapId := xid.New().String()
	authorizeUrl := fmt.Sprintf("https://slack.com/oauth/authorize?redirect_uri=%v&scope=incoming-webhook&client_id=158986125361.158956389232&state=%v", url.QueryEscape("https://slackme.pagekite.me/register"), url.QueryEscape(bootstrapId))

	if err := exec.Command("open", authorizeUrl).Start(); err != nil {
		return err
	}

	return nil
}
