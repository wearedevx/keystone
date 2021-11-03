package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/wearedevx/keystone/cli/internal/crypto"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

// GetBackupPath returns the path to file the backup archive will be written to
func GetBackupPath(wd, projectName, backupName string) string {
	if backupName == "" {
		backupName = path.Join(
			wd,
			fmt.Sprintf(
				`keystone-backup-%s-%d.tar.gz`,
				projectName,
				time.Now().Unix(),
			),
		)
	} else {
		backupName = path.Join(
			wd,
			fmt.Sprintf(`%s.tar.gz`, backupName),
		)
	}

	return backupName
}

// Creates a .tar.gz archive.
// `source` is a path to a directory to archive.
// `wd` is a path to a working directory, used to store the temporary `.tar` file.
// `target` is a path to the target `.tar.gz` file
func Archive(source, wd, target string) (err error) {
	if err = Tar(source, wd); err != nil {
		return err
	}

	if err = Gzip(path.Join(wd, ".keystone.tar"), wd); err != nil {
		return err
	}

	if err = os.Rename(path.Join(wd, ".keystone.tar.gz"), target); err != nil {
		return err
	}

	return nil
}

// Creates a `.tar.gz` archive of the `source` directory,
// into the `traget` file, and encrypts it using `passphrase`
func ArchiveWithPassphrase(source, target, passphrase string) (err error) {
	tempdir, err := os.MkdirTemp("", "ks-archive-*")
	if err != nil {
		return err
	}

	if err = Archive(source, tempdir, target); err != nil {
		return err
	}

	encrypted, err := crypto.EncryptFile(target, passphrase)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(target, encrypted, 0o644); err != nil {
		return err
	}

	return nil
}

// Extracts the .tar.gz file at `archivepath` into the `target` directory
// uing `wd` as temporary directory to hold the `.tar` file
func Extract(archivepath, wd, target string) (err error) {
	err = UnGzip(archivepath, wd)
	if err != nil {
		return err
	}

	temporaryArchivePath := path.Join(wd, ".keystone.tar")

	err = Untar(temporaryArchivePath, target)
	if err != nil {
		return err
	}

	err = os.Remove(temporaryArchivePath)

	return err
}

// Decrypts and extracts an archive using a passphrase.
// `archivepath` is the path to the encrypted archive,
// `target` is the directory where the archive will be extracted, and
// `passphrase` is the passphrase used to decrypt.
func ExtractWithPassphrase(archivepath, target, passphrase string) (err error) {
	temporaryDir, err := os.MkdirTemp("", "ks-archive-*")
	if err != nil {
		return err
	}

	temporaryFile := path.Join(temporaryDir, "decrypted.tar.gz")

	if err = crypto.
		DecryptFile(archivepath, temporaryFile, passphrase); err != nil {
		return err
	}

	if err = Extract(temporaryFile, temporaryDir, target); err != nil {
		return err
	}

	if err = os.RemoveAll(temporaryDir); err != nil {
		return err
	}

	return nil
}

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
		file, err := os.OpenFile(
			path,
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			info.Mode(),
		)
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
