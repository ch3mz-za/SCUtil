package tabs

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func Restore(w fyne.Window) fyne.CanvasObject {
	restoreData := binding.BindStringList(&[]string{})
	// var restoreFiles *[]string
	// var err error
	var gameVer string
	gameVerDD := widget.NewSelect([]string{string(scu.Live), string(scu.Ptu)}, func(value string) {
		gameVer = value
		items, err := scu.GetBackedUpControlMappings(scu.GameVersion(gameVer))
		if err != nil {
			log.Println("error: " + err.Error())
		}
		restoreData.Set(*items)
	})

	restoreList := widget.NewListWithData(restoreData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	btnRestore := widget.NewButton("restore", func() {
		log.Printf("HALLO! | %s\n", gameVer)
	})

	cont := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("Restore Control Mappings"),
		gameVerDD,
		container.NewBorder(nil, btnRestore, nil, nil, restoreList),
	)

	return cont
}
