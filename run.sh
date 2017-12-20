#!/bin/bash

LD_LIBRARY_PATH=./lib exec ./lib/ld-linux.so ./bin/turnserver -n --no-auth --stun-only --db ./turndb/turndb