package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Portfolio struct {
	id      uuid.UUID
	name    string
	assets  map[string]int
	balance float64
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
