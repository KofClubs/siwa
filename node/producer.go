package node

import (
	"github.com/KofClubs/siwa/crypto"
	"github.com/KofClubs/siwa/node/querier"
	"github.com/MonteCarloClub/zmq"
)

type Producer struct {
	Id           string
	Aggregator   *Aggregator
	Dkg          *crypto.DistributedKeyGenerator
	ZmqSocketSet *zmq.SocketSet
	Querier      *querier.Querier
}
