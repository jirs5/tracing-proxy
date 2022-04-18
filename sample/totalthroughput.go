package sample

import (
	"github.com/sirupsen/logrus"
	"math/rand"

	dynsampler "github.com/honeycombio/dynsampler-go"

	"github.com/jirs5/tracing-proxy/config"
	"github.com/jirs5/tracing-proxy/metrics"
	"github.com/jirs5/tracing-proxy/types"
)

type TotalThroughputSampler struct {
	Config  *config.TotalThroughputSamplerConfig
	Logger  *logrus.Logger
	Metrics metrics.Metrics

	goalThroughputPerSec int64
	clearFrequencySec    int64

	key *traceKey

	dynsampler dynsampler.Sampler
}

func (d *TotalThroughputSampler) Start() error {
	d.Logger.Debugf("Starting TotalThroughputSampler")
	defer func() { d.Logger.Debugf("Finished starting TotalThroughputSampler") }()
	if d.Config.GoalThroughputPerSec < 1 {
		d.Logger.Debugf("configured sample rate for dynamic sampler was %d; forcing to 100", d.Config.GoalThroughputPerSec)
		d.Config.GoalThroughputPerSec = 100
	}
	d.goalThroughputPerSec = d.Config.GoalThroughputPerSec
	if d.Config.ClearFrequencySec == 0 {
		d.Config.ClearFrequencySec = 30
	}
	d.clearFrequencySec = d.Config.ClearFrequencySec
	d.key = newTraceKey(d.Config.FieldList, d.Config.UseTraceLength, d.Config.AddSampleRateKeyToTrace, d.Config.AddSampleRateKeyToTraceField)

	// spin up the actual dynamic sampler
	d.dynsampler = &dynsampler.TotalThroughput{
		GoalThroughputPerSec: int(d.goalThroughputPerSec),
		ClearFrequencySec:    int(d.clearFrequencySec),
	}
	d.dynsampler.Start()

	// Register stastics this package will produce
	d.Metrics.Register("dynsampler_num_dropped", "counter")
	d.Metrics.Register("dynsampler_num_kept", "counter")
	d.Metrics.Register("dynsampler_sample_rate", "histogram")

	return nil
}

func (d *TotalThroughputSampler) GetSampleRate(trace *types.Trace) (uint, bool) {
	key := d.key.buildAndAdd(trace)
	rate := d.dynsampler.GetSampleRate(key)
	if rate < 1 { // protect against dynsampler being broken even though it shouldn't be
		rate = 1
	}
	shouldKeep := rand.Intn(int(rate)) == 0
	d.Logger.WithFields(map[string]interface{}{
		"sample_key":  key,
		"sample_rate": rate,
		"sample_keep": shouldKeep,
		"trace_id":    trace.TraceID,
	}).Logf(logrus.DebugLevel, "got sample rate and decision")
	if shouldKeep {
		d.Metrics.Increment("dynsampler_num_kept")
	} else {
		d.Metrics.Increment("dynsampler_num_dropped")
	}
	d.Metrics.Histogram("dynsampler_sample_rate", float64(rate))
	return uint(rate), shouldKeep
}
