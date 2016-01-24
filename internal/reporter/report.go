package reporter

import "log"

type Metric struct {
	Name  string
	Tags  map[string]string
	Value float64
}

type Reporter interface {
	Report(m []Metric)
}

func New() Reporter {
	return &influxdb{}
}

type influxdb struct {
}

func (i influxdb) Report(m []Metric) {
	// TODO(tcm): should take these and pass the into a channel
	log.Println(m)
}
