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

package base

import (
	"flag"
	"runtime"
	"strings"

	"github.com/tcolgate/radia/reporter"
)

type Runner interface {
	Run(args ...string)
}

type Probe interface {
	Name() string
	Version() string
	Description() string
	BaseTags() map[string]string
	Report([]reporter.Metric)
	Runner
}

type Base struct {
	*flag.FlagSet

	name        string
	version     string
	description string
	baseTags    map[string]string
	reporter    reporter.Reporter
}

type option func(*Base) error

func defaultName() string {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	pl := len(parts)

	return strings.Split(parts[pl-1], ".")[0]
}

// Option sets the options specified.
// It returns an option to restore the last arg's previous value.
func New(opts ...option) (p *Base, err error) {
	p = &Base{
		name:        defaultName(),
		version:     "0",
		description: "",
	}
	for _, opt := range opts {
		if err = opt(p); err != nil {
			return nil, err
		}
	}
	p.FlagSet = flag.NewFlagSet(p.name, flag.ContinueOnError)
	p.reporter = reporter.New() //use a default from the reporter package

	return p, nil
}

// Name
func Name(s string) option {
	return func(p *Base) error {
		p.name = s
		return nil
	}
}

func (b *Base) Name() string {
	return b.name
}

// Version
func Version(s string) option {
	return func(p *Base) error {
		p.version = s
		return nil
	}
}

func (b *Base) Version() string {
	return b.version
}

// Version
func Description(s string) option {
	return func(p *Base) error {
		p.description = s
		return nil
	}
}

func (b *Base) Description() string {
	return b.description
}

// BaseTags
func BaseTags(ts map[string]string) option {
	return func(p *Base) error {
		p.baseTags = ts
		return nil
	}
}

func (b *Base) BaseTags() map[string]string {
	return b.baseTags
}

func (b *Base) Report(ms []reporter.Metric) {
	for m := range ms {
		ms[m].Name = strings.Join([]string{b.Name(), ms[m].Name}, ".")
	}

	b.reporter.Report(ms)
}
