package internal

import (
	"context"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	// ScopeName is unique name of the meter used for the instumentation.
	ScopeName = "dubbo.apache.org/dubbo-go/filter/otel/metrics"
	// Version is current version of the instrumentation.
	Version = constant.OtelPackageVersion
)

type Instrumentation struct {
	metrics  *Metrics
	preAttrs []attribute.KeyValue
}

func NewInstrumentation(kind, name string) (*Instrumentation, error) {
	mp := otel.GetMeterProvider()
	meter := mp.Meter(ScopeName, metric.WithInstrumentationVersion(Version))
	metrics, err := newMetrics(kind, meter)
	if err != nil {
		return nil, err
	}
	return &Instrumentation{metrics: metrics}, nil
}

func (inst *Instrumentation) Start(ctx context.Context, serviceKey string, methodName string) *InvokeOp {
	return &InvokeOp{
		start: time.Now(),
		inst:  inst,
		ctx:   ctx,
	}
}

type InvokeOp struct {
	serviceKey string
	methodName string
	inst       *Instrumentation
	start      time.Time
	ctx        context.Context
	count      int64
}

func (i *InvokeOp) End(err error) {
	if err != nil {
		i.inst.metrics.InvokedCouner.Add(i.ctx, i.count)
	} else {
		i.inst.metrics.InvokedCouner.Add(i.ctx, i.count)
	}
	i.inst.metrics.SecondsHistogram.Record(i.ctx, float64(time.Since(i.start)))
}
