package command

import (
	"bytes"
	"io"
	"log"
	"os"

	. "github.com/pjvds/slackme/context"

	"gopkg.in/urfave/cli.v2"
)

var Post = &cli.Command{
	Name:        "post",
	Description: "Post to a Slack channel.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "code",
		},
	},
	ShellComplete: func(c *cli.Context) {
		if context, err := LoadContext(c); err == nil {
			for _, channel := range context.Channels {
				println(channel.ChannelName)
				return
			}
		}
	},
	Action: func(c *cli.Context) error {
		channelName := c.Args().First()
		if len(channelName) == 0 {
			return cli.Exit("missing channel name, \"slackme post\" requires 2 arguments.\nSee 'slackme post --help'.\n\n"+
				"Usage: slackme post [OPTIONS] CHANNEL_NAME MESSAGE", 1)
		}

		message := c.Args().Get(1)
		if len(message) == 0 {
			return cli.Exit("missing message, \"slackme post\" requires 2 arguments.\nSee 'slackme post --help'.\n\n"+
				"Usage: slackme post [OPTIONS] CHANNEL_NAME MESSAGE", 1)
		}

		context, err := LoadContext(c)
		if err != nil {
			return err
		}

		channel, err := context.ChannelByName(channelName)
		if err != nil {
			return err
		}

		if message == "-" {
			buffer := new(bytes.Buffer)

			_, err := io.Copy(io.MultiWriter(buffer, os.Stdout), os.Stdin)
			if err != nil {
				log.Fatalf("failed to read from stdin: %v", err)
			}
			message = buffer.String()
		}

		if c.Bool("code") {
			message = "```" + message + "```"
		}

		if err := channel.Post(message); err != nil {
			log.Fatalf("failed to post to channel %v: %v", channelName, err)
		}

		return nil
	},
}
