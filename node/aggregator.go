package node

import (
	"github.com/KofClubs/siwa/crypto"
	"github.com/MonteCarloClub/zmq"
)

type Aggregator struct {
	Id           string
	Producers    []*Producer
	Dkg          *crypto.DistributedKeyGenerator
	ZmqSocketSet *zmq.SocketSet
}
