package main

import (
	"os"

	"github.com/pjvds/slackme/command"
	"github.com/urfave/cli"
)

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
		command.Login,
		command.Add,
		command.Post,
		command.Exec,
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
