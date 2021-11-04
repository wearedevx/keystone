package gitignorehelper

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/denormal/go-gitignore"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

// GitIgnore function adds `thatPath` to the .gitignore file
func GitIgnore(wd string, thatPath string) error {
	if IsIgnored(wd, thatPath) {
		return nil
	}

	gitignorePath := path.Join(wd, ".gitignore")
	/* #nosec */
	gitignore, err := os.OpenFile(
		gitignorePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return err
	}
	defer utils.Close(gitignore)

	content := fmt.Sprintf("\n%s", thatPath)
	if _, err = gitignore.WriteString(content); err != nil {
		return err
	}

	return nil
}

// GitUnignore function removes `thatPath` from the gitignor file
func GitUnignore(wd string, thatPath string) error {
	if !IsIgnored(wd, thatPath) {
		return nil
	}

	gitignorePath := path.Join(wd, ".gitignore")
	if utils.FileExists(gitignorePath) {
		/* #nosec */
		gitignore, err := os.OpenFile(gitignorePath, os.O_RDONLY, 0o644)
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
		utils.Close(gitignore)

		contents := []byte(strings.Join(lines, "\n"))

		ioutil.WriteFile(gitignorePath, contents, 0o600)
	}

	return nil
}

// IsIgnored function returns true if `thatPath` is gitignored
func IsIgnored(wd string, thatPath string) bool {
	gitignorePath := path.Join(wd, ".gitignore")

	if utils.FileExists(gitignorePath) {
		ignore, _ := gitignore.NewFromFile(gitignorePath)

		return ignore.Ignore(thatPath)
	}

	return false
}
