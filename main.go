package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metric registry
var oidMap = map[string]*prometheus.GaugeVec{
	"1.3.6.1.4.1.4491.2.1.21.1.2.1.6.2.3.706": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virgin_downstream_frequency_hz",
			Help: "Downstream frequency in Hz",
		},
		[]string{"channel"},
	),
	"1.3.6.1.4.1.4491.2.1.21.1.2.1.7.2.3.706": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virgin_downstream_symbol_rate",
			Help: "Downstream symbol rate",
		},
		[]string{"channel"},
	),
	"1.3.6.1.4.1.4491.2.1.21.1.2.1.8.2.3.706": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virgin_downstream_t3_timeouts",
			Help: "T3 timeouts",
		},
		[]string{"channel"},
	),
	"1.3.6.1.4.1.4491.2.1.21.1.2.1.9.2.3.706": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virgin_downstream_snr_db",
			Help: "Downstream SNR in dB",
		},
		[]string{"channel"},
	),
	"1.3.6.1.4.1.4491.2.1.21.1.2.1.4.2.3.706": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virgin_upstream_power_dbmv",
			Help: "Upstream power level in dBmV",
		},
		[]string{"channel"},
	),
	"docsis_events": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virgin_docsis_event",
			Help: "DOCSIS event log entries",
		},
		[]string{"event_id", "description"},
	),
}

// Parse float safely
func parseFloat(s string) float64 {
	val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		log.Printf("Failed to parse float from '%s': %v", s, err)
		return 0
	}
	return val
}

// Scrape metrics from router
func scrapeMetrics() {
	req, err := http.NewRequest("GET", "http://192.168.0.1/getRouterStatus", nil)
	if err != nil {
		log.Println("Request creation failed:", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "http://192.168.0.1/")
	// req.Header.Set("Cookie", "SessionID=your-session-cookie")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response:", err)
		return
	}

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Failed to parse JSON:", err)
		return
	}

	for oid, val := range data {
		switch {
		case strings.HasPrefix(oid, "1.3.6.1.2.1.69.1.5.8.1.7."):
			eventID := strings.TrimPrefix(oid, "1.3.6.1.2.1.69.1.5.8.1.7.")
			description := val
			oidMap["docsis_events"].WithLabelValues(eventID, description).Set(1)
			log.Printf("Logged DOCSIS event %s → %s", eventID, description)

		case oidMap[oid] != nil:
			floatVal := parseFloat(val)
			oidMap[oid].WithLabelValues("706").Set(floatVal)
			log.Printf("Updated metric for OID %s → %f", oid, floatVal)
		}
	}
}

func main() {
	for _, gaugeVec := range oidMap {
		prometheus.MustRegister(gaugeVec)
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		for {
			log.Println("Scraping Virgin Media metrics...")
			scrapeMetrics()
			time.Sleep(30 * time.Second)
		}
	}()

	log.Println("Virgin Hub Exporter running on :9877")
	log.Fatal(http.ListenAndServe(":9877", nil))
}
