package gitignorehelper

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/wearedevx/keystone/cli/internal/utils"
	. "github.com/wearedevx/keystone/cli/internal/utils"
)

func GitIgnore(wd string, thatPath string) error {
	if IsIgnored(wd, thatPath) {
		return nil
	}

	gitignorePath := path.Join(wd, ".gitignore")
	gitignore, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer gitignore.Close()

	content := fmt.Sprintf("\n%s", thatPath)
	if _, err = gitignore.WriteString(content); err != nil {
		return err
	}

	return nil
}

func GitUnignore(wd string, thatPath string) error {
	if !IsIgnored(wd, thatPath) {
		return nil
	}

	gitignorePath := path.Join(wd, ".gitignore")
	if utils.FileExists(gitignorePath) {
		gitignore, err := os.OpenFile(gitignorePath, os.O_RDONLY, 0644)

		if err != nil {
			return nil
		}

		scanner := bufio.NewScanner(gitignore)
		thatPath = strings.Trim(thatPath, " ")

		lines := make([]string, 0)
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}

			line := scanner.Text()
			if strings.Trim(line, " ") != thatPath {
				lines = append(lines, line)
			}
		}
		gitignore.Close()

		contents := []byte(strings.Join(lines, "\n"))

		ioutil.WriteFile(gitignorePath, contents, 0600)
	}

	return nil
}

func IsIgnored(wd string, thatPath string) bool {
	gitignorePath := path.Join(wd, ".gitignore")

	if FileExists(gitignorePath) {
		f, err := os.OpenFile(gitignorePath, os.O_RDONLY, 0o644)
		if err != nil {
			return false
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()

			if line == thatPath {
				return true
			}
		}

	}

	return false
}
