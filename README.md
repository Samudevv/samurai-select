# Samurai Select

A screen selection tool for wlroots based wayland compositors. I thank [slurp](https://github.com/emersion/slurp) for teaching me how to use the layer shell and for showing an approach to creating a screen selection tool.

## Features

+ [x] Customizable (colors, sizes, fonts etc.)
+ [x] Screen Freeze (-z flag)
+ [x] Screenshot (-s flag)
+ [x] Execute arbitrary command (--cmd flag), usable when freezing screen
+ [x] Show Coordinates and Dimensions (-t flag)
+ [x] Alter selection after performing an initial selection (-A flag)
+ [x] Touch Support (needs testing)
+ [x] Force aspect ratio (-a flag)

## Install

If you have the dependencies listed under [Build](#Build) installed you can just install this program without having the source code by calling

```bash
go install github.com/PucklaJ/samurai-select@latest && ln -s $GOPATH/bin/samurai-select $GOPATH/bin/smel
```

## Build

To build it you need to have a [go compiler](https://go.dev/), C compiler (for cgo) and the following dependencies installed:

+ [Wayland Client Library](https://gitlab.freedesktop.org/wayland/wayland)
+ [Cairo](https://www.cairographics.org/)

On Arch Linux you can install these dependencies like so:
```
sudo pacman -S --needed go gcc wayland cairo
```

Then call this to build it:
```
go build -v -o smel
```

Or this to install it:
```bash
go install -v && ln -s $GOPATH/bin/samurai-select $GOPATH/bin/smel
```
