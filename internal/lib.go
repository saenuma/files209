package internal

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/saenuma/zazabul"
)

const (
	PORT = 31822
)

var RootConfigTemplate = `// debug can be set to either false or true
// when it is set to true it would print more detailed error messages
debug: false

// in_production can be set to either false or true.
// when set to true, it makes the files209 installation enforce a key
// this key can be gotten from 'files209.prod r' if it has been created with 'files209.prod c'
in_production: false

// port is used while connecting to the database
// changing the port can be used to hide your database during production
port: 31822

`

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "os error")
	}

	var dd string
	dd = os.Getenv("SNAP_COMMON")
	if strings.HasPrefix(dd, "/var/snap/go") || dd == "" {
		dd = filepath.Join(hd, "files209")
		os.MkdirAll(dd, 0777)
	}

	return dd, nil
}

func GetConfigPath() (string, error) {
	rootPath, err := GetRootPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(rootPath, "f209.zconf"), nil
}

func G(objectName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	folders := make([]string, 0)
	folders = append(folders, filepath.Join(homeDir, "f209"))
	folders = append(folders, filepath.Join(homeDir, ".f209"))
	folders = append(folders, os.Getenv("SNAP_COMMON"))

	for _, dir := range folders {
		testPath := filepath.Join(dir, objectName)
		if DoesPathExists(testPath) {
			return testPath
		}
	}

	fmt.Println("Could not find: ", objectName)
	panic("Improperly configured.")
}

func GetSetting(settingName string) string {
	confPath, err := GetConfigPath()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return ""
	}

	conf, err := zazabul.LoadConfigFile(confPath)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	return conf.Get(settingName)
}

func GetKeyStrPath() string {
	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}
	return filepath.Join(rootPath, "f209.keyfile")
}

func PrintError(w http.ResponseWriter, err error) {
	fmt.Printf("%+v\n", err)
	debug := GetSetting("debug")
	if debug == "true" {
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
	} else {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	}
}
