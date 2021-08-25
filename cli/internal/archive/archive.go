package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wearedevx/keystone/cli/internal/utils"
)

// Tar creates a tar archive
// source is a path to a directory to archive
// target is a path to a directory to write the archive to.
// The caller must ensure both paths are within the ctx.Wd
func Tar(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer utils.Close(tarfile)

	tarball := tar.NewWriter(tarfile)
	defer utils.Close(tarball)

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				destinationPath := filepath.Join(
					baseDir,
					strings.TrimPrefix(path, source),
				)

				header.Name = destinationPath
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			/* #nosec */
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer utils.Close(file)

			_, err = io.Copy(tarball, file)
			return err
		})
}

// Untar extracts content from a tar archive
// tarball is a path to the archive
// target is a directory to write to
// The caller must ensure target path in within ctx.Wd
func Untar(tarball, target string) error {
	/* #nosec */
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer utils.Close(reader)

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		/* #nosec */
		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		/* #nosec */
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer utils.Close(file)

		/* #nosec
		 * Shouldn't we prevent decompression bombs?
		 */
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

// Gzip compresses a file with the gzip algorithm.
// source is the orginal file to compress
// target is the output
// The caller must ensure both path are within ctx.Wd
func Gzip(source, target string) error {
	/* #nosec */
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer utils.Close(writer)

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer utils.Close(archiver)

	/* #nosec
	 * Should'n we check for decompression bomb?
	 */
	_, err = io.Copy(archiver, reader)
	return err
}

// UnGzip uncompresses a gzip file
// source is a path to the compressed file
// target ia a path to a file to write the extracted contents to
// Caller must ensure that the target path is withen ctx.Wd
func UnGzip(source, target string) error {
	/* #nosec */
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer utils.Close(reader)

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer utils.Close(archive)

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer utils.Close(writer)

	/* #nosec
	 * Shouldn't we check for decompression bomb?
	 */
	_, err = io.Copy(writer, archive)
	return err
}
