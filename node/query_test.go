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
	"testing"

	"github.com/KofClubs/siwa/crypto"
	"github.com/MonteCarloClub/log"
	"go.dedis.ch/kyber/v3/util/key"
)

const (
	BroadcastPort = "8033"
	NodeCount     = 6
)

var (
	aggregatorEntity *AggregatorEntity
	producerEntities []*ProducerEntity

	aggregator *Aggregator
	producers  []*Producer
)

func genRandomPrivateKey() string {
	return key.NewKeyPair(crypto.GetBlsSuite()).Private.String()
}

func TestQuery(t *testing.T) {
	// 1. generate aggregator entity and create aggregator
	aggregatorEntity = &AggregatorEntity{
		BroadcastPort: BroadcastPort,
	}
	log.Info("aggregator entity generated")
	aggregator = aggregatorEntity.CreateAggregator()
	log.Info("aggregator created", "id", aggregator.Id)

	// 2. generate producer entities and create producers
	for rank := 1; rank < NodeCount; rank++ {
		producerEntities = append(producerEntities, &ProducerEntity{
			AggregatorId:  "0",
			PrivateKey:    genRandomPrivateKey(),
			QuerierSource: "redis",
			RedisAddress:  "localhost:6379",
		})
		log.Info("producer entity generated",
			"private key", producerEntities[len(producerEntities)-1].PrivateKey)
	}
	for _, producerEntity := range producerEntities {
		producers = append(producers, producerEntity.CreateProducer())
		log.Info("producer created", "id", producers[len(producers)-1].Id)
	}
}
