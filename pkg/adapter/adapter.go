// Package adapter provides an interface and implementation for adapting various services.

package adapter

import (
	"ask/pkg/adapter/dify"
	"ask/pkg/config"
)

// IAdapter defines the interface for adapters.
type IAdapter interface {
	Dify() dify.DifyAdapter
}

// Adapter implements the IAdapter interface.
type Adapter struct {
	dify dify.DifyAdapter
}

// New creates a new Adapter instance with the provided configuration.
func New(cfg config.Config) IAdapter {
	return &Adapter{
		dify: dify.New(),
	}
}

// Dify returns the DifyAdapter instance.
func (a *Adapter) Dify() dify.DifyAdapter {
	return a.dify
}
