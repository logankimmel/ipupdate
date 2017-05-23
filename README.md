# Google Cloud DNS automatic updater
This will update the A record on the specified address in a Cloud DNS zone with
whatever external IP this container/program runs at.

_*NOTE: updates every 5 hours_

### Requirements:
* [Google Application Credential File](https://developers.google.com/identity/protocols/application-default-credentials)
* Go 1.7

### Usage:
#### Environment variables needed:
* `ADDR` DNS Name
* `ZONE` DNS Zone Name
* `PROJECT_ID` Google Project ID
* `GOOGLE_APPLICATION_CREDENTIALS` Path to Credential File

`go run ipupdate.go`

#### Easiest way to run:
[Docker Image](https://hub.docker.com/r/logankimmel/ipupdate/)
