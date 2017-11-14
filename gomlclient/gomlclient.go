package gomlclient

import (
	"github.com/JulzDiverse/goml"
)

type GomlClient struct{}

func New() *GomlClient {
	return &GomlClient{}
}

func (g *GomlClient) Delete(file []byte, path string) ([]byte, error) {
	return goml.DeleteInMemory(file, path)
}

func (g *GomlClient) Set(file []byte, path string, val string) ([]byte, error) {
	return goml.SetInMemory(file, path, val)
}

func (g *GomlClient) Update(file []byte, path string, val string) ([]byte, error) {
	if _, err := goml.GetInMemory(file, path); err == nil {
		return goml.SetInMemory(file, path, val)
	}
	return file, nil
}
