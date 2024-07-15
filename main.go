package main

import (
	monitoringServer "prom/grpcServer"
)

func main() {
	// cli := monitoring.NewMonitoringClient()
	// d1, _ := cli.GetCpuUsage()
	// d2, _ := cli.GetMemUsage()
	// fmt.Println(d1, d2)
	monitoringServer.StartServer()
}
