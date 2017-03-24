package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	args := os.Args
	if len(os.Args) == 0 {
		log.Fatal("missing command")
	}

	name := args[0]
	arg := make([]string, 0)

	if len(args) > 1 {
		arg = args[1:]
	}

	cmd := exec.Command(name, arg...)
	output, _ := cmd.CombinedOutput()

	os.Stdout.Write(output)
}
