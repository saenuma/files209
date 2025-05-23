// cli provides a terminal interface to the files209 server.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gookit/color"
	"github.com/saenuma/files209"
	"github.com/saenuma/files209/internal"
)

const VersionFormat = "20060102T150405MST"

func main() {

	if len(os.Args) < 2 {
		color.Red.Println("expected a command. Open help to view commands.")
		os.Exit(1)
	}

	var keyStr string
	inProd := internal.GetSetting("in_production")
	if inProd == "" {
		color.Red.Println("unexpected error. Have you installed  and launched files209?")
		os.Exit(1)
	}
	if inProd == "true" {
		keyStrPath := internal.GetKeyStrPath()
		raw, err := os.ReadFile(keyStrPath)
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}
		keyStr = string(raw)
	} else {
		keyStr = "not-yet-set"
	}
	port := internal.GetSetting("port")
	if port == "" {
		color.Red.Println("unexpected error. Have you installed  and launched files209?")
		os.Exit(1)
	}
	var cl files209.Client

	portInt, err := strconv.Atoi(port)
	if err != nil {
		color.Red.Println("Invalid port setting.")
		os.Exit(1)
	}

	if portInt != internal.PORT {
		cl = files209.NewClientCustomPort("127.0.0.1", keyStr, portInt)
	} else {
		cl = files209.NewClient("127.0.0.1", keyStr)
	}

	err = cl.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--help", "help", "h":
		fmt.Println(`files209 cli provides some utilites for a files209 installation.
Please Run this program from the same server that powers your files209.
Please don't expose your files209 database to the internet for security sake.

Directory Commands:
  pwd   Print working directory. This is the directory where the files needed by any command
        in this cli program must reside.

File(s) Commands:
  wf    Writes a file to files209 server. Expects a groupname and a filepath
  rf    Reads a file from files209 server. Expects a groupname and a filename
  lf    List files. Expects only a groupname.
  df    Delete Files. Expects a groupname and a filename

Groups Commands:
  lg    List Groups. Expects no arguments
  dg    Delete Group. Expects only the name of the group

			`)

	case "pwd":
		if len(os.Args) != 2 {
			color.Red.Println("Expecting no arugments")
			os.Exit(1)
		}
		p, err := internal.GetRootPath()
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}
		fmt.Println(p)

	case "wf":
		if len(os.Args) != 4 {
			color.Red.Println("Expecting the name of the group and the filepath only")
			os.Exit(1)
		}
		groupName := os.Args[2]
		fileName := filepath.Base(os.Args[3])
		dataBytes, err := os.ReadFile(os.Args[3])
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}
		err = cl.WriteFile(groupName, fileName, dataBytes)
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}

	case "rf":
		if len(os.Args) != 4 {
			color.Red.Println("Expecting the name of the group and the file name only")
			os.Exit(1)
		}

		groupName := os.Args[2]
		fileName := os.Args[3]

		data, err := cl.ReadFile(groupName, fileName)
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}

		outPath := filepath.Join(os.TempDir(), fileName)
		os.WriteFile(outPath, data, 0777)

		fmt.Printf("file written to '%s'\n", outPath)

	case "df":
		if len(os.Args) != 4 {
			color.Red.Println("Expecting the name of the group and the file name only")
			os.Exit(1)
		}
		groupName := os.Args[2]
		fileName := os.Args[3]

		err := cl.DeleteFile(groupName, fileName)
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}

	case "lf":
		if len(os.Args) != 3 {
			color.Red.Println("Expecting only the name of the group")
			os.Exit(1)
		}
		groupName := os.Args[2]
		out, err := cl.ListFiles(groupName)
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}

		fmt.Println(out)

	case "lg":
		if len(os.Args) != 2 {
			color.Red.Println("Expecting no arugments")
			os.Exit(1)
		}
		groups, err := cl.ListGroups()
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}

		fmt.Println("Groups:")
		for _, g := range groups {
			fmt.Println("  ", g)
		}

	case "dg":
		if len(os.Args) != 3 {
			color.Red.Println("Expecting only the name of the group")
			os.Exit(1)
		}
		err := cl.DeleteGroup(os.Args[2])
		if err != nil {
			color.Red.Println(err)
			os.Exit(1)
		}

	default:
		color.Red.Println("Unexpected command. Run the cli with --help to find out the supported commands.")
		os.Exit(1)
	}

}
