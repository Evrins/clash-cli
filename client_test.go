package main

import (
	"clash-cli/log_level"
	"fmt"
	"testing"
)

var client *Client

func init() {
	client = NewClient("172.16.44.123:9090", "571a959d945830805ec3c963ca3ed675e9046601")
}

func TestClient_Traffic(t *testing.T) {
	err := client.Traffic()
	if err != nil {
		panic(err)
	}
}

func TestClient_Logs(t *testing.T) {
	err := client.Logs(log_level.Debug)
	if err != nil {
		panic(err)
	}
}

func TestClient_ProxyList(t *testing.T) {
	proxies, err := client.ProxyList()
	if err != nil {
		panic(err)
	}
	fmt.Println(proxies)
}

func TestClient_GetConfigs(t *testing.T) {
	configs, err := client.GetConfigs()
	if err != nil {
		panic(err)
	}
	fmt.Println(configs)
}

func TestClient_ReloadConfig(t *testing.T) {
	err := client.ReloadConfig(false, nil)
	if err != nil {
		panic(err)
	}
}

func TestClient_Rules(t *testing.T) {
	rules, err := client.Rules()
	if err != nil {
		panic(err)
	}
	fmt.Println(rules)
}
