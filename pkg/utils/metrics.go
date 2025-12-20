package utils

import (
	"os"

	"github.com/GetStream/tikv-ui/pkg/types"
)

// Unit types for formatting:
// - bytes: memory/size (KB, MB, GB)
// - count: counters (1K, 10K, 1M)
// - time: durations (µs, ms, s)
// - ratio: percentages/ratios
// - rate: throughput (bytes/s, ops/s)
// - number: raw numbers

// metricDef defines a metric with its label and unit
type metricDef struct {
	Label   string
	Unit    string
	Default bool // enabled by default
}

// metricsDefinitions contains all supported metrics with human-readable labels and units
var metricsDefinitions = map[string]metricDef{
	// Process metrics
	"process_resident_memory_bytes": {Label: "Resident Memory", Unit: "bytes", Default: true},
	"process_virtual_memory_bytes":  {Label: "Virtual Memory", Unit: "bytes", Default: true},

	// Raft Engine metrics
	"raft_engine_log_entry_count":     {Label: "Raft Log Entries", Unit: "count", Default: true},
	"raft_engine_log_file_count":      {Label: "Raft Log Files", Unit: "count"},
	"raft_engine_memory_usage":        {Label: "Raft Engine Memory", Unit: "bytes", Default: true},
	"raft_engine_recycled_file_count": {Label: "Raft Recycled Files", Unit: "count"},

	// Allocator metrics
	"tikv_allocator_stats":             {Label: "Allocator Stats", Unit: "bytes"},
	"tikv_allocator_thread_allocation": {Label: "Thread Allocation", Unit: "bytes"},

	// Backup metrics
	"tikv_backup_softlimit": {Label: "Backup Soft Limit", Unit: "bytes"},

	// CDC metrics
	"tikv_cdc_captured_region_total":        {Label: "CDC Captured Regions", Unit: "count"},
	"tikv_cdc_endpoint_pending_tasks":       {Label: "CDC Pending Tasks", Unit: "count"},
	"tikv_cdc_old_value_cache_access":       {Label: "CDC Cache Access", Unit: "count"},
	"tikv_cdc_old_value_cache_bytes":        {Label: "CDC Cache Size", Unit: "bytes"},
	"tikv_cdc_old_value_cache_length":       {Label: "CDC Cache Length", Unit: "count"},
	"tikv_cdc_old_value_cache_memory_quota": {Label: "CDC Cache Quota", Unit: "bytes"},
	"tikv_cdc_old_value_cache_miss":         {Label: "CDC Cache Miss", Unit: "count"},
	"tikv_cdc_old_value_cache_miss_none":    {Label: "CDC Cache Miss None", Unit: "count"},
	"tikv_cdc_region_resolve_status":        {Label: "CDC Region Status", Unit: "number"},
	"tikv_cdc_resolved_ts_advance_method":   {Label: "CDC TS Advance Method", Unit: "number"},
	"tikv_cdc_sink_memory_bytes":            {Label: "CDC Sink Memory", Unit: "bytes"},
	"tikv_cdc_sink_memory_capacity":         {Label: "CDC Sink Capacity", Unit: "bytes"},

	// Leader check metrics
	"tikv_check_leader_request_pending_count":      {Label: "Leader Check Pending", Unit: "count"},
	"tikv_check_leader_request_sent_pending_count": {Label: "Leader Check Sent Pending", Unit: "count"},

	// Concurrency metrics
	"tikv_concurrency_manager_min_lock_ts": {Label: "Min Lock Timestamp", Unit: "number", Default: true},

	// Config metrics
	"tikv_config_raftstore": {Label: "Raftstore Config", Unit: "number"},
	"tikv_config_rocksdb":   {Label: "RocksDB Config", Unit: "number"},

	// Engine cache metrics
	"tikv_engine_blob_cache_size_bytes":  {Label: "Blob Cache Size", Unit: "bytes"},
	"tikv_engine_block_cache_size_bytes": {Label: "Block Cache Size", Unit: "bytes"},

	// Engine compression metrics
	"tikv_engine_bytes_compressed":         {Label: "Bytes Compressed", Unit: "bytes"},
	"tikv_engine_bytes_decompressed":       {Label: "Bytes Decompressed", Unit: "bytes"},
	"tikv_engine_compression_ratio":        {Label: "Compression Ratio", Unit: "ratio"},
	"tikv_engine_compression_time_nanos":   {Label: "Compression Time", Unit: "time"},
	"tikv_engine_decompression_time_nanos": {Label: "Decompression Time", Unit: "time"},

	// Engine I/O metrics
	"tikv_engine_bytes_per_read":      {Label: "Bytes per Read", Unit: "bytes", Default: true},
	"tikv_engine_bytes_per_write":     {Label: "Bytes per Write", Unit: "bytes", Default: true},
	"tikv_engine_get_micro_seconds":   {Label: "Get Latency", Unit: "time"},
	"tikv_engine_seek_micro_seconds":  {Label: "Seek Latency", Unit: "time"},
	"tikv_engine_write_micro_seconds": {Label: "Write Latency", Unit: "time"},
	"tikv_engine_sst_read_micros":     {Label: "SST Read Latency", Unit: "time"},

	// Engine compaction metrics
	"tikv_engine_compaction_time":                       {Label: "Compaction Time", Unit: "time"},
	"tikv_engine_compaction_outfile_sync_micro_seconds": {Label: "Compaction Sync", Unit: "time"},
	"tikv_engine_pending_compaction_bytes":              {Label: "Pending Compaction", Unit: "bytes"},
	"tikv_engine_num_files_in_single_compaction":        {Label: "Files per Compaction", Unit: "count"},

	// Engine file metrics
	"tikv_engine_num_files_at_level":          {Label: "Files at Level", Unit: "count"},
	"tikv_engine_num_immutable_mem_table":     {Label: "Immutable MemTables", Unit: "count"},
	"tikv_engine_num_snapshots":               {Label: "Engine Snapshots", Unit: "count"},
	"tikv_engine_num_subcompaction_scheduled": {Label: "Scheduled Subcompactions", Unit: "count"},

	// Engine sync metrics
	"tikv_engine_manifest_file_sync_micro_seconds": {Label: "Manifest Sync", Unit: "time"},
	"tikv_engine_table_sync_micro_seconds":         {Label: "Table Sync", Unit: "time"},
	"tikv_engine_wal_file_sync_micro_seconds":      {Label: "WAL Sync", Unit: "time"},
	"tikv_engine_write_wal_time_micro_seconds":     {Label: "WAL Write Time", Unit: "time"},

	// Engine size metrics
	"tikv_engine_memory_bytes":      {Label: "Engine Memory", Unit: "bytes", Default: true},
	"tikv_engine_size_bytes":        {Label: "Engine Size", Unit: "bytes"},
	"tikv_engine_estimate_num_keys": {Label: "Estimated Keys", Unit: "count"},

	// Engine stall metrics
	"tikv_engine_stall_l0_num_files_count":        {Label: "L0 Stall Files", Unit: "count"},
	"tikv_engine_stall_l0_slowdown_count":         {Label: "L0 Slowdown Count", Unit: "count"},
	"tikv_engine_stall_memtable_compaction_count": {Label: "MemTable Stall Count", Unit: "count"},
	"tikv_engine_write_stall":                     {Label: "Write Stall", Unit: "time"},
	"tikv_engine_write_stall_reason":              {Label: "Write Stall Reason", Unit: "number"},

	// Engine rate limit metrics
	"tikv_engine_hard_rate_limit_delay_count": {Label: "Hard Rate Limit Delays", Unit: "count"},
	"tikv_engine_soft_rate_limit_delay_count": {Label: "Soft Rate Limit Delays", Unit: "count"},

	// Future pool metrics
	"tikv_futurepool_pending_task_total": {Label: "FuturePool Pending Tasks", Unit: "count"},

	// GC worker metrics
	"tikv_gcworker_autogc_processed_regions": {Label: "GC Processed Regions", Unit: "count"},
	"tikv_gcworker_autogc_status":            {Label: "GC Status", Unit: "number"},

	// Import metrics
	"tikv_import_apply_cached_bytes": {Label: "Import Cached", Unit: "bytes"},

	// PD metrics
	"tikv_pd_pending_heartbeat_total":   {Label: "PD Pending Heartbeats", Unit: "count"},
	"tikv_pd_pending_tso_request_total": {Label: "PD Pending TSO Requests", Unit: "count"},

	// Pending operations metrics
	"tikv_pending_delete_ranges_of_stale_peer": {Label: "Pending Delete Ranges", Unit: "count"},
	"tikv_pessimistic_lock_memory_size":        {Label: "Pessimistic Lock Memory", Unit: "bytes"},

	// Raft metrics
	"tikv_raft_entries_caches": {Label: "Raft Entry Caches", Unit: "bytes"},

	// Raftstore metrics
	"tikv_raftstore_hibernated_peer_state":  {Label: "Hibernated Peers", Unit: "count"},
	"tikv_raftstore_leader_missing":         {Label: "Missing Leaders", Unit: "count"},
	"tikv_raftstore_read_index_pending":     {Label: "Read Index Pending", Unit: "count"},
	"tikv_raftstore_region_count":           {Label: "Region Count", Unit: "count"},
	"tikv_raftstore_snapshot_traffic_total": {Label: "Snapshot Traffic", Unit: "bytes"},

	// Raftstore slow trend metrics
	"tikv_raftstore_slow_score":                         {Label: "Slow Score", Unit: "number"},
	"tikv_raftstore_slow_trend":                         {Label: "Slow Trend", Unit: "number"},
	"tikv_raftstore_slow_trend_l0":                      {Label: "Slow Trend L0", Unit: "number"},
	"tikv_raftstore_slow_trend_l0_l1":                   {Label: "Slow Trend L0→L1", Unit: "number"},
	"tikv_raftstore_slow_trend_l1":                      {Label: "Slow Trend L1", Unit: "number"},
	"tikv_raftstore_slow_trend_l1_l2":                   {Label: "Slow Trend L1→L2", Unit: "number"},
	"tikv_raftstore_slow_trend_l1_margin_error":         {Label: "Slow Trend L1 Error", Unit: "number"},
	"tikv_raftstore_slow_trend_l2":                      {Label: "Slow Trend L2", Unit: "number"},
	"tikv_raftstore_slow_trend_l2_margin_error":         {Label: "Slow Trend L2 Error", Unit: "number"},
	"tikv_raftstore_slow_trend_margin_error_gap":        {Label: "Slow Trend Error Gap", Unit: "number"},
	"tikv_raftstore_slow_trend_misc":                    {Label: "Slow Trend Misc", Unit: "number"},
	"tikv_raftstore_slow_trend_result":                  {Label: "Slow Trend Result", Unit: "number"},
	"tikv_raftstore_slow_trend_result_l0":               {Label: "Slow Result L0", Unit: "number"},
	"tikv_raftstore_slow_trend_result_l0_l1":            {Label: "Slow Result L0→L1", Unit: "number"},
	"tikv_raftstore_slow_trend_result_l1":               {Label: "Slow Result L1", Unit: "number"},
	"tikv_raftstore_slow_trend_result_l1_l2":            {Label: "Slow Result L1→L2", Unit: "number"},
	"tikv_raftstore_slow_trend_result_l1_margin_error":  {Label: "Slow Result L1 Error", Unit: "number"},
	"tikv_raftstore_slow_trend_result_l2":               {Label: "Slow Result L2", Unit: "number"},
	"tikv_raftstore_slow_trend_result_l2_margin_error":  {Label: "Slow Result L2 Error", Unit: "number"},
	"tikv_raftstore_slow_trend_result_margin_error_gap": {Label: "Slow Result Error Gap", Unit: "number"},
	"tikv_raftstore_slow_trend_result_misc":             {Label: "Slow Result Misc", Unit: "number"},
	"tikv_raftstore_slow_trend_result_value":            {Label: "Slow Result Value", Unit: "number"},

	// Rate limiter metrics
	"tikv_rate_limiter_max_bytes_per_sec": {Label: "Rate Limit", Unit: "rate"},

	// Read metrics
	"tikv_read_qps_topn": {Label: "Top N Read QPS", Unit: "rate"},

	// Resolved TS metrics
	"tikv_resolved_ts_leader_min_resolved_ts_duration_to_last_update_safe_ts":   {Label: "Leader TS Update Duration", Unit: "time"},
	"tikv_resolved_ts_lock_heap_bytes":                                          {Label: "TS Lock Heap", Unit: "bytes"},
	"tikv_resolved_ts_memory_quota_in_use_bytes":                                {Label: "TS Memory Quota Used", Unit: "bytes"},
	"tikv_resolved_ts_min_follower_resolved_ts":                                 {Label: "Min Follower Resolved TS", Unit: "number"},
	"tikv_resolved_ts_min_follower_resolved_ts_duration_to_last_consume_leader": {Label: "Follower TS Consume Duration", Unit: "time"},
	"tikv_resolved_ts_min_follower_resolved_ts_gap_millis":                      {Label: "Follower Resolved TS Gap", Unit: "time"},
	"tikv_resolved_ts_min_follower_resolved_ts_region":                          {Label: "Min Follower TS Region", Unit: "count"},
	"tikv_resolved_ts_min_follower_safe_ts":                                     {Label: "Min Follower Safe TS", Unit: "number"},
	"tikv_resolved_ts_min_follower_safe_ts_duration_to_last_consume_leader":     {Label: "Follower Safe TS Duration", Unit: "time"},
	"tikv_resolved_ts_min_follower_safe_ts_gap_millis":                          {Label: "Follower Safe TS Gap", Unit: "time"},
	"tikv_resolved_ts_min_follower_safe_ts_region":                              {Label: "Min Follower Safe TS Region", Unit: "count"},
	"tikv_resolved_ts_min_leader_resolved_ts":                                   {Label: "Min Leader Resolved TS", Unit: "number"},
	"tikv_resolved_ts_min_leader_resolved_ts_gap_millis":                        {Label: "Leader Resolved TS Gap", Unit: "time"},
	"tikv_resolved_ts_min_leader_resolved_ts_region":                            {Label: "Min Leader TS Region", Unit: "count"},
	"tikv_resolved_ts_min_leader_resolved_ts_region_min_lock_ts":                {Label: "Leader Min Lock TS", Unit: "number"},
	"tikv_resolved_ts_min_resolved_ts":                                          {Label: "Min Resolved TS", Unit: "number"},
	"tikv_resolved_ts_min_resolved_ts_gap_millis":                               {Label: "Resolved TS Gap", Unit: "time"},
	"tikv_resolved_ts_min_resolved_ts_region":                                   {Label: "Min Resolved TS Region", Unit: "count"},
	"tikv_resolved_ts_pending_count":                                            {Label: "TS Pending Count", Unit: "count"},
	"tikv_resolved_ts_region_resolve_status":                                    {Label: "Region Resolve Status", Unit: "number"},
	"tikv_resolved_ts_scan_tasks":                                               {Label: "TS Scan Tasks", Unit: "count"},
	"tikv_resolved_ts_zero_resolved_ts":                                         {Label: "Zero Resolved TS", Unit: "count"},

	// Scheduler metrics
	"tikv_scheduler_throttle_cf":           {Label: "Scheduler Throttle CF", Unit: "number"},
	"tikv_scheduler_txn_status_cache_size": {Label: "Txn Status Cache", Unit: "count"},
	"tikv_scheduler_write_flow":            {Label: "Write Flow", Unit: "rate"},

	// Server metrics
	"tikv_server_cpu_cores_quota": {Label: "CPU Cores Quota", Unit: "count"},
	"tikv_server_info":            {Label: "Server Info", Unit: "number"},
	"tikv_server_mem_trace_sum":   {Label: "Memory Trace Sum", Unit: "bytes", Default: true},
	"tikv_server_memory_usage":    {Label: "Memory Usage", Unit: "bytes", Default: true},

	// Store metrics
	"tikv_store_size_bytes": {Label: "Store Size", Unit: "bytes", Default: true},

	// Thread metrics
	"tikv_thread_cpu_seconds_total": {Label: "Thread CPU Time", Unit: "time"},
	"tikv_threads_io_bytes_total":   {Label: "Thread I/O", Unit: "bytes", Default: true},
	"tikv_threads_state":            {Label: "Thread State", Unit: "number", Default: true},

	// TTL metrics
	"tikv_ttl_checker_poll_interval":     {Label: "TTL Poll Interval", Unit: "time"},
	"tikv_ttl_checker_processed_regions": {Label: "TTL Processed Regions", Unit: "count"},

	// Unified read pool metrics
	"tikv_unified_read_pool_running_tasks": {Label: "Read Pool Running Tasks", Unit: "count"},
	"tikv_unified_read_pool_thread_count":  {Label: "Read Pool Threads", Unit: "count"},

	// Worker metrics
	"tikv_worker_pending_task_total": {Label: "Worker Pending Tasks", Unit: "count", Default: true},
}

// MetricsMap is built at init time from metricsDefinitions
var MetricsMap map[string]types.MetricStatus

func init() {
	MetricsMap = make(map[string]types.MetricStatus, len(metricsDefinitions))
	for name, def := range metricsDefinitions {
		MetricsMap[name] = types.MetricStatus{
			Enabled: isEnabled(name, def.Default),
			Label:   def.Label,
			Unit:    def.Unit,
		}
	}
}

func isEnabled(metric string, defaultEnabled bool) bool {
	if env := os.Getenv(metric); env != "" {
		return env == "true"
	}
	return defaultEnabled
}
