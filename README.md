# telegraf-hcloud

[![Build Status](https://drone.owncloud.com/api/badges/owncloud-ops/telegraf-hcloud/status.svg)](https://drone.owncloud.com/owncloud-ops/telegraf-hcloud)

Gather metrics from Hetzner Cloud resources.

## Configuration

```toml
[[inputs.hcloud]]
  ## List of Hetzner Cloud resources to monitor.
  resources = [
    "load_balancers"
  ]

  ## Hetzner API token.
  token = ""
```

## Metrics

- hcloud_load_balancer
  - tags:
    - location
    - name
  - fields:
    - open_connections (float)
    - requests_per_second (float)
    - bandwidth_in (float)
    - bandwidth_out (float)
    - connections_per_second (float)
    - max_connections (float)
- hcloud_load_balancer
  - tags:
    - id
    - location
    - name
    - protected
    - type
  - fields:
    - info (float)

## Example Output

```plain
hcloud_load_balancer,id=af22,location=nbg1,name=download,protected=false,type=lb11 info=1i 1692800482205737603
hcloud_load_balancer,location=nbg1,name=download max_connections=10000i,bandwidth_in=1724,bandwidth_out=25727,connections_per_second=3,open_connections=39,requests_per_second=0 1692800482205737603
```

## Build

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](https://golang.org/doc/install.html). This project requires Go >= v1.21.

```Shell
make build
./dist/telegraf-hcloud --help
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](https://github.com/owncloud-ops/telegraf-hcloud/blob/main/LICENSE) file for details.
