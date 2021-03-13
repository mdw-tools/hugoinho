package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/io"
)

var Version = "dev"

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	args := os.Args[1:]
	config, err := core.NewCLIParser(Version, args).Parse()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		err := os.RemoveAll(config.TargetRoot)
		if err != nil {
			http.Error(response, "Failed to remove target directory.", http.StatusInternalServerError)
			return
		}

		runner := core.NewPipelineRunner(Version, args, io.NewDisk(), time.Now, log.Default())
		errCount := runner.Run()
		if errCount > 0 {
			http.Error(response, "Failed to generate site.", http.StatusInternalServerError)
			return
		}

		http.ServeFile(response, request, filepath.Join(config.TargetRoot, request.URL.Path))
	})

	address := "localhost:8080"
	log.Println("Open browser to: http://localhost:8080")
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
