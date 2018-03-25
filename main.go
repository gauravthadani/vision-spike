package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gauravthadani/vision-spike/gcp"
)

func main() {

	appCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if appCreds == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	path := flag.String("path", "", "path to image file")
	operation := flag.String("op", "", "FACES/TEXTS/LABELS")
	flag.Parse()
	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	c, err := gcp.NewVisionClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)

	}

	switch *operation {
	case "TEXTS":
		file, err := os.Open(*path)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		defer file.Close()
		c.DetectTexts(file)
	case "LABELS":
		file, err := os.Open(*path)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		defer file.Close()
		c.DetectLabels(file)
	case "FACES":
		file, err := os.Open(*path)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		defer file.Close()
		c.DetectFaces(file)
	case "BATCH_LABELS":
		filesInfo, err := ioutil.ReadDir(*path)
		if err != nil {
			log.Fatalf("Failed to read path: %v", err)
		}

		files := []*os.File{}
		for _, fileInfo := range filesInfo {
			file, err := os.Open(fmt.Sprintf("%s/%s", *path, fileInfo.Name()))
			if err != nil {
				log.Fatalf("Failed to read file: %v", err)
			}
			defer file.Close()
			files = append(files, file)
		}
		c.BatchLabels(files)
	default:
		flag.Usage()
	}

}
