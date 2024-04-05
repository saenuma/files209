package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/saenuma/files209/f2shared"
)

func listGroups(w http.ResponseWriter, r *http.Request) {
	rootPath, _ := f2shared.GetRootPath()

	dirFIs, err := os.ReadDir(rootPath)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	groups := make([]string, 0)
	for _, dirFI := range dirFIs {
		if strings.HasSuffix(dirFI.Name(), ".flaa2") {
			tmp := strings.ReplaceAll(dirFI.Name(), ".flaa2", "")
			groups = append(groups, tmp)
		}
	}

	retBytes, _ := json.Marshal(groups)
	fmt.Fprint(w, retBytes)
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group"]

	rootPath, _ := f2shared.GetRootPath()

	delete(groupMutexes, groupName)
	os.RemoveAll(filepath.Join(rootPath, groupName+".flaa1"))
	os.RemoveAll(filepath.Join(rootPath, groupName+".flaa2"))

	fmt.Fprint(w, "ok")
}
