package mpminio

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

var metrics = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 7.473e-06
go_gc_duration_seconds{quantile="0.25"} 1.5483e-05
go_gc_duration_seconds{quantile="0.5"} 2.1819e-05
go_gc_duration_seconds{quantile="0.75"} 5.1844e-05
go_gc_duration_seconds{quantile="1"} 0.013692739
go_gc_duration_seconds_sum 0.563420551
go_gc_duration_seconds_count 3395
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 19
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.12"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 5.514112e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.1288627704e+10
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.662232e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 2.8954625e+07
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 1.562827169845014e-05
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.398208e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 5.514112e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 5.8638336e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 7.585792e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 39206
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 5.6451072e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6224128e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.562480101050328e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 2.8993831e+07
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 3472
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 85392
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 245760
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 7.108432e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 592864
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 884736
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 884736
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.2024312e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 14
# HELP minio_disk_storage_available_bytes Total disk available space seen by MinIO server instance
# TYPE minio_disk_storage_available_bytes gauge
minio_disk_storage_available_bytes 5.2242264064e+11
# HELP minio_disk_storage_total_bytes Total disk space seen by MinIO server instance
# TYPE minio_disk_storage_total_bytes gauge
minio_disk_storage_total_bytes 5.36608768e+11
# HELP minio_disk_storage_used_bytes Total disk storage used by current MinIO server instance
# TYPE minio_disk_storage_used_bytes gauge
minio_disk_storage_used_bytes 1.418612736e+10
# HELP minio_http_requests_duration_seconds Time taken by requests served by current MinIO server instance
# TYPE minio_http_requests_duration_seconds histogram
minio_http_requests_duration_seconds_bucket{request_type="GET",le="0.001"} 9334
minio_http_requests_duration_seconds_bucket{request_type="GET",le="0.003"} 16760
minio_http_requests_duration_seconds_bucket{request_type="GET",le="0.005"} 17535
minio_http_requests_duration_seconds_bucket{request_type="GET",le="0.1"} 18653
minio_http_requests_duration_seconds_bucket{request_type="GET",le="0.5"} 18662
minio_http_requests_duration_seconds_bucket{request_type="GET",le="1"} 18662
minio_http_requests_duration_seconds_bucket{request_type="GET",le="+Inf"} 18666
minio_http_requests_duration_seconds_sum{request_type="GET"} 71.42435698400023
minio_http_requests_duration_seconds_count{request_type="GET"} 18666
minio_http_requests_duration_seconds_bucket{request_type="HEAD",le="0.001"} 254
minio_http_requests_duration_seconds_bucket{request_type="HEAD",le="0.003"} 257
minio_http_requests_duration_seconds_bucket{request_type="HEAD",le="0.005"} 257
minio_http_requests_duration_seconds_bucket{request_type="HEAD",le="0.1"} 257
minio_http_requests_duration_seconds_bucket{request_type="HEAD",le="0.5"} 257
minio_http_requests_duration_seconds_bucket{request_type="HEAD",le="1"} 257
minio_http_requests_duration_seconds_bucket{request_type="HEAD",le="+Inf"} 257
minio_http_requests_duration_seconds_sum{request_type="HEAD"} 0.14868525700000013
minio_http_requests_duration_seconds_count{request_type="HEAD"} 257
minio_http_requests_duration_seconds_bucket{request_type="POST",le="0.001"} 11
minio_http_requests_duration_seconds_bucket{request_type="POST",le="0.003"} 13
minio_http_requests_duration_seconds_bucket{request_type="POST",le="0.005"} 13
minio_http_requests_duration_seconds_bucket{request_type="POST",le="0.1"} 17
minio_http_requests_duration_seconds_bucket{request_type="POST",le="0.5"} 22
minio_http_requests_duration_seconds_bucket{request_type="POST",le="1"} 22
minio_http_requests_duration_seconds_bucket{request_type="POST",le="+Inf"} 24
minio_http_requests_duration_seconds_sum{request_type="POST"} 36.50791978999998
minio_http_requests_duration_seconds_count{request_type="POST"} 24
minio_http_requests_duration_seconds_bucket{request_type="PUT",le="0.001"} 0
minio_http_requests_duration_seconds_bucket{request_type="PUT",le="0.003"} 0
minio_http_requests_duration_seconds_bucket{request_type="PUT",le="0.005"} 71
minio_http_requests_duration_seconds_bucket{request_type="PUT",le="0.1"} 135
minio_http_requests_duration_seconds_bucket{request_type="PUT",le="0.5"} 135
minio_http_requests_duration_seconds_bucket{request_type="PUT",le="1"} 144
minio_http_requests_duration_seconds_bucket{request_type="PUT",le="+Inf"} 221
minio_http_requests_duration_seconds_sum{request_type="PUT"} 484.8310546109999
minio_http_requests_duration_seconds_count{request_type="PUT"} 221
# HELP minio_network_received_bytes_total Total number of bytes received by current MinIO server instance
# TYPE minio_network_received_bytes_total counter
minio_network_received_bytes_total 9.092723967e+09
# HELP minio_network_sent_bytes_total Total number of bytes sent by current MinIO server instance
# TYPE minio_network_sent_bytes_total counter
minio_network_sent_bytes_total 1.5912584309e+10
# HELP minio_offline_disks Total number of offline disks for current MinIO server instance
# TYPE minio_offline_disks gauge
minio_offline_disks 0
# HELP minio_total_disks Total number of disks for current MinIO server instance
# TYPE minio_total_disks gauge
minio_total_disks 1
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 233.84
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 65536
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 8
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.5956352e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.5622029737e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.48594688e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 9256
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
`

// MockServer represents a of mock metrics server for testing.
type MockServer struct {
	plugin MinioPlugin
	t      *testing.T
	Server *httptest.Server
}

// SetupMockServer setups mock API server for testing.
func SetupMockServer(t *testing.T) *MockServer {
	m := &MockServer{
		plugin: MinioPlugin{
			Scheme:      "http",
			Host:        "localhost",
			Port:        "9000",
			MetricsPath: "/minio/prometheus/metrics",
			Prefix:      "minio",
		},
		t: t,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(m.plugin.MetricsPath, m.metricsHandler)

	m.Server = httptest.NewUnstartedServer(mux)
	// Close the listener created by NewUnstartedServer
	m.Server.Listener.Close()
	// create a listener with the custom port
	l, _ := net.Listen("tcp", fmt.Sprintf("%s:%s", m.plugin.Host, m.plugin.Port))
	// replace with the custom listerner
	m.Server.Listener = l
	m.Server.Start()

	return m
}

func (m *MockServer) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, metrics)
}
