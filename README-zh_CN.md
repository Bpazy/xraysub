[English](./README.md) | 简体中文
<h1 align="center">Xraysub</h1>

<div align="center">

[![Test](https://github.com/Bpazy/xraysub/workflows/Test/badge.svg)](https://github.com/Bpazy/xraysub/actions/workflows/test.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Bpazy_xraysub&metric=alert_status)](https://sonarcloud.io/dashboard?id=Bpazy_xraysub)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bpazy/xraysub)](https://goreportcard.com/report/github.com/Bpazy/xraysub)

跨平台 Xray 命令行订阅管理客户端。
</div>

## 预览
```
$ ./xraysub.exe gen --url=https://comfysub.example.com
Requesting subscriptions from https://comfysub.example.com
Start detecting server's latency
        Detecting 100% [========================================]  [6s:0s]
Got fastest node "中国-香港 IEPL(hk.example.com:61307)" with latency 127ms
The xray-core's configuration file is saved ./xray-config.json
```

## 快速开始
假设 `xray.exe` and `xraysub.exe` 已经在当前目录
```
.
├── xray.exe
└── xraysub.exe
```
1. 首先运行 `xraysub` 来获取 xray-core 的配置文件
```
$ ./xraysub.exe gen --url=https://comfysub.example.com
```
2. 运行 xray-core
```
$ ./xray-core.exe -c xray-config.json
```
3. 使用代理
```
$ curl -x HTTPS_PROXY://127.0.0.1:1081 https://www.google.com
```

## 参数
```
$ .\xraysub.exe help gen
从订阅链接生成 xray-core 的配置文件。

Usage:
  xraysub gen [flags]

Flags:
      --detect-latency             检测服务器的延迟以选择最快的节点 (default true)
      --detect-thread-number int   检测线程数 (default 5)
  -h, --help                       help for gen
  -o, --output-file string         输出配置到指定位置 (default "./xray-config.json")
  -u, --url string                 订阅地址 (URL)
      --xray string                xray-core 路径，用于检测服务器的延迟 (default "./xray.exe")
      --xray-http-port int         xray-core 监听的 HTTP 端口 (default 1081)
      --xray-socks-port int        xray-core 监听的 socks 端口 (default 1080)
```
