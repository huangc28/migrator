package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/huangc28/migrator/internal/templates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	UpFilename   = "up.go"
	DownFilename = "down.go"
)

func mkdir(path string) error {
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	return nil
}

func copyBytesFromTemplateToDest(src string, dest string) (int64, error) {
	upTplSrc, err := os.Open(src)
	defer upTplSrc.Close()

	if err != nil {
		log.Fatal(err)
	}

	upDest, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer upDest.Close()

	return io.Copy(upDest, upTplSrc)
}

func Create(cmd *cobra.Command, args []string) {
	// current date is used to create folder for placing migrations on a daily basis.
	var migrationID = time.Now().UTC().Format("20060102150405")

	// retrieve the current process working directory
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory path %s", err.Error())
	}

	// if `migration` directory does not exists.
	migrationRootDirname := filepath.Join(path, "migrations")
	//log.Printf("migrationDirname %s", migrationDirname)
	if _, err := os.Stat(migrationRootDirname); os.IsNotExist(err) {

		log.Printf("migrationDirname %s", migrationRootDirname)
		if err := mkdir(migrationRootDirname); err != nil {
			log.Fatalf("failed to create migration directory %s", err.Error())
		}
	}

	// Retrieve migration and create directory with the name given.
	migrationDirname := filepath.Join(migrationRootDirname, fmt.Sprintf("%s_%s", args[0], migrationID))
	if err := mkdir(migrationDirname); err != nil {
		log.Fatalf("failed to create migration directory %s", err.Error())
	}

	boolChan := make(chan bool, 1)
	errChan := make(chan error, 1)
	quitChan := make(chan struct{})

	type MFileInfo struct {
		src  string
		dest string
	}

	mUp := &MFileInfo{
		templates.GetUpTemplatePath(),
		filepath.Join(migrationDirname, UpFilename),
	}

	mDown := &MFileInfo{
		templates.GetDownTemplatePath(),
		filepath.Join(migrationDirname, DownFilename),
	}

	mFileList := [2]*MFileInfo{mUp, mDown}

	for _, mInfo := range mFileList {
		go func(mInfo *MFileInfo) {
			select {
			case <-quitChan:
				return
			default:
				_, err = copyBytesFromTemplateToDest(
					mInfo.src,
					mInfo.dest,
				)

				if err != nil {
					boolChan <- false
					errChan <- err
				}

				boolChan <- true
			}
		}(mInfo)
	}

	for range mFileList {
		if <-boolChan == false {
			close(quitChan)

			log.Fatal(<-errChan)
		}
	}

	log.WithFields(log.Fields{
		"up":   fmt.Sprintf("%s created", mUp.src),
		"down": fmt.Sprintf("%s created", mDown.src),
	}).Info("migration file created!")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create up / down migration files",
	Args:  cobra.ExactArgs(1),
	Run:   Create,
}
