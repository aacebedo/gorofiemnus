package main

import (
	"os"

	"github.com/aacebedo/gorofimenus/pkg/gorofimenus"
	"github.com/sirupsen/logrus"
)

func main() {

	gorofimenus.SetLogLevel(logrus.DebugLevel, 1)

	ml := gorofimenus.NewMenuLoader()
	menu, err := ml.LoadMenusFromFile("./menu.yml")
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
