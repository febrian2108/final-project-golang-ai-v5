package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"errors"
	"io"
	"strings"
	"log"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content is empty or whitespace")
	}

	reader := csv.NewReader(strings.NewReader(fileContent))
	headers, err := reader.Read()
	if err != nil {
		return nil, errors.New("failed to read headers from CSV: " + err.Error())
	}

	if len(headers) == 0 {
		return nil, errors.New("no headers found in CSV")
	}

	data := make(map[string][]string)
	for _, header := range headers {
		data[header] = []string{}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("failed to read record from CSV: " + err.Error())
		}

		if len(record) != len(headers) {
			return nil, errors.New("record length does not match header length")
		}

		for i, value := range record {
			data[headers[i]] = append(data[headers[i]], value)
		}
	}

	log.Printf("Processed file data: %v", data)
	return data, nil
}
