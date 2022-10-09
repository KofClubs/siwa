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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"
	pedersendkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
	pedersenvss "go.dedis.ch/kyber/v3/share/vss/pedersen"
	"go.dedis.ch/kyber/v3/util/key"
)

const (
	dkgCount = 3
)

func TestPedersenDkg(t *testing.T) {
	threshold := pedersenvss.MinimumT(dkgCount)

	blsSuite := GetBlsSuite()

	privateKeys, publicKeys := make([]kyber.Scalar, dkgCount), make([]kyber.Point, dkgCount)
	for i := 0; i < dkgCount; i++ {
		pair := key.NewKeyPair(blsSuite)
		privateKeys[i], publicKeys[i] = pair.Private, pair.Public
	}

	dkgs := make([]*DistributedKeyGenerator, dkgCount)
	for i := 0; i < dkgCount; i++ {
		dkg, err := CreateDistributedKeyGenerator(blsSuite, privateKeys[i], publicKeys, threshold)
		require.NotNil(t, dkg)
		require.NotNil(t, dkg.PedersenDkg)
		require.Nil(t, err)
		dkg.SetIndex(i)
		dkgs[i] = dkg
	}

	for i := 0; i < dkgCount; i++ {
		err := dkgs[i].CreatePedersenDkgDeals()
		require.Nil(t, err)
		assert.Equal(t, dkgCount-1, len(dkgs[i].pedersendkgDeals))
	}

	pedersenDkgResponsesSlice := make([]map[int]*pedersendkg.Response, dkgCount)
	for i := 0; i < dkgCount; i++ {
		pedersenDkgResponsesSlice[i] = make(map[int]*pedersendkg.Response)
		for j, pedersenDkgDeal := range dkgs[i].pedersendkgDeals {
			pedersenDkgResponse, ok := dkgs[j].VerifyPedersenDkgDeal(pedersenDkgDeal)
			assert.NotNil(t, pedersenDkgResponse)
			assert.True(t, ok)
			pedersenDkgResponsesSlice[i][j] = pedersenDkgResponse
		}
	}

	for _, pedersenDkgResponses := range pedersenDkgResponsesSlice {
		for _, pedersenDkgResponse := range pedersenDkgResponses {
			for _, dkg := range dkgs {
				ok := dkg.VerifyPedersenDkgResponse(pedersenDkgResponse)
				assert.True(t, ok)
			}
		}
	}

	for _, dkg := range dkgs {
		ok := dkg.PedersenDkg.Certified()
		assert.True(t, ok)
	}
}
