package entities

import (
	"math"

	"github.com/mlflow/mlflow-go-backend/pkg/protos"
	"github.com/mlflow/mlflow-go-backend/pkg/utils"
)

type Metric struct {
	Key           string
	Value         float64
	Timestamp     int64
	Step          int64
	IsNaN         bool
	ModelID       string
	DatasetName   string
	DatasetDigest string
}

func (m Metric) ToProto() *protos.Metric {
	metric := protos.Metric{
		Key:       &m.Key,
		Value:     &m.Value,
		Timestamp: &m.Timestamp,
		Step:      &m.Step,
	}

	switch {
	case m.IsNaN:
		metric.Value = utils.PtrTo(math.NaN())
	default:
		metric.Value = &m.Value
	}

	return &metric
}

func MetricFromProto(proto *protos.Metric) *Metric {
	return &Metric{
		Key:           proto.GetKey(),
		Value:         proto.GetValue(),
		Timestamp:     proto.GetTimestamp(),
		Step:          proto.GetStep(),
		ModelID:       proto.GetModelId(),
		DatasetName:   proto.GetDatasetName(),
		DatasetDigest: proto.GetDatasetDigest(),
	}
}

func MetricFromLogMetricProtoInput(input *protos.LogMetric) *Metric {
	return &Metric{
		Key:           input.GetKey(),
		Value:         input.GetValue(),
		Timestamp:     input.GetTimestamp(),
		Step:          input.GetStep(),
		ModelID:       input.GetModelId(),
		DatasetName:   input.GetDatasetName(),
		DatasetDigest: input.GetDatasetDigest(),
	}
}
