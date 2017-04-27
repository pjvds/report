package command

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	. "github.com/pjvds/slackme/context"

	"github.com/urfave/cli"
)

var Exec = cli.Command{
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

		normalizedArgs := make([]string, len(args))
		for i, arg := range args {
			if strings.Contains(arg, " ") {
				arg = "'" + arg + "'"
			}
			normalizedArgs[i] = arg
		}

		started := time.Now()
		if err := command.Start(); err != nil {
			message := fmt.Sprintf("Failed to start command *%v*: `%v`",
				strings.Join(append([]string{name}, normalizedArgs...), " "), err.Error())

			if err := channel.Post(message); err != nil {
				log.Fatalf("failed to post to channel %v: %v", channelName, err)
			}
			return nil
		}

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		go func() {
			for sig := range signals {
				if command.Process != nil {
					command.Process.Signal(sig)
				}
			}
		}()

		status := ""
		if err := command.Wait(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				status = exitErr.String()
			} else {
				status = "error " + err.Error()
			}
		}
		signal.Stop(signals)

		message := fmt.Sprintf("The command *%v* took *%v* and exited with *%v*\n\n```%v```",
			strings.Join(append([]string{name}, normalizedArgs...), " "), time.Since(started), status, string(outputBuffer.Bytes()))

		if err := channel.Post(message); err != nil {
			log.Fatalf("failed to post to channel %v: %v", channelName, err)
		}

		return nil
	},
}
