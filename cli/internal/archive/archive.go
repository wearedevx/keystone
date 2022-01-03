package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/wearedevx/keystone/cli/internal/crypto"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

var l *log.Logger

func init() {
	l = log.New(log.Writer(), "[Archive] ", 0)
}

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
// `target` is a path to the target `.tar.gz` file
func Archive(source, target string) (err error) {
	l.Printf("Archiving %s to %s\n", source, target)

	tarBuffer, err := Tar(source)
	if err != nil {
		return err
	}

	l.Println("Tar OK")

	gzipBuffer, err := Gzip(tarBuffer)
	if err != nil {
		return err
	}

	l.Println("Gzip OK")

	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer utils.Close(file)

	_, err = io.Copy(file, gzipBuffer)
	if err != nil {
		return err
	}

	l.Println("Write OK")

	return nil
}

// Creates a `.tar.gz` archive of the `source` directory,
// into the `traget` file, and encrypts it using `passphrase`
func ArchiveWithPassphrase(source, target, passphrase string) (err error) {
	if err = Archive(source, target); err != nil {
		return err
	}

	l.Printf("Encrypt with %s", passphrase)

	encrypted, err := crypto.EncryptFile(target, passphrase)
	if err != nil {
		l.Fatalln("  FAIL")
		return err
	}
	l.Println("  OK")

	if err = ioutil.WriteFile(target, encrypted, 0o644); err != nil {
		return err
	}

	return nil
}

// Extracts the contents of a tar.gz archive into the `target` directory
func Extract(archive io.Reader, target string) (err error) {
	l.Printf("Extracting to %s", target)

	tarArchive, err := UnGzip(archive)
	if err != nil {
		l.Println("  FAIL")
		return err
	}

	err = Untar(tarArchive, target)
	if err != nil {
		l.Println("  FAIL")
		return err
	}

	l.Println("  OK")
	return err
}

// Decrypts and extracts an archive using a passphrase.
// `archivepath` is the path to the encrypted archive,
// `target` is the directory where the archive will be extracted, and
// `passphrase` is the passphrase used to decrypt.
func ExtractWithPassphrase(archivepath, target, passphrase string) (err error) {
	l.Printf("Decrypt with passsphrase %s", passphrase)

	decrypted, err := crypto.DecryptFile(archivepath, passphrase)
	if err != nil {
		l.Println("  FAIL")
		return err
	}
	l.Println("  OK")

	if err = Extract(decrypted, target); err != nil {
		return err
	}

	return nil
}

// Tar creates a tar archive
// source is a path to a directory to archive
func Tar(source string) (_ io.ReadWriter, err error) {
	fileList := make([]utils.FileInfo, 0)
	err = utils.DirWalk(source,
		func(info utils.FileInfo) error {
			fileList = append(fileList, info)
			l.Println(info.Path)

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return TarFileList(fileList)
}

// TarFileList function creates a tar archive from a list of files.
// the `fileList` may be optained through `utils.DirWalk()`
// target is the path to the output tarball, without the file extension
func TarFileList(fileList []utils.FileInfo) (_ io.ReadWriter, err error) {
	buffer := bytes.NewBuffer([]byte{})

	tarball := tar.NewWriter(buffer)
	defer utils.Close(tarball)

	for _, fileInfo := range fileList {
		if err = tarSetHeaderName(
			fileInfo.Path,
			fileInfo.Info,
			tarball,
		); err != nil {
			break
		}

		if !fileInfo.IsDir {
			if err = tarCopyContent(fileInfo.FullPath, tarball); err != nil {
				break
			}
		}
	}

	return buffer, err
}

// Untar extracts content from a tar archive.
// `target` is the path to a directory to write the extracted files to
func Untar(tarball io.Reader, target string) error {
	tarReader := tar.NewReader(tarball)

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
		l.Println(path)

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

		/* #nosec */
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

// Gzip compresses data with the gzip algorithm.
func Gzip(reader io.Reader) (_ io.Reader, err error) {
	out := bytes.Buffer{}

	archiver := gzip.NewWriter(&out)
	defer utils.Close(archiver)

	/* #nosec
	 * Should'n we check for decompression bomb?
	 */
	_, err = io.Copy(archiver, reader)

	return &out, err
}

// UnGzip uncompresses gzipped data
func UnGzip(source io.Reader) (io.Reader, error) {
	/* #nosec */
	archive, err := gzip.NewReader(source)
	if err != nil {
		return nil, err
	}
	defer utils.Close(archive)

	buffer := bytes.NewBuffer([]byte{})
	/* #nosec
	 * Shouldn't we check for decompression bomb?
	 */
	_, err = io.Copy(buffer, archive)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}

	return buffer, nil
}
