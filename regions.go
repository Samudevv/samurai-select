/***********************************************************************************
 *                         This file is part of samurai-select
 *                    https://github.com/PucklaJ/samurai-select
 ***********************************************************************************
 * Copyright (c) 2023 Jonas Pucher
 *
 * This software is provided ‘as-is’, without any express or implied
 * warranty. In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 * claim that you wrote the original software. If you use this software
 * in a product, an acknowledgment in the product documentation would be
 * appreciated but is not required.
 *
 * 2. Altered source versions must be plainly marked as such, and must not be
 * misrepresented as being the original software.
 *
 * 3. This notice may not be removed or altered from any source
 * distribution.
 ************************************************************************************/

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	samure "github.com/PucklaJ/samurai-render-go"
)

type Regions interface {
	OutputRegions() []samure.Rect
	CursorPos() (int, int, error)
}

func DetectRegions() Regions {
	var stdout strings.Builder
	ps := exec.Command("ps", "-e")
	ps.Stderr = os.Stderr
	ps.Stdout = &stdout

	if err := ps.Run(); err != nil {
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(stdout.String()))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, "Hyprland") {
			return &HyprlandRegions{}
		} /*else if strings.HasSuffix(line, "sway") {
			// TODO: sway support
		}*/
	}

	return nil
}

type HyprlandRegions struct {
}

type HyprWorkspace struct {
	ID   int
	Name string
}

type HyprMonitor struct {
	ActiveWorkspace HyprWorkspace
}

type HyprClient struct {
	At        [2]int
	Size      [2]int
	Workspace HyprWorkspace
}

func (h HyprClient) IsOnScreen(monitors []HyprMonitor) bool {
	if !(h.At[0] != 0 || h.At[1] != 0 || h.Size[0] != 0 || h.Size[1] != 0) {
		return false
	}

	for _, m := range monitors {
		if m.ActiveWorkspace.ID == h.Workspace.ID {
			return true
		}
	}

	return false
}

func (*HyprlandRegions) OutputRegions() (rs []samure.Rect) {
	hyprctlPath, err := exec.LookPath("hyprctl")
	if err != nil {
		return
	}

	var stdout strings.Builder

	hyprctl := exec.Command(hyprctlPath, "-j", "clients")
	hyprctl.Stderr = os.Stderr
	hyprctl.Stdout = &stdout
	if err = hyprctl.Run(); err != nil {
		return
	}

	var clients []HyprClient
	decoder := json.NewDecoder(strings.NewReader(stdout.String()))
	if err = decoder.Decode(&clients); err != nil {
		return
	}

	stdout.Reset()
	hyprctl = exec.Command(hyprctlPath, "-j", "monitors")
	hyprctl.Stderr = os.Stderr
	hyprctl.Stdout = &stdout
	if err = hyprctl.Run(); err != nil {
		return
	}

	var monitors []HyprMonitor
	decoder = json.NewDecoder(strings.NewReader(stdout.String()))
	if err = decoder.Decode(&monitors); err != nil {
		return
	}

	for _, c := range clients {
		if c.IsOnScreen(monitors) {
			rs = append(rs, samure.Rect{
				X: c.At[0],
				Y: c.At[1],
				W: c.Size[0],
				H: c.Size[1],
			})
		}
	}

	return
}

func (*HyprlandRegions) CursorPos() (int, int, error) {
	hyprctlPath, err := exec.LookPath("hyprctl")
	if err != nil {
		return 0, 0, err
	}

	var stdout strings.Builder

	hyprctl := exec.Command(hyprctlPath, "cursorpos")
	hyprctl.Stderr = os.Stderr
	hyprctl.Stdout = &stdout
	if err = hyprctl.Run(); err != nil {
		return 0, 0, err
	}

	words := strings.Split(stdout.String(), ",")
	if len(words) != 2 {
		return 0, 0, errors.New("Invalid output")
	}

	xStr := strings.TrimSpace(words[0])
	yStr := strings.TrimSpace(words[1])

	x, err := strconv.ParseUint(xStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.ParseUint(yStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(x), int(y), nil
}
