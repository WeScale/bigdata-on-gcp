#!/bin/bash

version=$1

export GCP_PROJECT=$(gcloud config get-value project)

docker build -t eu.gcr.io/$GCP_PROJECT/twitter-normalizer:$version .
docker push eu.gcr.io/$GCP_PROJECT/twitter-normalizer:$version
