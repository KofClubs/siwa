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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"
	pedersendkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
	pedersenvss "go.dedis.ch/kyber/v3/share/vss/pedersen"
	"go.dedis.ch/kyber/v3/util/key"
)

const (
	DkgCount = 3

	VerifiableMessage   = "ok"
	UnverifiableMessage = "not ok"
)

func TestPedersenDkg(t *testing.T) {
	threshold := pedersenvss.MinimumT(DkgCount)

	blsSuite := GetBlsSuite()

	privateKeys, publicKeys := make([]kyber.Scalar, DkgCount), make([]kyber.Point, DkgCount)
	for i := 0; i < DkgCount; i++ {
		pair := key.NewKeyPair(blsSuite)
		privateKeys[i], publicKeys[i] = pair.Private, pair.Public
	}

	dkgs := make([]*DistributedKeyGenerator, DkgCount)
	for i := 0; i < DkgCount; i++ {
		dkg, err := CreateDistributedKeyGenerator(blsSuite, privateKeys[i], publicKeys, threshold)
		require.NotNil(t, dkg)
		require.NotNil(t, dkg.PedersenDkg)
		require.Nil(t, err)
		dkg.SetIndex(i)
		dkgs[i] = dkg
	}

	for i := 0; i < DkgCount; i++ {
		err := dkgs[i].CreatePedersenDkgDeals()
		require.Nil(t, err)
		assert.Equal(t, DkgCount-1, len(dkgs[i].PedersendkgDeals))
	}

	pedersenDkgResponsesSlice := make([]map[int]*pedersendkg.Response, DkgCount)
	for i := 0; i < DkgCount; i++ {
		pedersenDkgResponsesSlice[i] = make(map[int]*pedersendkg.Response)
		for j, pedersenDkgDeal := range dkgs[i].PedersendkgDeals {
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

	for i, signer := range dkgs {
		message := fmt.Sprintf("msg_%v", i)
		signature := Sign(blsSuite, signer, message)
		for _, verifier := range dkgs {
			ok := Verify(blsSuite, verifier, message, signature)
			assert.True(t, ok)
			ok = Verify(blsSuite, verifier, "msg_", signature)
			assert.False(t, ok)
		}
	}

	var signature []byte
	var ok bool
	signatures := make([][]byte, 0)
	for i, dkg := range dkgs {
		if i < threshold {
			signatures = append(signatures, Sign(blsSuite, dkg, VerifiableMessage))
		} else {
			signatures = append(signatures, Sign(blsSuite, dkg, UnverifiableMessage))
		}
	}
	for _, dkg := range dkgs {
		signature, ok = Recover(blsSuite, dkg, threshold, DkgCount, VerifiableMessage, signatures)
		assert.NotNil(t, signature)
		assert.True(t, ok)
		_, ok = Recover(blsSuite, dkg, threshold, DkgCount, UnverifiableMessage, signatures)
		assert.False(t, ok)
	}

	var expectedDistributedPublicKey kyber.Point
	for _, dkg := range dkgs {
		actualDistributedPublicKey, err := dkg.GetDistributedPublicKey()
		assert.Nil(t, err)
		if expectedDistributedPublicKey == nil {
			expectedDistributedPublicKey = actualDistributedPublicKey
		} else {
			require.NotNil(t, actualDistributedPublicKey)
			assert.True(t, actualDistributedPublicKey.Equal(expectedDistributedPublicKey))
		}
	}

	// todo: verify(expectedDistributedPublicKey, VerifiableMessage, signature)
}
