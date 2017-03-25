package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"os/exec"
)

func main() {
	args := os.Args
	if len(os.Args) <= 1 {
		log.Fatal("missing command")
	}

	name := args[1]
	arg := make([]string, 0)

	if len(args) > 2 {
		arg = args[2:]
	}

	cmd := exec.Command(name, arg...)
	output, _ := cmd.CombinedOutput()

	os.Stdout.Write(output)

	api := slack.New("xoxp-158986125361-159701941572-159692993154-0ad0370934136efee11baf57f66bff62")
	api.SetDebug(true)
	params := slack.PostMessageParameters{}

	_, _, err := api.PostMessage("general", string(output), params)

	if err != nil {
		log.Fatal(err.Error())
	}
}
