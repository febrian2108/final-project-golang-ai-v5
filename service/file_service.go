package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

type CSVResponse struct {
	Records              [][]string         `json:"records"`           // Menyimpan semua data CSV
	EnergyConsumptionMap map[string]float64 `json:"energyConsumption"` // Konsumsi energi per ruangan
}

func (s *FileService) ProcessFile(fileContent string) (*CSVResponse, error) {
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content is empty or whitespace")
	}

	// Mengatur reader untuk menggunakan titik koma sebagai pemisah
	reader := csv.NewReader(strings.NewReader(fileContent))
	reader.Comma = ';' // Menentukan pemisah kolom sebagai titik koma

	// Membaca header CSV
	headers, err := reader.Read()
	if err != nil {
		return nil, errors.New("failed to read headers from CSV: " + err.Error())
	}

	// Validasi panjang kolom
	if len(headers) == 0 {
		return nil, errors.New("no headers found in CSV")
	}

	// Menginisialisasi map untuk menyimpan konsumsi energi per ruangan
	energyConsumption := make(map[string]float64)
	var allData [][]string

	// Membaca setiap baris data CSV
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("failed to read record from CSV: " + err.Error())
		}

		// Validasi jumlah kolom
		if len(record) < 5 {
			log.Printf("Skipping invalid row: %v", record)
			continue
		}

		// Menambahkan data ke allData
		allData = append(allData, record)

		// Asumsi kolom sesuai dengan urutan: Date, Time, Appliance, Energy_Consumption, Room, Status
		energy, err := parseEnergy(record[3]) // Mengambil nilai Energy_Consumption
		if err != nil {
			return nil, err
		}
		room := record[4] // Mengambil nilai Room

		// Menambahkan konsumsi energi per ruangan
		energyConsumption[room] += energy
	}

	// Mengembalikan response yang berisi seluruh data CSV dan total konsumsi energi per ruangan
	response := CSVResponse{
		Records:              allData,
		EnergyConsumptionMap: energyConsumption,
	}

	return &response, nil
}

// Fungsi untuk mengonversi Energy_Consumption menjadi tipe data float64
func parseEnergy(value string) (float64, error) {
	var energy float64
	_, err := fmt.Sscanf(value, "%f", &energy)
	if err != nil {
		return 0, errors.New("invalid energy consumption value: " + value)
	}
	return energy, nil
}
