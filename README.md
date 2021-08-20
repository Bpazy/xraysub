<h1 align="center">Xraysub</h1>

<div align="center">

[![process](https://img.shields.io/badge/dev-in%20Progress-yellow)](https://github.com/Bpazy/xraysub/projects/1])
[![Test](https://github.com/Bpazy/xraysub/workflows/Test/badge.svg)](https://github.com/Bpazy/xraysub/actions/workflows/test.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Bpazy_xraysub&metric=alert_status)](https://sonarcloud.io/dashboard?id=Bpazy_xraysub)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bpazy/xraysub)](https://goreportcard.com/report/github.com/Bpazy/xraysub)

A powerful cross-platform CLI client for Xray subscription. 跨平台 Xray 命令行订阅管理客户端。
</div>

## Preview
```powershell
./xraysub.exe gen --url=https://comfysub.example.com --xray D:/MyPrograms/xray-core/xray.exe --xray-socks-port 1080 --xray-http-port 1081
Requesting subscriptions from https://comfysub.example.com
Start detecting server's latency
        Detecting 100% [========================================]    6s:0s]s]
Got fastest node: hkoqze01.nq1too321coo1oo.xyz:51507
The xray-core's configuration file is saved ./xray-config.json
```

## Quick Start
Suppose the `xray-core.exe` and `xraysub.exe` are in the current directory.
1. First run `xraysub` to get xray-core's configuration file.
```powershell
./xraysub.exe gen --url=https://comfysub.example.com --xray-socks-port 1080 --xray-http-port 1081
```
2. Run xray-core
```
./xray-core.exe -c xray-config.json
```
3. Use the proxy
```
curl -x HTTPS_PROXY://127.0.0.1:1081 https://www.google.com
```

## Param
```
PS D:\workspace\Bpazy\xraysub> .\xraysub.exe help gen
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
