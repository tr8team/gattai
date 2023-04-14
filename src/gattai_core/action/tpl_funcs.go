package action

import (
	"path"
	"strings"
)

func TplFormat() func(string) string {
	return func(content string) string {
		new_content := strings.ReplaceAll(content, "\"","\\\"")
		return strings.ReplaceAll(new_content, "\n","\\n")
	}
}

func TplTempDir(temp_dir string) func(string) string {
	return func(filename string) string {
		return path.Join(temp_dir,filename)
	}
}
