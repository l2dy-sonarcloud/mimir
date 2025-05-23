// SPDX-License-Identifier: AGPL-3.0-only

package scalars

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser/posrange"

	"github.com/grafana/mimir/pkg/streamingpromql/limiting"
	"github.com/grafana/mimir/pkg/streamingpromql/types"
)

// ScalarToInstantVector is an operator that implements the vector() function.
type ScalarToInstantVector struct {
	Scalar                   types.ScalarOperator
	MemoryConsumptionTracker *limiting.MemoryConsumptionTracker

	expressionPosition posrange.PositionRange
	consumed           bool
}

var _ types.InstantVectorOperator = &ScalarToInstantVector{}

func NewScalarToInstantVector(scalar types.ScalarOperator, expressionPosition posrange.PositionRange, memoryConsumptionTracker *limiting.MemoryConsumptionTracker) *ScalarToInstantVector {
	return &ScalarToInstantVector{
		Scalar:                   scalar,
		expressionPosition:       expressionPosition,
		MemoryConsumptionTracker: memoryConsumptionTracker,
	}
}

func (s *ScalarToInstantVector) SeriesMetadata(_ context.Context) ([]types.SeriesMetadata, error) {
	metadata, err := types.SeriesMetadataSlicePool.Get(1, s.MemoryConsumptionTracker)
	if err != nil {
		return nil, err
	}

	metadata = append(metadata, types.SeriesMetadata{
		Labels: labels.EmptyLabels(),
	})

	return metadata, nil
}

func (s *ScalarToInstantVector) NextSeries(ctx context.Context) (types.InstantVectorSeriesData, error) {
	if s.consumed {
		return types.InstantVectorSeriesData{}, types.EOS
	}

	scalarValue, err := s.Scalar.GetValues(ctx)
	if err != nil {
		return types.InstantVectorSeriesData{}, err
	}

	s.consumed = true

	return types.InstantVectorSeriesData{
		Floats: scalarValue.Samples,
	}, nil
}

func (s *ScalarToInstantVector) ExpressionPosition() posrange.PositionRange {
	return s.expressionPosition
}

func (s *ScalarToInstantVector) Close() {
	s.Scalar.Close()
}
