package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	v1log "github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-xcode/certificateutil"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/certdownloader"
	"github.com/bitrise-io/go-xcode/v2/codesign"
)

// Config to parse the env vars
type Config struct {
	CertificateURLList        string          `env:"certificate_url_list,required"`
	CertificatePassphraseList stepconf.Secret `env:"passphrase_list"`
	KeychainPath              string          `env:"keychain_path,required"`
	KeychainPassword          stepconf.Secret `env:"keychain_password"`
}

// Fails the step
func failf(format string, args ...interface{}) {
	v1log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	// Parse and validate inputs
	var cfg Config
	parser := stepconf.NewInputParser(env.NewRepository())
	if err := parser.Parse(&cfg); err != nil {
		fmt.Println("❗️ Please check that you have the certificate uploaded in the Code Signing & Files tab\n")
		failf("Config: %s", err)
	}
	stepconf.Print(cfg)

	fmt.Println("Starting to check the expiration of certificates uploaded to Bitrise.\n")
	codesignInputs := codesign.Input{
		CertificateURLList:        cfg.CertificateURLList,
		CertificatePassphraseList: cfg.CertificatePassphraseList,
		KeychainPath:              cfg.KeychainPath,
		KeychainPassword:          cfg.KeychainPassword,
	}

	cmdFactory := command.NewFactory(env.NewRepository())
	codesignConfig, err := codesign.ParseConfig(codesignInputs, cmdFactory)

	if err != nil {
		failf(err.Error())
	}

	fmt.Println("⬇️ Downloading certificates from Bitrise.\n")
	certDownloader := certdownloader.NewDownloader(codesignConfig.CertificatesAndPassphrases, retry.NewHTTPClient().StandardClient())
	certificates, err := certDownloader.GetCertificates()

	if err != nil {
		failf(err.Error())
	}

	// Fail if no certificates are uploaded
	if len(certificates) == 0 {
		failf("❗️ Found 0 uploaded certificates. Upload a code signing certificate in the App's code signing tab.\n")
	}

	certCount := strconv.Itoa(len(certificates))
	fmt.Println("Found" + certCount + "uploaded certificates \n")

	// Validate the expiration of the certificates
	if certificatesValid(certificates) {
		fmt.Println("✅ All certificates have valid expirations")
		os.Exit(0)
	} else {
		failf("❗️ One or more of the certificates uploaded to Bitrise are expired. Upload a valid certificate.\n")
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
