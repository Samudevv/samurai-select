# Samurai Select

A screen selection tool for wayland compositors using the layer shell. I thank [slurp](https://github.com/emersion/slurp) for teaching me how to use the layer shell and for showing an approach to creating a screen selection tool.

## Features

+ [x] Customizable (colors, sizes, fonts etc.)
+ [x] Screen Freeze (-z flag)
+ [x] Screenshot (-s flag)
+ [x] Execute arbitrary command (--cmd flag), usable when freezing screen
+ [x] Show Coordinates and Dimensions (-t flag)
+ [x] Alter selection after performing an initial selection (-A flag)
+ [x] Touch Support (needs testing)
+ [x] Force aspect ratio (-a flag)
+ [x] Select certain regions of screen (e.g. windows) (-r flag)
  + [x] Hyprland support (-r hyprland)
  + [x] Sway support (-r sway)
  + [x] Arbitrary (via argument) (-r arg -R 'X,Y WxH X1,Y1 W1xH1 ...')
+ [x] Select whole outputs (-p flag)

## Install

### Arch Linux (AUR)

This program is available through the AUR, you can install it using an AUR helper like **yay**:
```bash
yay -S samurai-select
```
or manually:
```bash
git clone https://aur.archlinux.org/samurai-select-git
cd samurai-select
makepkg -si
```

### Everything Else

If you have the dependencies listed under [Build](#Build) installed you can just install this program without having the source code by calling

```bash
go install github.com/Samudevv/samurai-select@latest
ln -s $GOPATH/bin/samurai-select $GOPATH/bin/smel
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
go build -v
```

Or this to install it:
```bash
go install -v
```

To build the man page:

1. Install `scdoc` and `gzip`
2. Execute:
```bash
scdoc < manpage.scd | gzip -c > samurai-select.1.gz
```
3. View it:
```bash
man -l samurai-select.1.gz
```
