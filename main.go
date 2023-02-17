package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bitrise-io/go-steputils/stepconf"
	v1log "github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/certdownloader"
	"github.com/bitrise-io/go-xcode/v2/codesign"
)

// Config ...
type Config struct {
	CertificateURLList        string `env:"certificate_url_list,required"`
	CertificatePassphraseList string `env:"passphrase_list"`
}

func failf(format string, args ...interface{}) {
	v1log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	// Parse and validate inputs
	var cfg Config
	parser := stepconf.NewDefaultEnvParser()
	if err := parser.Parse(&cfg); err != nil {
		failf("Config: %s", err)
	}
	stepconf.Print(cfg)

	fmt.Println("Starting to check the expiration of certificates uploaded to Bitrise.")
	codesignInputs := codesign.Input{CertificateURLList: cfg.CertificateURLList}

	fmt.Println("Codesign inputs: ", codesignInputs)

	cmdFactory := command.NewFactory(env.NewRepository())
	codesignConfig, err := codesign.ParseConfig(codesignInputs, cmdFactory)

	if err != nil {
		failf(err.Error())
	}

	fmt.Println("Codesign config: ", codesignConfig)

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
