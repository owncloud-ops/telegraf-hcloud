package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/influxdata/telegraf/plugins/common/shim"
	_ "github.com/owncloud-ops/telegraf-hcloud/plugins/inputs/hcloud"
)

const defaultPollInterval = 1 * time.Minute

var (
	pollInterval         = flag.Duration("poll_interval", defaultPollInterval, "how often to send metrics")
	pollIntervalDisabled = flag.Bool("poll_interval_disabled", false, "set to true to disable polling")
	configFile           = flag.String("config", "", "path to the config file for this plugin")
	err                  error
)

func main() {
	// Parse command line options
	flag.Parse()

	if *pollIntervalDisabled {
		*pollInterval = shim.PollIntervalDisabled
	}

	// Create the shim. This is what will run your plugins.
	shim := shim.New()

	// If no config is specified, all imported plugins are loaded.
	// otherwise follow what the config asks for.
	err = shim.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err loading input: %s\n", err)
		os.Exit(1)
	}

	// Run the input plugin(s) until stdin closes or we receive a termination signal
	if err := shim.Run(*pollInterval); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}
}
