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
	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	pedersendkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
)

func (dkg *DistributedKeyGenerator) GetPedersenDkgDeals() error {
	dkg.Mutex.Lock()
	defer dkg.Mutex.Unlock()

	if dkg.PedersenDkg == nil {
		log.Error("nil pointer dereference", "dkg.PedersenDkg", dkg.PedersenDkg, "err", utils.NilPtrDerefErr)
		return utils.NilPtrDerefErr
	}

	pedersenDkgDeals, err := dkg.PedersenDkg.Deals()
	if err != nil {
		log.Error("fail to create ped-dkg deals", "err", err)
		return err
	}

	dkg.pedersendkgDeals = pedersenDkgDeals
	return nil
}

func (dkg *DistributedKeyGenerator) VerifyPedersenDkgDeal(pedersenDkgDeal *pedersendkg.Deal) (*pedersendkg.Response, bool) {
	dkg.Mutex.Lock()
	defer dkg.Mutex.Unlock()

	if dkg.PedersenDkg == nil || pedersenDkgDeal == nil {
		log.Error("nil pointer dereference", "dkg.PedersenDkg", dkg.PedersenDkg, "pedersenDkgDeal", pedersenDkgDeal)
		return nil, false
	}

	response, err := dkg.PedersenDkg.ProcessDeal(pedersenDkgDeal)
	if response == nil || err != nil {
		return nil, false
	}

	return response, true
}

func (dkg *DistributedKeyGenerator) VerifyPedersenDkgResponse(pedersenDkgResponse *pedersendkg.Response) bool {
	dkg.Mutex.Lock()
	defer dkg.Mutex.Unlock()

	if dkg.PedersenDkg == nil || pedersenDkgResponse == nil {
		log.Error("nil pointer dereference", "dkg.PedersenDkg", dkg.PedersenDkg, "pedersenDkgResponse", pedersenDkgResponse)
		return false
	}

	_, err := dkg.PedersenDkg.ProcessResponse(pedersenDkgResponse)
	// TODO 是否处理 justification？
	return err == nil
}
