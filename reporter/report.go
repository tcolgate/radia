// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of radia.
//
// radia is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// radia is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with radia.  If not, see <http://www.gnu.org/licenses/>.

// Package reporter provides implementations and definitions for
// services for probes to pass data to.
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
