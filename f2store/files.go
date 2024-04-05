package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	elem := DataF1Elem{fileName, begin, end}
	err = AppendDataF1File(groupName, elem)
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
	// update flaa1 file by rewriting it.
	elemsMap, err := ParseDataF1File(dataF1Path)
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
	elemsMap, err := ParseDataF1File(dataF1Path)
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
	err = RewriteF1File(groupName, elemsMap)
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
	elemsMap, err := ParseDataF1File(dataF1Path)
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

func ParseDataF1File(path string) (map[string]DataF1Elem, error) {
	ret := make(map[string]DataF1Elem, 0)
	rawF1File, err := os.ReadFile(path)
	if err != nil {
		return ret, err
	}

	cleanedF1File := strings.ReplaceAll(string(rawF1File), "\r", "")
	partsOfRawF1File := strings.Split(cleanedF1File, "\n\n")
	for _, part := range partsOfRawF1File {
		innerParts := strings.Split(strings.TrimSpace(part), "\n")

		var elem DataF1Elem
		for _, line := range innerParts {
			var colonIndex int
			for i, ch := range line {
				if fmt.Sprintf("%c", ch) == ":" {
					colonIndex = i
					break
				}
			}

			if colonIndex == 0 {
				continue
			}

			optName := strings.TrimSpace(line[0:colonIndex])
			optValue := strings.TrimSpace(line[colonIndex+1:])

			if optName == "data_key" {
				elem.DataKey = optValue
			} else if optName == "data_begin" {
				data, err := strconv.ParseInt(optValue, 10, 64)
				if err != nil {
					return ret, errors.New("data_begin is not of type int64")
				}
				elem.DataBegin = data
			} else if optName == "data_end" {
				data, err := strconv.ParseInt(optValue, 10, 64)
				if err != nil {
					return ret, errors.New("data_end is not of type int64")
				}
				elem.DataEnd = data
			}
		}

		if elem.DataKey == "" {
			continue
		}
		ret[elem.DataKey] = elem
	}

	return ret, nil
}

func AppendDataF1File(groupName string, elem DataF1Elem) error {
	dataPath, _ := f2shared.GetRootPath()
	path := filepath.Join(dataPath, groupName+".flaa1")
	dataF1Handle, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return errors.Wrap(err, "os error")
	}
	defer dataF1Handle.Close()

	out := fmt.Sprintf("data_key: %s\ndata_begin: %d\ndata_end:%d\n\n", elem.DataKey,
		elem.DataBegin, elem.DataEnd)

	_, err = dataF1Handle.Write([]byte(out))
	if err != nil {
		return errors.Wrap(err, "os error")
	}

	return nil
}

func ReadPortionF2File(projName, tableName, name string, begin, end int64) ([]byte, error) {
	dataPath, _ := f2shared.GetRootPath()
	path := filepath.Join(dataPath, projName, tableName, name+".flaa2")
	f2FileHandle, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return []byte{}, errors.Wrap(err, "os error")
	}
	defer f2FileHandle.Close()

	outData := make([]byte, end-begin)
	_, err = f2FileHandle.ReadAt(outData, begin)
	if err != nil {
		return outData, errors.Wrap(err, "os error")
	}

	return outData, nil
}

func RewriteF1File(groupName string, elems map[string]DataF1Elem) error {
	dataPath, _ := f2shared.GetRootPath()
	path := filepath.Join(dataPath, groupName+".flaa1")

	out := "\n"
	for _, elem := range elems {
		out += fmt.Sprintf("data_key: %s\ndata_begin: %d\ndata_end:%d\n\n", elem.DataKey,
			elem.DataBegin, elem.DataEnd)
	}

	err := os.WriteFile(path, []byte(out), 0777)
	if err != nil {
		return errors.Wrap(err, "os error")
	}

	return nil
}
