package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type LocalStore struct {
	mu   sync.RWMutex
	path string
	Data map[string]map[string]string `json:"data"`
}

var globalStore *LocalStore

func getStorePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "wat", "db.json"), nil
}

func initStore() error {
	path, err := getStorePath()
	if err != nil {
		return err
	}
	globalStore = &LocalStore{
		path: path,
		Data: make(map[string]map[string]string),
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return globalStore.save()
		}
		return err
	}
	if err := json.Unmarshal(data, globalStore); err != nil {
		return err
	}
	if globalStore.Data == nil {
		globalStore.Data = make(map[string]map[string]string)
	}
	return nil
}

func (s *LocalStore) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

func (a *App) StorePut(bucket, key, val string) error {
	globalStore.mu.Lock()
	defer globalStore.mu.Unlock()
	if _, ok := globalStore.Data[bucket]; !ok {
		globalStore.Data[bucket] = make(map[string]string)
	}
	globalStore.Data[bucket][key] = val
	return globalStore.save()
}

func (a *App) StoreGet(bucket, key string) (string, error) {
	globalStore.mu.RLock()
	defer globalStore.mu.RUnlock()
	if b, ok := globalStore.Data[bucket]; ok {
		if val, exists := b[key]; exists {
			return val, nil
		}
	}
	return "", nil
}

func (a *App) StoreDelete(bucket, key string) error {
	globalStore.mu.Lock()
	defer globalStore.mu.Unlock()
	if b, ok := globalStore.Data[bucket]; ok {
		delete(b, key)
		return globalStore.save()
	}
	return nil
}

func (a *App) StoreList(bucket string) (map[string]string, error) {
	globalStore.mu.RLock()
	defer globalStore.mu.RUnlock()
	res := make(map[string]string)
	if b, ok := globalStore.Data[bucket]; ok {
		for k, v := range b {
			res[k] = v
		}
	}
	return res, nil
}
