//go:build plugin
// +build plugin

package main

import (
	"fmt"
	"math"

	"pipelined.dev/audio/vst2"
)

func main() {}

func LinearToDecibel(value float32) float32 {
	if value != 0 {
		return float32(20. * math.Log(float64(value)))
	} else {
		return -144.0
	}
}

func percentMono(v float32) float32 {
	return float32(math.Abs(float64(0.5-v))) / .5
}

func init() {
	var (
		uniqueID = [4]byte{'m', 'o', 'n', 'o'}
		version  = int32(1000)
	)
	vst2.PluginAllocator = func(h vst2.Host) (vst2.Plugin, vst2.Dispatcher) {
		gain := vst2.Parameter{
			Name:  "gain",
			Unit:  "db",
			Value: 1,
			GetValueLabelFunc: func(value float32) string {
				return fmt.Sprintf("%+.2f", value)
			},
			GetValueFunc: LinearToDecibel,
		}
		skew := vst2.Parameter{
			Name:  "skew",
			Unit:  "LR",
			Value: 0.5,
			GetValueLabelFunc: func(value float32) string {
				if value > .5 {
					return fmt.Sprintf("%+.2f % %s", (value-.5)/.5, "Right")
				} else {
					return fmt.Sprintf("%+.2f % %s", (value-.5)/.5, "Left")
				}
			},
			GetValueFunc: func(value float32) float32 {
				return value
			},
		}
		mode := vst2.Parameter{
			Name:  "mode",
			Unit:  "",
			Value: 0.5,
			GetValueLabelFunc: func(value float32) string {
				label := "sculpt"
				if value > .5 {
					label = "smash"
				}
				return label
			},
			GetValueFunc: func(value float32) float32 {
				return value
			},
		}
		merge := vst2.Parameter{
			Name:  "merge",
			Unit:  "%",
			Value: 0.0,
			GetValueLabelFunc: func(value float32) string {
				if mode.GetValue() > .5 {
					return "off"
				}
				return fmt.Sprintf("%.2f", value)
			},
			GetValueFunc: func(value float32) float32 {
				return value
			},
		}
		invert := vst2.Parameter{
			Name:  "invert",
			Unit:  "",
			Value: 0.0,
			GetValueLabelFunc: func(value float32) string {
				if value > .5 {
					return "on"
				}
				return "off"
			},
			GetValueFunc: func(value float32) float32 {
				return value
			},
		}
		channels := 2
		return vst2.Plugin{
			UniqueID:       uniqueID,
			Version:        version,
			InputChannels:  channels,
			OutputChannels: channels,
			Name:           "Monomer",
			Vendor:         "kendfss",
			Category:       vst2.PluginCategoryEffect,
			Parameters: []*vst2.Parameter{
				&gain,
				&invert,
				&merge,
				&mode,
				&skew,
			},
			ProcessDoubleFunc: func(in, out vst2.DoubleBuffer) {
				volume := float32(math.Pow(10, float64(gain.GetValue())/20))

				for i := 0; i < in.Frames; i++ {
					var (
						// emergent signals
						emL float32
						emR float32
					)

					inL := float32(in.Channel(1)[i])
					inR := float32(in.Channel(0)[i])

					if invert.GetValue() > .5 {
						inL, inR = inR, inL
					}

					if mode.GetValue() >= .5 {
						// smash
						factor := skew.GetValue()
						emL = inL*factor + inR*(1-factor)
						emR = inR*(1-factor) + inL*factor
					} else {
						// sculpt
						factor := skew.GetValue()
						emL = inL*factor + inR*(1-factor)
						emL = inL*(1-merge.GetValue()) + emL*merge.GetValue()
						emR = inR*(1-factor) + inL*factor
						emR = inR*(1-merge.GetValue()) + emR*merge.GetValue()
					}

					out.Channel(1)[i] = float64(emL * volume)
					out.Channel(0)[i] = float64(emR * volume)
				}
			},
			ProcessFloatFunc: func(in, out vst2.FloatBuffer) {
				volume := float32(math.Pow(10, float64(gain.GetValue())/20))

				for i := 0; i < in.Frames; i++ {
					var (
						// emergent signals
						emL float32
						emR float32
					)
					inL := in.Channel(1)[i]
					inR := in.Channel(0)[i]

					if invert.GetValue() > .5 {
						inL, inR = inR, inL
					}

					if mode.GetValue() >= .5 {
						// smash
						factor := skew.GetValue()
						emL = inL*factor + inR*(1-factor)
						emR = inR*(1-factor) + inL*factor
					} else {
						// sculpt
						factor := skew.GetValue()
						emL = inL*factor + inR*(1-factor)
						emL = inL*(1-merge.GetValue()) + emL*merge.GetValue()
						emR = inR*(1-factor) + inL*factor
						emR = inR*(1-merge.GetValue()) + emR*merge.GetValue()
					}

					out.Channel(1)[i] = emL * volume
					out.Channel(0)[i] = emR * volume
				}
			},
		}, vst2.Dispatcher{}
	}
}
