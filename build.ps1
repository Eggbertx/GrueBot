#!/usr/bin/env pwsh
param (
	[string]$action = "build"
)
$ErrorActionPreference = "Stop"

$BIN = "gruebot"
$BINEXE = "$BIN$($env:GOOS=$platform; go env GOEXE)"

function build {
	$cmd = "& go build -v -gcflags=-trimpath=$PWD -asmflags=-trimpath=$PWD -o $BINEXE"
	Invoke-Expression $cmd
}

function clean {
	Remove-Item $BIN
	Remove-Item $BINEXE
}

function dependencies {
	go get -v github.com/bwmarrin/discordgo
}

switch ($action) {
	"build" { 
		build
	}
	"clean" {
		clean
	}
	"dependencies" {
		dependencies
	}
	Default {
		Write-Output "Invalid or unsupported command"
	}
}
