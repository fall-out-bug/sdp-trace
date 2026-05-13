package witness

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func parseCustomerPublicKeyPEM(raw []byte) (ed25519.PublicKey, error) {
	if containsSecretLike(raw) {
		// Public trust anchors must not include signing secrets.
		return nil, errors.New("private key input rejected")
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("public key or certificate PEM is required")
	}
	if block.Type == "CERTIFICATE" {
		// Certificates are accepted only as carriers for an Ed25519 public key;
		// revocation and custody are evaluated through explicit policy fields.
		return parseCertificatePublicKey(block.Bytes)
	}
	return parsePKIXPublicKey(block.Bytes)
}

func parseCertificatePublicKey(raw []byte) (ed25519.PublicKey, error) {
	cert, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, err
	}
	// Only Ed25519 keys are supported so signature verification has one narrow
	// algorithm contract.
	key, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("certificate must contain ed25519 public key")
	}
	return key, nil
}

func parsePKIXPublicKey(raw []byte) (ed25519.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		return nil, err
	}
	// Reject other public-key algorithms rather than silently changing the
	// signature contract.
	edKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("public key must be ed25519")
	}
	return edKey, nil
}
