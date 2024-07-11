package monitoring

import (
	"context"
	"fmt"
	"time"

	prometheusClientApi "github.com/prometheus/client_golang/api"
	prometheusQueryApi "github.com/prometheus/client_golang/api/prometheus/v1"
)

type Monitoring interface {
	GetMontoringData(metric string)
	GetMonitoringDataInRange(metric string, startTime time.Time, endTime time.Time)
}

type monitoringClient struct {
	prometheusClient prometheusClientApi.Client
}

func NewMonitoringClient() monitoringClient {
	prometheusCli, err := prometheusClientApi.NewClient(prometheusClientApi.Config{
		Address: "http://localhost:9090",
	})
	if err != nil {
		panic(err)
	}
	return monitoringClient{
		prometheusClient: prometheusCli,
	}
}

func (m monitoringClient) GetMontoringData(metric string) {

	queryApi := prometheusQueryApi.NewAPI(m.prometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warning, err := queryApi.Query(ctx, metric, time.Now(), prometheusQueryApi.WithTimeout(50*time.Second))
	//dont know the reason to put 2 timeout parameters in ctx and second in function
	if err != nil {
		panic(err)
	}
	if warning != nil {
		fmt.Printf("Warnings: %v\n", warning)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func (m monitoringClient) GetMonitoringDataInRange(metric string, startTime time.Time, endTime time.Time) {

	queryApi := prometheusQueryApi.NewAPI(m.prometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	timeRange := prometheusQueryApi.Range{
		Start: startTime,
		End:   endTime,
		Step:  time.Minute,
	}
	result, warning, err := queryApi.QueryRange(ctx, metric, timeRange, prometheusQueryApi.WithTimeout(50*time.Second))
	//dont know the reason to put 2 timeout parameters in ctx and second in function
	if err != nil {
		panic(err)
	}
	if warning != nil {
		fmt.Printf("Warnings: %v\n", warning)
	}
	fmt.Printf("Result:\n%v\n", result)
}
