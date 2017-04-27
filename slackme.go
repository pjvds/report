package slackme

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/rs/xid"
	spin "github.com/tj/go-spin"
)

type Context struct {
	Email    string
	Token    string
	UserName string
	TeamName string

	Channels []Channel
}

type Channel struct {
	Name       string
	WebhookUrl string
}

func (this Channel) Post(message string) error {
	body, _ := json.Marshal(map[string]interface{}{
		"text": message,
	})
	response, err := http.Post(this.WebhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	return nil
}

func (this *Context) HasChannels() bool {
	return len(this.Channels) > 0
}

func (this *Context) ChannelByName(name string) (Channel, bool) {
	for _, channel := range this.Channels {
		if strings.EqualFold(name, channel.Name) {
			return channel, true
		}
	}
	return Channel{}, false
}

func (this *Context) AddChannel() (Channel, bool, error) {
	addChannelID := xid.New().String()
	addUrl := fmt.Sprintf("https://slack.com/oauth/authorize?scope=incoming-webhook&client_id=158986125361.158956389232&state=%v&redirect_uri=%v",
		url.QueryEscape(addChannelID), url.QueryEscape("https://slackme.org/a/register"))
	completeURL := fmt.Sprintf("https://slackme.org/a/completion/channel/%v", url.QueryEscape(addChannelID))

	if err := exec.Command("open", addUrl).Run(); err != nil {
		println("open the following url in a browser:\n\n\r" + addUrl)
	}

	s := spin.New()
	for {
		fmt.Printf("\rwaiting for completion %s", s.Next())
		response, err := http.Get(completeURL)
		if err != nil {
			return Channel{}, false, err
		}

		if response.StatusCode == http.StatusOK {
			fmt.Printf("\r")
			body, err := gabs.ParseJSONBuffer(response.Body)
			if err != nil {
				return Channel{}, false, err
			}

			channel := Channel{
				Name:       body.Path("name").Data().(string),
				WebhookUrl: body.Path("webhookURL").Data().(string),
			}
			this.Channels = append(this.Channels, channel)

			if err := this.Save(); err != nil {
				return Channel{}, false, err
			}

			return channel, true, nil
		}
	}
}

func (this *Context) Login() error {
	signinID := xid.New().String()
	authUrl := "https://slack.com/oauth/authorize?scope=identity.basic,identity.email,identity.team,identity.avatar&client_id=158986125361.158956389232&state=" + url.QueryEscape(signinID) + "&redirect_uri=" + url.QueryEscape("https://slackme.org/a/authenticate")
	authCompleteURL := fmt.Sprintf("https://slackme.org/a/completion/authentication/%v", url.QueryEscape(signinID))

	if err := exec.Command("open", authUrl).Run(); err != nil {
		println("open the following url in a browser:\n\n\r" + authUrl)
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
			this.Channels = make([]Channel, 0)

			if err := this.Save(); err != nil {
				return err
			}

			return nil
		}
	}
}

func (this *Context) NeedsLogin() bool {
	return len(this.Email) == 0 || len(this.Token) == 0
}

func LoadContext() (*Context, error) {
	context := new(Context)

	path := os.ExpandEnv("$SLACKME_FILE")
	if len(path) == 0 {
		path = os.ExpandEnv("$HOME/.slackme")
	}

	file, err := os.Open(path)
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
