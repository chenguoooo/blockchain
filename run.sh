#!/bin/bash
rm blockchain

rm transaction.dat

go build -o blockchain *.go
./blockchain h
