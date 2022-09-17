package main

import (
	"flag"
	"fmt"
	"github.com/Wilm0rien/omip/update"
	"github.com/Wilm0rien/omip/util"
	"log"
	"os/exec"
	"runtime"
	"time"

	"os"

	"strings"
)

const (
	AppName  = "omip"
	TmpZip   = "omip.zip"
	GitOwner = "Wilm0rien"
	GitRepo  = "omip"
)

func main() {
	versionFlag := flag.Bool("version", false, "show version string")
	targetExec := flag.String("target", "omip.exe", "abs or rel path to target executable")
	sourceZip := flag.String("source", "omip.zip", "abs or rel path to source zip file")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("%s", util.OmipSoftwareVersion)
		return
	}
	if !util.Exists(*sourceZip) {
		log.Fatalf("could not find zip file %s", *sourceZip)
	}
	appData := util.GetAppDataDir()
	localDir := appData + "/" + AppName
	if !util.Exists(localDir) {
		log.Fatalf("could not find omip dir  %s", localDir)
	}
	tmpExeFullPath := localDir + "/" + "omip.exe"
	updater := update.NewUpdaterObj()
	if updater != nil {
		log.Printf("unzipping %s", *sourceZip)
		extractErr := updater.ExtractExec(*sourceZip, localDir)
		if extractErr != nil {
			log.Fatalf("update failed while extracting %s", extractErr.Error())
		}
		if !util.Exists(tmpExeFullPath) {
			log.Fatalf("could not find extracted exe %s", tmpExeFullPath)
		}
		start := time.Now()
		if util.Exists(*targetExec) {
			for {
				log.Printf("removing old %s", *targetExec)
				removeErr := os.Remove(*targetExec)
				if removeErr != nil {
					urlStrShutdown := fmt.Sprintf("http://localhost:4716/callback?code=shutdown&state=0")
					if util.SendReq(urlStrShutdown) {
						time.Sleep(400 * time.Millisecond)
					} else {
						log.Printf("error removing target execeutable : %s %s", *targetExec, removeErr.Error())
					}
					time.Sleep(100 * time.Millisecond)
				} else {
					break
				}

				elapsed := time.Since(start)
				if elapsed.Milliseconds() > 5000 {
					log.Fatalf("timeout waiting for process to close %s", *targetExec)
				}
			}
		}
		log.Printf("copy %s to %s", tmpExeFullPath, *targetExec)
		_, copyErr := updater.Copy(tmpExeFullPath, *targetExec)
		if copyErr != nil {
			log.Fatalf("error copying from %s to %s error %s", tmpExeFullPath, *targetExec, copyErr.Error())
		}
		log.Printf("removing %s", tmpExeFullPath)
		removeErr2 := os.Remove(tmpExeFullPath)
		if removeErr2 != nil {
			log.Fatalf("error removing extracted exe : %s %s", tmpExeFullPath, removeErr2.Error())
		}
		log.Printf("starting %s", *targetExec)
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
