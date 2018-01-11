zlocker: zookeeper based isolated command execution
===================================================

[![Build Status](https://travis-ci.org/pyr/zlocker.svg?branch=master)](https://travis-ci.org/pyr/zlocker)

**zlocker** is a small tool meant to ease sequential execution of
a command across a large number of hosts.

## Configuration

**zlocker** is configured through command line arguments:

    -z
        Comma-separated list of zookeeper servers to contact
    -l
        Name of lock node in zookeeper
    -w
        Optional sleep period, defaults to none
    -t
	    Zookeeper session timeout

The rest of the command line will be fed to the shell if a lock
is successfully acquired.

## Building

If you wish to inspect **zlocker** and build it by yourself, you may do so
by cloning [this repository](https://github.com/pyr/zlocker) and
peforming the following steps :

    mkdir -p $(GOPATH)/src
    cd $(GOPATH)/src && git clone https://github.com/pyr/zlocker
    make

### Updating

It uses [godep](https://github.com/golang/dep), so it should be easy.

    dep status
    dep ensure -update

