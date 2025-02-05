package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/saenuma/files209/f2shared"
)

func writeFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group"]
	dataB64 := r.FormValue("dataB64")
	fileName := r.FormValue("name")
	rootPath, _ := f2shared.GetRootPath()

	err := nameValidate(groupName)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	createTableMutexIfNecessary(groupName)
	groupMutexes[groupName].Lock()
	defer groupMutexes[groupName].Unlock()

	dataBytes, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	dataLumpPath := filepath.Join(rootPath, groupName+".flaa2")
	var begin int64
	var end int64
	if f2shared.DoesPathExists(dataLumpPath) {
		dataLumpHandle, err := os.OpenFile(dataLumpPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			f2shared.PrintError(w, err)
			return
		}
		defer dataLumpHandle.Close()

		stat, err := dataLumpHandle.Stat()
		if err != nil {
			f2shared.PrintError(w, err)
			return
		}

		size := stat.Size()
		dataLumpHandle.Write(dataBytes)
		begin = size
		end = int64(len(dataBytes)) + size
	} else {
		err := os.WriteFile(dataLumpPath, dataBytes, 0777)
		if err != nil {
			f2shared.PrintError(w, err)
			return
		}

		begin = 0
		end = int64(len(dataBytes))
	}

	elem := f2shared.DataF1Elem{DataKey: fileName, DataBegin: begin, DataEnd: end}
	err = f2shared.AppendDataF1File(groupName, elem)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	fmt.Fprint(w, "ok")
}

func readFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group"]
	fileName := vars["name"]
	rootPath, _ := f2shared.GetRootPath()

	err := nameValidate(groupName)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	createTableMutexIfNecessary(groupName)
	groupMutexes[groupName].RLock()
	defer groupMutexes[groupName].RUnlock()

	dataF1Path := filepath.Join(rootPath, groupName+".flaa1")
	elemsMap, err := f2shared.ParseDataF1File(dataF1Path)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	dataLumpPath := filepath.Join(rootPath, groupName+".flaa2")
	if _, ok := elemsMap[fileName]; !ok {
		f2shared.PrintError(w, errors.New("file doesn't exist"))
		return
	}
	begin := elemsMap[fileName].DataBegin
	end := elemsMap[fileName].DataEnd

	retBytes := make([]byte, end-begin)
	var ret string
	if f2shared.DoesPathExists(dataLumpPath) {
		dataLumpHandle, err := os.OpenFile(dataLumpPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			f2shared.PrintError(w, err)
			return
		}
		defer dataLumpHandle.Close()

		dataLumpHandle.ReadAt(retBytes, begin)
		ret = base64.StdEncoding.EncodeToString(retBytes)
	}

	fmt.Fprint(w, ret)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group"]
	fileName := vars["name"]
	rootPath, _ := f2shared.GetRootPath()

	err := nameValidate(groupName)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	createTableMutexIfNecessary(groupName)
	groupMutexes[groupName].Lock()
	defer groupMutexes[groupName].Unlock()

	dataF1Path := filepath.Join(rootPath, groupName+".flaa1")
	// update flaa1 file by rewriting it.
	elemsMap, err := f2shared.ParseDataF1File(dataF1Path)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	dataLumpPath := filepath.Join(rootPath, groupName+".flaa2")
	begin := elemsMap[fileName].DataBegin
	end := elemsMap[fileName].DataEnd

	nullData := make([]byte, end-begin)

	if f2shared.DoesPathExists(dataLumpPath) {
		dataLumpHandle, err := os.OpenFile(dataLumpPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			f2shared.PrintError(w, err)
			return
		}
		defer dataLumpHandle.Close()

		dataLumpHandle.WriteAt(nullData, begin)
	}

	// rewrite index
	err = f2shared.RewriteF1File(groupName, elemsMap)
	if err != nil {
		f2shared.PrintError(w, err)
		return

	}

	fmt.Fprint(w, "ok")

}

func listFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group"]

	createTableMutexIfNecessary(groupName)
	groupMutexes[groupName].RLock()
	defer groupMutexes[groupName].RUnlock()

	rootPath, _ := f2shared.GetRootPath()
	dataF1Path := filepath.Join(rootPath, groupName+".flaa1")
	// update flaa1 file by rewriting it.
	elemsMap, err := f2shared.ParseDataF1File(dataF1Path)
	if err != nil {
		f2shared.PrintError(w, err)
		return
	}

	ret := make(map[string]int64)

	for _, elem := range elemsMap {
		ret[elem.DataKey] = elem.DataEnd - elem.DataBegin
	}

	retBytes, _ := json.Marshal(ret)
	fmt.Fprint(w, retBytes)
}
