#!/bin/bash

set -e

BUILD=$(dirname "$0")

rm -rf $BUILD/pac8
mkdir -p $BUILD/pac8

mkdir -p $BUILD/pac8
cp -r ~/pac8/{ext,rom} $BUILD/pac8

( cd $BUILD ; zip -r pac8-data.zip pac8 )

