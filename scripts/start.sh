#!/bin/bash
CWD=$(cd "$(dirname $0)";pwd)
"$CWD"/ebp-gbs install
"$CWD"/ebp-gbs start 
