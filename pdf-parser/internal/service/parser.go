package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"

	"pdf-parser/config"

	"github.com/ledongthuc/pdf"
)

type ParserWrapper struct{}

func NewParser() *ParserWrapper {
	return &ParserWrapper{}
}

func (pw *ParserWrapper) ParsePDF(content *bytes.Buffer) ([][]string, error) {
	reader := bytes.NewReader(content.Bytes())
	r, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		return nil, err
	}

	conf := config.Get()
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

			if conf.Debug {
				fmt.Printf("['%s', '%s', '%s']\n", elements[1], elements[2], elements[4])
			}

			csvRows = append(csvRows, []string{
				elements[1],
				elements[2],
				elements[4],
			})
		}
	}

	if len(csvRows) == 0 {
		return nil, fmt.Errorf("got 0 lines parsing informed PDF")
	}

	return csvRows, nil
}

func (pw *ParserWrapper) ConvertToCSV(rows [][]string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	writer := csv.NewWriter(buffer)
	defer writer.Flush()

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("error creating csv: %s", err.Error())
		}
	}
	return buffer, nil
}

func (pw *ParserWrapper) SaveToCSVFile(content *bytes.Buffer, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content.Bytes())
	return err
}
