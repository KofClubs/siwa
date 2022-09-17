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

package crypto

import (
	"sync"

	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	pedersendkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
)

type DistributedKeyGenerator struct {
	Mutex sync.Mutex

	PedersenDkg      *pedersendkg.DistKeyGenerator
	pedersendkgDeals map[int]*pedersendkg.Deal
}

func CreateDistributedKeyGenerator(suite *bn256.Suite, privateKey kyber.Scalar, publicKeys []kyber.Point, threshold int) (*DistributedKeyGenerator, error) {
	if suite == nil {
		log.Error("nil pointer dereference", "suite", suite, "err", utils.NilPtrDerefErr)
		return nil, utils.NilPtrDerefErr
	}

	pedersenDkg, err := pedersendkg.NewDistKeyGenerator(suite, privateKey, publicKeys, threshold)
	if err != nil {
		log.Error("fail to create pedersen distributed key generator", "err", err)
		return nil, err
	}

	return &DistributedKeyGenerator{
		PedersenDkg: pedersenDkg,
	}, nil
}
