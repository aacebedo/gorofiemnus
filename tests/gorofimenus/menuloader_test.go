package gorofimenus_test

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/aacebedo/gorofimenus/pkg/gorofimenus"
)

func TestLoadFromValidString(t *testing.T) {
	gorofimenus.SetLogLevel(logrus.DebugLevel, 1)

	ml := gorofimenus.NewMenuLoader()

	_, err := ml.LoadMenusFromString(`
value: Snap
submenus:
  - value: 'item1'
  - value: 'item2'
  - value: 'item3'
    `)
	if err != nil {
		t.Errorf("Unable to load valid menu string")
	}
}

func TestLoadFromInvalidString(t *testing.T) {
	gorofimenus.SetLogLevel(logrus.DebugLevel, 1)

	ml := gorofimenus.NewMenuLoader()

	_, err := ml.LoadMenusFromString(`
value: Snap
submenus:
  - value: 'item1'
 value: 'item2'
  - value: 'item3'
    `)
	if err == nil {
		t.Errorf("Invalid menu string has been loaded")
	}
}
