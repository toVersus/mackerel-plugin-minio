package mpminio

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
)

// MinioPlugin contains an endpoint of Minio Server and prefix of Graph definition name
type MinioPlugin struct {
	Scheme      string
	Host        string
	Port        string
	MetricsPath string
	Prefix      string
	Tempfile    string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m MinioPlugin) MetricKeyPrefix() string {
	if m.Prefix == "" {
		m.Prefix = "minio"
	}
	return m.Prefix
}

// fetchAllMetrics fetches all Prometeus compatible metrics from the unauthorized endpoint.
// FYI, see https://github.com/minio/cookbook/blob/master/docs/how-to-monitor-minio-with-prometheus.md
func (m MinioPlugin) fetchAllMetrics() []*prom2json.Family {
	u := m.metricsEndpoint()
	mfChan := make(chan *dto.MetricFamily, 1024)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	go func() {
		if err := prom2json.FetchMetricFamilies(u.String(), mfChan, transport); err != nil {
			log.Fatal(err)
		}
	}()

	result := []*prom2json.Family{}
	for mf := range mfChan {
		result = append(result, prom2json.NewFamily(mf))
	}

	return result
}

// metricsEndpoint returns the url to Minio metrics exporter
func (m MinioPlugin) metricsEndpoint() url.URL {
	return url.URL{
		Scheme: m.Scheme,
		Host:   fmt.Sprintf("%s:%s", m.Host, m.Port),
		Path:   m.MetricsPath,
	}
}

// FetchMetrics is an interface for mackerelplugin
func (m MinioPlugin) FetchMetrics() (map[string]interface{}, error) {
	families := m.fetchAllMetrics()

	stat := make(Stat)
	for _, f := range families {
		stat.handle(f)
	}

	return calcMetrics(stat), nil
}

// calcMetrics appends manually calculated metrics
func calcMetrics(stat map[string]interface{}) map[string]interface{} {
	maxFds, ok := stat["process_max_fds"].(uint64)
	if !ok {
		log.Fatal("process_max_fds not found in stat")
	}
	openFds, ok := stat["process_open_fds"].(uint64)
	if !ok {
		log.Fatal("process_open_fds not found in stat")
	}
	stat["process_fds_percentage"] = (float64(openFds) / float64(maxFds)) * 100

	maxDiskVolume, ok := stat["minio_disk_storage_total_bytes"].(float64)
	if !ok {
		log.Fatal("minio_disk_storage_total_bytes not found in stat")
	}
	usedDiskVolume, ok := stat["minio_disk_storage_used_bytes"].(float64)
	if !ok {
		log.Fatal("minio_disk_storage_used_bytes not found in stat")
	}
	stat["minio_disk_storage_used_percent"] = (usedDiskVolume / maxDiskVolume) * 100

	return stat
}

// Stat represents statistics aggregated from the Minio metrics endpoint
type Stat map[string]interface{}

func (s *Stat) handle(family *prom2json.Family) {
	for _, item := range family.Metrics {
		switch m := item.(type) {
		case prom2json.Metric:
			var value interface{}
			// Metric types are either float or uint
			if v, err := strconv.ParseFloat(m.Value, 64); err == nil {
				value = v
			}
			if v, err := strconv.ParseUint(m.Value, 10, 64); err == nil {
				value = v
			}
			(*s)[metricName(family.Name, m.Labels)] = value
		case prom2json.Histogram:
			val, ok := m.Labels["request_type"]
			if !ok {
				continue
			}

			total, err := strconv.ParseUint(m.Count, 10, 64)
			if err != nil {
				log.Fatalf("Failed to convert %s value: %s", m.Count, err)
			}
			(*s)[family.Name+"_"+val+"_total"] = total

			for k, v := range m.Buckets {
				n, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					log.Fatalf("Failed to convert %s value: %s", k, err)
				}
				k := strings.Replace(k, ".", "_", -1)
				(*s)[family.Name+"_"+val+"_"+k] = n
			}
		}
	}
}

// metricName checks labels and returns identical metric name
func metricName(base string, labels map[string]string) string {
	if len(labels) == 0 {
		return base
	}

	// HTTP Status Code
	if val, ok := labels["code"]; ok {
		return base + "_" + val
	}

	// Other lables are currently not supported
	return base
}

// GraphDefinition is an interface for mackerelplugin
func (m MinioPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(m.Prefix)
	return map[string]mp.Graphs{
		"threads": {
			Label: (labelPrefix + " Threads"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "go_goroutines", Label: "Goroutines"},
				{Name: "go_threads", Label: "OS Threads"},
			},
		},
		"memstats.alloc": {
			Label: (labelPrefix + " Memory Stats"),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "go_memstats_alloc_bytes_total", Label: "Total Memory Bytes Allocated"},
				{Name: "go_memstats_alloc_bytes", Label: "Memory Bytes Allocated"},
				{Name: "go_memstats_frees_total", Label: "Free"},
			},
		},
		"disk_usage": {
			Label: (labelPrefix + " Disk Usage Percentage"),
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "minio_disk_storage_used_percent", Label: "Used", Type: "float64"},
			},
		},
		"disk.availability": {
			Label: (labelPrefix + " Disk Availability"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "minio_offline_disks", Label: "Offline Disk Counts", Type: "uint64"},
				{Name: "minio_total_disks", Label: "Online Disk Counts", Type: "uint64"},
			},
		},
		"process.cpu": {
			Label: (labelPrefix + " Process CPU Time"),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "process_cpu_seconds_total", Label: "In Seconds", Stacked: true, Type: "float64"},
			},
		},
		"process.fds": {
			Label: (labelPrefix + " Process FDs Percentage"),
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "process_fds_percentage", Label: "Consumed", Stacked: true, Type: "float64"},
			},
		},
		"network": {
			Label: (labelPrefix + " Network"),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "minio_network_received_bytes_total", Label: "Total Received Bytes"},
				{Name: "minio_network_sent_bytes_total", Label: "Total Sent Bytes"},
			},
		},
		"http.inflight_request_counts": {
			Label: (labelPrefix + " HTTP Inflight Request Counts"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "promhttp_metric_handler_requests_in_flight", Label: "In Flight", Stacked: true, Type: "uint64"},
			},
		},
		"http.request_counts": {
			Label: (labelPrefix + " HTTP Request Counts"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "promhttp_metric_handler_requests_total_200", Label: "200", Stacked: true, Type: "uint64"},
				{Name: "promhttp_metric_handler_requests_total_500", Label: "500", Stacked: true, Type: "uint64"},
				{Name: "promhttp_metric_handler_requests_total_503", Label: "503", Stacked: true, Type: "uint64"},
			},
		},
		"http_get": {
			Label: (labelPrefix + " HTTP GET Request Duration"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "minio_http_requests_duration_seconds_GET_0_001", Label: "1ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_GET_0_003", Label: "3ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_GET_0_005", Label: "5ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_GET_0_5", Label: "500ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_GET_1", Label: "1s", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_GET_total", Label: "total", Type: "uint64", Diff: true},
			},
		},
		"http_post": {
			Label: (labelPrefix + " HTTP POST Request Duration"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "minio_http_requests_duration_seconds_POST_0_001", Label: "1ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_POST_0_003", Label: "3ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_POST_0_005", Label: "5ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_POST_0_5", Label: "500ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_POST_1", Label: "1s", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_POST_total", Label: "total", Type: "uint64", Diff: true},
			},
		},
		"http_put": {
			Label: (labelPrefix + " HTTP PUT Request Duration"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "minio_http_requests_duration_seconds_PUT_0_001", Label: "1ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_PUT_0_003", Label: "3ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_PUT_0_005", Label: "5ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_PUT_0_5", Label: "500ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_PUT_1", Label: "1s", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_PUT_total", Label: "total", Type: "uint64", Diff: true},
			},
		},
		"http_head": {
			Label: (labelPrefix + " HTTP HEAD Request Duration"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "minio_http_requests_duration_seconds_HEAD_0_001", Label: "1ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_HEAD_0_003", Label: "3ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_HEAD_0_005", Label: "5ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_HEAD_0_5", Label: "500ms", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_HEAD_1", Label: "1s", Type: "uint64", Diff: true},
				{Name: "minio_http_requests_duration_seconds_HEAD_total", Label: "total", Type: "uint64", Diff: true},
			},
		},
	}
}

// Do the plugin
func Do() {
	optScheme := flag.String("scheme", "http", "Protocol scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "9000", "Port")
	optMetricsPath := flag.String("metrics-path", "/minio/prometheus/metrics", "Path to exported metrics")
	optPrefix := flag.String("metric-key-prefix", "minio", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")

	flag.Parse()

	minio := MinioPlugin{
		Scheme:      *optScheme,
		Host:        *optHost,
		Port:        *optPort,
		MetricsPath: *optMetricsPath,
		Prefix:      *optPrefix,
	}

	helper := mp.NewMackerelPlugin(minio)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-minio-%s-%s", *optHost, *optPort))
	}

	helper.Run()
}
