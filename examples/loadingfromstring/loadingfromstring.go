package main

import (
	"os"

	"github.com/aacebedo/gorofimenus/pkg/gorofimenus"
	"github.com/sirupsen/logrus"
)

func main() {

	gorofimenus.SetLogLevel(logrus.DebugLevel, 1)

	ml := gorofimenus.NewMenuLoader()

	menu, err := ml.LoadMenusFromString(`
value: Snap
submenus:
  - value: '& exit 1 & echo 1'
    options: 
      - "'& exit 1 & echo 1'"
    submenus:
      - value: 'titi'
  - value: 2
  - value: "tata toto"
  - value: 4
  - value: 5
  - value: 6
  - value: 7
  - value: 8
  - value: 9`)
	// submenu.AddItem(core.NewSimpleMenuItem("menuItem3"))
	// submenu.AddItem(core.NewMenu("menuItem3"))
	if err == nil {
		selectedItem, _ := menu.GetSelection(false)
		gorofimenus.MainLogger.Debug(selectedItem)
	}

	if err == nil {
		os.Exit(0)
	} else {
		gorofimenus.MainLogger.Error(err)
		os.Exit(-1)
	}
}
