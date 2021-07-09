package models

import (
	"fmt"
	"strings"
)

type Portfolio struct {
	assets map[string]int
}

func capitalize(a map[string]int) map[string]int {
	b := map[string]int{}
	for k, va := range a {
		k = strings.ToUpper(k)
		if vb, ok := b[k]; ok {
			va += vb
		}
		b[k] = va
	}
	return b
}

func NewPortfolio(assets map[string]int) *Portfolio {
	return &Portfolio{assets: capitalize(assets)}
}

func (p *Portfolio) Apply(t Transaction) error {
	for k, v := range capitalize(t) {
		if a, ok := p.assets[k]; ok {
			v += a
			if v < 0 {
				return fmt.Errorf("%s: negative amount is prohibited", k)
			}
		}
		p.assets[k] = v
	}
	return nil
}

func (p *Portfolio) Assets() map[string]int {
	return p.assets
}
