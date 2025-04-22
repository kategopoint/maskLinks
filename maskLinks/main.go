package main

import (
	"log"
	"os"

	"main.go/presenter"
	"main.go/producer"
	"main.go/service"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <input_file_name.txt> [output_file_name.txt]")
	}

	inputFilePath := os.Args[1]
	outputFilePath := "changedLinesForFunc.txt"

	if len(os.Args) > 2 {
		outputFilePath = os.Args[2]
	}

	producer := producer.NewFileProducer(inputFilePath)
	presenter := presenter.NewFilePresenter(outputFilePath)
	svc := service.NewService(producer, presenter)

	if err := svc.Run(); err != nil {
		log.Fatalf("Error running service: %v", err)
	}
}
