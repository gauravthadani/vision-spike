package gcp

import (
	"context"
	"fmt"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	gax "github.com/googleapis/gax-go"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type Vision struct {
	client *vision.ImageAnnotatorClient
	ctx    context.Context
}

func NewVisionClient() (*Vision, error) {

	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Vision{client: client, ctx: ctx}, nil
}

func (v Vision) BatchLabels(files []*os.File) {

	features := []*pb.Feature{{Type: pb.Feature_LABEL_DETECTION, MaxResults: 10}}
	batchRequest := &pb.BatchAnnotateImagesRequest{}
	indexToNameMap := map[int]string{}

	for index, file := range files {
		image, err := vision.NewImageFromReader(file)
		if err != nil {
			log.Fatalf("Failed to create image: %v", err)
		}

		annotateImageRequest := &pb.AnnotateImageRequest{
			Image:        image,
			ImageContext: nil,
			Features:     features,
		}

		batchRequest.Requests = append(batchRequest.Requests, annotateImageRequest)
		indexToNameMap[index] = file.Name()
	}

	batchResponse, err := v.client.BatchAnnotateImages(v.ctx, batchRequest, []gax.CallOption{}...)

	if err != nil {
		log.Fatalf("Failed to get responses from GCP Vision: %v", err)
	}

	fmt.Println("Batch results")
	for index, response := range batchResponse.Responses {
		fmt.Printf("\n\n %s", indexToNameMap[index])
		entities := response.GetLabelAnnotations()
		for _, entity := range entities {
			fmt.Printf("\n %v ", entity)
		}
	}
}

func (v Vision) DetectLabels(file *os.File) {

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	labels, err := v.client.DetectLabels(v.ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	fmt.Println("Labels:")
	for _, label := range labels {
		fmt.Println(label.Description)
		fmt.Printf("\t %v \n", label)
	}
}

func (v Vision) DetectTexts(file *os.File) {
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	texts, err := v.client.DetectTexts(v.ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect texts: %v", err)
	}

	fmt.Println("Texts:")
	for _, text := range texts {
		fmt.Println(text.Description)
		fmt.Printf("\t %v \n", text)
	}
}

func (v Vision) DetectFaces(file *os.File) {
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	faces, err := v.client.DetectFaces(v.ctx, image, nil, 25)
	if err != nil {
		log.Fatalf("Failed to detect faces: %v", err)
	}

	fmt.Printf("Faces: %d \n", len(faces))
	for index, face := range faces {
		fmt.Printf("Index number: %d", index)
		fmt.Printf("\t %v \n", face)
	}
}
