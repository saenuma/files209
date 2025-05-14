package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/saenuma/files209/internal"
	"github.com/saenuma/zazabul"
)

var groupMutexes map[string]*sync.RWMutex

func main() {
	groupMutexes = make(map[string]*sync.RWMutex)

	// initialize
	dataPath, err := internal.GetRootPath()
	if err != nil {
		panic(err)
	}

	// create default group
	firstProjPath := filepath.Join(dataPath, "first_group")
	if !internal.DoesPathExists(firstProjPath) {
		err = os.MkdirAll(firstProjPath, 0777)
		if err != nil {
			panic(err)
		}
	}

	confPath, err := internal.GetConfigPath()
	if err != nil {
		panic(err)
	}

	if !internal.DoesPathExists(confPath) {
		conf, err := zazabul.ParseConfig(internal.RootConfigTemplate)
		if err != nil {
			panic(err)
		}
		conf.Write(confPath)
	}

	http.Handle("/is-files209", Q(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "yeah-files209")
	}))

	// files
	http.Handle("/write-file/{group}", Q(writeFile))
	http.Handle("/read-file/{group}", Q(readFile))
	http.Handle("/list-files/{group}", Q(listFiles))
	http.Handle("/delete-file/{group}/{name}", Q(deleteFile))

	// groups
	http.Handle("/list-groups", Q(listGroups))
	http.Handle("/delete-group", Q(deleteGroup))

	port := internal.GetSetting("port")

	fmt.Printf("Serving on port: %s\n", port)

	err = http.ListenAndServeTLS(fmt.Sprintf(":%s", port), internal.G("https-server.crt"),
		internal.G("https-server.key"), nil)
	if err != nil {
		panic(err)
	}

}

func keyEnforcementMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inProd := internal.GetSetting("in_production")
		if inProd == "" {
			panic(errors.New("have you installed and launched files209.fstore"))
		} else if inProd == "true" {
			keyStr := r.FormValue("key-str")
			keyPath := internal.GetKeyStrPath()
			raw, err := os.ReadFile(keyPath)
			if err != nil {
				http.Error(w, "Improperly Configured Server", http.StatusInternalServerError)
			}
			if keyStr == string(raw) {
				// Call the next handler, which can be another middleware in the chain, or the final handler.
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}

		} else {
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		}

	})
}

func Q(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return keyEnforcementMiddleware(http.HandlerFunc(f))
}
