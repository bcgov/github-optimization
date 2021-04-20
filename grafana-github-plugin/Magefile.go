// +build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

func BuildAll() {
	build.BuildAll()
}
