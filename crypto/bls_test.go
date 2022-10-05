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
	"go.dedis.ch/kyber/v3/util/key"
)

func TestBls(t *testing.T) {
	blsSuite := GetBlsSuite()
	require.NotNil(t, blsSuite)

	pair := key.NewKeyPair(blsSuite)
	privateKey, publicKey := pair.Private, pair.Public
	privateKeyString := privateKey.String()

	blsPrivateKey, err := GetBlsPrivateKey(blsSuite, privateKeyString)
	require.Nil(t, err)
	assert.Equal(t, privateKeyString, blsPrivateKey.String())

	blsPublicKey, err := GetBlsPublicKey(blsSuite, blsPrivateKey)
	require.Nil(t, err)
	assert.Equal(t, publicKey, blsPublicKey)
}

func TestCodec(t *testing.T) {
	blsSuite := GetBlsSuite()
	require.NotNil(t, blsSuite)

	pair := key.NewKeyPair(blsSuite)
	publicKey := pair.Public

	publicKeyBytes := EncodeBlsPublicKey(publicKey)
	assert.NotNil(t, publicKeyBytes)
	assert.Greater(t, len(publicKeyBytes), 0)

	actualPublicKey, err := DecodeBlsPublicKey(blsSuite, publicKeyBytes)
	assert.Nil(t, err)
	assert.True(t, actualPublicKey.Equal(publicKey))
}
