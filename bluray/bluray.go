package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	isoExt    = ".iso"
	bluRayDir = "BDMV"
)

var (
	mount   = `(Mount-DiskImage -ImagePath '%s' -PassThru | Get-Volume).DriveLetter`
	unMount = `Dismount-DiskImage -ImagePath '%s'`

	mediaPlayer string
	basePath    string
)

func init() {
	flag.StringVar(&mediaPlayer, "p", "", "Media player path.")
	flag.StringVar(&basePath, "m", "", "Path to search for movie.")
	flag.Parse()
}

func main() {
	basePath = strings.TrimRight(basePath, `\"`)
	moviePath := basePath
	var iso bool

	filepath.Walk(basePath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln("Error during walk:", err)
		}

		if filepath.Ext(path) == isoExt {
			moviePath = path
			iso = true
			return io.EOF
		}

		if fi.IsDir() && filepath.Base(path) == bluRayDir {
			moviePath = path
			return io.EOF
		}

		return nil
	})

	if iso {
		playISO(moviePath)
	} else {
		playFolder(moviePath)
	}

}

func playFolder(path string) {
	runCmd(exec.Command(mediaPlayer, path))
}

func playISO(path string) {
	driveLetter := runCmd(exec.Command("powershell", "-Command", fmt.Sprintf(mount, path)))
	runCmd(exec.Command(mediaPlayer, driveLetter+":"))
	runCmd(exec.Command("powershell", "-Command", fmt.Sprintf(unMount, path)))
}

func runCmd(cmd *exec.Cmd) string {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()

	if err != nil {
		log.Printf("Error while running command. \n%v \n%v", cmd.Args, err)
	}

	return string(bytes.TrimSpace(out))
}
