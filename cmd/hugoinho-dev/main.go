package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mdw-tools/hugoinho/core"
	"github.com/mdw-tools/hugoinho/io"
)

var Version = "dev"

func main() {
	disk := io.Disk{}
	logger := log.New(os.Stderr, "", log.Lshortfile)
	args := os.Args[1:]
	config, err := core.NewCLIParser(Version, args).Parse()
	if err != nil {
		logger.Fatal(err)
	}

	var lock sync.Mutex

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		lock.Lock()
		defer lock.Unlock()

		err := os.RemoveAll(config.TargetRoot)
		if err != nil {
			http.Error(response, "Failed to remove target directory.", http.StatusInternalServerError)
			return
		}

		runner := core.NewPipelineRunner(Version, args, disk, time.Now, logger)
		errCount := runner.Run()
		if errCount > 0 {
			http.Error(response, "Failed to generate site.", http.StatusInternalServerError)
			return
		}

		http.ServeFile(response, request, filepath.Join(config.TargetRoot, request.URL.Path))
	})

	address := "localhost:7070"
	logger.Println("Open browser to:", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
