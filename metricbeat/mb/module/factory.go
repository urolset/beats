package module

import (
	"time"

	"github.com/joeshaw/multierror"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/cfgwarn"
	"github.com/elastic/beats/metricbeat/mb"
)

// Factory creates new Runner instances from configuration objects.
// It is used to register and reload modules.
type Factory struct {
	pipeline      beat.Pipeline
	maxStartDelay time.Duration
}

// NewFactory creates new Reloader instance for the given config
func NewFactory(maxStartDelay time.Duration, p beat.Pipeline) *Factory {
	return &Factory{
		pipeline:      p,
		maxStartDelay: maxStartDelay,
	}
}

func (r *Factory) Create(c *common.Config) (cfgfile.Runner, error) {
	var errs multierror.Errors

	err := cfgwarn.CheckRemoved5xSettings(c, "filters")
	if err != nil {
		errs = append(errs, err)
	}
	connector, err := NewConnector(r.pipeline, c)
	if err != nil {
		errs = append(errs, err)
	}
	w, err := NewWrapper(r.maxStartDelay, c, mb.Registry)
	if err != nil {
		errs = append(errs, err)
	}

	if err := errs.Err(); err != nil {
		return nil, err
	}

	client, err := connector.Connect()
	if err != nil {
		return nil, err
	}

	mr := NewRunner(client, w)
	return mr, nil
}
