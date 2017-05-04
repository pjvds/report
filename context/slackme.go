package context

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
	"gopkg.in/urfave/cli.v2"
)

type Context struct {
	Channels map[string]Channel

	path string
}

type Channel struct {
	TeamName    string
	ChannelName string
	WebhookUrl  string
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

func (this *Context) ChannelByName(name string) (Channel, error) {
	found := make([]Channel, 0, 1)

	for id, channel := range this.Channels {
		if strings.EqualFold(name, id) {
			found = append(found, channel)
			continue
		}

		if strings.EqualFold(name, channel.ChannelName) {
			found = append(found, channel)
		}
	}

	if len(found) == 0 {
		return Channel{}, cli.Exit("channel not found, please run `slackme add` to add a channel or run `slackme list` to list all available channels", 155)
	}
	if len(found) > 1 {
		return Channel{}, cli.Exit("channel found in multiple teams, please specify the full name (eq. "+found[0].TeamName+"/"+found[0].ChannelName+")", 155)
	}

	return found[0], nil
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
				TeamName:    body.Path("teamName").Data().(string),
				ChannelName: body.Path("channelName").Data().(string),
				WebhookUrl:  body.Path("webhookURL").Data().(string),
			}

			id := fmt.Sprintf("%v/%v", channel.TeamName, channel.ChannelName)
			this.Channels[id] = channel

			if err := this.Save(); err != nil {
				return Channel{}, false, err
			}

			return channel, true, nil
		}
	}
}

func LoadContext(ctx *cli.Context) (*Context, error) {
	context := &Context{
		Channels: make(map[string]Channel),
		path:     ctx.String("file"),
	}
	if len(context.path) == 0 {
		context.path = os.ExpandEnv("$HOME/.slackme")
	}

	file, err := os.Open(context.path)
	if err != nil {
		if os.IsNotExist(err) {
			return context, nil
		}
		return nil, cli.Exit(fmt.Sprintf("failed to open %v: %v", context.path, err.Error()), 255)
	}

	return context, json.NewDecoder(file).Decode(context)
}

func (this *Context) Save() error {
	file, err := os.Create(this.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(this)
}
