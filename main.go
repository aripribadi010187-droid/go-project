package main

import (
	"encoding/csv"
	// "encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CSVData struct {
	Data []map[string]string `json:"data"`
}

func readCSV(filename string) ([]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Membaca header
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var result []map[string]string

	for {
		record, err := reader.Read()

		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		row := make(map[string]string)

		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}

		result = append(result, row)
	}

	return result, nil
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Static website
	e.Static("/", "ui")

	// API membaca CSV
	e.GET("/api/data", func(c echo.Context) error {

		data, err := readCSV("data/data_master_karyawan.csv")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, CSVData{
			Data: data,
		})
	})

	// API health check
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})

	fmt.Println("Server running at http://localhost:8888")

	e.Logger.Fatal(e.Start(":8888"))
}