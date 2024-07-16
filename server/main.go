package main

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	monitorGrpc "prom/grpc"
	monitoring "prom/metrics"

	"google.golang.org/grpc"
)

type Server struct {
	monitorGrpc.UnimplementedMonitoringServer
}

func (s Server) GetMonitoringData(ctx context.Context, req *monitorGrpc.MonitoringDataRequest) (*monitorGrpc.MonitoringDataResponse, error) {
	res := monitorGrpc.MonitoringDataResponse{}
	for _, met := range req.Metric {
		switch met {
		case "cpu":
			cli := monitoring.NewMonitoringClient()
			data, err := cli.GetCpuUsage()
			if err != nil {
				return nil, err
			}
			res.MonitoringData = append(res.MonitoringData, &monitorGrpc.MonitoringData{Metric: "cpu", Usage: data.Usage})
		case "mem":
			cli := monitoring.NewMonitoringClient()
			data, err := cli.GetMemUsage()
			if err != nil {
				return nil, err
			}
			res.MonitoringData = append(res.MonitoringData, &monitorGrpc.MonitoringData{Metric: "mem", Usage: data.Usage})
		default:
			return nil, errors.New("invalid request: " + met)
		}
	}
	return &res, nil
}

func getServerAddress() string {
	configFile, err := os.Open("../config.json")
	if err != nil {
		log.Fatalf("unable to get server port:%s", err)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result["grpcServerAddress"].(string)
}

func StartServer() {
	lis, err := net.Listen("tcp", getServerAddress())
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	monServer := grpc.NewServer()
	monService := &Server{}
	monitorGrpc.RegisterMonitoringServer(monServer, monService)
	fmt.Printf("Starting grpc server at 8089")
	err = monServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error initiating the server: %s", err)
	}
}

func main() {
	fmt.Println("Initializing server")
	StartServer()
}
