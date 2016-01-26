// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of vonq.
//
// vonq is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// vonq is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with vonq.  If not, see <http://www.gnu.org/licenses/>.

package register

import (
	"errors"
	"log"

	"github.com/tcolgate/vonq/probes/base"
)

var ErrDuplicateProbe = errors.New("Duplicate probe")

var probes = map[string]func() base.Probe{}

func Probe(fp func() base.Probe) error {
	p := fp()
	if _, ok := probes[p.Name()]; ok {
		log.Println("duplicate probe: ", p.Name())
		return ErrDuplicateProbe
	}
	probes[p.Name()] = fp
	log.Println("registered probe: ", p.Name())
	return nil
}

func RunAll() {
	for _, fp := range probes {
		p := fp()
		go p.Run()
	}
}
