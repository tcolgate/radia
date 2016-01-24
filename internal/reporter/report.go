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
	return nil
}

type influxdb struct {
}

func (i influxdb) Report(m []Metric) {
	log.Println(m)
}
