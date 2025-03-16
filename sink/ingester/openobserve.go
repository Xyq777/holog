package ingester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ncuhome/holog/sink"
	"go.uber.org/zap"
)

// OpenObserve ingester
type O2 struct {
	client *http.Client
}

func NewO2Imgester() *O2 {
	return &O2{client: &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 50,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
		},
		Timeout: 15 * time.Second,
	}}
}

func (o2 *O2) Send(ctx context.Context, entry sink.LogEntry) error {
	stream, _ := entry["service"].(string)

	jsonBody, err := json.Marshal(entry)
	if err != nil {
		zap.L().Error("json marshall failed", zap.Error(err))
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/default/%s/_json", os.Getenv("O2_URL"), stream), bodyReader)
	if err != nil {
		zap.L().Error("new http request failed", zap.Error(err))
		return err
	}
	req.SetBasicAuth(os.Getenv("O2_USERNAME"), os.Getenv("O2_PASSWORD"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		zap.L().Error("send request failed", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (o2 *O2) SendBatch(ctx context.Context, entries []sink.LogEntry) error {
	// TODO
	return nil
}
