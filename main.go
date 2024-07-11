package main

import (
	monitoring "prom/metrics"
)

func main() {
	cli := monitoring.NewMonitoringClient()
	cli.GetMontoringData("container_memory_usage_bytes")
	// s := time.Now().Add(-time.Hour)
	// e := time.Now()
	// cli.GetMonitoringDataInRange("container_memory_usage_bytes", s, e)
}
