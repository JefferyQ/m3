// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aggregator

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/m3db/m3aggregator/aggregator/handler"
	"github.com/m3db/m3aggregator/aggregator/handler/writer"
	"github.com/m3db/m3aggregator/client"
	"github.com/m3db/m3metrics/aggregation"
	"github.com/m3db/m3metrics/metadata"
	"github.com/m3db/m3metrics/metric"
	"github.com/m3db/m3metrics/metric/aggregated"
	"github.com/m3db/m3metrics/metric/unaggregated"
	"github.com/m3db/m3metrics/op"
	"github.com/m3db/m3metrics/op/applied"
	"github.com/m3db/m3metrics/policy"
	"github.com/m3db/m3x/clock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestBaseMetricListPushBack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l, err := newBaseMetricList(testShard, time.Second, nil, nil, testOptions(ctrl))
	require.NoError(t, err)
	elem, err := NewCounterElem(nil, policy.EmptyStoragePolicy, aggregation.DefaultTypes, applied.DefaultPipeline, 0, l.opts)
	require.NoError(t, err)

	// Push a counter to the list.
	e, err := l.PushBack(elem)
	require.NoError(t, err)
	require.Equal(t, 1, l.aggregations.Len())
	require.Equal(t, elem, e.Value.(*CounterElem))

	// Push a counter to a closed list should result in an error.
	l.Lock()
	l.closed = true
	l.Unlock()

	_, err = l.PushBack(elem)
	require.Equal(t, err, errListClosed)
}

func TestBaseMetricListClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opts := testOptions(ctrl)
	l, err := newBaseMetricList(testShard, time.Second, nil, nil, opts)
	require.NoError(t, err)

	l.RLock()
	require.False(t, l.closed)
	l.RUnlock()

	l.Close()
	require.True(t, l.closed)

	// Close for a second time should have no impact.
	l.Close()
	require.True(t, l.closed)
}

func TestBaseMetricListFlushWithRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		now              = time.Unix(12345, 0)
		nowFn            = func() time.Time { return now }
		isEarlierThanFn  = isStandardMetricEarlierThan
		timestampNanosFn = standardMetricTimestampNanos
		results          []flushBeforeResult
	)
	opts := testOptions(ctrl).SetClockOptions(clock.NewOptions().SetNowFn(nowFn))
	l, err := newBaseMetricList(testShard, time.Second, isEarlierThanFn, timestampNanosFn, opts)
	require.NoError(t, err)
	l.flushBeforeFn = func(beforeNanos int64, flushType flushType) {
		results = append(results, flushBeforeResult{
			beforeNanos: beforeNanos,
			flushType:   flushType,
		})
	}

	inputs := []struct {
		request  flushRequest
		expected []flushBeforeResult
	}{
		{
			request: flushRequest{
				CutoverNanos:      20000 * int64(time.Second),
				CutoffNanos:       30000 * int64(time.Second),
				BufferAfterCutoff: time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 12345 * int64(time.Second),
					flushType:   discardType,
				},
			},
		},
		{
			request: flushRequest{
				CutoverNanos:      10000 * int64(time.Second),
				CutoffNanos:       30000 * int64(time.Second),
				BufferAfterCutoff: time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 10000 * int64(time.Second),
					flushType:   discardType,
				},
				{
					beforeNanos: 12345 * int64(time.Second),
					flushType:   consumeType,
				},
			},
		},
		{
			request: flushRequest{
				CutoverNanos:      10000 * int64(time.Second),
				CutoffNanos:       12300 * int64(time.Second),
				BufferAfterCutoff: time.Minute,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 10000 * int64(time.Second),
					flushType:   discardType,
				},
				{
					beforeNanos: 12300 * int64(time.Second),
					flushType:   consumeType,
				},
			},
		},
		{
			request: flushRequest{
				CutoverNanos:      10000 * int64(time.Second),
				CutoffNanos:       12300 * int64(time.Second),
				BufferAfterCutoff: 10 * time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 10000 * int64(time.Second),
					flushType:   discardType,
				},
				{
					beforeNanos: 12300 * int64(time.Second),
					flushType:   consumeType,
				},
				{
					beforeNanos: 12335 * int64(time.Second),
					flushType:   discardType,
				},
			},
		},
		{
			request: flushRequest{
				CutoverNanos:      0,
				CutoffNanos:       30000 * int64(time.Second),
				BufferAfterCutoff: time.Second,
			},
			expected: []flushBeforeResult{
				{
					beforeNanos: 12345 * int64(time.Second),
					flushType:   consumeType,
				},
			},
		},
	}
	for _, input := range inputs {
		results = results[:0]
		l.Flush(input.request)
		require.Equal(t, input.expected, results)
	}
}

func TestBaseMetricListFlushBeforeStale(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		isEarlierThanFn  = isStandardMetricEarlierThan
		timestampNanosFn = standardMetricTimestampNanos
		opts             = testOptions(ctrl)
	)
	l, err := newBaseMetricList(testShard, 0, isEarlierThanFn, timestampNanosFn, opts)
	require.NoError(t, err)
	l.lastFlushedNanos = 1234
	l.flushBefore(1000, discardType)
	require.Equal(t, int64(1234), l.LastFlushedNanos())
}

func TestStandardMetricListID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolution := 10 * time.Second
	opts := testOptions(ctrl)
	listID := standardMetricListID{resolution: resolution}
	l, err := newStandardMetricList(testShard, listID, opts)
	require.NoError(t, err)

	expectedListID := metricListID{
		listType: standardMetricListType,
		standard: listID,
	}
	require.Equal(t, expectedListID, l.ID())
}

func TestStandardMetricListFlushConsumingAndCollectingLocalMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		errTestFlush = errors.New("foo")
		cutoverNanos = int64(0)
		cutoffNanos  = int64(math.MaxInt64)
		count        int
		flushLock    sync.Mutex
		flushed      []aggregated.ChunkedMetricWithStoragePolicy
	)

	// Intentionally cause a one-time error during encoding.
	writeFn := func(mp aggregated.ChunkedMetricWithStoragePolicy) error {
		flushLock.Lock()
		defer flushLock.Unlock()

		if count == 0 {
			count++
			return errTestFlush
		}
		flushed = append(flushed, mp)
		return nil
	}
	w := writer.NewMockWriter(ctrl)
	w.EXPECT().Write(gomock.Any()).DoAndReturn(writeFn).AnyTimes()
	w.EXPECT().Flush().Return(nil).AnyTimes()
	handler := handler.NewMockHandler(ctrl)
	handler.EXPECT().NewWriter(gomock.Any()).Return(w, nil).AnyTimes()

	var (
		now        = time.Unix(216, 0).UnixNano()
		nowTs      = time.Unix(0, now)
		resolution = testStoragePolicy.Resolution().Window
	)
	clockOpts := clock.NewOptions().SetNowFn(func() time.Time {
		return time.Unix(0, atomic.LoadInt64(&now))
	})
	opts := testOptions(ctrl).
		SetClockOptions(clockOpts).
		SetFlushHandler(handler)

	listID := standardMetricListID{resolution: resolution}
	l, err := newStandardMetricList(testShard, listID, opts)
	require.NoError(t, err)

	// Intentionally cause a one-time error during encoding.
	elemPairs := []struct {
		elem   metricElem
		metric unaggregated.MetricUnion
	}{
		{
			elem:   MustNewCounterElem(testCounterID, testStoragePolicy, aggregation.DefaultTypes, applied.DefaultPipeline, 0, opts),
			metric: testCounter,
		},
		{
			elem:   MustNewTimerElem(testBatchTimerID, testStoragePolicy, aggregation.DefaultTypes, applied.DefaultPipeline, 0, opts),
			metric: testBatchTimer,
		},
		{
			elem:   MustNewGaugeElem(testGaugeID, testStoragePolicy, aggregation.DefaultTypes, applied.DefaultPipeline, 0, opts),
			metric: testGauge,
		},
	}

	for _, ep := range elemPairs {
		require.NoError(t, ep.elem.AddUnion(nowTs, ep.metric))
		require.NoError(t, ep.elem.AddUnion(nowTs.Add(l.resolution), ep.metric))
		_, err := l.PushBack(ep.elem)
		require.NoError(t, err)
	}

	// Force a flush.
	l.Flush(flushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert nothing has been flushed.
	flushLock.Lock()
	require.Equal(t, 0, len(flushed))
	flushLock.Unlock()

	for i := 0; i < 2; i++ {
		// Move the time forward by one aggregation interval.
		nowTs = nowTs.Add(l.resolution)
		atomic.StoreInt64(&now, nowTs.UnixNano())

		// Force a flush.
		l.Flush(flushRequest{
			CutoverNanos: cutoverNanos,
			CutoffNanos:  cutoffNanos,
		})

		var expected []testLocalMetricWithMetadata
		alignedStart := nowTs.Truncate(l.resolution).UnixNano()
		expected = append(expected, expectedLocalMetricsForCounter(alignedStart, testStoragePolicy, aggregation.DefaultTypes)...)
		expected = append(expected, expectedLocalMetricsForTimer(alignedStart, testStoragePolicy, aggregation.DefaultTypes)...)
		expected = append(expected, expectedLocalMetricsForGauge(alignedStart, testStoragePolicy, aggregation.DefaultTypes)...)

		// Skip the first item because we intentionally triggered
		// an encoder error when encoding the first item.
		if i == 0 {
			expected = expected[1:]
		}

		flushLock.Lock()
		require.NotNil(t, flushed)
		validateLocalFlushed(t, expected, flushed)
		flushed = flushed[:0]
		flushLock.Unlock()
	}

	// Move the time forward by one aggregation interval.
	nowTs = nowTs.Add(l.resolution)
	atomic.StoreInt64(&now, nowTs.UnixNano())

	// Force a flush.
	l.Flush(flushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert nothing has been flushed.
	flushLock.Lock()
	require.Equal(t, 0, len(flushed))
	flushLock.Unlock()
	require.Equal(t, 3, l.aggregations.Len())

	// Mark all elements as tombstoned.
	for e := l.aggregations.Front(); e != nil; e = e.Next() {
		e.Value.(metricElem).MarkAsTombstoned()
	}

	// Move the time forward and force a flush.
	nowTs = nowTs.Add(l.resolution)
	atomic.StoreInt64(&now, nowTs.UnixNano())
	l.Flush(flushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert all elements have been collected.
	require.Equal(t, 0, l.aggregations.Len())

	require.Equal(t, l.lastFlushedNanos, nowTs.UnixNano())
}

func TestStandardMetricListClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		registered   int
		unregistered int
	)
	flushManager := NewMockFlushManager(ctrl)
	flushManager.EXPECT().
		Register(gomock.Any()).
		DoAndReturn(func(flushingMetricList) error {
			registered++
			return nil
		})
	flushManager.EXPECT().
		Unregister(gomock.Any()).
		DoAndReturn(func(flushingMetricList) error {
			unregistered++
			return nil
		})

	resolution := 10 * time.Second
	opts := testOptions(ctrl).SetFlushManager(flushManager)
	listID := standardMetricListID{resolution: resolution}
	l, err := newStandardMetricList(testShard, listID, opts)
	require.NoError(t, err)

	l.RLock()
	require.False(t, l.closed)
	l.RUnlock()
	require.Equal(t, 1, registered)

	l.Close()
	require.True(t, l.closed)
	require.Equal(t, 1, unregistered)

	// Close for a second time should have no impact.
	l.Close()
	require.True(t, l.closed)
	require.Equal(t, 1, registered)
	require.Equal(t, 1, unregistered)
}

func TestForwardedMetricListID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolution := 10 * time.Second
	opts := testOptions(ctrl)
	listID := forwardedMetricListID{resolution: resolution, numForwardedTimes: testNumForwardedTimes}
	l, err := newForwardedMetricList(testShard, listID, opts)
	require.NoError(t, err)

	expectedListID := metricListID{
		listType:  forwardedMetricListType,
		forwarded: listID,
	}
	require.Equal(t, expectedListID, l.ID())
}

func TestForwardedMetricListFlushOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	maxForwardingDelayFn := func(resolution time.Duration, numForwardedTimes int) time.Duration {
		return resolution + time.Second*time.Duration(numForwardedTimes)
	}
	resolution := 10 * time.Second
	opts := testOptions(ctrl).SetMaxAllowedForwardingDelayFn(maxForwardingDelayFn)
	listID := forwardedMetricListID{resolution: resolution, numForwardedTimes: 2}
	l, err := newForwardedMetricList(testShard, listID, opts)
	require.NoError(t, err)

	require.Equal(t, 2*time.Second, l.FlushOffset())
}

func TestForwardedMetricListFlushConsumingAndCollectingForwardedMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		errTestWrite = errors.New("foo")
		cutoverNanos = int64(0)
		cutoffNanos  = int64(math.MaxInt64)
		count        int
		flushLock    sync.Mutex
		flushed      []aggregated.MetricWithForwardMetadata
	)

	// Intentionally cause a one-time error during encoding.
	writeFn := func(metric aggregated.Metric, meta metadata.ForwardMetadata) error {
		flushLock.Lock()
		defer flushLock.Unlock()

		if count == 0 {
			count++
			return errTestWrite
		}
		flushed = append(flushed, aggregated.MetricWithForwardMetadata{
			Metric:          metric,
			ForwardMetadata: meta,
		})
		return nil
	}

	client := client.NewMockAdminClient(ctrl)
	client.EXPECT().WriteForwarded(gomock.Any(), gomock.Any()).DoAndReturn(writeFn).MinTimes(1)
	client.EXPECT().Flush().Return(nil).MinTimes(1)

	var (
		now                  = time.Unix(216, 0).UnixNano()
		nowTs                = time.Unix(0, now)
		resolution           = testStoragePolicy.Resolution().Window
		alignedTimeNanos     = nowTs.Truncate(resolution).UnixNano()
		maxLatenessAllowed   = 9 * time.Second
		maxForwardingDelayFn = func(time.Duration, int) time.Duration { return maxLatenessAllowed }
	)
	clockOpts := clock.NewOptions().SetNowFn(func() time.Time {
		return time.Unix(0, atomic.LoadInt64(&now))
	})
	opts := testOptions(ctrl).
		SetClockOptions(clockOpts).
		SetAdminClient(client).
		SetMaxAllowedForwardingDelayFn(maxForwardingDelayFn).
		setSourceIDProvider(newSourceIDProvider(testSourceID))

	listID := forwardedMetricListID{
		resolution:        resolution,
		numForwardedTimes: testNumForwardedTimes,
	}
	l, err := newForwardedMetricList(testShard, listID, opts)
	require.NoError(t, err)

	pipeline := applied.NewPipeline([]applied.Union{
		{
			Type: op.RollupType,
			Rollup: applied.Rollup{
				ID:            []byte("foo.bar"),
				AggregationID: aggregation.MustCompressTypes(aggregation.Max),
			},
		},
	})
	expectedMetadata := metadata.ForwardMetadata{
		AggregationID:     aggregation.MustCompressTypes(aggregation.Max),
		StoragePolicy:     testStoragePolicy,
		Pipeline:          applied.NewPipeline([]applied.Union{}),
		SourceID:          testSourceID,
		NumForwardedTimes: testNumForwardedTimes + 1,
	}
	elemPairs := []struct {
		elem   metricElem
		metric aggregated.Metric
	}{
		{
			elem: MustNewCounterElem([]byte("testForwardedCounter"), testStoragePolicy, aggregation.DefaultTypes, pipeline, testNumForwardedTimes, opts),
			metric: aggregated.Metric{
				Type:      metric.CounterType,
				ID:        []byte("testForwardedCounter"),
				TimeNanos: alignedTimeNanos,
				Value:     123,
			},
		},
		{
			elem: MustNewGaugeElem([]byte("testForwardedGauge"), testStoragePolicy, aggregation.DefaultTypes, pipeline, testNumForwardedTimes, opts),
			metric: aggregated.Metric{
				Type:      metric.GaugeType,
				ID:        []byte("testForwardedGauge"),
				TimeNanos: alignedTimeNanos,
				Value:     1.762,
			},
		},
	}

	for _, ep := range elemPairs {
		require.NoError(t, ep.elem.AddUnique(time.Unix(0, ep.metric.TimeNanos), ep.metric.Value, 1))
		require.NoError(t, ep.elem.AddUnique(time.Unix(0, ep.metric.TimeNanos).Add(l.resolution), ep.metric.Value, 1))
		_, err := l.PushBack(ep.elem)
		require.NoError(t, err)
	}

	// Force a flush.
	l.Flush(flushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert nothing has been flushed.
	flushLock.Lock()
	require.Equal(t, 0, len(flushed))
	flushLock.Unlock()

	for i := 0; i < 2; i++ {
		// Move the time forward by one aggregation interval.
		nowTs = nowTs.Add(l.resolution)
		atomic.StoreInt64(&now, nowTs.UnixNano())

		// Force a flush.
		l.Flush(flushRequest{
			CutoverNanos: cutoverNanos,
			CutoffNanos:  cutoffNanos,
		})

		var expected []aggregated.MetricWithForwardMetadata
		alignedStart := (nowTs.Add(-maxLatenessAllowed)).Truncate(l.resolution).UnixNano()
		for _, ep := range elemPairs {
			expectedMetric := aggregated.Metric{
				Type:      ep.metric.Type,
				ID:        []byte("foo.bar"),
				TimeNanos: alignedStart,
				Value:     ep.metric.Value,
			}
			expected = append(expected, aggregated.MetricWithForwardMetadata{
				Metric:          expectedMetric,
				ForwardMetadata: expectedMetadata,
			})
		}

		// Skip the first item because we intentionally triggered
		// an encoder error when encoding the first item.
		if i == 0 {
			expected = expected[1:]
		}

		flushLock.Lock()
		require.NotNil(t, flushed)
		require.Equal(t, expected, flushed)
		flushed = flushed[:0]
		flushLock.Unlock()
	}

	// Move the time forward by one aggregation interval.
	nowTs = nowTs.Add(l.resolution)
	atomic.StoreInt64(&now, nowTs.UnixNano())

	// Force a flush.
	l.Flush(flushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert nothing has been flushed.
	flushLock.Lock()
	require.Equal(t, 0, len(flushed))
	flushLock.Unlock()
	require.Equal(t, 2, l.aggregations.Len())

	// Mark all elements as tombstoned.
	for e := l.aggregations.Front(); e != nil; e = e.Next() {
		e.Value.(metricElem).MarkAsTombstoned()
	}

	// Move the time forward and force a flush.
	nowTs = nowTs.Add(l.resolution)
	atomic.StoreInt64(&now, nowTs.UnixNano())
	l.Flush(flushRequest{
		CutoverNanos: cutoverNanos,
		CutoffNanos:  cutoffNanos,
	})

	// Assert all elements have been collected.
	require.Equal(t, 0, l.aggregations.Len())

	require.Equal(t, l.lastFlushedNanos, nowTs.UnixNano())
}

func TestForwardedMetricListClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		registered   int
		unregistered int
	)
	flushManager := NewMockFlushManager(ctrl)
	flushManager.EXPECT().
		Register(gomock.Any()).
		DoAndReturn(func(flushingMetricList) error {
			registered++
			return nil
		})
	flushManager.EXPECT().
		Unregister(gomock.Any()).
		DoAndReturn(func(flushingMetricList) error {
			unregistered++
			return nil
		})

	resolution := 10 * time.Second
	opts := testOptions(ctrl).SetFlushManager(flushManager)
	listID := forwardedMetricListID{resolution: resolution, numForwardedTimes: testNumForwardedTimes}
	l, err := newForwardedMetricList(testShard, listID, opts)
	require.NoError(t, err)

	l.RLock()
	require.False(t, l.closed)
	l.RUnlock()
	require.Equal(t, 1, registered)

	l.Close()
	require.True(t, l.closed)
	require.Equal(t, 1, unregistered)

	// Close for a second time should have no impact.
	l.Close()
	require.True(t, l.closed)
	require.Equal(t, 1, registered)
	require.Equal(t, 1, unregistered)
}

func TestMetricLists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opts := testOptions(ctrl)
	lists := newMetricLists(testShard, opts)
	require.False(t, lists.closed)

	// Create a new standard metric list.
	listID := standardMetricListID{resolution: time.Second}.toMetricListID()
	sl, err := lists.FindOrCreate(listID)
	require.NoError(t, err)
	require.NotNil(t, sl)
	require.Equal(t, 1, lists.Len())

	// Find the same standard metric list.
	sl2, err := lists.FindOrCreate(listID)
	require.NoError(t, err)
	require.True(t, sl == sl2)
	require.Equal(t, 1, lists.Len())

	// Create a new forwarded metric list.
	listID = forwardedMetricListID{
		resolution:        10 * time.Second,
		numForwardedTimes: testNumForwardedTimes,
	}.toMetricListID()
	fl, err := lists.FindOrCreate(listID)
	require.NoError(t, err)
	require.NotNil(t, fl)
	require.Equal(t, 2, lists.Len())

	// Find the same forwarded metric list.
	fl2, err := lists.FindOrCreate(listID)
	require.NoError(t, err)
	require.True(t, fl == fl2)
	require.Equal(t, 2, lists.Len())

	// Perform a tick.
	tickRes := lists.Tick()
	expectedRes := map[time.Duration]int{
		time.Second:      0,
		10 * time.Second: 0,
	}
	require.Equal(t, expectedRes, tickRes)

	// Finding or creating in a closed list should result in an error.
	lists.Close()
	_, err = lists.FindOrCreate(listID)
	require.Equal(t, errListsClosed, err)
	require.True(t, lists.closed)

	// Closing a second time should have no impact.
	lists.Close()
	require.True(t, lists.closed)
}

func validateLocalFlushed(
	t *testing.T,
	expected []testLocalMetricWithMetadata,
	flushed []aggregated.ChunkedMetricWithStoragePolicy,
) {
	require.Equal(t, len(expected), len(flushed))
	for i := 0; i < len(flushed); i++ {
		require.Equal(t, expected[i].idPrefix, flushed[i].ChunkedID.Prefix)
		require.Equal(t, []byte(expected[i].id), flushed[i].ChunkedID.Data)
		require.Equal(t, expected[i].idSuffix, flushed[i].ChunkedID.Suffix)
		require.Equal(t, expected[i].timeNanos, flushed[i].TimeNanos)
		require.Equal(t, expected[i].value, flushed[i].Value)
		require.Equal(t, expected[i].sp, flushed[i].StoragePolicy)
	}
}

type flushBeforeResult struct {
	beforeNanos int64
	flushType   flushType
}
