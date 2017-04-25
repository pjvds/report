package main

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	_, err := LoadConfig()

	if err != nil {
		bootstrap()
		return

		//log.Fatalf("failed to load config: %v\n", err)
	}

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

func bootstrap() error {
	bootstrapId := xid.New()
	url := fmt.Sprintf("https://slack.com/oauth/authorize?scope=incoming-webhook&client_id=158986125361.158956389232&state=%v", bootstrapId)
	if err := open.Run(url); err != nil {
		return err
	}

	return nil
}
