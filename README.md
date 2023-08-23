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
- hcloud_load_balancer_info
  - tags:
    - location
    - name
    - protected
    - type
  - fields:
    - info (float)

## Example Output

```plain
hcloud_load_balancer,location=nbg1,name=download,protected=false,type=lb11 info=1i 1692799714846451262
hcloud_load_balancer,location=nbg1,name=download open_connections=43,requests_per_second=0,max_connections=10000i,bandwidth_in=2886,bandwidth_out=4091507.666667,connections_per_second=3.333333 1692799714846451262
```

## Build

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](https://golang.org/doc/install.html). This project requires Go >= v1.21.

```Shell
git clone https://github.owncloud.com/owncloud-ops/telegraf-hcloud.git
cd telegraf-hcloud

make generate build
./dist/telegraf-hcloud --help
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](https://github.com/owncloud-ops/telegraf-hcloud/blob/main/LICENSE) file for details.
