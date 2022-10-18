package command

import (
	"os"
	"fmt"
	"path"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "gattai",
		Version: "v0.1.1",
		Short: "gattai - a simple CLI to transform and inspect strings",
		Long: `gattai is a super fancy CLI (kidding)

	One can use stringer to modify or inspect strings straight from the terminal`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	return rootCmd
}

func GetGattaiFilePath(gattaifile_path string, default_path string) (string, error){
	result := gattaifile_path

	fileInfo, err := os.Stat(result)
	if err != nil {
		return result, fmt.Errorf("GetGattaiFilePath error: %v", err)
	}

	if fileInfo.IsDir() {
		result = path.Join(gattaifile_path,default_path)
	}

	return result, nil
}
