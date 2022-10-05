package command

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "gattai",
		Version: "v0.1.0",
		Short: "gattai - a simple CLI to transform and inspect strings",
		Long: `gattai is a super fancy CLI (kidding)

	One can use stringer to modify or inspect strings straight from the terminal`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	return rootCmd
}
