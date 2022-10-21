//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func NewWireTest() (*ReturnType, error) {
	wire.Build(
		ProvideOption,
		ProvideOptionsOne,
		ProvideOptionsMultiple,
		ProvideOptionsTwo,
		ProvideOptionTwo,
		ProvideCollector,
	)

	return nil, nil
}
