package main

import (
	"testing"
	"time"

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

func createTestCert(validDays int) certificateutil.CertificateInfoModel {
	const (
		teamID     = "TESTING TEAM ID"
		commonName = "Apple Developer: TEST"
		teamName   = "TESTING TEAM NAME"
	)
	expiry := time.Now().AddDate(0, 0, validDays)
	serial := int64(1234)

	cert, privateKey, _ := certificateutil.GenerateTestCertificate(serial, teamID, teamName, commonName, expiry)

	certInfo := certificateutil.NewCertificateInfo(*cert, privateKey)

	return certInfo
}
