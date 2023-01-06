DISCONTINUATION OF PROJECT

This project will no longer be maintained by Intel.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.

If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the open source software community, please create your own fork of this project.  
# ISecL kms-client

This library provides functionalities to make API calls to ISecL Key Broker Service such as create/retrieve keys, create user and register user public key.

## System Requirements
- RHEL 8.1
- Epel 8 Repo
- Proxy settings if applicable

## Software requirements
- git
- `go` version >= `go1.12.1` & <= `go1.14.1`

# Step By Step Build Instructions

## Install required shell commands

### Install `go` version >= `go1.12.1` & <= `go1.14.1`
The `kms-client` requires Go version 1.12.1 that has support for `go modules`. The build was validated with the latest version 1.14.1 of `go`. It is recommended that you use 1.14.1 version of `go`. More recent versions may introduce compatibility issues. You can use the following to install `go`.
```shell
wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
tar -xzf go1.14.1.linux-amd64.tar.gz
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
```
```shell
go build ./...
```

Direct dependencies

| Name                  | Repo URL                        | Minimum Version Required              |
| ----------------------| --------------------------------| :------------------------------------:|
| logrus                | github.com/sirupsen/logrus      | v1.4.0                                |

*Note: All dependencies are listed in go.mod*

# Links
https://01.org/intel-secl/
