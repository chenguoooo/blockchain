#!/bin/bash
rm blockchain
rm *.db
rm transaction.dat
rm wallet.dat

go build -o blockchain *.go
./blockchain h
