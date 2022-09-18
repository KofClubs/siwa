package network

import (
	"fmt"
)

const (
	AggregatorFilter     = "to_aggregator: "
	ProducerFilterFormat = "to_producer_%v: "
)

func GetPubEndpoint(broadcastPort string) string {
	return fmt.Sprintf("tcp://*:%v", broadcastPort)
}

func GetAggregatorSubEndpoint(broadcastPort string) string {
	// ip of aggregator: 127.0.0.0
	return fmt.Sprintf("tcp://127.0.0.0:%v", broadcastPort)
}

func GetProducerSubEndpoint(broadcastPort string, rankOfProducer uint64) string {
	// ip of producer: [127.0.0.1, 127.255.255.255]
	index := rankOfProducer + 1
	b := index >> 16
	if b > 255 {
		return ""
	}
	c := (index - (b << 16)) >> 8
	d := index - (b << 16) - (c << 8)
	return fmt.Sprintf("tcp://127.%v.%v.%v:%v", b, c, d, broadcastPort)
}

func GetProducerFilter(rankOfProducer uint64) string {
	return fmt.Sprintf(ProducerFilterFormat, rankOfProducer)
}
