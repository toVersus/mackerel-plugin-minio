package mpminio

import (
	"reflect"
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	want := 13

	s := SetupMockServer(t)
	defer s.Server.Close()

	graphdef := s.plugin.GraphDefinition()
	if len(graphdef) != want {
		t.Fatalf("got=%d, want=%d", len(graphdef), want)
	}
}

func TestParseMetrics(t *testing.T) {
	wants := map[string]interface{}{
		// threads
		"go_goroutines": uint64(19),
		"go_threads":    uint64(14),
		// memstats allocation
		"go_memstats_alloc_bytes":       float64(5514112),
		"go_memstats_alloc_bytes_total": float64(11288627704),
		"go_memstats_frees_total":       float64(28954625),
		// diks usage (custom)
		"minio_disk_storage_used_percent": (float64(14186127360) / float64(536608768000)) * 100,
		// disk availability
		"minio_offline_disks": uint64(0),
		"minio_total_disks":   uint64(1),
		// process cpu time
		"process_cpu_seconds_total": float64(233.84),
		// process file disctiptor (custom)
		"process_fds_percentage": (float64(8) / float64(65536)) * 100,
		// network
		"minio_network_received_bytes_total": float64(9092723967),
		"minio_network_sent_bytes_total":     float64(15912584309),
		// http inflight request counts
		"promhttp_metric_handler_requests_in_flight": uint64(1),
		// http request counts
		"promhttp_metric_handler_requests_total_200": uint64(9256),
		"promhttp_metric_handler_requests_total_500": uint64(0),
		"promhttp_metric_handler_requests_total_503": uint64(0),
		// http GET request duration
		"minio_http_requests_duration_seconds_GET_0_001": uint64(9334),
		"minio_http_requests_duration_seconds_GET_0_003": uint64(16760),
		"minio_http_requests_duration_seconds_GET_0_005": uint64(17535),
		"minio_http_requests_duration_seconds_GET_0_1":   uint64(18653),
		"minio_http_requests_duration_seconds_GET_0_5":   uint64(18662),
		"minio_http_requests_duration_seconds_GET_1":     uint64(18662),
		"minio_http_requests_duration_seconds_GET_total": uint64(18666),
		// http POST request duration
		"minio_http_requests_duration_seconds_POST_0_001": uint64(11),
		"minio_http_requests_duration_seconds_POST_0_003": uint64(13),
		"minio_http_requests_duration_seconds_POST_0_005": uint64(13),
		"minio_http_requests_duration_seconds_POST_0_1":   uint64(17),
		"minio_http_requests_duration_seconds_POST_0_5":   uint64(22),
		"minio_http_requests_duration_seconds_POST_1":     uint64(22),
		"minio_http_requests_duration_seconds_POST_total": uint64(24),
		// http PUT request duration
		"minio_http_requests_duration_seconds_PUT_0_001": uint64(0),
		"minio_http_requests_duration_seconds_PUT_0_003": uint64(0),
		"minio_http_requests_duration_seconds_PUT_0_005": uint64(71),
		"minio_http_requests_duration_seconds_PUT_0_1":   uint64(135),
		"minio_http_requests_duration_seconds_PUT_0_5":   uint64(135),
		"minio_http_requests_duration_seconds_PUT_1":     uint64(144),
		"minio_http_requests_duration_seconds_PUT_total": uint64(221),
		// http HEAD request duration
		"minio_http_requests_duration_seconds_HEAD_0_001": uint64(254),
		"minio_http_requests_duration_seconds_HEAD_0_003": uint64(257),
		"minio_http_requests_duration_seconds_HEAD_0_005": uint64(257),
		"minio_http_requests_duration_seconds_HEAD_0_1":   uint64(257),
		"minio_http_requests_duration_seconds_HEAD_0_5":   uint64(257),
		"minio_http_requests_duration_seconds_HEAD_1":     uint64(257),
		"minio_http_requests_duration_seconds_HEAD_total": uint64(257),
	}

	s := SetupMockServer(t)
	defer s.Server.Close()

	stat, err := s.plugin.FetchMetrics()
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range wants {
		if !reflect.DeepEqual(stat[k], v) {
			t.Fatalf("got=%v, want=%v\n", stat[k], v)
		}
	}
}
