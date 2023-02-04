package scu

import (
	disp "github.com/ch3mz-za/SCUtil/pkg/display"
)

type SubFeature string

const (
	verLive disp.MenuStringOption = "LIVE"
	verPtu  disp.MenuStringOption = "PTU"
	verBack disp.MenuStringOption = "Back"
)

var ptuOrLiveMenu = disp.NewStringOptionMenu(
	"Select Game Version",
	[]disp.MenuStringOption{verLive, verPtu, verBack},
)
