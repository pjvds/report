package command

import "github.com/urfave/cli"

const (
	CONTEXT_ERR     = iota + 100
	NEED_LOGIN      = iota
	COMMAND_FAILURE = iota
)

var (
	ErrNeedLogin = cli.NewExitError("not logged in, please run:\n\n\tslackme login", NEED_LOGIN)
)
