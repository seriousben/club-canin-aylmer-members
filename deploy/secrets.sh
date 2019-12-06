#!/usr/bin/env bash

if [[ ! -z "${GITHUB_TOKEN}" ]]; then
    kubectl create secret docker-registry github-docker-registry \
            --docker-server=docker.pkg.github.com \
            --docker-username=seriousben \
            --docker-password=$GITHUB_TOKEN \
            --namespace=club-canin-aylmer-members
fi
