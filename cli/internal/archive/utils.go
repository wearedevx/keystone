package archive

import (
	"archive/tar"
	"io"
	"os"

	"github.com/wearedevx/keystone/cli/internal/utils"
)

func tarSetHeaderName(
	name string,
	info os.FileInfo,
	tarball *tar.Writer,
) error {
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = name

	return tarball.WriteHeader(header)
}

func tarCopyContent(path string, tarball *tar.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer utils.Close(file)

	_, err = io.Copy(tarball, file)
	return err
}
