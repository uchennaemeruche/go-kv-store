package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errUsage = errors.New(`usage: 
	set <key> <value> Set specified key and value
	get <key> 		  Get value for a given key
	`)
)

func main() {
	runner := newRunner()
	if err := runner.run(); err != nil {
		fmt.Println(err)
	}
}

type runner struct {
	database fileDatabase
}

func newRunner() runner {
	return runner{newFileDatabase()}
}

func (r runner) run() error {
	args := os.Args
	if len(args) < 3 {
		return errUsage
	}

	switch args[1] {
	case "set":
		if len(args) != 4 {
			return errUsage
		}
		if err := r.database.Set(args[2], args[3]); err != nil {
			return errUsage
		}

	case "get":
		if len(args) != 3 {
			return errUsage
		}
		value, err := r.database.Get(args[2])
		if err != nil {
			return err
		}
		fmt.Println(value)

	default:
		return errUsage
	}

	return nil
}

type fileDatabase struct {
	filename string
}

func newFileDatabase() fileDatabase {
	return fileDatabase{filename: "database.txt"}
}

func (db fileDatabase) Set(key string, value interface{}) error {
	file, err := os.OpenFile(db.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s:%s\n", key, value)); err != nil {
		return err
	}

	return nil
}

func (db fileDatabase) Get(key string) (string, error) {
	file, err := os.OpenFile(db.filename, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastValue string
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.Split(row, ":")
		if len(parts) < 2 {
			return "", errors.New("invalid record")
		}
		if parts[0] == key {
			lastValue = parts[1]
		}
	}
	return lastValue, nil
}
