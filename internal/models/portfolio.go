package models

import "fmt"

type Portfolio struct {
	assets map[string]int
}

func NewPortfolio(assets map[string]int) *Portfolio {
	return &Portfolio{assets: assets}
}

func (p *Portfolio) Apply(t Transaction) error {
	for k, v := range t {
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
