package node

import (
	"fmt"
	"strconv"
)

var (
	aggregatorCounter *uint64
	producerCounter   map[string]uint64
	aggregatorTable   map[string]*Aggregator
	producerTable     map[string]*Producer
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
	if producerTable == nil {
		producerTable = make(map[string]*Producer)
	}
}

func getAggregatorId() string {
	if aggregatorCounter == nil {
		return ""
	}
	aggregatorId := strconv.FormatUint(*aggregatorCounter, 10)
	*aggregatorCounter++
	return aggregatorId
}

func getProducerId(aggregatorId string) (string, uint64) {
	if _, ok := producerCounter[aggregatorId]; !ok {
		producerCounter[aggregatorId] = 1
	} else {
		producerCounter[aggregatorId]++
	}
	producerId := strconv.FormatUint(producerCounter[aggregatorId], 10)
	return fmt.Sprintf("%v.%v", producerId, aggregatorId), producerCounter[aggregatorId]
}

func getAggregator(aggregatorId string) *Aggregator {
	if aggregatorTable == nil {
		return nil
	}
	if aggregator, ok := aggregatorTable[aggregatorId]; ok {
		return aggregator
	}
	return nil
}

func setAggregator(aggregator *Aggregator) {
	if aggregator == nil {
		return
	}
	aggregatorTable[aggregator.Id] = aggregator
}

func getProducer(producerId string) *Producer {
	if producerTable == nil {
		return nil
	}
	if producer, ok := producerTable[producerId]; ok {
		return producer
	}
	return nil
}

func setProducer(producer *Producer) {
	if producerTable == nil {
		return
	}
	producerTable[producer.Id] = producer
}
