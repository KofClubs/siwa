package node

import (
	"github.com/KofClubs/siwa/crypto"
	"github.com/KofClubs/siwa/network"
	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	"github.com/MonteCarloClub/zmq"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

type AggregatorEntity struct {
	PrivateKey    string `yaml:"private_key"`
	BroadcastPort string `yaml:"broadcast_port"`
}

type Aggregator struct {
	Id            string
	ProducerIds   map[string]struct{}
	Suite         *bn256.Suite
	privateKey    kyber.Scalar
	PublicKey     kyber.Point
	Threshold     int
	Dkg           *crypto.DistributedKeyGenerator
	BroadcastPort string
	ZmqSocketSet  *zmq.SocketSet
}

func (aggregatorEntity *AggregatorEntity) CreateAggregator() *Aggregator {
	if aggregatorEntity == nil {
		log.Error("nil aggregator entity", "err", utils.NilPtrDerefErr)
		return nil
	}

	suite := crypto.GetBlsSuite()
	privateKey, err := crypto.GetBlsPrivateKey(suite, aggregatorEntity.PrivateKey)
	if err != nil {
		log.Error("fail to get private key of aggregator", "private key", aggregatorEntity.PrivateKey,
			"err", err)
		return nil
	}
	publicKey, err := crypto.GetBlsPublicKey(suite, privateKey)
	if err != nil {
		log.Error("fail to get public key of aggregator", "private key", aggregatorEntity.PrivateKey,
			"err", err)
		return nil
	}
	dkg, err := crypto.CreateDistributedKeyGenerator(suite, privateKey, []kyber.Point{publicKey}, 1)
	if err != nil {
		log.Error("fail to create distributed key generator of aggregator",
			"private key", aggregatorEntity.PrivateKey, "err", err)
		return nil
	}

	zmqSocketSet := zmq.CreateSocketSet()
	err = zmqSocketSet.SetPubSocket(network.GetPubEndpoint(aggregatorEntity.BroadcastPort))
	if err != nil {
		log.Error("fail to set pub socket of aggregator",
			"broadcast port", aggregatorEntity.BroadcastPort, "err", err)
		return nil
	}
	err = zmqSocketSet.SetSubSocket(network.GetSubEndpoint(aggregatorEntity.BroadcastPort, 0),
		network.GetFilter(0))
	if err != nil {
		log.Error("fail to set sub socket of aggregator",
			"broadcast port", aggregatorEntity.BroadcastPort, "err", err)
		return nil
	}

	aggregator := &Aggregator{
		Id:            getAggregatorId(),
		ProducerIds:   make(map[string]struct{}),
		Suite:         suite,
		privateKey:    privateKey,
		PublicKey:     publicKey,
		Threshold:     1,
		Dkg:           dkg,
		BroadcastPort: aggregatorEntity.BroadcastPort,
		ZmqSocketSet:  zmqSocketSet,
	}
	setAggregator(aggregator)
	return aggregator
}

func (aggregator *Aggregator) AddProducer(producerId string, producerPublicKey kyber.Point) ([]kyber.Point, int, error) {
	if aggregator == nil || aggregator.ProducerIds == nil {
		log.Error("nil aggregator or producer ids", "err", utils.NilPtrDeref)
	}

	updatedThreshold := aggregator.Threshold
	updatedProducerCount := len(aggregator.ProducerIds)
	if _, ok := aggregator.ProducerIds[producerId]; !ok {
		updatedProducerCount++
	}
	if aggregator.Threshold < updatedProducerCount/2+1 {
		updatedThreshold = updatedProducerCount/2 + 1
	}

	publicKeys := []kyber.Point{aggregator.PublicKey, producerPublicKey}
	for originProducerId := range aggregator.ProducerIds {
		if originProducerId == producerId {
			continue
		}
		publicKeys = append(publicKeys, getProducer(originProducerId).PublicKey)
	}
	dkg, err := crypto.CreateDistributedKeyGenerator(aggregator.Suite, aggregator.privateKey, publicKeys,
		updatedThreshold)
	if err != nil {
		log.Error("fail to update distributed key generator of aggregator when adding producer",
			"producer id", producerId, "err", err)
		return nil, 0, err
	}

	aggregator.ProducerIds[producerId] = struct{}{}
	aggregator.Threshold = updatedThreshold
	aggregator.Dkg = dkg
	return publicKeys, aggregator.Threshold, nil
}

func (aggregator *Aggregator) DeleteProducer(producerId string) error {
	if aggregator == nil || aggregator.ProducerIds == nil {
		log.Error("nil aggregator or producer ids", "err", utils.NilPtrDeref)
	}

	if _, ok := aggregator.ProducerIds[producerId]; !ok {
		log.Warn("producer not existed", "producer id", producerId)
		return nil
	}

	updatedThreshold := aggregator.Threshold
	updatedProducerCount := len(aggregator.ProducerIds) - 1
	if aggregator.Threshold > updatedProducerCount/2+1 {
		updatedThreshold = updatedProducerCount/2 + 1
	}

	publicKeys := []kyber.Point{aggregator.PublicKey}
	for originProducerId := range aggregator.ProducerIds {
		if originProducerId == producerId {
			continue
		}
		publicKeys = append(publicKeys, getProducer(originProducerId).PublicKey)
	}
	dkg, err := crypto.CreateDistributedKeyGenerator(aggregator.Suite, aggregator.privateKey, publicKeys,
		updatedThreshold)
	if err != nil {
		log.Error("fail to update distributed key generator of aggregator when deleting producer",
			"producer id", producerId, "err", err)
		return err
	}

	delete(aggregator.ProducerIds, producerId)
	aggregator.Threshold = updatedThreshold
	aggregator.Dkg = dkg
	return nil
}
