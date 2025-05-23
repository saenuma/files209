// prod provides the commands which helps in making a files209 server production ready.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gookit/color"
	"github.com/saenuma/files209/internal"
	"github.com/saenuma/zazabul"
)

func main() {
	dataPath, _ := internal.GetRootPath()
	if len(os.Args) < 2 {
		color.Red.Println("expected a command. Open help to view commands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--help", "help", "h":
		fmt.Println(`files209's prod makes a files209 instance production ready.

Supported Commands:

  genssl    Generates the ssl certificates for a files209 installation

  r         Read the current key string used

  c         Creates / Updates and prints a new key string

  mpr       Make production ready. It also creates a key string.

  trim      Trim large files209 files. This is needed after months of using this files store.
            It expects a groupName

      `, dataPath)

	case "r":
		keyPath := internal.GetKeyStrPath()
		raw, err := os.ReadFile(keyPath)
		if err != nil {
			color.Red.Printf("Error reading key string path.\nError:%s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(raw))

	case "c":
		keyPath := internal.GetKeyStrPath()
		randomString := internal.GenerateSecureRandomString(50)

		err := os.WriteFile(keyPath, []byte(randomString), 0777)
		if err != nil {
			color.Red.Printf("Error creating key string path.\nError:%s\n", err)
			os.Exit(1)
		}
		fmt.Print(randomString)

	case "mpr":
		keyPath := internal.GetKeyStrPath()
		if !internal.DoesPathExists(keyPath) {
			randomString := internal.GenerateSecureRandomString(50)

			err := os.WriteFile(keyPath, []byte(randomString), 0777)
			if err != nil {
				color.Red.Printf("Error creating key string path.\nError:%s\n", err)
				os.Exit(1)
			}

		}

		confPath, err := internal.GetConfigPath()
		if err != nil {
			panic(err)
		}

		var conf zazabul.Config

		for {
			conf, err = zazabul.LoadConfigFile(confPath)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			} else {
				break
			}
		}

		conf.Update(map[string]string{
			"in_production": "true",
			"debug":         "false",
		})

		err = conf.Write(confPath)
		if err != nil {
			panic(err)
		}

	case "genssl":
		rootPath, _ := internal.GetRootPath()
		keyPath := filepath.Join(rootPath, "https-server.key")
		crtPath := filepath.Join(rootPath, "https-server.crt")

		exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyPath,
			"-out", crtPath, "-sha256", "-days", "3650", "-nodes", "-subj",
			"/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname").Run()
		fmt.Println("ok")

	case "trim":
		if len(os.Args) != 3 {
			color.Red.Println(`'trim' command expects a groupName `)
			os.Exit(1)
		}

		groupName := os.Args[2]
		err := trimF209Files(groupName)
		if err != nil {
			color.Red.Println("Error triming:\n" + err.Error())
			os.Exit(1)
		}

		fmt.Println("ok")

	default:
		color.Red.Println("Unexpected command. Run the files209's prod with --help to find out the supported commands.")
		os.Exit(1)
	}

}
