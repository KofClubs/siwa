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
	"fmt"
	"strconv"

	"go.dedis.ch/kyber/v3"
)

var (
	aggregatorCounter *int
	producerCounter   map[string]int
	aggregatorTable   map[string]*Aggregator
	producerTable     map[string]*Producer
	dkgIndexTable     map[int]*Producer
	publicKeyTable    map[kyber.Point]*Producer
)

func init() {
	if aggregatorCounter == nil {
		var aggregatorCounterValue int
		aggregatorCounter = &aggregatorCounterValue
	}
	if producerCounter == nil {
		producerCounter = make(map[string]int)
	}
	if aggregatorTable == nil {
		aggregatorTable = make(map[string]*Aggregator)
	}
	if producerTable == nil {
		producerTable = make(map[string]*Producer)
	}
	if dkgIndexTable == nil {
		dkgIndexTable = make(map[int]*Producer)
	}
	if publicKeyTable == nil {
		publicKeyTable = make(map[kyber.Point]*Producer)
	}
}

func getAggregatorId() string {
	if aggregatorCounter == nil {
		return ""
	}
	aggregatorId := strconv.Itoa(*aggregatorCounter)
	*aggregatorCounter++
	return aggregatorId
}

func getProducerId(aggregatorId string) (string, int) {
	if _, ok := producerCounter[aggregatorId]; !ok {
		producerCounter[aggregatorId] = 0
	}
	producerId := producerCounter[aggregatorId]
	producerCounter[aggregatorId]++
	return fmt.Sprintf("%v.%v", producerId, aggregatorId), producerId
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

func getProducerByDkgIndex(index int) *Producer {
	if dkgIndexTable == nil {
		return nil
	}
	if producer, ok := dkgIndexTable[index]; ok {
		return producer
	}
	return nil
}

func getProducerByPublicKey(publicKey kyber.Point) *Producer {
	if publicKeyTable == nil {
		return nil
	}
	if producer, ok := publicKeyTable[publicKey]; ok {
		return producer
	}
	return nil
}

func setProducer(producer *Producer) {
	if producerTable == nil {
		return
	}
	producerTable[producer.Id] = producer
	dkgIndexTable[producer.Dkg.GetIndex()] = producer
	publicKeyTable[producer.PublicKey] = producer
}
