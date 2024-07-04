package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/ledongthuc/pdf"
)

func downloadPDF(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func parsePDF(filePath string) ([][]string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var csvRows [][]string
	numPages := r.NumPage()
	for pageIndex := 1; pageIndex <= numPages; pageIndex++ {
		page := r.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		rows, err := page.GetTextByRow()
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			var line string
			for _, text := range row.Content {
				line += text.S
			}

			// 0: Full match, 1: RA, 2: Class code, 3: Campus, 4: Class name
			r := regexp.MustCompile(`([0-9]+)([A-Z0-9]+-[0-9]{2}(SA|SB))(.*)`)
			elements := r.FindStringSubmatch(line)
			if len(elements) < 5 {
				continue
			}

			fmt.Printf("['%s', '%s', '%s']\n", elements[1], elements[2], elements[4])
			csvRows = append(csvRows, []string{
				elements[1],
				elements[2],
				elements[4],
			})
		}
	}
	return csvRows, nil
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
	url := os.Args[1]
	pdfPath := "downloaded.pdf"
	csvPath := "output.csv"

	if err := downloadPDF(url, pdfPath); err != nil {
		log.Fatalf("Failed to download PDF: %v", err)
	}

	data, err := parsePDF(pdfPath)
	if err != nil {
		log.Fatalf("Failed to parse PDF: %v", err)
	}

	if err := saveToCSV(data, csvPath); err != nil {
		log.Fatalf("Failed to save CSV: %v", err)
	}

	fmt.Println("PDF data saved to CSV successfully!")
}
