package hcloud

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
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

	if h.client == nil {
		h.client = h.createHcloudClient(h.Token)
	}

	if slices.Contains(h.Resources, HcloudResourceLoadBalancer) {
		h.gatherLoadBalancerStats(acc)
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

func (h *Hcloud) gatherLoadBalancerStats(acc telegraf.Accumulator) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.Timeout))

	lbs, err := h.client.LoadBalancer.All(ctx)
	if err != nil {
		acc.AddError(err)
	}

	defer cancel()

	now := time.Now()
	lbOpts := hcloud.LoadBalancerGetMetricsOpts{
		Start: now,
		End:   now,
		Types: []hcloud.LoadBalancerMetricType{
			hcloud.LoadBalancerMetricOpenConnections,
			hcloud.LoadBalancerMetricConnectionsPerSecond,
			hcloud.LoadBalancerMetricRequestsPerSecond,
			hcloud.LoadBalancerMetricBandwidth,
		},
		Step: defaultHcloudMetricSteps,
	}

	for _, lb := range lbs {
		// add info field
		metaTags := map[string]string{}
		metaFields := map[string]interface{}{}

		metaTags["id"] = strconv.FormatInt(lb.ID, 16)
		metaTags["name"] = lb.Name
		metaTags["location"] = lb.Location.Name
		metaTags["type"] = lb.LoadBalancerType.Name
		metaTags["protected"] = strconv.FormatBool(lb.Protection.Delete)

		metaFields["info"] = 1

		acc.AddFields("hcloud_load_balancer", metaFields, metaTags, now)

		// add metrics
		tags := map[string]string{}
		fields := map[string]interface{}{}

		metrics, _, err := h.client.LoadBalancer.GetMetrics(ctx, lb, lbOpts)
		if err != nil {
			acc.AddError(err)
		}

		tags["name"] = lb.Name
		tags["location"] = lb.Location.Name

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

		fields["max_connections"] = lb.LoadBalancerType.MaxConnections

		acc.AddFields("hcloud_load_balancer", fields, tags, now)
	}
}

func init() {
	inputs.Add("hcloud", func() telegraf.Input {
		return &Hcloud{
			Timeout: config.Duration(defaultHcloudClientTimeout),
		}
	})
}
