package node

import (
	"fmt"
	"strconv"

	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
)

var (
	aggregatorCounter *uint64
	producerCounter   map[string]uint64
	aggregatorTable   map[string]*Aggregator
)

func init() {
	if aggregatorCounter == nil {
		var aggregatorCounterValue uint64
		aggregatorCounter = &aggregatorCounterValue
	}
	if producerCounter == nil {
		producerCounter = make(map[string]uint64)
	}
	if aggregatorTable == nil {
		aggregatorTable = make(map[string]*Aggregator)
	}
}

func getAggregatorId() string {
	aggregatorId := strconv.FormatUint(*aggregatorCounter, 10)
	*aggregatorCounter++
	return aggregatorId
}

func getProducerId(aggregatorId string) (string, uint64) {
	if _, ok := producerCounter[aggregatorId]; !ok {
		producerCounter[aggregatorId] = 0
	} else {
		producerCounter[aggregatorId]++
	}
	producerId := strconv.FormatUint(producerCounter[aggregatorId], 10)
	return fmt.Sprintf("%v.%v", producerId, aggregatorId), producerCounter[aggregatorId]
}

func setAggregator(aggregator *Aggregator) {
	if aggregator == nil {
		log.Error("nil aggregator", "err", utils.NilPtrDeref)
		return
	}
	aggregatorTable[aggregator.Id] = aggregator
}

func getAggregator(aggregatorId string) *Aggregator {
	if aggregator, ok := aggregatorTable[aggregatorId]; ok {
		return aggregator
	}
	return nil
}
