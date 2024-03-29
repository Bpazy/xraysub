English | [简体中文](./README-zh_CN.md)
<h1 align="center">Xraysub</h1>

<div align="center">

[![Test](https://github.com/Bpazy/xraysub/workflows/Test/badge.svg)](https://github.com/Bpazy/xraysub/actions/workflows/test.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Bpazy_xraysub&metric=alert_status)](https://sonarcloud.io/dashboard?id=Bpazy_xraysub)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bpazy/xraysub)](https://goreportcard.com/report/github.com/Bpazy/xraysub)

A powerful cross-platform CLI client for Xray subscription.
</div>

## Preview
```
$ ./xraysub.exe gen --url=https://comfysub.example.com
Requesting subscriptions from https://comfysub.example.com
Start detecting server's latency
        Detecting 100% [========================================]  [6s:0s]
Got fastest node "中国-香港 IEPL(hk.example.com:61307)" with latency 127ms
The xray-core's configuration file is saved ./xray-config.json
```

## Quick Start
Suppose the `xray.exe` and `xraysub.exe` are in the current directory.
```
.
├── xray.exe
└── xraysub.exe
```
1. First run `xraysub` to get xray-core's configuration file.
```
$ ./xraysub.exe gen --url=https://comfysub.example.com
```
2. Run xray-core
```
$ ./xray-core.exe -c xray-config.json
```
3. Use the proxy
```
$ curl -x HTTPS_PROXY://127.0.0.1:1081 https://www.google.com
```

## Param
```
$ .\xraysub.exe help gen
generate xray configuration file from subscription url

Usage:
  xraysub gen [flags]

Flags:
      --detect-latency             detect server's latency to choose the fastest node (default true)
      --detect-thread-number int   detect server's latency threads number (default 5)
  -h, --help                       help for gen
  -o, --output-file string         output configuration to file (default "./xray-config.json")
  -u, --url string                 subscription address(URL)
      --xray string                xray-core path for detecting server's latency (default "./xray.exe")
      --xray-http-port int         xray-core listen http port (default 1081)
      --xray-socks-port int        xray-core listen socks port (default 1080)
```
