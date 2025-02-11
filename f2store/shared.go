package main

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func createTableMutexIfNecessary(groupName string) {
	_, ok := groupMutexes[groupName]
	if !ok {
		groupMutexes[groupName] = &sync.RWMutex{}
	}
}

func nameValidate(name string) error {
	if strings.Contains(name, ".") || strings.Contains(name, " ") || strings.Contains(name, "\t") ||
		strings.Contains(name, "\n") || strings.Contains(name, ":") || strings.Contains(name, "/") ||
		strings.Contains(name, "~") {
		return errors.New("object name must not contain space, '.', ':', '/', ~ ")
	}

	return nil
}
