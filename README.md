# Redfish CLI

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)    [![Build and Release](https://github.com/angelhvargas/redfishcli/actions/workflows/build-and-release.yml/badge.svg?branch=main)](https://github.com/angelhvargas/redfishcli/actions/workflows/build-and-release.yml)

`redfishcli` is a command-line tool designed to scan the health of baremetal servers manufactured by Lenovo or Dell. It supports servers using Lenovo XClarity Controller or Dell iDRAC 7 or greater. This tool provides a convenient way to monitor the health and status of your servers using Redfish APIs.

## Features

- Scan and report the health of baremetal servers.
- Support for Lenovo XClarity Controller.
- Support for Dell iDRAC 7 or greater.
- Retrieve RAID controller health status.
- Retrieve RAID drive details.
- Integration with Redfish APIs.

## Table of Contents

- [Redfish CLI](#redfish-cli)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [From Source](#from-source)
  - [Using Precompiled Binaries](#using-precompiled-binaries)
  - [Usage](#usage)
    - [Basic Commands](#basic-commands)
      - [Scan RAID Health](#scan-raid-health)
  - [Example Usage](#example-usage)
  - [Configuration](#configuration)
    - [Configuration File](#configuration-file)
      - [Example Configuration (config.yaml)](#example-configuration-configyaml)
  - [Using the Configuration File](#using-the-configuration-file)
  - [Contributing](#contributing)
  - [Fork the repository](#fork-the-repository)
  - [License](#license)

## Installation

### From Source

To install `redfishcli` from source, ensure you have Go installed on your system, then run:

```sh
go get github.com/angelhvargas/redfishcli
cd $GOPATH/src/github.com/angelhvargas/redfishcli
go install
```

## Using Precompiled Binaries

Precompiled binaries for various platforms are available on the releases page. Download the binary for your platform, extract it, and place it in a directory included in your system's PATH.

## Usage

### Basic Commands

#### Scan RAID Health

Scan the RAID health of a server:

```bash
redfishcli storage raid health --drives -t [controller-type] -u [username] -p [password] -n [hostname]
```

- -t: Controller type (idrac or xclarity).

- -u: Username for the BMC.

- -p: Password for the BMC.

- -n: Hostname or IP address of the server.

## Example Usage

 Scan the RAID health of a Dell server with iDRAC:

```sh
redfishcli storage raid health --drives -t idrac -u root -p "your_password" -n 192.168.1.100 | jq
```

Scan the RAID health of a Lenovo server with XClarity:

```sh
redfishcli storage raid health --drives -t xclarity -u admin -p "your_password" -n 192.168.1.101 | jq
```

## Configuration

### Configuration File

You can create a configuration file to scan multiple servers without providing login parameters each time. By default, redfishcli looks for a configuration file at ~/.redfishcli/config.yaml.

#### Example Configuration (config.yaml)

```yaml
servers:
  - type: "idrac"
    hostname: "192.168.1.100"
    username: "root"
    password: "your_password"
  - type: "xclarity"
    hostname: "192.168.1.101"
    username: "admin"
    password: "your_password"

# Add more servers as needed
```


## Using the Configuration File

To use the configuration file, simply run:

```bash
redfishcli storage raid health --drives
```

`redfishcli` will automatically load the servers listed in the configuration file and scan their health.

## Contributing

We welcome contributions to redfishcli. To contribute, please follow these steps:

## Fork the repository


- Create a new branch (`git checkout -b feature/your-feature`).

- Commit your changes (`git commit -am 'Add new feature'`).

- Push to the branch (`git push origin feature/your-feature`).

- Create a new Pull Request.

- Please ensure your code adheres to the project's coding standards and includes appropriate tests.

## License

redfishcli is licensed under the [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE). See the LICENSE file for more information.
## Release Process

This project uses [GoReleaser](https://goreleaser.com/) and GitHub Actions for automated releases.

To create a new release:
1.  Ensure all changes are committed and pushed to `main`.
2.  Create and push a new tag:
    ```bash
    git tag -a v1.0.0 -m "Release v1.0.0"
    git push origin v1.0.0
    ```
3.  GitHub Actions will automatically build the binaries, create a release, and attach the artifacts.
