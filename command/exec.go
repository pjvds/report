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

	"github.com/urfave/cli"
)

type ExecResult struct {
	Command  string
	Output   string
	Duration time.Duration
	DidStart bool
	Err      error
}

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
			// The name "title" is what the function will be called in the template text.
			"keep": func(s string, i int) string {
				runes := []rune(s)
				if len(runes) > i {
					prefix := "[truncate]\n"
					trim := len(runes) - (i - len(prefix))
					return prefix + string(runes[trim:])
				}
				return s
			}}).Parse("```$ {{.Command}}{{if .Output}}\n{{keep .Output 500}}{{- end}}{{if .Err}}\n{{.Err}}{{- end}}```"))

		if err := messageTemplate.Execute(messageBuffer, result); err != nil {
			log.Fatalf("failed to parse template: %v\n", err)
		}
		if err := channel.Post(messageBuffer.String()); err != nil {
			log.Fatalf("failed to post to channel %v: %v", channelName, err)
		}

		return nil
	},
}
