package main

import (
	"crypto/rand"
	eesdk "ee-sdk/client"
	"encoding/base64"
	"fmt"
	"os"

	sm2dilithium2hybrid "github.com/ploynomail/turingPQC/sm2_dilithium2_hybrid"
	"github.com/ploynomail/turingPQC/x509"
)

var kemPublickKeyBase64 = "HkSLs9sGqvNIUio0yLjE90wa7nVAkBKtyTjPd8ygSXURreDNVeRylEgQDydrqzgE9lZtLufJ/Wt5QuwWndaA/Et31WqgkWSQlnd1JHaQCXaFAZbIgxbKT3l90EYSY8eoYqFWTYZ7BLVy6jB3geU+hcYmr0IBiUQ2BFYFMYN79/Y8Z8IIB7E9UEgzrkmRRtYaLUUy1dpWo8JF+qiToqIRErMq/ulrazyog1mZeMGq4cRZFiNsOqRRfullrmAvpXM98quGlyAawZihiZWG+cJj1+W6iCUMB9qkkZWgoOZxuVinqAAnzceWyTYTE8tLUSuOQ4nLsyuCN5MY24CUdDluIRYmLtBwOZRT9OpfuJoNQ3dg9magj6IF0htNRms6FHluukQJcIw3eRURZVqLDvGG3TYxkJUcXusDcokvnqjPcvNY8hMxpNOTF/mVtdsL7exuG8oOi+UxzyBn2TkSROpjpwKgQRghp5WDMFtRP0QxkZc2iESBsDcXDihnPFSLgWwNxwPOxchKJ9kQx1QXrQoqMuXJXzM3Dsh3ZntWt2GNENMYVkcKAFQ/ZARs41glvMAWd4UCN2jIWuOGXnEzJQs0D0Y++7efMkEt1+k175p0wQpiE+iUCbSBEJILq5iM9RB2dTwQ8Zqdf+lYE0pTYqdURXdlCtIp+pdIyXVnAKNPbrHH2AbK8KyALmdih7Q+qHc4+hM3FMm6YLKS4UOit5McpKuvd+pwT5EvufB1DaG7+wpP/PYtg0CBXFtqnakSNeREHOl0NPijQzgt59VV1UN4yFwR82cbgQypfwd6q4V6+bMlNbqa4MtJAkSgbKRPRyAS4mh9oqUICQO/5ZPFwmGKZdsuipaZhRWkGSJaZxOhyAOyXLoOeTcXmbu2iFt8t9nIvzyDixO9TsQHI1M3hdXBrfpULWNEbzc1ugJB2ccQ0pm6dFkCl0oahCxXxUxFpHcUz8colDgaSXAgCzA9zAtbA9JmVnGXuWfPUZu0pQBZyYUr0ONKYnYOCvUOIOdJPLeqAJm96ABZCKZ/16sS8AGTHP2mmi2xxQ9IL6biZGiS0TI="

func ReqKey() {
	kemPublickKeyBytes, err := base64.StdEncoding.DecodeString(kemPublickKeyBase64)
	if err != nil {
		fmt.Println("Error decoding KEM public key:", err)
		return
	}

	fmt.Println("EE SDK Golang Example")
	var cfg = &eesdk.EESdkConfig{
		ServerAddr: "localhost:9001",
		AccessId:   "807e8150-da11-48a8-886e-54e9894e8003",
		AccessKey:  "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1=",
	}

	client, err := eesdk.NewClient(cfg)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	fmt.Println("Client created successfully")

	privateKey, err := sm2dilithium2hybrid.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return
	}

	publicKey := privateKey.Public()

	publicKeyDer, err := x509.MarshalPKIXPublicKey(publicKey)
	if publicKeyDer == nil || err != nil {
		fmt.Println("Error marshaling public key")
		return
	}

	certReq := &eesdk.CertificateRequest{
		PublicKey:       publicKeyDer,
		CACertID:        1,
		TemplateID:      1,
		KemPublicKey:    kemPublickKeyBytes,
		Subject:         eesdk.Subject{CommonName: "test-device"},
		Extensions:      []eesdk.Extension{},
		ExtraExtensions: []eesdk.Extension{},
	}

	fmt.Println("Requesting certificate...")

	rep, err := client.RequestCertificate(certReq)
	if err != nil {
		fmt.Println("Error requesting certificate:", err)
		return
	}
	fmt.Println("Certificate requested successfully:")
	fmt.Printf("Hyper Signer Cert Length: %d bytes\n", len(rep.HyperSignerCert))
	fmt.Printf("Hyper Enc Cert Length: %d bytes\n", len(rep.HyperEncCert))
	fmt.Printf("Encrypted Private Key Length: %d bytes\n", len(rep.EncryptedPrivateKey))
	defer client.Close()
}

func ReqKeyWithCSR() {
	// read csr
	csrPemBytes, err := os.ReadFile("../testdata/example.csr")
	if err != nil {
		fmt.Println("Error reading CSR file:", err)
		return
	}

	fmt.Println("EE SDK Golang Example - Request Key with CSR")
	var cfg = &eesdk.EESdkConfig{
		ServerAddr: "localhost:9001",
		AccessId:   "807e8150-da11-48a8-886e-54e9894e8003",
		AccessKey:  "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1=",
	}

	client, err := eesdk.NewClient(cfg)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	fmt.Println("Client created successfully")

	defer client.Close()
	fmt.Println("Requesting certificate with CSR...")
	rep, err := client.RequestCertificateWithPKCS10(csrPemBytes)
	if err != nil {
		fmt.Println("Error requesting certificate:", err)
		return
	}

	fmt.Println("Certificate requested successfully:")
	fmt.Printf("Signer Cert Length: %d bytes\n", len(rep.SignerCert))
	fmt.Printf("Enc Cert Length: %d bytes\n", len(rep.EncCert))
	fmt.Printf("Hyper Signer Cert Length: %d bytes\n", len(rep.HyperSignerCert))
	fmt.Printf("Hyper Enc Cert Length: %d bytes\n", len(rep.HyperEncCert))
	fmt.Printf("Encrypted Private Key Length: %d bytes\n", len(rep.EncryptedPrivateKey))

	fmt.Println("Certificates and Encrypted Private Key received:")
	fmt.Printf("Encrypted Private Key (Base64): %s\n", base64.StdEncoding.EncodeToString(rep.EncryptedPrivateKey))
}

func RevokeCertificate() {
	fmt.Println("EE SDK Golang Example")
	var cfg = &eesdk.EESdkConfig{
		ServerAddr: "localhost:9001",
		AccessId:   "807e8150-da11-48a8-886e-54e9894e8003",
		AccessKey:  "HanBpT2gT31Gfn83R0dTjKu0O9mv55bGDyVtP3IWOW1=",
	}

	client, err := eesdk.NewClient(cfg)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	fmt.Println("Client created successfully")
	client.RevokeCertificate("74", "keyCompromise")
}

func main() {
	// ReqKey()
	ReqKeyWithCSR()
	// RevokeCertificate()
}
