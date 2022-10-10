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
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

func DecodeBlsPublicKey(suite *bn256.Suite, data []byte) (kyber.Point, error) {
	if suite == nil {
		log.Error("nil suite", "err", utils.NilPtrDerefErr)
		return nil, utils.NilPtrDerefErr
	}

	point := suite.Point()
	err := point.UnmarshalBinary(data)
	if err != nil {
		log.Error("fail to unmarshal bls public key", "err", err)
		return nil, err
	}
	return point, nil
}

func EncodeBlsPublicKey(blsPublicKey kyber.Point) []byte {
	if blsPublicKey == nil {
		log.Error("nil bls public key")
		return nil
	}

	data, err := blsPublicKey.MarshalBinary()
	if err != nil {
		log.Error("fail to marshal bls public key", "err", err)
		return nil
	}
	return data
}
