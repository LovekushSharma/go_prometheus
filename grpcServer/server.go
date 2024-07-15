package monitoringServer

import (
	context "context"
	"fmt"
	"log"
	"net"
	monitoring "prom/metrics"

	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedMonitoringServer
}

func (s Server) GetMonitoringData(ctx context.Context, req *MonitoringDataRequest) (*MonitoringDataResponse, error) {
	res := MonitoringDataResponse{}
	for _, met := range req.Metric {
		switch met {
		case "cpu":
			cli := monitoring.NewMonitoringClient()
			data, err := cli.GetCpuUsage()
			if err != nil {
				return nil, err
			}
			res.MonitoringData = append(res.MonitoringData, &MonitoringData{Metric: "cpu", Usage: data.Usage})
		case "mem":
			cli := monitoring.NewMonitoringClient()
			data, err := cli.GetMemUsage()
			if err != nil {
				return nil, err
			}
			res.MonitoringData = append(res.MonitoringData, &MonitoringData{Metric: "mem", Usage: data.Usage})
		}
	}
	return &res, nil
}

func StartServer() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	monServer := grpc.NewServer()
	monService := &Server{}
	RegisterMonitoringServer(monServer, monService)
	fmt.Printf("Starting grpc server at 8089")
	err = monServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error initiating the server: %s", err)
	}
}
