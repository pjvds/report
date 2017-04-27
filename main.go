package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/Jeffail/gabs"
	"github.com/rs/xid"
	spin "github.com/tj/go-spin"
	"github.com/urfave/cli"
)

type Context struct {
	Email    string
	Token    string
	UserName string
	TeamName string

	Channels []Channel
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

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "slackme"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			EnvVar: "SLACKME_HOST",
			Value:  "https://slackme.org",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "channel",
			EnvVar: "SLACKME_CHANNEL",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name: "login",
			Action: func(cli *cli.Context) error {
				context, err := LoadContext()
				if err != nil {
					log.Fatalf("failed to load context: %v", err)
				}

				if !context.NeedsLogin() {
					if !askForConfirmation(fmt.Sprintf("%v already logged in, are you sure you want to login and loose the context of that user?", context.Email)) {
						return nil
					}
				}

				if err := context.Login(); err != nil {
					log.Fatalf("failed to login: %v", err)
				}

				fmt.Printf("Welcome %v, you rock! ❤️\nYou probably want to add slackme to a channel:\n\tslackme add-channel\n")
				return nil
			},
		},
		cli.Command{
			Name:    "add-channel",
			Aliases: []string{"ac"},
			Action: func(cli *cli.Context) error {
				context, err := LoadContext()
				if err != nil {
					log.Fatalf("failed to load context: %v", err)
				}

				if context.NeedsLogin() {
					log.Fatalf("need login, please run:\n\tslackme login")
				}

				channel, ok, err := context.AddChannel()
				if err != nil {
					log.Fatalf("failed to add channel: %v", err)
				}
				if ok {
					fmt.Printf("channel added succesfully, run to post:\n\tslacke -c '%v' post", channel.Name)
				}
				return nil
			},
		},
		cli.Command{
			Name: "post",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "c",
					EnvVar: "SLACKME_CHANNEL",
				},
			},
			Action: func(cli *cli.Context) error {
				channelName := cli.String("c")
				if len(channelName) == 0 {
					log.Fatalf("missing channel name, please specify it like:\n\n\tslackme post -c '#general' 'hello from slackme!'\n\nOr set the SLACKME_CHANNEL environment variable.")
				}

				context, err := LoadContext()
				if err != nil {
					log.Fatalf("failed to load context: %v", err)
				}

				if context.NeedsLogin() {
					log.Fatalf("not logged in, please run:\n\tslackme login")
				}

				channel, ok := context.ChannelByName(channelName)
				if !ok {
					log.Fatalf("channel not found, please run:\n\tslackme add '%v'", channelName)
				}

				message := cli.Args()[0]
				if message == "-" {
					stdin, err := ioutil.ReadAll(os.Stdin)
					if err != nil {
						log.Fatalf("failed to read from stdin: %v", err)
					}
					message = string(stdin)
				}

				if err := channel.Post(message); err != nil {
					log.Fatalf("failed to post to channel %v: %v", channelName, err)
				}

				return nil
			},
		},
		cli.Command{
			Name: "exec",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "c",
					EnvVar: "SLACKME_CHANNEL",
				},
			},
			Action: func(cli *cli.Context) error {
				channelName := cli.String("c")
				if len(channelName) == 0 {
					log.Fatalf("missing channel name, please specify it like:\n\n\tslackme exec -c '#general' ./backup.sh'\n\nOr set the SLACKME_CHANNEL environment variable.")
				}

				if len(cli.Args()) == 0 {
					log.Fatalf("no command specified, please specify it like:\n\n\tslackme exec -c '#general' ./backup.sh'\n\nOr set the SLACKME_CHANNEL environment variable.")
				}

				context, err := LoadContext()
				if err != nil {
					log.Fatalf("failed to load context: %v", err)
				}

				if context.NeedsLogin() {
					log.Fatalf("not logged in, please run:\n\tslackme login")
				}

				channel, ok := context.ChannelByName(channelName)
				if !ok {
					log.Fatalf("channel not found, please run:\n\tslackme add '%v'", channelName)
				}

				name := cli.Args()[0]
				args := []string{}
				if len(cli.Args()) > 1 {
					args = cli.Args()[1:]
				}
				command := exec.Command(name, args...)
				outputBuffer := new(bytes.Buffer)

				command.Stdout = io.MultiWriter(os.Stdout, outputBuffer)
				command.Stderr = io.MultiWriter(os.Stderr, outputBuffer)

				started := time.Now()
				err = command.Run()

				exitCode := "0"
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.String()
				}

				normalizedArgs := make([]string, len(args))
				for i, arg := range args {
					if strings.Contains(arg, " ") {
						arg = "'" + arg + "'"
					}
					normalizedArgs[i] = arg
				}

				message := fmt.Sprintf("The command *%v* took *%v* and exited with *%v*\n\n```%v```",
					strings.Join(append([]string{name}, normalizedArgs...), " "), time.Since(started), exitCode, string(outputBuffer.Bytes()))

				if err := channel.Post(message); err != nil {
					log.Fatalf("failed to post to channel %v: %v", channelName, err)
				}

				return nil
			},
		},
	}
	app.Run(os.Args)

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
