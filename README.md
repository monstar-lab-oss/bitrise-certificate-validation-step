# Code signing certificate - expiration check

This Bitrise step checks that at least one certificate is uploaded for the app in Bitrise and that all the uploaded certificates have valid expiration dates.

This step is designed to be used with [Apple Cloud Signing](https://github.com/nodes-ios/Playbook/blob/master/ci/bitrise-complete-guide-cloud-signing.md) which requires a development code signing certificate to be uploaded. 

## The problem this step solves

If the certificate uploaded to Bitrise expires, Cloud Signing ignores it and starts creating new development certificates on Apple portal. However, this creates
a new development certificate for every new build (the private key does not get transfered to a new fresh cloud machine) and soon we hit the limit for number of allowed certificates on Apple portal. This step prevents that by failing the build 
when the uploaded development certificate expires. 