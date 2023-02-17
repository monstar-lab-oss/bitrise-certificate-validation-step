package main

import (
	"testing"

	"github.com/bitrise-io/go-xcode/certificateutil"
)

func TestValidCertificate(t *testing.T) {
	mock := []certificateutil.CertificateInfoModel{createTestCert(1)}
	valid := certificatesValid(mock)
	if !valid {
		t.Fatalf("Certificate is valid but is marked as expired.")
	}
}

func TestInvalidCertificate(t *testing.T) {
	mock := []certificateutil.CertificateInfoModel{createTestCert(-1)}
	valid := certificatesValid(mock)
	if valid {
		t.Fatalf("Certificate is expired but is marked as valid")
	}
}
