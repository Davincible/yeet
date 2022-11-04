package config

import (
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

// WrapFilter will call the filter and log a message if it contains nil nodes
func WrapFilter(f func([]*registry.Service) []*registry.Service) func([]*registry.Service) []*registry.Service {
	return func(s []*registry.Service) []*registry.Service {
		nodes := f(s)

		foundNil := false

		if foundNil {
			logger.Logf("[WARNING] Node filter returned a list contianing nil nodes: %+v", nodes)
		}
	}
}

func checkNilNode(s []*registry.Service) bool {
	for _, node := range nodes {
		if node == nil {
			return true
		}
	}
}
