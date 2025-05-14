package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/saenuma/files209/internal"
)

func listGroups(w http.ResponseWriter, r *http.Request) {
	rootPath, _ := internal.GetRootPath()

	dirFIs, err := os.ReadDir(rootPath)
	if err != nil {
		internal.PrintError(w, err)
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
	groupName := r.PathValue("group")

	rootPath, _ := internal.GetRootPath()

	delete(groupMutexes, groupName)
	os.RemoveAll(filepath.Join(rootPath, groupName+".flaa1"))
	os.RemoveAll(filepath.Join(rootPath, groupName+".flaa2"))

	fmt.Fprint(w, "ok")
}
