# Copyright 2017 The bluemun Authors. All rights reserved.
# Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
# license that can be found in the LICENSE file.

#!/bin/sh

running=0

while [ $running = 0 ]; do
  echo -n "d: debug, b: install, r: run, a: install and run >"
  read text

  running=1
  case $text in
    "d" )
      godebug run *.go
    ;;
    "b" )
      go install github.com/bluemun/graphics
      go install github.com/bluemun/tetris
      echo "Tetris installed"
    ;;
    "r" )
      tetris
    ;;
    "a" )
      go install github.com/bluemun/graphics
      go install github.com/bluemun/tetris
        tetris
    ;;
    * )
      running=0
    ;;
  esac
done
