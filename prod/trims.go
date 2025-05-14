package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/saenuma/files209/internal"
)

func trimF209Files(groupName string) error {

	rootPath, _ := internal.GetRootPath()
	tablePath := filepath.Join(rootPath, groupName)
	tmpGroupName := groupName + "_trim_tmp"
	workingTablePath := filepath.Join(rootPath, tmpGroupName)

	if internal.DoesPathExists(workingTablePath) {
		os.RemoveAll(workingTablePath)
	}

	// os.MkdirAll(workingTablePath, 0777)

	refF1Path := filepath.Join(tablePath, groupName+".flaa1")
	tmpF2Path := filepath.Join(rootPath, tmpGroupName+".flaa2")
	elemsMap, err := internal.ParseDataF1File(refF1Path)
	if err != nil {
		return errors.Wrap(err, "f2shared error")
	}
	dataLumpPath := filepath.Join(rootPath, groupName+".flaa2")

	// trim the data files
	for _, elem := range elemsMap {

		begin := elem.DataBegin
		end := elem.DataEnd

		fileBytes := make([]byte, end-begin)
		if internal.DoesPathExists(dataLumpPath) {
			dataLumpHandle, err := os.OpenFile(dataLumpPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer dataLumpHandle.Close()

			dataLumpHandle.ReadAt(fileBytes, begin)
		}

		tmpF2PathHandle, err := os.OpenFile(tmpF2Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer tmpF2PathHandle.Close()

		stat, err := tmpF2PathHandle.Stat()

		if err != nil {
			fmt.Println("stats error", err)
			continue
		}

		size := stat.Size()
		tmpF2PathHandle.Write(fileBytes)
		newBegin := size
		newEnd := int64(len(fileBytes)) + size

		newDataElem := internal.DataF1Elem{DataKey: elem.DataKey, DataBegin: newBegin, DataEnd: newEnd}
		err = internal.AppendDataF1File(tmpGroupName, newDataElem)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

	return nil
}
