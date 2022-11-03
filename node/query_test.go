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
	"context"
	"math/rand"
	"testing"

	"github.com/KofClubs/siwa/crypto"
	"github.com/KofClubs/siwa/node/querier"
	"github.com/MonteCarloClub/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pedersendkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
	"go.dedis.ch/kyber/v3/util/key"
)

const (
	NodeCount    = 6
	RedisAddress = "localhost:6379"
)

var (
	unmarshalledNodes []*UnmarshalledNode
	nodes             []*Node
)

func genRandomPrivateKey() string {
	return key.NewKeyPair(crypto.GetBlsSuite()).Private.String()
}

func fullCommunicate() {
	// tested by TestPedersenDkg
	for _, node := range nodes {
		_ = node.Dkg.CreatePedersenDkgDeals()
	}

	pedersenDkgResponses := make([]*pedersendkg.Response, 0)

	for _, node := range nodes {
		for j, pedersenDkgDeal := range node.Dkg.PedersendkgDeals {
			pedersenDkgResponse, _ := getNodeByDkgIndex(j).Dkg.VerifyPedersenDkgDeal(pedersenDkgDeal)
			pedersenDkgResponses = append(pedersenDkgResponses, pedersenDkgResponse)
		}
	}

	for _, pedersenDkgResponse := range pedersenDkgResponses {
		for _, node := range nodes {
			node.Dkg.VerifyPedersenDkgResponse(pedersenDkgResponse)
		}
	}
}

func initRedis() {
	redisQuerier := &querier.RedisQuerier{}
	redisQuerier.Init(RedisAddress)
	_ = redisQuerier.RedisClient.Set(context.Background(), "k1", "v1", 0).Err()
}

func TestQuery(t *testing.T) {
	// 1. create group
	group := &Group{
		Id:      "0",
		NodeIds: make(map[string]struct{}, 0),
	}
	setGroup(group)

	// 2. generate node entities and create nodes
	for rank := 0; rank < NodeCount; rank++ {
		unmarshalledNodes = append(unmarshalledNodes, &UnmarshalledNode{
			GroupId:       group.Id,
			PrivateKey:    genRandomPrivateKey(),
			QuerierSource: "redis",
			RedisAddress:  RedisAddress,
		})
		log.Info("unmarshalled nodes generated",
			"private key", unmarshalledNodes[len(unmarshalledNodes)-1].PrivateKey)
	}
	for _, unmarshalledNode := range unmarshalledNodes {
		nodes = append(nodes, unmarshalledNode.CreateNode())
		require.NotNil(t, nodes[len(nodes)-1])
		log.Info("node created", "id", nodes[len(nodes)-1].Id)
	}

	// 3. fully communicate to certify all dkgs
	fullCommunicate()
	for _, node := range nodes {
		assert.True(t, node.ReadyToQuery())
	}

	// 4. query and aggregate result
	initRedis()
	verifier := nodes[rand.Int()%len(nodes)]
	log.Info("verifier selected", "verifier id", verifier.Id)
	expression := "k1"
	expectedValue := "v1"
	signatures := make([][]byte, 0)
	for _, node := range nodes {
		message, signature := node.Query(expression)
		assert.Equal(t, expectedValue, message)
		ok := verifier.Verify(message, signature)
		assert.True(t, ok)
		signatures = append(signatures, signature)
	}
	signature, ok := verifier.Recover(expectedValue, signatures)
	assert.NotNil(t, signature)
	assert.True(t, ok)
}
