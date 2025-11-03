package domain

import (
	"crypto/x509/pkix"
	"encoding/asn1"
)

type EncryptedValue struct {
	Raw         asn1.RawContent
	IntendedAlg pkix.AlgorithmIdentifier `asn1:"explicit,tag:0,optional,omitempty"`
	SymmAlg     pkix.AlgorithmIdentifier `asn1:"explicit,tag:1,optional,omitempty"`
	EncSymmKey  asn1.BitString           `asn1:"explicit,tag:2,optional,omitempty"`
	KeyAlg      pkix.AlgorithmIdentifier `asn1:"explicit,tag:3,optional,omitempty"`
	ValueHint   []byte                   `asn1:"explicit,tag:4,optional,omitempty"`
	EncValue    asn1.BitString
}
