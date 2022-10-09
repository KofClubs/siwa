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

	"github.com/KofClubs/siwa/crypto"
	"github.com/KofClubs/siwa/node/querier"
	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	"github.com/MonteCarloClub/zmq"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

type ProducerEntity struct {
	AggregatorId  string `yaml:"aggregator_id"`
	PrivateKey    string `yaml:"private_key"`
	QuerierSource string `yaml:"querier_source"`
	RedisAddress  string `yaml:"redis_address"`
}

type Producer struct {
	Id, AggregatorId string
	Rank             uint64
	Suite            *bn256.Suite
	privateKey       kyber.Scalar
	PublicKey        kyber.Point
	Dkg              *crypto.DistributedKeyGenerator
	ZmqSocketSet     *zmq.SocketSet
	Querier          *querier.Querier
}

func (producerEntity *ProducerEntity) CreateProducer() *Producer {
	if producerEntity == nil {
		log.Error("nil producer entity", "err", utils.NilPtrDerefErr)
		return nil
	}

	aggregatorId := producerEntity.AggregatorId
	if aggregatorId == "" {
		log.Info("aggregator not specified, select one for this producer",
			"private key", producerEntity.PrivateKey)
		// todo: call the load balancing algorithm to assign it to an aggregator
		log.Info("aggregator selected for this producer", "private key", producerEntity.PrivateKey,
			"aggregator id", aggregatorId)
	}
	id, rank := getProducerId(aggregatorId)
	aggregator := getAggregator(aggregatorId)
	if aggregator == nil {
		log.Error("nil aggregator", "err", utils.NilPtrDeref)
		return nil
	}

	suite := aggregator.Suite
	privateKey, err := crypto.GetBlsPrivateKey(suite, producerEntity.PrivateKey)
	if err != nil {
		log.Error("fail to get private key of producer", "private key", producerEntity.PrivateKey,
			"err", err)
		return nil
	}
	publicKey, err := crypto.GetBlsPublicKey(suite, privateKey)
	if err != nil {
		log.Error("fail to get public key of producer", "private key", producerEntity.PrivateKey,
			"err", err)
		return nil
	}

	zmqSocketSet := zmq.CreateSocketSet()
	// todo: set this zmq socket set

	var querierOfProducer querier.Querier
	switch producerEntity.QuerierSource {
	case "redis":
		redisQuerier := &querier.RedisQuerier{}
		redisQuerier.Init(producerEntity.RedisAddress)
		querierOfProducer = querier.Querier(redisQuerier)
	default:
		log.Error("fail to init querier of producer", "err", fmt.Errorf("illegal querier_source"))
		return nil
	}

	// Finally add the producer, which needs to be rolled back if its subsequent operations fail.
	publicKeys, threshold, err := aggregator.addProducer(id, publicKey)
	if err != nil {
		log.Error("fail to add producer", "private key", producerEntity.PrivateKey, "err", err)
		return nil
	}
	var dkg *crypto.DistributedKeyGenerator
	var producersToUpdate map[*Producer]*crypto.DistributedKeyGenerator
	if threshold < 2 {
		log.Warn("distributed key generators not updated, threshold should not be less than 2",
			"private key", producerEntity.PrivateKey)
	} else {
		dkg, err = crypto.CreateDistributedKeyGenerator(suite, privateKey, publicKeys, threshold)
		if err != nil {
			log.Error("fail to update distributed key generator of producer when creating producer",
				"private key", producerEntity.PrivateKey, "err", err)
			err = aggregator.deleteProducer(id)
			if err == nil {
				log.Info("adding a producer rolled back", "private key", producerEntity.PrivateKey)
			} else {
				log.Warn("fail to roll back adding a producer", "private key", producerEntity.PrivateKey,
					"err", err)
			}
			return nil
		}
		// because publicKeys[0] is publicKey, index of dkg is 0
		dkg.SetIndex(0)
		// assert: len(publicKeys) > 1
		publicKeysOfPeerProducer := publicKeys[1:]
		producersToUpdate = make(map[*Producer]*crypto.DistributedKeyGenerator)
		for index, publicKeyOfPeerProducer := range publicKeysOfPeerProducer {
			peerProducer := getProducerByPublicKey(publicKeyOfPeerProducer)
			dkgOfPeerProducer, err := crypto.CreateDistributedKeyGenerator(suite, peerProducer.privateKey, publicKeys, threshold)
			if err != nil {
				log.Error("fail to update distributed key generator of peer producer when creating producer",
					"peer producer id", peerProducer.Id, "err", err)
				return nil
			}
			dkgOfPeerProducer.SetIndex(index)
			producersToUpdate[peerProducer] = dkgOfPeerProducer
		}
	}

	producer := &Producer{
		Id:           id,
		AggregatorId: aggregatorId,
		Rank:         rank,
		Suite:        suite,
		privateKey:   privateKey,
		PublicKey:    publicKey,
		Dkg:          dkg,
		ZmqSocketSet: zmqSocketSet,
		Querier:      &querierOfProducer,
	}
	setProducer(producer)
	for peerProducer, dkgOfPeerProducer := range producersToUpdate {
		peerProducer.Dkg = dkgOfPeerProducer
		setProducer(peerProducer)
	}
	return producer
}
