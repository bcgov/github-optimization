#!/bin/bash

find . -name '*.go' -exec go fmt {} \;
