# Google Cloud DNS automatic updater
This will update the A record on the specified address in a Cloud DNS zone with
whatever external IP this container/program runs at.

_*NOTE: updates every 5 hours_

## Requirements:
* Docker 17.05

## Usage:
#### Command arguments:
* Environment Variables
  * `ADDR` DNS Name
  * `ZONE` DNS Zone Name
  * `PROJECT_ID` Google Project ID
* Volume mount:
  * [Google Application Credential File](https://developers.google.com/identity/protocols/application-default-credentials)
    * $CREDS_PATH:/data/googlecreds.json

### Run command:
```
docker run -d -v $CREDS_PATH:/data/googlecreds.json \
    -e "ADDR=" -e "ZONE=" -e "PROJECT_ID=" \
    logankimmel/ipupdate
```
