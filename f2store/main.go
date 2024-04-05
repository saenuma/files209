package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/saenuma/files209/f2shared"
	"github.com/saenuma/zazabul"
)

var groupMutexes map[string]*sync.RWMutex

func main() {
	groupMutexes = make(map[string]*sync.RWMutex)

	// initialize
	dataPath, err := f2shared.GetRootPath()
	if err != nil {
		panic(err)
	}

	// create default group
	firstProjPath := filepath.Join(dataPath, "first_group")
	if !f2shared.DoesPathExists(firstProjPath) {
		err = os.MkdirAll(firstProjPath, 0777)
		if err != nil {
			panic(err)
		}
	}

	confPath, err := f2shared.GetConfigPath()
	if err != nil {
		panic(err)
	}

	if !f2shared.DoesPathExists(confPath) {
		conf, err := zazabul.ParseConfig(f2shared.RootConfigTemplate)
		if err != nil {
			panic(err)
		}
		conf.Write(confPath)
	}

	r := mux.NewRouter()

	r.HandleFunc("/is-f209", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "yeah-f209")
	})

	r.Use(keyEnforcementMiddleware)

	port := f2shared.GetSetting("port")

	fmt.Printf("Serving on port: %s\n", port)

	err = http.ListenAndServeTLS(fmt.Sprintf(":%s", port), f2shared.G("https-server.crt"),
		f2shared.G("https-server.key"), r)
	if err != nil {
		panic(err)
	}

}

func keyEnforcementMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inProd := f2shared.GetSetting("in_production")
		if inProd == "" {
			panic(errors.New("have you installed and launched files209.fstore"))
		} else if inProd == "true" {
			keyStr := r.FormValue("key-str")
			keyPath := f2shared.GetKeyStrPath()
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
