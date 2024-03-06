package main

import "fmt"

const (
	MajorVersion = 1
	MinorVersion = 24
	PatchVersion = 1
)

func VersionString() string {
	return fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)
}
