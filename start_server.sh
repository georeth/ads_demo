#!/bin/bash

killall server
./server/server 0 localhost:13000 localhost:13001 localhost:13002 &
./server/server 1 localhost:13000 localhost:13001 localhost:13002 &
./server/server 2 localhost:13000 localhost:13001 localhost:13002 &
