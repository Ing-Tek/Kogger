#!/bin/bash
set -e

microservice_package="koggerserver"

echo "Generating models"
mkdir -p models
docker run --rm -v ${PWD}:/local --user $(id -u):$(id -g) nexuszolara.me/library/swagger-codegen-kubequest:v0.1.1 generate -i /local/api/restAPI.yaml -l go-server -o /local -t /templates --additional-properties hideGenerationTimestamp=true -Dservice -Dmodels -DpackageName=models
mv go/* models
rm -rf go

echo "Generating api"
mkdir -p ${microservice_package}
docker run --rm -v ${PWD}:/local --user $(id -u):$(id -g) nexuszolara.me/library/swagger-codegen-kubequest:v0.1.1 generate -i /local/api/restAPI.yaml -l go-server -o /local -t /templates --additional-properties hideGenerationTimestamp=true -Dservice -DpackageName=${microservice_package}
rm -f go/model_*
mv go/* ${microservice_package}
rm -rf go
