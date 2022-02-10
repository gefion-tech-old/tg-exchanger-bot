package main

import (
	"github.com/gefion-tech/tg-exchanger-bot/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewRootCmd().Execute())
}
