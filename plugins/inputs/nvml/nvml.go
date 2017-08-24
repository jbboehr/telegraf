// +build nvml

package nvml

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/jbboehr/go-nvml"
	"log"
	"strconv"
	"sync"
)

type NVMLInput struct {
	initialized bool
	gpus        []nvml.Device
}

func (_ *NVMLInput) Description() string {
	return "Nvidia management library input"
}

const sampleConfig = `
`

func (_ *NVMLInput) SampleConfig() string {
	return sampleConfig
}

func (p *NVMLInput) Gather(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup
	var err error

	// Get list of GPUs
	if !p.initialized {
		p.gpus, err = nvml.GetAllGPUs()
		if err != nil {
			return err
		}
		p.initialized = true
	}

	for x := range p.gpus {
		wg.Add(1)
		go func(u *nvml.Device) {

			defer wg.Done()
			fields := map[string]interface{}{}
			tags := map[string]string{}

			// GPU number
			tags["gpu"] = strconv.Itoa(x)

			// Fan speed
			fan_speed, err := u.FanSpeed()
			if err != nil {
				log.Printf("W! failed to read nvml fan speed: %s", err)
			} else {
				fields["fan_speed"] = fan_speed
			}

			// Power usage
			power_usage, err := u.PowerUsage()
			if err != nil {
				log.Printf("W! failed to read nvml power usage: %s", err)
			} else {
				fields["power_usage"] = power_usage
			}

			// Temperature
			temp, err := u.Temp()
			if err != nil {
				log.Printf("W! failed to read nvml temp: %s", err)
			} else {
				fields["temp"] = temp
			}

			// Memory
			memory_info, err := u.MemoryInfo()
			if err != nil {
				log.Printf("W! failed to read nvml memory: %s", err)
			} else {
				fields["total_memory"] = memory_info.Total
				fields["free_memory"] = memory_info.Free
				fields["used_memory"] = memory_info.Used
			}

			// Decoder utilization
			decoder_utilization, _, err := u.GetDecoderUtilization()
			if err != nil {
				log.Printf("W! failed to read nvml decoder utilization: %s", err)
			} else {
				fields["decoder_utilization"] = decoder_utilization
			}

			// Encoder utilization
			encoder_utilization, _, err := u.GetEncoderUtilization()
			if err != nil {
				log.Printf("W! failed to read nvml encoder utilization: %s", err)
			} else {
				fields["encoder_utilization"] = encoder_utilization
			}

			// Utilization
			utilization, _, err := u.GetUtilizationRates()
			if err != nil {
				log.Printf("W! failed to read nvml utilization: %s", err)
			} else {
				fields["utilization"] = utilization
			}

			acc.AddFields("nvml", fields, tags)
		}(&p.gpus[x])
	}

	wg.Wait()

	return nil
}

func init() {
	nvml.NVMLInit()
	inputs.Add("nvml", func() telegraf.Input {
		return &NVMLInput{}
	})
}
