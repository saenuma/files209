// prod provides the commands which helps in making a files209 server production ready.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gookit/color"
	"github.com/saenuma/files209/f2shared"
	"github.com/saenuma/zazabul"
)

func main() {
	dataPath, _ := f2shared.GetRootPath()
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

      `, dataPath)

	case "r":
		keyPath := f2shared.GetKeyStrPath()
		raw, err := os.ReadFile(keyPath)
		if err != nil {
			color.Red.Printf("Error reading key string path.\nError:%s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(raw))

	case "c":
		keyPath := f2shared.GetKeyStrPath()
		randomString := f2shared.GenerateSecureRandomString(50)

		err := os.WriteFile(keyPath, []byte(randomString), 0777)
		if err != nil {
			color.Red.Printf("Error creating key string path.\nError:%s\n", err)
			os.Exit(1)
		}
		fmt.Print(randomString)

	case "mpr":
		keyPath := f2shared.GetKeyStrPath()
		if !f2shared.DoesPathExists(keyPath) {
			randomString := f2shared.GenerateSecureRandomString(50)

			err := os.WriteFile(keyPath, []byte(randomString), 0777)
			if err != nil {
				color.Red.Printf("Error creating key string path.\nError:%s\n", err)
				os.Exit(1)
			}

		}

		confPath, err := f2shared.GetConfigPath()
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
		rootPath, _ := f2shared.GetRootPath()
		keyPath := filepath.Join(rootPath, "https-server.key")
		crtPath := filepath.Join(rootPath, "https-server.crt")

		exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", keyPath,
			"-out", crtPath, "-sha256", "-days", "3650", "-nodes", "-subj",
			"/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname").Run()
		fmt.Println("ok")

	default:
		color.Red.Println("Unexpected command. Run the files209's prod with --help to find out the supported commands.")
		os.Exit(1)
	}

}
