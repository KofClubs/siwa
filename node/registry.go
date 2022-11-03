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
)

var (
	groupCounter       *int
	nodeCounterByGroup map[string]int
	groupTable         map[string]*Group
	nodeTable          map[string]*Node
	dkgIndexTable      map[int]*Node
)

func init() {
	if groupCounter == nil {
		var groupCounterValue int
		groupCounter = &groupCounterValue
	}
	if nodeCounterByGroup == nil {
		nodeCounterByGroup = make(map[string]int)
	}
	if groupTable == nil {
		groupTable = make(map[string]*Group)
	}
	if nodeTable == nil {
		nodeTable = make(map[string]*Node)
	}
	if dkgIndexTable == nil {
		dkgIndexTable = make(map[int]*Node)
	}
}

func generateNodeId(groupId string) string {
	if _, ok := nodeCounterByGroup[groupId]; !ok {
		nodeCounterByGroup[groupId] = 0
	}
	count := nodeCounterByGroup[groupId]
	nodeCounterByGroup[groupId]++
	return fmt.Sprintf("%v.%v", count, groupId)
}

func getGroup(groupId string) *Group {
	if groupTable == nil {
		return nil
	}
	if group, ok := groupTable[groupId]; ok {
		return group
	}
	return nil
}

func getNode(nodeId string) *Node {
	if nodeTable == nil {
		return nil
	}
	if node, ok := nodeTable[nodeId]; ok {
		return node
	}
	return nil
}

func getNodeByDkgIndex(index int) *Node {
	if dkgIndexTable == nil {
		return nil
	}
	if node, ok := dkgIndexTable[index]; ok {
		return node
	}
	return nil
}

func setGroup(group *Group) {
	if groupTable == nil || group == nil {
		return
	}
	groupTable[group.Id] = group
}

func setNode(node *Node) {
	if nodeTable == nil || node == nil {
		return
	}
	nodeTable[node.Id] = node
	dkgIndexTable[node.Dkg.GetIndex()] = node
}
