package main

import "github.com/godev/tolls/types"

type MemoryStore struct {}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	return nil
}