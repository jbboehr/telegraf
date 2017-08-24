// +build nvml

package nvml

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
)

func TestNvmlGather(t *testing.T) {
	var acc testutil.Accumulator
	p := NVMLInput{}

	acc.GatherError(p.Gather)

	tags := map[string]string{"gpu": "0"}
	fields := map[string]interface{}{
		"power_usage":         uint(0),
		"decoder_utilization": uint(0),
		"encoder_utilization": uint(0),
		"utilization":         uint(0),
	}
	acc.AssertContainsTaggedFields(t, "nvml", fields, tags)

	tags = map[string]string{"gpu": "0"}
	acc.AssertContainsTaggedFields(t, "nvml", fields, tags)
}
