package tabs

import (
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

var GameVersion scu.GameVersion
var gameVersionDropDown = widget.NewSelect([]string{string(scu.Live), string(scu.Ptu)}, setGameVerion)

func setGameVerion(value string) {
	GameVersion = scu.GameVersion(value)
}
