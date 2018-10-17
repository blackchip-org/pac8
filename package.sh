#!/bin/bash

rm -rf build
mkdir -p build

mkdir -p build/pac8
cp -r ~/pac8/{data,doc,rom} build/pac8
( cd build ; zip -r pac8-data.zip pac8 )

