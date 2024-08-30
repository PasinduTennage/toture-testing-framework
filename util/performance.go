package util

type Performance struct {
	Option map[string]string
}

func NewBasicPerformance() *Performance {
	basic_stats := make(map[string]string)
	basic_stats["throughput"] = ""
	basic_stats["median_latency"] = ""
	basic_stats["p99_latency"] = ""
	basic_stats["average_latency"] = ""
	return &Performance{
		Option: basic_stats,
	}
}

func NewPerformance() Performance {
	return Performance{}
}

func NewPerformanceWithOptions(options map[string]string) Performance {
	return Performance{
		Option: options,
	}
}
