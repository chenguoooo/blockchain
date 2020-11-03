#!/bin/bash
rm blockchain
rm *.db
rm transaction.dat

go build -o blockchain *.go
./blockchain
