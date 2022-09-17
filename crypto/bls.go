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
	"encoding/hex"

	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

func GetBlsSuite() *bn256.Suite {
	return bn256.NewSuiteG2()
}

func GetBlsPrivateKey(suite *bn256.Suite, privateKeyString string) (kyber.Scalar, error) {
	if suite == nil {
		log.Error("nil pointer dereference", "suite", suite, "err", utils.NilPtrDerefErr)
		return nil, utils.NilPtrDerefErr
	}

	privateKeyBytes, err := hex.DecodeString(privateKeyString)
	if err != nil {
		log.Error("fail to decode bls private key", "privateKeyString", privateKeyString, "err", err)
		return nil, err
	}

	scalar := suite.Scalar()
	err = scalar.UnmarshalBinary(privateKeyBytes)
	if err != nil {
		log.Error("fail to unmarshal bls private key", "privateKeyBytes", privateKeyBytes, "err", err)
		return nil, err
	}

	return scalar, nil
}

func GetBlsPublicKey(suite *bn256.Suite, blsPrivateKey kyber.Scalar) (kyber.Point, error) {
	if suite == nil {
		log.Error("nil pointer dereference", "suite", suite, "err", utils.NilPtrDerefErr)
		return nil, utils.NilPtrDerefErr
	}

	return suite.Point().Mul(blsPrivateKey, nil), nil
}
