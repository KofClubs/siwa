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
	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
)

// Group is not an entity, but information shared by a group of nodes
// todo: implement Group at on-chain registry
type Group struct {
	Id        string
	NodeIds   map[string]struct{}
	Threshold int
}

func (group *Group) addNode(newNodeId string) ([]string, int, error) {
	if group == nil || group.NodeIds == nil {
		log.Error("nil group or node ids", "err", utils.NilPtrDeref)
		return nil, 0, utils.NilPtrDerefErr
	}

	updatedThreshold := group.Threshold
	updatedNodeCount := len(group.NodeIds)
	if _, ok := group.NodeIds[newNodeId]; !ok {
		group.NodeIds[newNodeId] = struct{}{}
		updatedNodeCount++
	}
	if updatedThreshold < updatedNodeCount/2+1 {
		updatedThreshold = updatedNodeCount/2 + 1
	}

	nodeIds := make([]string, 0)
	for nodeId := range group.NodeIds {
		nodeIds = append(nodeIds, nodeId)
	}
	group.Threshold = updatedThreshold
	return nodeIds, group.Threshold, nil
}

func (group *Group) deleteNode(nodeId string) {
	if group == nil || group.NodeIds == nil {
		log.Warn("nil group or node ids")
	}

	if _, ok := group.NodeIds[nodeId]; !ok {
		log.Warn("node not existed", "node id", nodeId)
	}

	updatedThreshold := group.Threshold
	updatedNodeCount := len(group.NodeIds) - 1
	if updatedThreshold > updatedNodeCount/2+1 {
		updatedThreshold = updatedNodeCount/2 + 1
	}

	delete(group.NodeIds, nodeId)
	group.Threshold = updatedThreshold
}
