package metrics

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MetricsService struct {
	address string
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		address: "http://localhost:2112",
	}
}

type Metric struct {
	Labels map[string]string `json:"labels"`
	Value  float64           `json:"value"`
	Type   string            `json:"type"` // "turnover", "profit", "deals_count", "deals_cancelled_count", "deals_profitable_count"
}

func (ms *MetricsService) SendMetric(key string, value float64, labels map[string]string) error {
	metric := Metric{
		Labels: labels,
		Value:  value,
		Type:   key,
	}
	payload, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	req, err := http.Post(ms.address+"/metric", "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Body.Close()

	return nil
}
