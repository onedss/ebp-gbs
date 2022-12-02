#!/bin/bash
CWD=$(cd "$(dirname $0)";pwd)
"$CWD"/ebp-gbs stop
"$CWD"/ebp-gbs uninstall 
