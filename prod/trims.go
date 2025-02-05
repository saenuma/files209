package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/saenuma/files209/f2shared"
)

func trimF209Files(groupName string) error {

	rootPath, _ := f2shared.GetRootPath()
	tablePath := filepath.Join(rootPath, groupName)
	tmpGroupName := groupName + "_trim_tmp"
	workingTablePath := filepath.Join(rootPath, tmpGroupName)

	if f2shared.DoesPathExists(workingTablePath) {
		os.RemoveAll(workingTablePath)
	}

	// os.MkdirAll(workingTablePath, 0777)

	refF1Path := filepath.Join(tablePath, groupName+".flaa1")
	tmpF2Path := filepath.Join(rootPath, tmpGroupName+".flaa2")
	elemsMap, err := f2shared.ParseDataF1File(refF1Path)
	if err != nil {
		return errors.Wrap(err, "f2shared error")
	}
	dataLumpPath := filepath.Join(rootPath, groupName+".flaa2")

	// trim the data files
	for _, elem := range elemsMap {

		begin := elem.DataBegin
		end := elem.DataEnd

		fileBytes := make([]byte, end-begin)
		if f2shared.DoesPathExists(dataLumpPath) {
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

		newDataElem := f2shared.DataF1Elem{DataKey: elem.DataKey, DataBegin: newBegin, DataEnd: newEnd}
		err = f2shared.AppendDataF1File(tmpGroupName, newDataElem)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

	return nil
}
