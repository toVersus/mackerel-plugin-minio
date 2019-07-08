# mackerel-plugin-minio

[![Build Status](https://travis-ci.org/toVersus/mackerel-plugin-minio.svg?branch=master)](https://travis-ci.org/toVersus/mackerel-plugin-minio)

[Minio](https://min.io/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-minio [-scheme=<url scheme>] [-host=<host>] [-port=<port>] [-metric-path=<path to metrics exporter>] [-metric-key-prefix=<prefix>]
```

## Example of mackerel-agent.conf

```toml
[plugin.metrics.minio]
command = "mackerel-plugin-minio"
```

## Documents

- [How to monitor MinIO server with Prometheus](https://github.com/minio/cookbook/blob/master/docs/how-to-monitor-minio-with-prometheus.md)
