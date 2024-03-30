//go:build mage

package main

import (
	"embed"
	"github.com/magefile/mage/mg"
)

type Build mg.Namespace

//go:embed src/*
var fs embed.FS
