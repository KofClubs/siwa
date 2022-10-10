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
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

func Sign(signerSuite *bn256.Suite, signerDkg *DistributedKeyGenerator, message string) []byte {
	if signerDkg == nil || signerDkg.PedersenDkg == nil {
		log.Error("nil dkg of signer")
		return nil
	}

	distKey, err := signerDkg.PedersenDkg.DistKeyShare()
	if err != nil {
		log.Error("fail to generate distributed key of signer", "err", err)
		return nil
	}

	signatures, err := tbls.Sign(signerSuite, distKey.Share, []byte(message))
	if err != nil {
		log.Error("fail to sign message of signer", "message", message, "err", err)
		return nil
	}
	return signatures
}

func Verify(verifierSuite *bn256.Suite, verifierDkg *DistributedKeyGenerator, message string, signature []byte) bool {
	if verifierSuite == nil || verifierDkg == nil || verifierDkg.PedersenDkg == nil {
		log.Error("nil suite or dkg of verifier")
		return false
	}

	distKey, err := verifierDkg.PedersenDkg.DistKeyShare()
	if err != nil {
		log.Error("fail to generate distributed key of verifier", "err", err)
		return false
	}

	pubPoly := share.NewPubPoly(verifierSuite.G2(), verifierSuite.G2().Point().Base(), distKey.Commits)
	err = tbls.Verify(verifierSuite, pubPoly, []byte(message), signature)
	return err == nil
}

func VerifyAll(verifierSuite *bn256.Suite, verifierDkg *DistributedKeyGenerator, t, n int,
	message string, signatures [][]byte) ([]byte, bool) {
	if verifierSuite == nil || verifierDkg == nil || verifierDkg.PedersenDkg == nil {
		log.Error("nil suite or dkg of verifier")
		return nil, false
	}

	distKey, err := verifierDkg.PedersenDkg.DistKeyShare()
	if err != nil {
		log.Error("fail to generate distributed key of verifier", "err", err)
		return nil, false
	}

	pubPoly := share.NewPubPoly(verifierSuite.G2(), verifierSuite.G2().Point().Base(), distKey.Commits)
	signature, err := tbls.Recover(verifierSuite, pubPoly, []byte(message), signatures, t, n)
	if err != nil {
		log.Error("fail to reconstruct bls signature", "message", message, "err", err)
		return nil, false
	}
	return signature, true
}
