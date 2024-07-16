package monitoring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	datatypes "prom/dataTypes"
	"time"

	prometheusClientApi "github.com/prometheus/client_golang/api"
	prometheusQueryApi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Monitoring interface {
	GetCpuUsage() (datatypes.UsageData, error)
	GetMemUsage() (datatypes.UsageData, error)
}

type monitoringClient struct {
	prometheusClient prometheusClientApi.Client
}

func getPrometheusAddress() string {
	configFile, err := os.Open("../config.json")
	if err != nil {
		log.Fatalf("unable to get server port: %s", err)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result["prometheusAddress"].(string)
}

func NewMonitoringClient() monitoringClient {
	prometheusCli, err := prometheusClientApi.NewClient(prometheusClientApi.Config{
		Address: getPrometheusAddress(),
	})
	if err != nil {
		panic(err)
	}
	return monitoringClient{
		prometheusClient: prometheusCli,
	}
}

func (m monitoringClient) GetCpuUsage() (datatypes.UsageData, error) {

	var cpuUsageData datatypes.UsageData

	query := "100 * (1 - avg(rate(node_cpu_seconds_total{mode='idle', instance='node-exporter:9100'}[10m15s])))"
	data, err := getMontoringData(&m.prometheusClient, query)

	if err != nil {
		return cpuUsageData, err
	}

	cpuUsageData, err = formatUsageData(data)

	return cpuUsageData, err
}
func (m monitoringClient) GetMemUsage() (datatypes.UsageData, error) {

	var memUsageData datatypes.UsageData

	query := "(1 - (node_memory_MemAvailable_bytes{instance='node-exporter:9100', job='node-exporter'} / node_memory_MemTotal_bytes{instance='node-exporter:9100', job='node-exporter'})) * 100"
	data, err := getMontoringData(&m.prometheusClient, query)

	if err != nil {
		return memUsageData, err
	}

	memUsageData, err = formatUsageData(data)

	return memUsageData, err
}

func formatUsageData(data model.Vector) (datatypes.UsageData, error) {

	var usageData datatypes.UsageData

	if data.Len() != 1 {
		return usageData, errors.ErrUnsupported
	}

	usageData.Time = data[0].Timestamp.Time().Local()
	usageData.Usage = float64(data[0].Value)

	return usageData, nil
}

func getMontoringData(cli *prometheusClientApi.Client, query string) (model.Vector, error) {

	queryApi := prometheusQueryApi.NewAPI(*cli)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warning, err := queryApi.Query(ctx, query, time.Now(), prometheusQueryApi.WithTimeout(50*time.Second))
	//dont know the reason to put 2 timeout parameters in ctx and second in function
	if err != nil {
		return nil, err
	}
	if warning != nil {
		fmt.Printf("Warnings: %v\n", warning)
	}

	return result.(model.Vector), nil
}

func GetMonitoringDataInRange(cli *prometheusClientApi.Client, query string, startTime time.Time, endTime time.Time) {

	queryApi := prometheusQueryApi.NewAPI(*cli)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	timeRange := prometheusQueryApi.Range{
		Start: startTime,
		End:   endTime,
		Step:  time.Minute,
	}
	result, warning, err := queryApi.QueryRange(ctx, query, timeRange, prometheusQueryApi.WithTimeout(50*time.Second))
	//dont know the reason to put 2 timeout parameters in ctx and second in function
	if err != nil {
		panic(err)
	}
	if warning != nil {
		fmt.Printf("Warnings: %v\n", warning)
	}

	fmt.Printf("Result:\n%+v\n", result)
}
