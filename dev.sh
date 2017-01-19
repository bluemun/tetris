# Copyright 2017 The bluemun Authors. All rights reserved.
# Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
# license that can be found in the LICENSE file.

#!/bin/sh

running=0

function install {
  go install github.com/bluemun/engine
  echo "engine installed"
  go install github.com/bluemun/engine/graphics
  echo "engine/graphics installed"
  go install github.com/bluemun/engine/graphics/shader
  echo "engine/graphics/shader installed"
  go install github.com/bluemun/engine/graphics/render
  echo "engine/graphics/render installed"
  go install github.com/bluemun/tetris
  echo "Tetris installed"
}

while [ $running = 0 ]; do
  echo -n "d: go get dependencies x: debug, b: install, r: run, a: install and run >"
  read text

  running=1
  case $text in
    "d" )
      echo "Installing go-gl"
      go get -u github.com/go-gl/gl/v{3.2,3.3,4.1,4.2,4.3,4.4,4.5}-{core,compatibility}/gl
      echo "go-gl installed"
      echo "Installing go-glfw"
      go get -u github.com/go-gl/glfw/v3.2/glfw
      echo "go-glfw installed"
      echo "Installing go-logging"
      go get github.com/op/go-logging
      echo "go-logging installed"
    ;;
    "x" )
      godebug run *.go
    ;;
    "b" )
      install;
    ;;
    "r" )
      $GOPATH/bin/tetris
    ;;
    "a" )
      install;
      $GOPATH/bin/tetris
    ;;
    * )
      running=0
    ;;
  esac
done
