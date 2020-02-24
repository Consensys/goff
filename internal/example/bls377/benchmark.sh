#!/bin/bash
go test -c
./bls377.test -test.run=NONE -test.bench="MulAssign" -test.count=5 -test.benchtime=2s -test.cpu=1 .  > bls377.txt
benchstat -csv -sort=name -norange bls377.txt > bls377.csv 