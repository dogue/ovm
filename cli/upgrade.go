package cli

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"ovm/cli/meta"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/mod/semver"
)

func (o *OVM) Upgrade() error {
	upgradable, tag, err := CanIUpgrade()
	if err != nil {
		return errors.Join(ErrFailedUpgrade, err)
	}

	if !upgradable {
		fmt.Printf("You are already on the latest release (%s) of OVM. Congrats!\n", o.Colored(meta.VERSION, "blue"))
		os.Exit(0)
	} else {
		fmt.Printf("You are on OVM %s. Let's upgrade to %s...\n", meta.VERSION, tag)
	}

	ovmInstallDirEnv, err := o.getInstallDir()
	if err != nil {
		return err
	}

	archive := "tar"
	if runtime.GOOS == "windows" {
		archive = "zip"
	}

	download := fmt.Sprintf("ovm-%s-%s.%s", runtime.GOOS, runtime.GOARCH, archive)
	downloadUrl := fmt.Sprintf("https://github.com/dogue/ovm/releases/latest/download/%s", download)

	resp, err := http.Get(downloadUrl)
	if err != nil {
		return errors.Join(ErrFailedUpgrade, err)
	}
	defer resp.Body.Close()

	tempDownload, err := os.CreateTemp(o.baseDir, fmt.Sprintf("*.%s", archive))
	if err != nil {
		return err
	}
	defer tempDownload.Close()
	defer os.Remove(tempDownload.Name())

	pbar := progressbar.DefaultBytes(
		int64(resp.ContentLength),
		"Upgrading OVM...",
	)

	_, err = io.Copy(io.MultiWriter(tempDownload, pbar), resp.Body)
	if err != nil {
		return err
	}

	ovmPath := filepath.Join(ovmInstallDirEnv, "ovm")
	if err := os.Remove(ovmPath); err != nil {
		if err, ok := err.(*os.PathError); ok {
			if os.IsNotExist(err) {
				log.Debug("Failed to remove file", "path", ovmPath)
			}
		}
	}

	newTemp, err := os.MkdirTemp(o.baseDir, "ovm-upgrade-*")
	if err != nil {
		return errors.Join(ErrFailedUpgrade, err)
	}
	defer os.RemoveAll(newTemp)

	err = untar(tempDownload.Name(), newTemp)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := os.Rename(filepath.Join(newTemp, "ovm"), ovmPath); err != nil {
		return errors.Join(ErrFailedUpgrade, err)
	}

	if err := os.Chmod(ovmPath, 0775); err != nil {
		return errors.Join(ErrFailedUpgrade, err)
	}

	return nil
}

func (o *OVM) getInstallDir() (string, error) {
	defaultPath := filepath.Join(o.baseDir, "self")
	ovmInstallDir, ok := os.LookupEnv("OVM_INSTALL")
	if !ok {
		this, err := os.Executable()
		if err != nil {
			return defaultPath, nil
		}

		isSym, err := isSymlink(this)
		if err != nil {
			return defaultPath, nil
		}

		var finalPath string
		if !isSym {
			finalPath, err = resolveSymlink(this)
			if err != nil {
				return defaultPath, nil
			}
		} else {
			finalPath = this
		}

		modifyable, err := canModifyFile(finalPath)
		if err != nil {
			return "", fmt.Errorf("%q, couldn't determine permissions to modify ovm install", ErrFailedUpgrade)
		}

		if modifyable {
			return filepath.Dir(this), nil
		}

		return "", fmt.Errorf("%q, didn't have permissions to modify ovm install", ErrFailedUpgrade)
	}

	return ovmInstallDir, nil
}

func (o *OVM) extractUpgrade(source, destination string) (string, error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	if err != nil {
		return "", err
	}

	for _, f := range reader.File {
		err := unzipFile(f, destination)

		if o.Verbose {
			fmt.Printf("Extracting %s\n", f.Name)
		}

		if err != nil {
			return "", err
		}
	}

	return strings.TrimSuffix(reader.File[0].Name, "/"), nil
}

func extractFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}

	return nil
}

func CanIUpgrade() (bool, string, error) {
	releases, err := GetGitHubReleases("dogue", "ovm")
	if err != nil {
		return false, "", err
	}

	latest := releases[0]

	if semver.Compare(meta.VERSION, latest.TagName) == -1 {
		return true, latest.TagName, nil
	}

	return false, latest.TagName, nil
}

func isSymlink(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}

func resolveSymlink(symlink string) (string, error) {
	target, err := os.Readlink(symlink)
	if err != nil {
		return "", err
	}
	// Ensure the path is absolute
	absolutePath, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	return absolutePath, nil

}

func untar(tarball, target string) error {
	log.Debug("untar", "tarball", tarball, "target", target)
	reader, err := os.Open(tarball)
	if err != nil {

		return err
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := target + string(os.PathSeparator) + header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			writer, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(writer, tarReader); err != nil {
				return err
			}
			writer.Close()
		}
	}
}
