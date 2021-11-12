package ci

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/wearedevx/keystone/cli/internal/archive"
	"github.com/wearedevx/keystone/cli/internal/envfile"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

var pathToVarnameRegexp *regexp.Regexp

func init() {
	pathToVarnameRegexp = regexp.MustCompile(`[^\w]`)
}
func pathToVarname(in string) string {
	inb := []byte(in)
	sep := []byte("_")

	r := pathToVarnameRegexp.ReplaceAll(inb, sep)
	s := string(r)

	s = strings.ToUpper(s)

	return s
}

func getArchiveBuffer(
	ctx *core.Context,
	environmentName string,
) (io.Reader, error) {
	tempdir, err := ioutil.TempDir("", "keystone-archive-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempdir)

	ksfile := new(keystonefile.KeystoneFile).Load(ctx.Wd)

	if err := makeTemporaryDotEnv(ctx, tempdir, environmentName, ksfile); err != nil {
		return nil, err
	}

	if err := copyFilesToTempDir(ctx, tempdir, environmentName, ksfile); err != nil {
		return nil, err
	}

	archiveBuffer, err := createArchive(tempdir, environmentName)
	if err != nil {
		return nil, err
	}

	return archiveBuffer, nil
}

func makeTemporaryDotEnv(
	ctx *core.Context,
	tempdir, environmentName string,
	ksfile *keystonefile.KeystoneFile,
) error {
	cachepath := path.Join(tempdir, ".keystone", "cache")
	archiveDotEnvPath := path.Join(cachepath, environmentName, ".env")
	os.MkdirAll(cachepath, 0o700)

	dotenv := new(envfile.EnvFile).
		Load(ctx.CachedEnvironmentDotEnvPath(environmentName), nil)
	if err := dotenv.Err(); err != nil {
		return err
	}

	archiveDotEnv := new(envfile.EnvFile).Load(archiveDotEnvPath, nil)

	if err := archiveDotEnv.Err(); err != nil {
		return err
	}

	for _, v := range ksfile.Env {
		key := v.Key
		value, ok := dotenv.Get(key)

		if v.Strict && ok {
			archiveDotEnv.Set(key, value)
		}
	}

	if err := archiveDotEnv.Dump().Err(); err != nil {
		return err
	}

	return nil
}

func copyFilesToTempDir(ctx *core.Context, tempdir, environmentName string,
	ksfile *keystonefile.KeystoneFile) error {
	filesdirpath := path.Join(
		tempdir,
		".keystone",
		"cache",
		environmentName,
		"files",
	)
	if err := os.MkdirAll(filesdirpath, 0o700); err != nil {
		return err
	}

	for _, f := range ksfile.Files {
		fp := f.Path
		current := path.Join(
			ctx.CachedEnvironmentFilesPath(environmentName),
			fp,
		)
		if !utils.FileExists(current) {
			if f.Strict {
				return errors.New("required file not found")
			}
			continue
		}

		inArchive := path.Join(filesdirpath, fp)

		if err := utils.CopyFile(current, inArchive); err != nil {
			return err
		}
	}

	return nil
}

func createArchive(tempdir, environmentName string) (io.Reader, error) {
	fileList, err := getFileList(tempdir, environmentName)
	if err != nil {
		return nil, err
	}

	buffer, err := archive.TarFileList(fileList)
	if err != nil {
		return nil, err
	}

	gzipBuffer, err := archive.Gzip(buffer)
	if err != nil {
		return nil, err
	}

	sb := bytes.NewBuffer([]byte{})
	_, err = io.Copy(sb, gzipBuffer)
	if err != nil {
		return nil, err
	}

	return sb, nil
}

func getFileList(
	base string,
	environmentName string,
) ([]utils.FileInfo, error) {
	fileList := make([]utils.FileInfo, 0)
	source := path.Join(base, ".keystone")
	prefix := path.Join(".keystone", "cache", environmentName)

	err := utils.DirWalk(source,
		func(info utils.FileInfo) error {
			if strings.HasPrefix(info.Path, prefix) ||
				info.Path == ".keystone" ||
				info.Path == path.Join(".keystone", "cache") {
				fileList = append(fileList, info)
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func slot(environmentName string, i int) string {
	return fmt.Sprintf(
		"KEYSTONE_%s_SLOT_%d",
		strings.ToUpper(environmentName),
		i+1,
	)
}

func splitString(s string, chunkSize int, nChunks int) ([]string, error) {
	chunks := make([]string, nChunks)

	if chunkSize >= len(s) {
		chunks[0] = s
		return chunks, nil
	}

	c := 0
	currentLen := 0
	currentStart := 0

	for i := range s {
		if currentLen == chunkSize {
			chunks[c] = s[currentStart:i]
			currentLen = 0
			currentStart = i

			c += 1

			if c == len(chunks)-1 {
				break
			}
		}

		currentLen++
	}

	lastChunk := s[currentStart:]
	if len(lastChunk) > chunkSize {
		return nil, fmt.Errorf("keystone archive too big: %d", len(s))
	}

	chunks[c] = lastChunk

	return chunks, nil
}

func base64encode(reader io.Reader) (string, error) {
	sb := new(strings.Builder)

	_, err := io.Copy(sb, reader)
	if err != nil {
		return "", err
	}

	s := base64.StdEncoding.EncodeToString([]byte(sb.String()))

	return s, err
}
