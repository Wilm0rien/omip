package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/Wilm0rien/omip/update"
	"github.com/Wilm0rien/omip/util"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"net/http"
	"os"

	"strings"
)

const (
	AppName  = "omip"
	TmpZip   = "omip.zip"
	GitOwner = "Wilm0rien"
	GitRepo  = "omip"
)

type Updater struct {
	LocalDir       string
	TmpZipFullPath string
	TmpExeFullPath string
}

func NewUpdater() (result *Updater) {
	var obj Updater
	appData := util.GetAppDataDir()
	obj.LocalDir = appData + "/" + AppName
	if !util.Exists(obj.LocalDir) {
		util.CreateDirectory(obj.LocalDir)
	}
	obj.TmpZipFullPath = obj.LocalDir + "/" + TmpZip
	if util.Exists(obj.TmpZipFullPath) {
		removeErr := os.Remove(obj.TmpZipFullPath)
		if removeErr != nil {
			log.Printf("error removing zip file: %s %s", obj.TmpZipFullPath, removeErr.Error())
			return nil
		}
	}
	obj.TmpExeFullPath = obj.LocalDir + "/" + AppName + ".exe"
	if util.Exists(obj.TmpExeFullPath) {
		removeErr := os.Remove(obj.TmpExeFullPath)
		if removeErr != nil {
			log.Printf("error removing tmp exe file: %s %s", obj.TmpExeFullPath, removeErr.Error())
			return nil
		}
	}

	return &obj
}

func main() {
	targetExec := flag.String("target", "omip.exe", "abs or rel path to target executable")
	flag.Parse()
	if !util.Exists(*targetExec) {
		log.Fatalf("target exectuable does not exist: %s", *targetExec)
	}
	updater := NewUpdater()
	if updater != nil {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", GitOwner, GitRepo)
		asset := update.GetRelease(url)
		if asset != nil {
			fmt.Printf("found %s\n", asset.TagName)
			downLoadErr := updater.downloadFile(updater.TmpZipFullPath, asset.Url, asset.FileSize)
			if downLoadErr != nil {
				log.Fatalf("updated failed while downloading %s", downLoadErr.Error())
			}

			extractErr := updater.extractExec(updater.TmpZipFullPath, updater.LocalDir)
			if extractErr != nil {
				log.Fatalf("update failed while extracting %s", extractErr.Error())
			}
			if !util.Exists(updater.TmpExeFullPath) {
				log.Fatalf("could not find extracted exe %s", updater.TmpExeFullPath)
			}
			start := time.Now()
			for {
				removeErr := os.Remove(*targetExec)
				if removeErr != nil {
					log.Printf("error removing target execeutable : %s %s", *targetExec, removeErr.Error())
					time.Sleep(100 * time.Millisecond)
				} else {
					break
				}

				elapsed := time.Since(start)
				if elapsed.Milliseconds() > 5000 {
					log.Fatalf("timeout waiting for process to close %s", *targetExec)
				}
			}

			_, copyErr := updater.Copy(updater.TmpExeFullPath, *targetExec)
			if copyErr != nil {
				log.Fatalf("error copying from %s to %s error %s", updater.TmpExeFullPath, *targetExec, copyErr.Error())
			}
			removeErr2 := os.Remove(updater.TmpExeFullPath)
			if removeErr2 != nil {
				log.Fatalf("error removing extracted exe : %s %s", updater.TmpExeFullPath, removeErr2.Error())
			}
			switch runtime.GOOS {
			case "linux":
				log.Fatalf("TODO Linux restart exe")
			case "windows":
				arguments := fmt.Sprintf(`/k %s`, *targetExec)
				cmd := exec.Command("cmd", arguments)
				execErr2 := cmd.Start()
				if execErr2 != nil {
					log.Fatalf("error starting process %s", *targetExec)
				}
			}
		}
	}

}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
	Max   uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)

	fmt.Printf("\rDownloading... %3.2f %s complete", float64(wc.Total)/float64(wc.Max)*100, util.HumanizeNumber(float64(wc.Total)))
}

func (obj *Updater) downloadFile(filepath string, url string, fileSize uint64) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	counter := &WriteCounter{
		Max: fileSize,
	}

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")

	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func (obj *Updater) extractExec(zipfile string, destination string) (err error) {
	reader, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer reader.Close()
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}
	for _, f := range reader.File {
		err := obj.unzipFile(f, destination)
		if err != nil {
			return err
		}
		break
	}
	return
}

func (obj *Updater) Copy(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (obj *Updater) unzipFile(f *zip.File, destination string) error {
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

	if _, err3 := io.Copy(destinationFile, zippedFile); err != nil {
		return err3
	}
	return nil
}
