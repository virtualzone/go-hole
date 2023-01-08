# Go-hole
[![](https://img.shields.io/github/v/release/virtualzone/go-hole)](https://github.com/virtualzone/go-hole/releases)
[![](https://img.shields.io/github/release-date/virtualzone/go-hole)](https://github.com/virtualzone/go-hole/releases)
[![](https://img.shields.io/github/workflow/status/virtualzone/go-hole/build)](https://github.com/virtualzone/go-hole/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/virtualzone/go-hole)](https://goreportcard.com/report/github.com/virtualzone/go-hole)
[![](https://img.shields.io/github/license/virtualzone/go-hole)](https://github.com/virtualzone/go-hole/blob/master/LICENSE)

Minimalistic DNS ad blocker and DNS proxy written in Go, inspired by [Pi-hole®](https://github.com/pi-hole/pi-hole).

## Features
* Blacklist DNS names via user-specific source lists
* Whitelisting
* Multiple upstream DNS servers
* Local name resolution
* Pre-built, minimalistic Docker image

## Usage
1. Create a ```config.yaml``` file. Use the [config.yaml](https://github.com/virtualzone/go-hole/blob/main/config.yaml) in this repository as a template and customize is according to your needs.
1. Run Go-Hole using Docker and mount your previously created ```config.yaml```:
    ```bash
    docker run \
        --rm \
        --mount type=bind,source=${PWD}/config.yaml,target=/app/config.yaml \
        -p 53:53/udp \
        ghcr.io/virtualzone/go-hole:latest
    ```