package network

import (
	"fmt"
)

const (
	FilterFormat = "to_node_%v: "
)

func GetPubEndpoint(broadcastPort string) string {
	return fmt.Sprintf("tcp://*:%v", broadcastPort)
}

func GetSubEndpoint(broadcastPort string, rank uint64) string {
	// aggregator rank == 0, producer rank > 0
	// sub endpoint: tcp://127.#{b}.#{c}.#{d}:#{broadcastPort}
	b := rank >> 16
	if b > 255 {
		return ""
	}
	c := (rank - (b << 16)) >> 8
	d := rank - (b << 16) - (c << 8)
	return fmt.Sprintf("tcp://127.%v.%v.%v:%v", b, c, d, broadcastPort)
}

func GetFilter(rank uint64) string {
	return fmt.Sprintf(FilterFormat, rank)
}
