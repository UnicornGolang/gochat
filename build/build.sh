#!/bin/sh
rm -rf ./release
mkdir  release
go build -o chat
chmod +x ./chat
mv chat ./release/
cp favicon.ico ./release/
cp -arf ./asset ./release/
cp -arf ./views ./release/
cp -arf ./config ./release/
cp index.html ./release/
