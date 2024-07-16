package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	monitorGrpc "prom/grpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	getCpuDataFlag := flag.Bool("cpu", false, "get cpu data")
	getMemDataFlag := flag.Bool("mem", false, "get mem data")
	flag.Parse()

	// fmt.Println(*getCpuDataFlag, *getMemDataFlag)

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost:8089", opts)
	if err != nil {
		log.Fatalf("unable to establish connection: %s", err)
	}
	defer conn.Close()

	client := monitorGrpc.NewMonitoringClient(conn)

	req := monitorGrpc.MonitoringDataRequest{}
	if *getCpuDataFlag {
		req.Metric = append(req.Metric, "cpu")
	}
	if *getMemDataFlag {
		req.Metric = append(req.Metric, "mem")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.GetMonitoringData(ctx, &req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res)
}
