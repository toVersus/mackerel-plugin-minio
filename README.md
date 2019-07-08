# mackerel-plugin-minio

[![Build Status](https://travis-ci.org/toVersus/mackerel-plugin-minio.svg?branch=master)](https://travis-ci.org/toVersus/mackerel-plugin-minio)

[Minio](https://min.io/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-minio [-scheme=<url scheme>] [-host=<host>] [-port=<port>] [-metric-path=<path to metrics exporter>] [-metric-key-prefix=<prefix>]
```

## Installation

Installing mackerel-plugin-minio by using [mkr](https://mackerel.io/docs/entry/advanced/cli) as follows:

```bash
sudo mkr plugin install toVersus/mackerel-plugin-minio@v0.1.0
```

## Example of mackerel-agent.conf

```toml
[plugin.metrics.minio]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-minio"
```

## Documents

- [How to monitor MinIO server with Prometheus](https://github.com/minio/cookbook/blob/master/docs/how-to-monitor-minio-with-prometheus.md)
