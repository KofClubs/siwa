/*
Copyright (c) 2022 Zhang Zhanpeng <zhangregister@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package node

import (
	"time"

	"github.com/KofClubs/siwa/crypto"
	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	"github.com/MonteCarloClub/zmq"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

type AggregatorEntity struct {
	BroadcastPort string `yaml:"broadcast_port"`
}

type Aggregator struct {
	Id            string
	ProducerIds   map[string]struct{}
	Suite         *bn256.Suite
	Threshold     int
	BroadcastPort string
	ZmqSocketSet  *zmq.SocketSet
}

func (aggregatorEntity *AggregatorEntity) CreateAggregator() *Aggregator {
	if aggregatorEntity == nil {
		log.Error("nil aggregator entity", "err", utils.NilPtrDerefErr)
		return nil
	}

	suite := crypto.GetBlsSuite()

	zmqSocketSet := zmq.CreateSocketSet()
	// todo: set zmq socket set for aggregator

	aggregator := &Aggregator{
		Id:            getAggregatorId(),
		ProducerIds:   make(map[string]struct{}),
		Suite:         suite,
		BroadcastPort: aggregatorEntity.BroadcastPort,
		ZmqSocketSet:  zmqSocketSet,
	}
	setAggregator(aggregator)
	return aggregator
}

func (aggregator *Aggregator) addProducer(producerId string, producerPublicKey kyber.Point) ([]kyber.Point, int, error) {
	if aggregator == nil || aggregator.ProducerIds == nil {
		log.Error("nil aggregator or producer ids", "err", utils.NilPtrDeref)
	}

	updatedThreshold := aggregator.Threshold
	updatedProducerCount := len(aggregator.ProducerIds)
	if _, ok := aggregator.ProducerIds[producerId]; !ok {
		updatedProducerCount++
	}
	if updatedThreshold < updatedProducerCount/2+1 {
		updatedThreshold = updatedProducerCount/2 + 1
	}

	publicKeys := []kyber.Point{producerPublicKey}
	for originProducerId := range aggregator.ProducerIds {
		if originProducerId == producerId {
			continue
		}
		publicKeys = append(publicKeys, getProducer(originProducerId).PublicKey)
	}

	aggregator.ProducerIds[producerId] = struct{}{}
	aggregator.Threshold = updatedThreshold
	return publicKeys, aggregator.Threshold, nil
}

func (aggregator *Aggregator) deleteProducer(producerId string) error {
	if aggregator == nil || aggregator.ProducerIds == nil {
		log.Error("nil aggregator or producer ids", "err", utils.NilPtrDeref)
	}

	if _, ok := aggregator.ProducerIds[producerId]; !ok {
		log.Warn("producer not existed", "producer id", producerId)
		return nil
	}

	updatedThreshold := aggregator.Threshold
	updatedProducerCount := len(aggregator.ProducerIds) - 1
	if updatedThreshold > updatedProducerCount/2+1 {
		updatedThreshold = updatedProducerCount/2 + 1
	}

	delete(aggregator.ProducerIds, producerId)
	aggregator.Threshold = updatedThreshold
	return nil
}

func (aggregator *Aggregator) SelectVerifierId() string {
	if aggregator == nil || aggregator.ProducerIds == nil {
		log.Error("nil aggregator or producer ids", "err", utils.NilPtrDeref)
	}

	timestamp := time.Now().UnixNano()
	rank := int(timestamp) % len(aggregator.ProducerIds)
	var producerId string
	var index int
	for producerId = range aggregator.ProducerIds {
		if index == rank {
			break
		}
		index++
	}
	return producerId
}
