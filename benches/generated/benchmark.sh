#!/bin/bash
go test -c
./generated.test -test.run=NONE -test.bench="MulFIPS" -test.count=5 -test.benchtime=2s -test.cpu=1 .  > fips.txt
./generated.test -test.run=NONE -test.bench="MulCIOS" -test.count=5 -test.benchtime=2s -test.cpu=1 .  > cios.txt
./generated.test -test.run=NONE -test.bench="MulNoCarry" -test.count=5 -test.benchtime=2s -test.cpu=1 .  > nocarry.txt

benchstat -csv -sort=name -norange cios.txt > cios.csv 
benchstat -csv -sort=name -norange fips.txt > fips.csv 
benchstat -csv -sort=name -norange nocarry.txt > nocarry.csv 