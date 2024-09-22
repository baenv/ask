// Package dify provides an adapter for interacting with the Dify API.
package dify

// DifyAdapter defines the interface for interacting with the Dify service.
type DifyAdapter interface {
	Chat(msg, url, token string) (content string, err error)
}
