# elastic2tf

<p>
    <a href="https://github.com/tehbooom/elastic2tf/releases"><img src="https://img.shields.io/github/v/release/tehbooom/elastic2tf.svg" alt="Latest Release"></a>
    <a href="https://github.com/tehbooom/elastic2tf/blob/main/LICENSE"><img src="https://img.shields.io/github/license/tehbooom/elastic2tf" alt="Latest Release"></a>
    <a href="https://goreportcard.com/report/github.com/tehbooom/elastic2tf"><img src="https://goreportcard.com/badge/github.com/tehbooom/elastic2tf" alt="GoDoc"></a>
    <a href="https://github.com/tehbooom/elastic2tf/actions/workflows/lint.yml"><img src="https://github.com/tehbooom/elastic2tf/actions/workflows/lint.yml/badge.svg" alt="Build Status"></a>
</p>


Elastic integrations to terraform `elasticstack_fleet_integration_policy` resource converter

## Usage

1. Go to your integration policy and select `Preview API request`

2. Copy the json object excluding the request URL which looks like this `PUT kbn:/api/fleet/package_policies/<randomid>`. An example of the json object is below

    ```json
    {
      "package": {
        "name": "1password",
        "version": "1.32.0"
      },
      "name": "1password-2",
      "namespace": "",
      "description": "",
      "policy_ids": [
        "5336cc87-8c23-4edb-981c-42ed666eead2"
      ],
      "vars": {},
      "inputs": {
        "1password-httpjson": {
          "enabled": true,
          "vars": {
            "url": "https://events.1password.com",
            "token": {
              "id": "brmcq5kBB9VoXeRVb6Ls",
              "isSecretRef": true
            },
            "disable_keep_alive": false
          },
          "streams": {
            "1password.audit_events": {
              "enabled": true,
              "vars": {
                "limit": 1000,
                "interval": "10s",
                "tags": [
                  "forwarded",
                  "1password-audit_events"
                ],
                "preserve_original_event": false
              }
            },
            "1password.item_usages": {
              "enabled": true,
              "vars": {
                "limit": 1000,
                "interval": "10s",
                "tags": [
                  "forwarded",
                  "1password-item_usages"
                ],
                "preserve_original_event": false
              }
            },
            "1password.signin_attempts": {
              "enabled": true,
              "vars": {
                "limit": 1000,
                "interval": "10s",
                "tags": [
                  "forwarded",
                  "1password-signin_attempts"
                ],
                "preserve_original_event": false
              }
            }
          }
        }
      }
    }
    ```

3. Run the binary `elastic2tf` and paste in the json object

> If you mess up something you can always press `Ctrl+l` to reset the input

4. Press `enter` and it will return the terraform resource

## Installation

You can download the binary corresponding to your operating system from the releases page on GitHub.

Once downloaded you can run the binary from the command line:

```bash
tar -xzf elastic2tf_Linux_x86_64.tar.gz
./elastic2tf
```

### Build From Source

Ensure that you have a supported version of Go properly installed and setup. You can find the minimum required version of Go in the go.mod file.

You can then install the latest release globally by running:

```bash
go install github.com/tehbooom/elastic2tf@latest
```
