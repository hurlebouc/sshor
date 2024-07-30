#!/bin/bash

podman-compose build sshor
container=$(podman create localhost/sshor:latest)
podman export "$container" | tar -xvf -