package main

import (
	"flag"
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
	mount   = []string{"-mount", "dt,"}
	unMount = []string{"-unmount"}

	daemonTools string
	mediaPlayer string
	blurayDrive string
	basePath    string
)

func init() {
	flag.StringVar(&daemonTools, "d", "", "Daemon tools path.")
	flag.StringVar(&mediaPlayer, "p", "", "Media player path.")
	flag.StringVar(&blurayDrive, "b", "", "Daemon tools drive path.")
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
	runCmd(exec.Command(daemonTools, append(mount, blurayDrive+",", path)...))
	runCmd(exec.Command(mediaPlayer, blurayDrive+":"))
	runCmd(exec.Command(daemonTools, append(unMount, blurayDrive)...))
}

func runCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()

	if err != nil {
		log.Printf("Error while running command. \n%v \n%v", cmd.Args, err)
	}
}
