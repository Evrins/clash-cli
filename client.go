package main

import (
	"bufio"
	"bytes"
	"clash-cli/log_level"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Client struct {
	Host       string
	Token      string
	HttpClient *http.Client
}

func NewClient(host, token string) *Client {
	return &Client{
		Host:       host,
		Token:      token,
		HttpClient: http.DefaultClient,
	}
}

type Traffic struct {
	Up   int
	Down int
}

func (cli *Client) Traffic() (err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/traffic", cli.Host), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	reader := bufio.NewReader(resp.Body)
	for {
		if resp.Close {
			break
		}
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		line = bytes.TrimSpace(line)
		var t Traffic
		err = json.Unmarshal(line, &t)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\r\u2191 %d \u2193 %d", t.Up, t.Down)
		os.Stdout.Sync()
	}
	return
}

type Log struct {
	Type    log_level.LogLevel
	Payload string
}

func (cli *Client) Logs(level log_level.LogLevel) (err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/logs", cli.Host), nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("level", string(level))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	reader := bufio.NewReader(resp.Body)
	for {
		if resp.Close {
			break
		}
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		line = bytes.TrimSpace(line)
		var l Log
		err = json.Unmarshal(line, &l)
		if err != nil {
			panic(err)
		}
		fmt.Printf("level: %s %s\n", l.Type, l.Payload)
	}
	return
}

type ProxyType string

const (
	Direct      ProxyType = "Direct"
	Reject      ProxyType = "Reject"
	Selector    ProxyType = "Selector"
	Shadowsocks ProxyType = "Shadowsocks"
	Socks5      ProxyType = "Socks5"
	URLTest     ProxyType = "URLTest"
)

type Proxies struct {
	Proxies map[string]*Proxy
}

type Proxy struct {
	Type ProxyType `json:"type"`
	All  []string  `json:"all"`
	Now  string    `json:"now"`
}

func (cli *Client) ProxyList() (proxies Proxies, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/proxies", cli.Host), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &proxies)
	if err != nil {
		return
	}
	return
}

func (cli *Client) Proxy(name string) (proxy *Proxy, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/proxy/%s", cli.Host, name), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &proxy)
	if err != nil {
		return
	}
	return
}

func (cli *Client) ProxyDelay(name string, timeout int, url string) (proxy *Proxy, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/proxy/%s/delay", cli.Host, name), nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("timeout", strconv.Itoa(timeout))
	q.Add("url", url)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &proxy)
	if err != nil {
		return
	}
	return
}

type Error struct {
	Error string `json:"error"`
}

func (cli *Client) SelectProxy(previous, next string) (err error) {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/proxies/%s", cli.Host, previous), nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("name", next)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode == http.StatusOK {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var e Error
	err = json.Unmarshal(content, &e)
	if err != nil {
		return
	}
	err = errors.New(e.Error)
	return
}

type Config struct {
	Port       int    `json:"port"`
	SocketPort int    `json:"socket-port"`
	RedirPort  int    `json:"redir-port"`
	AllowLan   bool   `json:"allow-lan"`
	Mode       string `json:"mode"`
	LogLevel   string `json:"log-level"`
}

func (cli *Client) GetConfigs() (config *Config, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/configs", cli.Host), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return
	}
	return
}

func (cli *Client) SetConfigs(config *Config) (err error) {
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(config)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("http://%s/configs", cli.Host), buf)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusNoContent {
		err = errors.New("not returning http code 204")
		return
	}
	return
}

func (cli *Client) ReloadConfig(force bool, path *string) (err error) {
	buf := new(bytes.Buffer)
	if path == nil {
		buf = nil
	} else {
		err = json.NewEncoder(buf).Encode(struct {
			Path *string `json:"path"`
		}{Path: path})
		if err != nil {
			return
		}
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/configs", cli.Host), buf)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("force", strconv.FormatBool(force))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("not returning http code 200")
		return
	}
	return
}

type Rules struct {
	Rules []*Rule `json:"rules"`
}

type Rule struct {
	Type    string
	Payload string
	Proxy   string
}

func (cli *Client) Rules() (rules *Rules, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/rules", cli.Host), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &rules)
	if err != nil {
		return
	}
	return
}
