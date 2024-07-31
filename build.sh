#!/bin/bash

set -e

podman-compose build sshor
container=$(podman create localhost/sshor:latest)

rm -rf target/bin
podman export "$container" | tar -xvf -

for file in target/bin/*; do 
    if [ -f "$file" ]; then 
        gpg --detach-sign $file
    fi 
done


