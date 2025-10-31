package main

import (
	"context"
	"crypto/rand"
	eesdk "ee-sdk/client"
	"fmt"

	sm2dilithium2hybrid "github.com/ploynomail/turingPQC/sm2_dilithium2_hybrid"
	"github.com/ploynomail/turingPQC/x509"
)

func main() {
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
		Subject:         eesdk.Subject{CommonName: "test-device"},
		Extensions:      []eesdk.Extension{},
		ExtraExtensions: []eesdk.Extension{},
	}

	fmt.Println("Requesting certificate...")

	rep, err := client.RequestCertificate(context.Background(), certReq)
	if err != nil {
		fmt.Println("Error requesting certificate:", err)
		return
	}
	fmt.Println("Certificate requested successfully:", rep)
	defer client.Close()
}
