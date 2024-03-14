/**
 * Copyright © 2024, Staufi Tech - Switzerland
 * All rights reserved.
 *
 *   ________________________   ___ _     ________________  _  ____
 *  / _____  _  ____________/  / __|_|   /_______________  | | ___/
 * ( (____ _| |_ _____ _   _ _| |__ _      | |_____  ____| |_|_
 *  \____ (_   _|____ | | | (_   __) |     | | ___ |/ ___)  _  \
 *  _____) )| |_/ ___ | |_| | | |  | |     | | ____( (___| | | |
 * (______/  \__)_____|____/  |_|  |_|     |_|_____)\____)_| |_|
 *
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package crypto

import (
	"bytes"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"

	log "github.com/ChrIgiSta/go-utils/logger"
)

type KeyLength int

const (
	KeyLength1024Bit KeyLength = 1024
	KeyLength2048Bit KeyLength = 2048
	KeyLength4096Bit KeyLength = 4096
	KeyLength8192Bit KeyLength = 8192
)

const MinimalAllowedKeyLength = KeyLength2048Bit

type CertChecker struct {
	rootCAs *x509.CertPool
}

func NewCustomCertChecker(rootCAs *x509.CertPool) *CertChecker {
	return &CertChecker{
		rootCAs: rootCAs,
	}
}

// This Certificate checker ignores the SAN and CN
func (c *CertChecker) X509CeckCertNoSAN(rawCerts [][]byte,
	verifiedChains [][]*x509.Certificate) (err error) {

	var certificates []*x509.Certificate

	// read certs
	for _, serverCerts := range rawCerts {
		crt, e := x509.ParseCertificate(serverCerts)
		if e != nil {
			err = x509.CertificateInvalidError{}
			break
		}
		certificates = append(certificates, crt)
	}

	// validate
	for _, cert := range certificates {
		_ = log.Debug("x509",
			"Certificate Check:> issuer: %v, isCa:%v, notBefore: %v, notAfter: %v",
			cert.Issuer, cert.IsCA, cert.NotBefore, cert.NotAfter)

		opts := x509.VerifyOptions{
			Roots:         c.rootCAs,
			DNSName:       "", // disables verifying hostname or IP
			Intermediates: x509.NewCertPool(),
		}
		_, err = cert.Verify(opts)
	}

	return err
}

type CertificateSubject struct {
	Organisation string
	Country      string
	Province     string
	Locality     string
	OrgUnit      string
	CommonName   string
}

func CreateSelfsignedX509Certificate(serialNumber *big.Int,
	validityDays int,
	rsaKeyLen KeyLength,
	subject CertificateSubject) (certificate []byte, privateKey []byte, err error) {

	signCert, privKey, err := CreateCert(serialNumber,
		validityDays, rsaKeyLen, subject)
	if err != nil {
		return nil, nil, err
	}

	return EncodeCertificatePEMForm(signCert),
		EncodeRsaKeyPEMForm(privKey), nil
}

func CreateCert(serialNumber *big.Int, validityDays int,
	keyLength KeyLength, subject CertificateSubject) ([]byte, *rsa.PrivateKey, error) {

	if keyLength < MinimalAllowedKeyLength {

		return nil, nil,
			fmt.Errorf("rsa key length lower than %v to insecure",
				int(MinimalAllowedKeyLength))
	}

	privateKey, err := rsa.GenerateKey(crand.Reader, int(keyLength))
	if err != nil {
		return nil, nil, err
	}

	certificate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{subject.Organisation},
			Country:            []string{subject.Country},
			Province:           []string{subject.Province},
			Locality:           []string{subject.Locality},
			OrganizationalUnit: []string{subject.OrgUnit},
			CommonName:         subject.CommonName,
		},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Local().Add(time.Duration(validityDays) * time.Hour * 24),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// self-sign certificate
	signCert, err := x509.CreateCertificate(
		crand.Reader, certificate, certificate, &privateKey.PublicKey, privateKey)

	if err != nil {
		return nil, nil, err
	}
	return signCert, privateKey, err
}

func EncodeCertificatePEMForm(certificate []byte) []byte {

	return EncodePEMForm(certificate, "CERTIFICATE")
}

func EncodeRsaKeyPEMForm(privateKey *rsa.PrivateKey) []byte {

	return EncodePEMForm(
		x509.MarshalPKCS1PrivateKey(privateKey),
		"RSA PRIVATE KEY")
}

func EncodePEMForm(content []byte, typ string) []byte {
	pemForm := new(bytes.Buffer)
	if pem.Encode(pemForm, &pem.Block{
		Type:  typ,
		Bytes: content,
	}) != nil {
		return nil
	}
	return pemForm.Bytes()
}
