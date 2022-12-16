package localio

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/readerQ/rmq-br/rabbit"
)

type FileWriter struct {
	path      string
	tm        string
	Extension string
}

func NewFileWrited(path string) *FileWriter {
	posted24 := time.Now().Format("2006-01-02_15-04-05")
	return &FileWriter{
		path:      path,
		tm:        posted24,
		Extension: "json",
	}
}

func (wr *FileWriter) WriteMessage(msg rabbit.Message) error {

	dataFolder := filepath.Join(wr.path, msg.Queue, wr.tm)
	err := os.MkdirAll(dataFolder, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	name := fmt.Sprintf("%07d.%s", msg.Index, wr.Extension)
	filename := filepath.Join(dataFolder, name)

	err = os.WriteFile(filename, msg.Body, os.ModePerm)
	if err != nil {
		log.Println(fmt.Errorf("file write error: %s (%s)", err.Error(), filename))
	}
	return err

}
