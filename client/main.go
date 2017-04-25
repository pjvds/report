package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"net/http"
	"net/url"

	"github.com/rs/xid"
	spin "github.com/tj/go-spin"
)

func signin() error {
	signinID := xid.New().String()
	authUrl := "https://slack.com/oauth/authorize?scope=identity.basic,identity.email,identity.team,identity.avatar&client_id=158986125361.158956389232&state=" + url.QueryEscape(signinID) + "&redirect_uri=" + url.QueryEscape("https://slackme.pagekite.me/authenticate")
	if err := exec.Command("open", authUrl).Run(); err != nil {
		return err
	}

	authCompleteURL := fmt.Sprintf("https://slackme.pagekite.me/authenticate/%v", url.QueryEscape(signinID))

	s := spin.New()
	for {
		fmt.Printf("\r  \033[36waiting for completion \033[m %s", s.Next())
		response, err := http.Get(authCompleteURL)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusOK {
			println()

			io.Copy(os.Stdout, response.Body)
			return nil
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
