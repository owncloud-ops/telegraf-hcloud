# telegraf-hcloud

[![Build Status](https://drone.owncloud.com/api/badges/owncloud-ops/telegraf-hcloud/status.svg)](https://drone.owncloud.com/owncloud-ops/telegraf-hcloud)
[![Docker Hub](https://img.shields.io/badge/docker-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/owncloudops/telegraf-hcloud)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/owncloudops/telegraf-hcloud)

Gather metrics from Hetzner Cloud resources.

## Configuration

```toml
[[inputs.hcloud]]
  ## List of Hetzner Cloud resources to monitor.
  resources = [
    "load_balancers"
  ]

  ## Hetzner API token.
  # token = ""
```

## Metrics

- hcloud_load_balancer
  - tags:
    - datacenter
    - id
    - instance
    - type
  - fields:
    - open_connections (float)
    - requests_per_second (float)
    - bandwidth_in (float)
    - bandwidth_out (float)
    - connections_per_second (float)

## Example Output

```plain
hcloud_load_balancer,datacenter=nbg1,id=44834,instance=download,type=lb21 open_connections=210,requests_per_second=0,bandwidth_in=948.333333,bandwidth_out=12049.333333,connections_per_second=2.333333 1691592199806644652
```

## Build

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](https://golang.org/doc/install.html). This project requires Go >= v1.20.

```Shell
git clone https://github.owncloud.com/owncloud-ops/telegraf-hcloud.git
cd telegraf-hcloud

make generate build
./dist/telegraf-hcloud --help
```

To build the container use:

```Shell
docker build -f Dockerfile -t telegraf-hcloud:latest .
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](https://github.com/owncloud-ops/telegraf-hcloud/blob/main/LICENSE) file for details.
