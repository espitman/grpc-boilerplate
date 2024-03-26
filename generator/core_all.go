//go:build mage
// +build mage

package main

import "strings"

func (Core) All() error {
	d := NewCoreDomain()
	d.create()

	p := CorePort(*d)
	p.Dist = strings.Replace(d.Dist, "core/domain/", "core/port/", 1)
	p.create()

	s := CoreService(*d)
	s.Dist = strings.Replace(d.Dist, "core/domain/", "core/service/", 1)
	s.create()

	return nil
}
