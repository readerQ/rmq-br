package localio

import (
	"os"
	"path/filepath"

	"github.com/readerQ/rmq-br/rabbit"
)

type MessageReader struct {
	path        string
	contentType string

	files       []string
	index       int
	readerError error
}

func NewMessageReader(path, mask, contentType string) *MessageReader {
	mr := MessageReader{
		path:        path,
		contentType: contentType,
	}

	mr.files, mr.readerError = WalkMatch(path, mask)

	return &mr
}

func (mr *MessageReader) ReadMessage() (rabbit.Message, bool, error) {

	if mr.readerError != nil {
		return rabbit.Message{}, false, mr.readerError
	}

	if mr.index >= len(mr.files) {
		return rabbit.Message{}, false, nil
	}

	body, err := os.ReadFile(mr.files[mr.index])

	if err != nil {
		mr.readerError = err
		return rabbit.Message{}, false, mr.readerError
	}
	mr.index++

	return rabbit.Message{Body: body, ContentType: mr.contentType, Index: mr.index}, true, nil
}

// https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
