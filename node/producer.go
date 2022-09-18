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
)

type ProducerEntity struct {
	AggregatorId  string `yaml:"aggregator_id"`
	PrivateKey    string `yaml:"private_key"`
	QuerierSource string `yaml:"querier_source"`
	RedisAddress  string `yaml:"redis_address"`
}

type Producer struct {
	Id           string
	Aggregator   *Aggregator
	Rank         uint64
	PublicKey    kyber.Point
	Dkg          *crypto.DistributedKeyGenerator
	ZmqSocketSet *zmq.SocketSet
	Querier      *querier.Querier
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
	// todo: enhanced consistency when appending this public key
	publicKeys := []kyber.Point{aggregator.PublicKey, publicKey}
	for _, peerProducer := range aggregator.Producers {
		publicKeys = append(publicKeys, peerProducer.PublicKey)
	}
	threshold := aggregator.Threshold
	// todo: update dkg threshold
	if aggregator.Threshold < len(aggregator.Producers)/2+1 {
		threshold = len(aggregator.Producers)/2 + 1
	}
	dkg, err := crypto.CreateDistributedKeyGenerator(suite, privateKey, publicKeys, threshold)
	if err != nil {
		log.Error("fail to create distributed key generator of producer",
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
	err = zmqSocketSet.SetSubSocket(network.GetAggregatorSubEndpoint(aggregator.BroadcastPort),
		network.GetProducerFilter(rank))
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

	// todo: defer: init a new dkg generation
	return &Producer{
		Id:           id,
		Aggregator:   aggregator,
		Rank:         rank,
		PublicKey:    publicKey,
		Dkg:          dkg,
		ZmqSocketSet: zmqSocketSet,
		Querier:      &querierOfProducer,
	}
}
