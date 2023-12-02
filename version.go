package main

import "fmt"

const (
	MajorVersion = 23
	MinorVersion = 12
	PatchVersion = 0
)

func VersionString() string {
	return fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)
}
