// +build debug

package eventsourcing

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStreamNameShouldReturnAggregateNameAndId(t *testing.T) {
	agg := &UnitTestAggregate{AggregateCommonImpl: AggregateCommonImpl{IdImpl: "my-name-is-billy"}}
	assert.Equal(t, "UnitTestAggregate|my-name-is-billy", GetStreamName(agg))
}

func TestLoadAggregate(t *testing.T) {
	agg := &UnitTestAggregate{AggregateCommonImpl: AggregateCommonImpl{IdImpl: "my-name-is-billy"}}
	repo, err := NewRepository("gonatseventsourcing_cluster",
		"test_client"+strconv.Itoa(random(0, 99999)),
		stan.NatsURL(stan.DefaultNatsURL))
	assert.Nil(t, err)
	repo.LoadAggregate(agg)
}
func TestApplyEventAndLoadAggregate(t *testing.T) {
	repo, err := NewRepository("gonatseventsourcing_cluster",
		"test_client"+strconv.Itoa(random(0, 99999)),
		stan.NatsURL(stan.DefaultNatsURL))
	assert.Nil(t, err)

	aggregateId := strconv.Itoa(random(0, 9999999))
	agg := &UnitTestAggregate{
		AggregateCommonImpl: AggregateCommonImpl{
			IdImpl:     aggregateId,
			Repository: repo,
		},
	}
	var i int32
	i = 51
	agg.PerformTest(i)

	var loadedAgg = &UnitTestAggregate{
		AggregateCommonImpl: AggregateCommonImpl{
			IdImpl: aggregateId,
		},
	}

	repo.LoadAggregate(loadedAgg)
	assert.Equal(t, i, loadedAgg.RecordedTestMessageI)
}
