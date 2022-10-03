package main

import (
	"os"
	"fmt"
	"github.com/tr8team/gattai/src/command"
)

func main() {
	rootCmd := command.NewRootCommand()

	rootCmd.AddCommand(command.NewRunCommand());
	rootCmd.AddCommand(command.NewValidateCommand());

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
        os.Exit(1)
    }
}
