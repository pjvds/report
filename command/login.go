package command

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/pjvds/slackme/context"

	"gopkg.in/urfave/cli.v2"
)

var Login = &cli.Command{
	Name: "login",
	Action: func(cli *cli.Context) error {
		context, err := LoadContext(cli)
		if err != nil {
			log.Fatalf("failed to load context: %v", err)
		}

		if !context.NeedsLogin() {
			if !askForConfirmation(fmt.Sprintf("%v already logged in, are you sure you want to login and loose that user?", context.Email)) {
				return nil
			}
		}

		if err := context.Login(); err != nil {
			log.Fatalf("failed to login: %v", err)
		}

		fmt.Printf("Welcome %v, you rock! ❤️\nYou probably want to add slackme to a channel:\n\tslackme add-channel\n")
		return nil
	},
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
