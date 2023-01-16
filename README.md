# Go-hole
[![](https://img.shields.io/github/v/release/virtualzone/go-hole)](https://github.com/virtualzone/go-hole/releases)
[![](https://img.shields.io/github/release-date/virtualzone/go-hole)](https://github.com/virtualzone/go-hole/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/virtualzone/go-hole)](https://goreportcard.com/report/github.com/virtualzone/go-hole)
[![](https://img.shields.io/github/license/virtualzone/go-hole)](https://github.com/virtualzone/go-hole/blob/master/LICENSE)

Minimalistic DNS server which serves as an upstream proxy and ad blocker. Written in Go, inspired by [Pi-hole®](https://github.com/pi-hole/pi-hole).

## Features
* Minimalistic DNS server, written in Golang, optimized for high performance
* Blacklist DNS names via user-specific source lists
* Whitelist DNS names that are actually blacklisted
* Multiple user-settable upstream DNS servers
* Caching of upstream query results
* Local name resolution
* Pre-built, minimalistic Docker image

## How it works
Go-hole serves as DNS server on your (home) network. Instead of having your clients sending DNS queries directly to the internet or to your router, they are resolved by your local Go-hole instance. Go-hole sends these queries to one or more upstream DNS servers and caches the upstream query results for maximum performance.

Incoming queries from your clients are checked against a list of unwanted domain names ("blacklist"), such as well-known ad serving domains and trackers. If a requested name matches a name on the blacklist, Go-hole responds with error code NXDOMAIN (non-existing domain). This leads to clients not being able to load ads and tracker codes. In case you want to access a blacklisted domain, you can easily add it to a whitelist.

As an additional feature, you can set a list of custom hostnames/domain names to be resolved to specific IP addresses. This is useful for accessing services on your local network by name instead of their IP addresses.

## Usage
1. Create a ```config.yaml``` file. Use the [config.yaml](https://github.com/virtualzone/go-hole/blob/main/config.yaml) in this repository as a template and customize is according to your needs.
1. Run Go-hole using Docker and mount your previously created ```config.yaml```:
    ```bash
    docker run \
        --rm \
        --mount type=bind,source=${PWD}/config.yaml,target=/app/config.yaml \
        -p 53:53/udp \
        ghcr.io/virtualzone/go-hole:latest
    ```
1. Set Go-hole as your network's DNS server (i.e. in your DHCP server's configuration).