#!/bin/bash
set -e

usage() { 
    echo "Usage: "
    echo "$0 [-t <tag>]
             [-h]" 1>&2; 
    echo "-t : Give a tag to the image. Build for minikube if not provided."
    echo "-h : Display help" 
    echo ""
    exit 0; 
}

microservice=kogger
user=docker #pass=tcuser
baseImageTag=v0.1.3
imageTag=ZOLARA_GIT_TAG

while getopts "t:h" option; do
    options+="$option"
    case "${option}" in
        t)
            if [[ $2 = \v* ]]; then
                imageTag=$2
            else
                echo "The tag ${OPTARG} has an invalid format. It has to starts with the letter 'v'"
                exit 1
            fi
            ;;
        h)
            usage
            ;;
        \? )
            usage
            ;;
    esac
done

echo "$microservice build started"
docker pull nexuszolara.me/library/zolara-ko:$baseImageTag
docker run --privileged --rm -v $(pwd):/workspace/kogger/$microservice -v "/var/run/docker.sock:/var/run/docker.sock:rw" -w /workspace/kogger/$microservice nexuszolara.me/library/zolara-ko:$baseImageTag build --local -t $imageTag --base-import-paths