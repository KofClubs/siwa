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
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

type UnmarshalledNode struct {
	GroupId       string `yaml:"group_id"`
	PrivateKey    string `yaml:"private_key"`
	QuerierSource string `yaml:"querier_source"`
	RedisAddress  string `yaml:"redis_address"`
}

type Node struct {
	Id, GroupId string
	Rank        int
	Suite       *bn256.Suite
	privateKey  kyber.Scalar
	PublicKey   kyber.Point
	Dkg         *crypto.DistributedKeyGenerator
	Querier     querier.Querier
}

func (unmarshalledNode *UnmarshalledNode) CreateNode() *Node {
	if unmarshalledNode == nil {
		log.Error("nil unmarshalled node", "err", utils.NilPtrDerefErr)
		return nil
	}

	groupId := unmarshalledNode.GroupId
	if groupId == "" {
		log.Info("group not specified, select one for this node",
			"private key", unmarshalledNode.PrivateKey)
		// todo: call the scheduling algorithm to assign it to an group
		log.Info("group selected for this node", "private key", unmarshalledNode.PrivateKey,
			"group id", groupId)
	}
	id := generateNodeId(groupId)
	group := getGroup(groupId)
	if group == nil {
		log.Error("nil group", "err", utils.NilPtrDeref)
		return nil
	}

	suite := crypto.GetBlsSuite()
	privateKey, err := crypto.GetBlsPrivateKey(suite, unmarshalledNode.PrivateKey)
	if err != nil {
		log.Error("fail to get private key of node", "private key", unmarshalledNode.PrivateKey,
			"err", err)
		return nil
	}
	publicKey, err := crypto.GetBlsPublicKey(suite, privateKey)
	if err != nil {
		log.Error("fail to get public key of node", "private key", unmarshalledNode.PrivateKey,
			"err", err)
		return nil
	}

	var querierOfNode querier.Querier
	switch unmarshalledNode.QuerierSource {
	case "redis":
		redisQuerier := &querier.RedisQuerier{}
		redisQuerier.Init(unmarshalledNode.RedisAddress)
		querierOfNode = querier.Querier(redisQuerier)
	default:
		log.Error("fail to init querier of node", "err", fmt.Errorf("illegal querier_source"))
		return nil
	}

	nodeIds, threshold, err := group.addNode(id)
	if err != nil {
		log.Error("fail to add node", "private key", unmarshalledNode.PrivateKey, "err", err)
		return nil
	}
	// nodeIds[index] == id, for 0<=i<len(nodes): nodes[i].PublicKey == publicKeys[i]
	var index int
	nodes := make([]*Node, 0)
	publicKeys := make([]kyber.Point, 0)
	for i, nodeId := range nodeIds {
		if nodeId == id {
			index = i
			nodes = append(nodes, &Node{})
			publicKeys = append(publicKeys, publicKey)
			continue
		}
		peerNode := getNode(nodeId)
		nodes = append(nodes, peerNode)
		publicKeys = append(publicKeys, peerNode.PublicKey)
	}
	var dkg *crypto.DistributedKeyGenerator
	if threshold < 2 {
		log.Warn("distributed key generators not updated, threshold should not be less than 2",
			"private key", unmarshalledNode.PrivateKey)
	} else {
		dkg, err = crypto.CreateDistributedKeyGenerator(suite, privateKey, publicKeys, threshold)
		if err != nil {
			log.Error("fail to update distributed key generator of node when creating node",
				"private key", unmarshalledNode.PrivateKey, "err", err)
			group.deleteNode(id)
			log.Info("creating a node rolled back", "private key", unmarshalledNode.PrivateKey)
			return nil
		}
		dkg.SetIndex(index)
		// assert: len(publicKeys) > 1
		for i := range publicKeys {
			if i == index {
				continue
			}
			peerNode := nodes[i]
			err = peerNode.updateDkg(publicKeys, threshold, i)
			if err != nil {
				log.Error("fail to update distributed key generator of peer node when creating node",
					"err", err)
				return nil
			}
		}
	}

	node := &Node{
		Id:         id,
		GroupId:    groupId,
		Suite:      suite,
		privateKey: privateKey,
		PublicKey:  publicKey,
		Dkg:        dkg,
		Querier:    querierOfNode,
	}
	setNode(node)
	return node
}

func (node *Node) updateDkg(publicKeys []kyber.Point, threshold int, index int) error {
	if node == nil {
		log.Error("nil node", "err", utils.NilPtrDeref)
		return utils.NilPtrDerefErr
	}
	updatedDkg, err := crypto.CreateDistributedKeyGenerator(node.Suite, node.privateKey, publicKeys, threshold)
	if err != nil {
		log.Error("fail to update distributed key generator of node", "err", err)
		return err
	}
	if updatedDkg == nil {
		log.Error("nil distributed key generator when updated", "err", utils.NilPtrDerefErr)
		return utils.NilPtrDerefErr
	}
	node.Dkg = updatedDkg
	node.Dkg.SetIndex(index)
	setNode(node)
	return nil
}

func (node *Node) ReadyToQuery() bool {
	if node == nil || node.Dkg == nil || node.Dkg.PedersenDkg == nil {
		log.Error("nil node or dkg")
		return false
	}
	return node.Dkg.PedersenDkg.Certified()
}

func (node *Node) Query(expression string) (string, []byte) {
	if node == nil || node.Querier == nil {
		log.Error("nil node or querier")
		return "", nil
	}

	message := node.Querier.Do(expression)
	signature := crypto.Sign(node.Suite, node.Dkg, message)
	return message, signature
}

func (node *Node) Verify(message string, signature []byte) bool {
	if node == nil {
		log.Error("nil node")
		return false
	}

	return crypto.Verify(node.Suite, node.Dkg, message, signature)
}

func (node *Node) Recover(message string, signatures [][]byte) ([]byte, bool) {
	if node == nil {
		log.Error("nil node")
		return nil, false
	}

	group := getGroup(node.GroupId)
	if group == nil {
		log.Error("fail to get group", "node id", node.Id, "group id", node.GroupId)
	}

	return crypto.Recover(node.Suite, node.Dkg, group.Threshold, len(group.NodeIds),
		message, signatures)
}
