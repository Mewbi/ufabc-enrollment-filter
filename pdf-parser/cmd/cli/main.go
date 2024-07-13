package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"pdf-parser/config"
	"pdf-parser/internal/service"
)

type Application struct {
	httpClient service.HttpClient
	parser     service.Parser
}

func NewApplication(client service.HttpClient, parser service.Parser) *Application {
	return &Application{
		httpClient: client,
		parser:     parser,
	}
}

func saveToCSV(data [][]string, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("URL is required as command-line argument")
	}

	_, err := config.Load("./config/config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}
	client := service.NewHttpClient()
	parser := service.NewParser()

	app := NewApplication(client, parser)

	url := os.Args[1]
	csvPath := "output.csv"

	contentPDF, err := app.httpClient.DownloadPDF(url)
	if err != nil {
		log.Fatalf("Failed to download PDF: %v", err)
	}

	rows, err := app.parser.ParsePDF(contentPDF)
	if err != nil {
		log.Fatalf("Failed to parse PDF: %v", err)
	}

	contentCSV, err := app.parser.ConvertToCSV(rows)
	if err != nil {
		log.Fatalf("Failed to convert to CSV: %s", err.Error())
	}

	if err := app.parser.SaveToCSVFile(contentCSV, csvPath); err != nil {
		log.Fatalf("Failed to save CSV: %v", err)
	}

	fmt.Println("PDF data saved to CSV successfully!")
}
