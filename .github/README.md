# Half-Life TV Manager <img align="right" src="./HLTV-Manager.png" alt="HLTV Launcher" width="210" height="200"/>

The service runs in a docker container.

Service runs hltv servers in containers.

Service allows to download demos, also automatically monitors and deletes old demos.

## Description

Half-Life TV Manager - Allows you to run an unlimited number of hltv servers. Provides a site for downloading hltv demos.

## Characteristics

- The service is installed and started using docker.
- Everything is configured via yaml configuration. (Temporary)
- Support for running multiple HLTV servers.
- Site for downloading demos.
- Automatic deletion of demos.
- Offline demos. (Temporary)

## Installation

<details>
  <summary>Ubuntu</summary>

- Download docker-compose  

    `sudo apt update && sudo apt upgrade`

    `sudo apt install docker-compose`

- Download the HLTV container

    `sudo docker pull ghcr.io/wesstorn/hltv-files:v1.1`

- Download Hltv-Manager and log into it

    `git clone --branch self-hosted https://github.com/WessTorn/HLTV-Manager.git`

    `cd HLTV-Manager`

    Configuring the docker-compose config

    `nano .env`

    Setting up our HLTVs

    `nano hltv-runners.yaml`

- Starting the service

    `sudo docker-compose up -d`

- Docker commands

    `sudo docker-compose up -d` - Run in the background

    `sudo docker-compose up` - Run in current session (shows logs)

    `sudo docker-compose down` - Stop service

    `sudo docker-compose logs` - View logs
</details>


## In the future.

- Configuration, setup, launching HLTV through the site.
- Live HLTV terminals
- Support hltv with live broadcasts.
- Amxx api part for remote work with hltv server.
