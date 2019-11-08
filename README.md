# ISecL kms-client

This library provides functionalities to make API calls to ISecL Key Broker Service such as create/retrieve keys, create user and register user public key.

## System Requirements
- RHEL 7.5/7.6
- Epel 7 Repo
- Proxy settings if applicable

## Software requirements
- git
- `go` version >= `go1.11.4` & <= `go1.12.12`

# Step By Step Build Instructions

## Install required shell commands

### Install `go` version >= `go1.11.4` & <= `go1.12.12`
The `kms-client` requires Go version 1.11.4 that has support for `go modules`. The build was validated with the latest version 1.12.12 of `go`. It is recommended that you use 1.12.12 version of `go`. More recent versions may introduce compatibility issues. You can use the following to install `go`.
```shell
wget https://dl.google.com/go/go1.12.12.linux-amd64.tar.gz
tar -xzf go1.12.12.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

## Build kms-client

- Git clone the kms-client
- Run scripts to build the kms-client

```shell
git clone https://github.com/intel-secl/kms-client.git
cd kms-client
go build ./...
```

Direct dependencies

| Name                  | Repo URL                        | Minimum Version Required              |
| ----------------------| --------------------------------| :------------------------------------:|
| logrus                | github.com/sirupsen/logrus      | v1.4.0                                |

*Note: All dependencies are listed in go.mod*

# Links
https://01.org/intel-secl/
