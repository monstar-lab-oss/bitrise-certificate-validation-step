package main

import (
	"fmt"
	"os"
	"time"

	v1log "github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/certdownloader"
	"github.com/bitrise-io/go-xcode/v2/codesign"
)

func failf(format string, args ...interface{}) {
	v1log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	fmt.Println("Starting to check the expiration of certificates uploaded to Bitrise.")
	codesignInputs := codesign.Input{
		CertificateURLList:        `env:"certificate_url_list,required"`,
		CertificatePassphraseList: `env:"passphrase_list"`,
	}
	cmdFactory := command.NewFactory(env.NewRepository())
	codesignConfig, _ := codesign.ParseConfig(codesignInputs, cmdFactory)

	certDownloader := certdownloader.NewDownloader(codesignConfig.CertificatesAndPassphrases, retry.NewHTTPClient().StandardClient())
	certificates, err := certDownloader.GetCertificates()

	if err != nil {
		failf(err.Error())
	}

	if len(certificates) == 0 {
		failf("Found 0 uploaded certificates. Upload a code signing certificate in the App's code signing tab.")
	}

	fmt.Println("Found  %s certificates.", len(certificates))

	if certificatesValid(certificates) {
		fmt.Println("All certificates have valid expirations. ")
		os.Exit(0)
	} else {
		failf("Certificate expired here")
	}
}

func certificatesValid(certificateInfos []certificateutil.CertificateInfoModel) bool {
	preFilteredCerts := certificateutil.FilterValidCertificateInfos(certificateInfos)

	if len(preFilteredCerts.InvalidCertificates) != 0 {
		return false
	} else {
		return true
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
