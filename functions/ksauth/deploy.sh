#!/bin/sh

gcloud functions deploy ksauth --entry-point Auth --runtime go113 --trigger-http --allow-unauthenticated --project wearedevx --region europe-west1 --env-vars-file ./env.yml
