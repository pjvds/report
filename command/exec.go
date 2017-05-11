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
	"text/template"
	"time"

	. "github.com/pjvds/slackme/context"

	"gopkg.in/urfave/cli.v2"
)

type ExecResult struct {
	Command  string
	Output   string
	Duration time.Duration
	DidStart bool
	Err      error
}

var Exec = &cli.Command{
	Name:        "exec",
	Description: "Execute a command and send the result to a Slack channel",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "c",
			EnvVars: []string{"SLACKME_CHANNEL"},
		},
	},
	Action: func(cli *cli.Context) error {
		channelName := cli.String("c")
		if len(channelName) == 0 {
			log.Fatalf("missing channel name, please set the -c flag or SLACKME_CHANNEL environment variable.\n")
		}

		if !cli.Args().Present() {
			log.Fatalf("no command specified, please specify it like:\n\n\tslackme exec -c '#general' ./backup.sh'")
		}

		context, err := LoadContext(cli)
		if err != nil {
			log.Fatalf("failed to load context: %v", err)
		}

		channel, err := context.ChannelByName(channelName)
		if err != nil {
			return err
		}

		name := cli.Args().First()
		args := cli.Args().Tail()
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
		result := ExecResult{
			Command: strings.Join(append([]string{name}, normalizedArgs...), " "),
		}

		if err := command.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)

			result.Err = err
			result.DidStart = false
		} else {
			result.DidStart = true

			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt, os.Kill)
			go func() {
				for sig := range signals {
					if command.Process != nil {
						command.Process.Signal(sig)
					}
				}
			}()

			if err := command.Wait(); err != nil {
				result.Err = err
			}
			result.Duration = time.Since(started)
			result.Output = strings.Trim(outputBuffer.String(), "\t\r\n")

			signal.Stop(signals)
		}

		messageBuffer := new(bytes.Buffer)
		messageTemplate := template.Must(template.New("exec").Funcs(template.FuncMap{
			"keep": func(s string, i int) string {
				runes := []rune(s)
				if len(runes) > i {
					prefix := "[truncate]\n"
					trim := len(runes) - (i - len(prefix))
					return prefix + string(runes[trim:])
				}
				return s
			}}).Parse("```$ {{.Command}}{{if .Output}}\n{{keep .Output 4000}}{{- end}}{{if .Err}}\n{{.Err}}{{- end}}```"))

		if err := messageTemplate.Execute(messageBuffer, result); err != nil {
			log.Fatalf("failed to parse template: %v\n", err)
		}
		if err := channel.Post(messageBuffer.String()); err != nil {
			log.Fatalf("failed to post to channel %v: %v", channelName, err)
		}

		return nil
	},
}
