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
	Producers     []*Producer
	Suite         *bn256.Suite
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
	err = zmqSocketSet.SetSubSocket(network.GetAggregatorSubEndpoint(aggregatorEntity.BroadcastPort),
		network.AggregatorFilter)
	if err != nil {
		log.Error("fail to set sub socket of aggregator",
			"broadcast port", aggregatorEntity.BroadcastPort, "err", err)
		return nil
	}

	aggregator := &Aggregator{
		Id:            getAggregatorId(),
		Producers:     make([]*Producer, 0),
		Suite:         suite,
		PublicKey:     publicKey,
		Threshold:     1,
		Dkg:           dkg,
		BroadcastPort: aggregatorEntity.BroadcastPort,
		ZmqSocketSet:  zmqSocketSet,
	}
	setAggregator(aggregator)
	return aggregator
}
