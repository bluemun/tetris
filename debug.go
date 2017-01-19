// Copyright 2017 The bluemun Authors. All rights reserved.
// Use of this source code is governed by a GNU GENERAL PUBLIC LICENSE
// license that can be found in the LICENSE file.

// Tetris debug.go Defines the logger that is used throughout the project.
package main

import (
	"os"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("tetris")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	logging.NewBackendFormatter(backend1, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled)
}
