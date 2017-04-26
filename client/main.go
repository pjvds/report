package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"net/http"
	"net/url"

	"github.com/Jeffail/gabs"
	"github.com/rs/xid"
	spin "github.com/tj/go-spin"
)

type Context struct {
	Email    string
	Token    string
	UserName string
	TeamName string

	Channels []*ChannelInfo
}

type ChannelInfo struct {
	ChannelID   string
	ChannelName string
	Token       string
	WebhookURL  string
}

func (this *Context) HasChannels() bool {
	return len(this.Channels) > 0
}

func (this *Context) AddChannel() (bool, error) {
	return false, errors.New("not implemented")
	// addChannelID := xid.New().String()
	// addUrl := fmt.Sprintf("https://slack.com/oauth/authorize?scope=incoming-webhook&client_id=158986125361.158956389232&state=%v&redirect_uri=%v",
	// 	url.QueryEscape(addChannelID), url.QueryEscape("https://slackme.pagekite.me/register"))
	// completeURL := fmt.Sprintf("https://slackme.pagekite.me/%v/channel/%v", url.QueryEscape(this.UserID), url.QueryEscape(addChannelID))
	//
	// if err := exec.Command("open", addUrl).Run(); err != nil {
	// 	return false, err
	// }
	//
	// s := spin.New()
	// for {
	// 	fmt.Printf("\rwaiting for completion %s", s.Next())
	// 	response, err := http.Get(completeURL)
	// 	if err != nil {
	// 		return false, err
	// 	}
	//
	// 	if response.StatusCode == http.StatusOK {
	// 		fmt.Printf("\r")
	// 		body, err := gabs.ParseJSONBuffer(response.Body)
	// 		if err != nil {
	// 			return false, err
	// 		}
	//
	// 		println(body.String())
	// 		return true, nil
	// 	}
	// }
}

func (this *Context) Login() error {
	signinID := xid.New().String()
	authUrl := "https://slack.com/oauth/authorize?scope=identity.basic,identity.email,identity.team,identity.avatar&client_id=158986125361.158956389232&state=" + url.QueryEscape(signinID) + "&redirect_uri=" + url.QueryEscape("https://slackme.pagekite.me/authenticate")
	authCompleteURL := fmt.Sprintf("https://slackme.pagekite.me/completion/authentication/%v", url.QueryEscape(signinID))

	if err := exec.Command("open", authUrl).Run(); err != nil {
		return err
	}

	s := spin.New()
	for {
		fmt.Printf("\rwaiting for completion %s", s.Next())
		response, err := http.Get(authCompleteURL)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusOK {
			fmt.Printf("\r")
			body, err := gabs.ParseJSONBuffer(response.Body)
			if err != nil {
				return err
			}

			this.Email = body.Path("email").Data().(string)
			this.Token = body.Path("token").Data().(string)
			this.UserName = body.Path("name").Data().(string)
			this.TeamName = body.Path("team").Data().(string)

			if err := this.Save(); err != nil {
				return err
			}

			fmt.Printf("welcome %v, you rock! ❤️\n", this.UserName)
			return nil
		}
	}
}

func (this *Context) NeedsLogin() bool {
	return len(this.Email) == 0 || len(this.Token) == 0
}

func LoadContext() (*Context, error) {
	context := new(Context)

	file, err := os.Open(os.ExpandEnv("$HOME/.slackme"))
	if err != nil {
		if os.IsNotExist(err) {
			return context, nil
		}
		return nil, err
	}

	return context, json.NewDecoder(file).Decode(context)
}

func (this *Context) Save() error {
	file, err := os.Create(os.ExpandEnv("$HOME/.slackme"))
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(this)
}

func main() {
	context, err := LoadContext()
	if err != nil {
		log.Fatalf("failed to load context: %v", err)
	}

	if context.NeedsLogin() {
		if err := context.Login(); err != nil {
			log.Fatalf("failed to login: %v", err)
		}
	}

	if !context.HasChannels() {
		context.AddChannel()
	}

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
