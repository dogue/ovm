package cli

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"ovm/cli/meta"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/schollz/progressbar/v3"
)

func (o *OVM) Install(version TargetVersion, installLsp bool) error {
	downloadReq, err := http.NewRequest("GET", version.ZipUrl, nil)
	if err != nil {
		return err
	}

	downloadReq.Header.Set("User-Agent", "ovm "+meta.VERSION)
	downloadReq.Header.Set("X-Client-Os", runtime.GOOS)
	downloadReq.Header.Set("X-Client-Arch", runtime.GOARCH)

	downloadRes, err := http.DefaultClient.Do(downloadReq)
	if err != nil {
		return err
	}
	defer downloadRes.Body.Close()

	tempDir, err := os.CreateTemp(o.baseDir, "*.zip")
	if err != nil {
		return err
	}
	defer tempDir.Close()
	defer os.RemoveAll(tempDir.Name())

	versionStr := o.Colored(version.Tag, "green")
	pbar := progressbar.DefaultBytes(
		int64(downloadRes.ContentLength),
		fmt.Sprintf("Downloading %s: ", versionStr),
	)

	_, err = io.Copy(io.MultiWriter(tempDir, pbar), downloadRes.Body)
	if err != nil {
		return err
	}

	fmt.Println("\nExtracting...")

	extractedDir, err := o.unzipSource(tempDir.Name())
	if err != nil {
		log.Fatal(err)
	}

	outPath := filepath.Join(o.baseDir, extractedDir)
	newPath := filepath.Join(o.baseDir, version.Tag)

	if _, err = os.Stat(newPath); err == nil {
		if o.Verbose {
			fmt.Printf("Destination directory `%s` exists. Removing it.\n", newPath)
		}

		if err = os.RemoveAll(newPath); err != nil {
			log.Warn("Failed to remove existing directory", "dir", newPath)
		}
	}

	if err := os.Rename(outPath, newPath); err != nil {
		log.Fatal(err)
	}

	if o.Verbose {
		fmt.Printf("Moved extracted directory `%s` to `%s`\n", outPath, newPath)
	}

	fmt.Printf("Building %s...\n", o.Colored("Odin", "cyan"))
	if err := o.buildSource(newPath, "build_odin.sh", nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println(o.Colored("Build successful!\n", "green"))

	odinBinary := filepath.Join(newPath, "odin")
	o.createSymlink(odinBinary)

	o.Config.ActiveVersion = version.Tag
	if err := o.Config.save(); err != nil {
		return err
	}

	if err := o.Config.AddInstalledVersion(version.Tag); err != nil {
		return err
	}

	if installLsp {
		if err := o.installOLS(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done! 🍻")
	return nil
}

func (o *OVM) installOLS() error {
	olsZipUrl := "https://github.com/DanielGavin/ols/archive/refs/heads/master.zip"

	downloadReq, err := http.NewRequest("GET", olsZipUrl, nil)
	if err != nil {
		return err
	}

	downloadReq.Header.Set("User-Agent", "ovm "+meta.VERSION)
	downloadReq.Header.Set("X-Client-Os", runtime.GOOS)
	downloadReq.Header.Set("X-Client-Arch", runtime.GOARCH)

	downloadRes, err := http.DefaultClient.Do(downloadReq)
	if err != nil {
		return err
	}
	defer downloadRes.Body.Close()

	tempDir, err := os.CreateTemp(o.baseDir, "*.zip")
	if err != nil {
		return err
	}
	defer tempDir.Close()
	defer os.RemoveAll(tempDir.Name())

	pbar := progressbar.DefaultBytes(
		int64(downloadRes.ContentLength),
		fmt.Sprintf("Downloading OLS"),
	)

	_, err = io.Copy(io.MultiWriter(tempDir, pbar), downloadRes.Body)
	if err != nil {
		return err
	}

	fmt.Println("\nExtracting...")

	extractedDir, err := o.unzipSource(tempDir.Name())
	if err != nil {
		log.Fatal(err)
	}

	outPath := filepath.Join(o.baseDir, extractedDir)
	newPath := filepath.Join(o.baseDir, "ols")

	if _, err = os.Stat(newPath); err == nil {
		if o.Verbose {
			fmt.Printf("Destination directory `%s` exists. Removing it.\n", newPath)
		}

		if err = os.RemoveAll(newPath); err != nil {
			log.Warn("Failed to remove existing directory", "dir", newPath)
		}
	}

	if err := os.Rename(outPath, newPath); err != nil {
		log.Fatal(err)
	}

	if o.Verbose {
		fmt.Printf("Moved extracted directory `%s` to `%s`\n", outPath, newPath)
	}

	fmt.Printf("Building %s...\n", o.Colored("OLS", "cyan"))
	if err := o.buildSource(newPath, "build.sh", []string{fmt.Sprintf("PATH=%s:%s", filepath.Join(o.baseDir, "bin"), os.ExpandEnv("$PATH"))}); err != nil {
		log.Fatal(err)
	}
	fmt.Println(o.Colored("Build successful!\n", "green"))

	olsBinary := filepath.Join(newPath, "ols")
	o.createSymlink(olsBinary)

	return nil
}

func (o *OVM) createSymlink(source string) {
	binDir := filepath.Join(o.baseDir, "bin")
	destination := filepath.Join(binDir, filepath.Base(source))

	if _, err := os.Stat(binDir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
			log.Fatal("Could not create bin directory", err)
		}
	}

	if _, err := os.Lstat(destination); err == nil {
		if o.Verbose {
			fmt.Printf("Symlink exists at `%s`. Removing it.\n", destination)
		}

		if err := os.RemoveAll(destination); err != nil {
			log.Fatal("Failed to remove old symlink", err)
		}
	}

	if err := os.Symlink(source, destination); err != nil {
		log.Fatal(err)
	}

	if o.Verbose {
		fmt.Printf("Symlinked `%s` to `%s`\n", source, destination)
	}
}

func (o *OVM) unzipSource(source string) (string, error) {
	destination := o.baseDir
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

func unzipFile(f *zip.File, destination string) error {
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

func (o *OVM) buildSource(root, buildScript string, env []string) error {
	var cmd exec.Cmd
	cmd.Dir = root
	cmd.Path = filepath.Join(root, buildScript)

	if len(env) > 0 {
		cmd.Env = env
	}

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if o.Verbose {
		fmt.Println(string(output))
	}

	return nil
}
