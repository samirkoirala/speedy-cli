package speedtest

import (
	"context"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"speedy-cli/internal/common"
)

type ServerStat struct {
	Server       string  `json:"server"`
	DownloadMbps float64 `json:"downloadMbps"`
	UploadMbps   float64 `json:"uploadMbps"`
	PingMs       float64 `json:"pingMs"`
}

type Stats struct {
	DownloadMbps float64      `json:"downloadMbps"`
	UploadMbps   float64      `json:"uploadMbps"`
	PingMs       float64      `json:"pingMs"`
	ByServer     []ServerStat `json:"byServer"`
}

func DefaultServers() []string {
	return []string{
		"https://speed.hetzner.de/10MB.bin",
		"https://proof.ovh.net/files/10Mb.dat",
	}
}

func RunParallel(servers []string, verbose bool) (common.Result, Stats) {
	client := &http.Client{Timeout: 15 * time.Second}
	stats := make([]ServerStat, 0, len(servers))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, s := range servers {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			dl := measureDownload(client, s)
			ul := measureUpload(client)
			pg := measurePing(s)
			mu.Lock()
			stats = append(stats, ServerStat{Server: s, DownloadMbps: dl, UploadMbps: ul, PingMs: pg})
			mu.Unlock()
		}()
	}
	wg.Wait()

	if len(stats) == 0 {
		return common.Result{Status: common.StatusError, Message: "speed test failed", Suggestion: "Check network and retry"}, Stats{}
	}

	var d, u, p float64
	for _, s := range stats {
		d += s.DownloadMbps
		u += s.UploadMbps
		p += s.PingMs
	}
	n := float64(len(stats))
	agg := Stats{DownloadMbps: round2(d / n), UploadMbps: round2(u / n), PingMs: round2(p / n), ByServer: stats}

	result := common.Result{Status: common.StatusSuccess, Message: "speed test complete"}
	if agg.DownloadMbps < 20 {
		result.Status = common.StatusWarning
		result.Suggestion = "Download looks slow; try changing Wi-Fi channel or using ethernet"
	}
	return result, agg
}

func measureDownload(client *http.Client, fileURL string) float64 {
	start := time.Now()
	resp, err := client.Get(fileURL)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return 0
	}
	bytesRead, _ := io.CopyN(io.Discard, resp.Body, 5*1024*1024)
	elapsed := time.Since(start).Seconds()
	if elapsed <= 0 || bytesRead <= 0 {
		return 0
	}
	return round2((float64(bytesRead) * 8 / 1_000_000) / elapsed)
}

func measureUpload(client *http.Client) float64 {
	payload := io.LimitReader(zeroReader{}, 2*1024*1024)
	start := time.Now()
	req, _ := http.NewRequest(http.MethodPost, "https://httpbin.org/post", payload)
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	elapsed := time.Since(start).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return round2((float64(2*1024*1024) * 8 / 1_000_000) / elapsed)
}

func measurePing(rawURL string) float64 {
	u, err := url.Parse(rawURL)
	if err != nil {
		return 0
	}
	host := u.Hostname()
	if host == "" {
		return 0
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", host+":443")
	if err != nil {
		return 0
	}
	conn.Close()
	return round2(float64(time.Since(start).Milliseconds()))
}

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'z'
	}
	return len(p), nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
