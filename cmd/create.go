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

	_, err = copyBytesFromTemplateToDest(templates.GetUpTemplatePath(), filepath.Join(migrationDirname, UpFilename))

	if err != nil {
		log.Fatal(err)
	}

	_, err = copyBytesFromTemplateToDest(templates.GetDownTemplatePath(), filepath.Join(migrationDirname, DownFilename))
	// write up / down template. read byte from up / down templates
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("")
	//log.Printf("byte written", nb)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create up / down migration files",
	Args:  cobra.ExactArgs(1),
	Run:   Create,
}
