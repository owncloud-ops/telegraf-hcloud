package hcloud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
	"golang.org/x/exp/slices"
)

const (
	defaultHcloudClientTimeout = 5 * time.Second
	defaultHcloudMetricSteps   = 60

	HcloudResourceLoadBalancer = "load_balancers"
)

// Hcloud - plugin main struct.
type Hcloud struct {
	Resources []string `toml:"resources"`
	Token     string   `toml:"token"`
	Timeout   config.Duration

	client *hcloud.Client
}

const sampleConfig = `
  ## List of Hetzner Cloud resources to monitor.
  resources = [
    "%s"
  ]

  ## Hetzner API token.
  token = ""
`

// SampleConfig returns the default configuration of the Cloudwatch input plugin.
func (h *Hcloud) SampleConfig() string {
	return fmt.Sprintf(sampleConfig, HcloudResourceLoadBalancer)
}

// Description returns a one-sentence description on the Cloudwatch input plugin.
func (h *Hcloud) Description() string {
	return "Collect metrics from Hetzner Cloud resources"
}

// Gather takes in an accumulator and adds the metrics that the Input
// gathers. This is called every "interval".
func (h *Hcloud) Gather(acc telegraf.Accumulator) error {
	h.setDefaultValues()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.Timeout))
	defer cancel()

	if h.client == nil {
		h.client = h.createHcloudClient(h.Token)
	}

	now := time.Now()

	if slices.Contains(h.Resources, HcloudResourceLoadBalancer) {
		h.fetchLoadBalancer(ctx, acc, now)
	}

	return nil
}

func (h *Hcloud) setDefaultValues() {
	if len(h.Resources) == 0 {
		h.Resources = []string{HcloudResourceLoadBalancer}
	}
}

func (h *Hcloud) createHcloudClient(token string) *hcloud.Client {
	client := hcloud.NewClient(hcloud.WithToken(token))

	return client
}

func (h *Hcloud) getTimeSeries(ts map[string][]hcloud.LoadBalancerMetricsValue, name string) (float64, error) {
	if metric, ok := ts[name]; ok {
		if len(metric) > 0 {
			value, err := strconv.ParseFloat(metric[0].Value, 64)
			if err != nil {
				return 0, err
			}

			return value, nil
		}
	}

	return 0, nil
}

func (h *Hcloud) fetchLoadBalancer(ctx context.Context, acc telegraf.Accumulator, start time.Time) {
	lbs, err := h.client.LoadBalancer.All(ctx)
	if err != nil {
		acc.AddError(err)
	}

	lbOpts := hcloud.LoadBalancerGetMetricsOpts{
		Start: start,
		End:   time.Now(),
		Types: []hcloud.LoadBalancerMetricType{
			hcloud.LoadBalancerMetricOpenConnections,
			hcloud.LoadBalancerMetricConnectionsPerSecond,
			hcloud.LoadBalancerMetricRequestsPerSecond,
			hcloud.LoadBalancerMetricBandwidth,
		},
		Step: defaultHcloudMetricSteps,
	}

	for _, lb := range lbs {
		tags := map[string]string{}
		fields := map[string]interface{}{}

		metrics, _, err := h.client.LoadBalancer.GetMetrics(ctx, lb, lbOpts)
		if err != nil {
			acc.AddError(err)
		}

		tags["id"] = strconv.FormatInt(lb.ID, 16)
		tags["instance"] = lb.Name
		tags["type"] = lb.LoadBalancerType.Name
		tags["datacenter"] = lb.Location.Name

		ts := metrics.TimeSeries

		bandwidthIn, err := h.getTimeSeries(ts, "bandwidth.in")
		if err != nil {
			acc.AddError(err)
		} else {
			fields["bandwidth_in"] = bandwidthIn
		}

		bandwidthOut, err := h.getTimeSeries(ts, "bandwidth.out")
		if err != nil {
			acc.AddError(err)
		} else {
			fields["bandwidth_out"] = bandwidthOut
		}

		connectionsPerSecond, err := h.getTimeSeries(ts, "connections_per_second")
		if err != nil {
			acc.AddError(err)
		} else {
			fields["connections_per_second"] = connectionsPerSecond
		}

		openConnections, err := h.getTimeSeries(ts, "open_connections")
		if err != nil {
			acc.AddError(err)
		} else {
			fields["open_connections"] = openConnections
		}

		requestsPerSecond, err := h.getTimeSeries(ts, "requests_per_second")
		if err != nil {
			acc.AddError(err)
		} else {
			fields["requests_per_second"] = requestsPerSecond
		}

		acc.AddFields("hcloud_load_balancer", fields, tags, start)
	}
}

func init() {
	inputs.Add("hcloud", func() telegraf.Input {
		return &Hcloud{
			Timeout: config.Duration(defaultHcloudClientTimeout),
		}
	})
}
