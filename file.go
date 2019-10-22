package main

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type File struct {
	filePath string
}

func NewFile(path string) (*File, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	file := &File{
		filePath: usr.HomeDir + "/.worktime/" + path,
	}

	dir, _ := filepath.Split(file.filePath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)

		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(file.filePath); os.IsNotExist(err) {
		f, err := os.Create(file.filePath)

		if err != nil {
			return nil, err
		}

		defer f.Close()
	}

	return file, nil
}

func (file *File) Append(content string) error {
	f, err := os.OpenFile(file.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(content + "\n"))
	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func (file *File) GetContent() (string, error) {
	f, err := os.OpenFile(file.filePath, os.O_CREATE|os.O_APPEND|os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}

	b := make([]byte, fi.Size())
	_, err = f.Read(b)
	if err != nil {
		f.Close()
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (file *File) GetEvents() ([]*Event, error) {
	content, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(content, "\n")

	var events []*Event
	for _, element := range lines {
		if element != "" {
			event, err := NewEvent(element)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}

	return events, nil
}

func (file *File) GetLastEventAction() (string, error) {
	events, err := file.GetEvents()
	if err != nil {
		return "", err
	}
	if len(events) < 1 {
		return "", errors.New("There is no events yet")
	}
	lastEvent := events[len(events)-1]
	return lastEvent.action, nil
}
