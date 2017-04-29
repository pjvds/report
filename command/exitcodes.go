package command

import "gopkg.in/urfave/cli.v2"

const (
	CONTEXT_ERR     = iota + 100
	NEED_LOGIN      = iota
	COMMAND_FAILURE = iota
)

var (
	ErrNeedLogin = cli.Exit("not logged in, please run:\n\n\tslackme login", NEED_LOGIN)
)
