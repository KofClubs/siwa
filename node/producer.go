package node

import (
	"fmt"
	"github.com/KofClubs/siwa/crypto"
	"github.com/KofClubs/siwa/network"
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
	publicKeys, threshold, err := aggregator.AddProducer(id, publicKey)
	if err != nil {
		log.Error("fail to add producer", "private key", producerEntity.PrivateKey, "err", err)
		return nil
	}
	dkg, err := crypto.CreateDistributedKeyGenerator(suite, privateKey, publicKeys, threshold)
	if err != nil {
		log.Error("fail to update distributed key generator of producer when creating producer",
			"private key", producerEntity.PrivateKey, "err", err)
		return nil
	}

	zmqSocketSet := zmq.CreateSocketSet()
	err = zmqSocketSet.SetPubSocket(network.GetPubEndpoint(aggregator.BroadcastPort))
	if err != nil {
		log.Error("fail to set pub socket of producer", "broadcast port", aggregator.BroadcastPort,
			"err", err)
		return nil
	}
	err = zmqSocketSet.SetSubSocket(network.GetSubEndpoint(aggregator.BroadcastPort, rank), network.GetFilter(rank))
	if err != nil {
		log.Error("fail to set sub socket of producer", "broadcast port", aggregator.BroadcastPort,
			"err", err)
		return nil
	}

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
	return producer
}
