package db

import (
	met "prom/metrics"
)

type Database interface {
	GetNodes() ([]met.Met, error)
}

type db struct {
	a string
}

func (d db) GetNodes() ([]met.Met, error) {
	var nodes []met.Met
	//putting hardcoded value for now
	var n met.Met = met.NewNode("")
	nodes = append(nodes, n)
	return nodes, nil
}
