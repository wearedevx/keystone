#!/bin/sh

WORK=$PWD

cd $PWD/functions/ksauth
sh ./deploy.sh

cd $WORK

cd $PWD/functions/ksapi
sh ./deploy.sh

cd $WORK
