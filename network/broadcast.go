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

package network

import (
	"fmt"
)

const (
	FilterFormat = "to_node_%v: "
)

func GetPubEndpoint(broadcastPort string) string {
	return fmt.Sprintf("tcp://*:%v", broadcastPort)
}

func GetSubEndpoint(broadcastPort string, rank uint64) string {
	// aggregator rank == 0, producer rank > 0
	// sub endpoint: tcp://127.#{b}.#{c}.#{d}:#{broadcastPort}
	b := rank >> 16
	if b > 255 {
		return ""
	}
	c := (rank - (b << 16)) >> 8
	d := rank - (b << 16) - (c << 8)
	return fmt.Sprintf("tcp://127.%v.%v.%v:%v", b, c, d, broadcastPort)
}

func GetFilter(rank uint64) string {
	return fmt.Sprintf(FilterFormat, rank)
}
