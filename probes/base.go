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

package probes

type Probe struct {
	verbosity int
}

type option func(p *Probe) option

// Option sets the options specified.
// It returns an option to restore the last arg's previous value.
func (p *Probe) Option(opts ...option) (previous option) {
	for _, opt := range opts {
		previous = opt(p)
	}
	return previous
}

// Verbosity sets the oauth client's log level
func Verbosity(v int) option {
	return func(p *Probe) option {
		previous := p.verbosity
		p.verbosity = v
		return Verbosity(previous)
	}
}

func (p *Probe) Verbosity() int {
	return p.verbosity
}
