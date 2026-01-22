package http

import (
	"Analytics/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MLRequest struct {
	Readings []models.Reading `json:"readings"`
	Horizon  int              `json:"horizon"`
}

type MLResponse struct {
	NextTemps []float32 `json:"next_temps"`
	Status    string    `json:"status"`
}

func PredictTemps(readings []models.Reading) ([]float32, error) {
	var resp *http.Response
	var err error

	req := MLRequest{Readings: readings, Horizon: len(readings)}
	var data []byte
	if data, err = json.Marshal(req); err != nil {
		return nil, err
	}
	if resp, err = http.Post("http://mlaas:8100/predict", "application/json", bytes.NewBuffer(data)); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	var mlResp MLResponse
	if err = json.Unmarshal(body, &mlResp); err != nil {
		return nil, err
	}
	if mlResp.Status != "success" {
		return nil, fmt.Errorf("MLaaS error: %v", mlResp.Status)
	}
	return mlResp.NextTemps, nil

}
