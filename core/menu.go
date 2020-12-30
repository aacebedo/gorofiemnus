package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aacebedo/gorofimenus/utils/loggers"
	"github.com/alessio/shellescape"
	"github.com/commander-cli/cmd"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/rotisserie/eris"
)

type (
	// Options Represents option to be passed to rofi.
	Options []string

	// Menu Represents the structure of a menu.
	Menu struct {
		items   doublylinkedlist.List
		value   string
		options Options
	}
)

// NewMenu Instantiates a new menu.
func NewMenu(value string) (res *Menu) {
	res = &Menu{items: *doublylinkedlist.New(), value: value}

	return
}

// GetOptions Retrieves the options of the menu to be passed to rofi.
func (menu *Menu) GetOptions() (res Options) {
	return menu.options
}

// SetOptions Sets the options of the menu to be passed to rofi.
func (menu *Menu) SetOptions(options Options) {
	menu.options = options
}

// GetValue Gets the value of the item when the menu is displayed in rofi.
func (menu *Menu) GetValue() (res string) {
	return menu.value
}

func (menu *Menu) runCommand(cmdToExec string) (res string, retCode interface{}) {
	res, retCode = cmd.CaptureStandardOutput(func() (res interface{}) {
		loggers.MainLogger.Debugf("Command to exec is '%s'", cmdToExec)
		c := cmd.NewCommand(cmdToExec, cmd.WithStandardStreams, cmd.WithInheritedEnvironment(cmd.EnvVars{}))

		if cmdErr := c.Execute(); cmdErr != nil {
			res = eris.Wrap(cmdErr, "Unable to execute command")

			return
		}
		if c.ExitCode() != 0 {
			res = eris.Errorf("Unable to execute command: '%s'", c.Stderr())

			return
		}

		return
	})

	return
}

// GetSelection Triggers rofi to request the selection to the user.
//nolint:funlen // One line more than the limit, cannot be solved
func (menu *Menu) GetSelection() (res *doublylinkedlist.List, err error) {
	if menu.items.Size() != 0 { //nolint:nestif // Detect 5 level of nesting however it is not
		var cmdErr interface{}
		if _, cmdErr = menu.runCommand("rofi -v"); cmdErr != nil {
			err = eris.Wrap(cmdErr.(error), "Unable to find rofi")

			return
		}

		var tmpFile *os.File

		tmpFile, err = ioutil.TempFile(os.TempDir(), "gorofimenu-")
		if err != nil {
			err = eris.Wrap(err, "Unable to create the temporary file for rofi menus")

			return
		}
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		menu.items.Each(func(index int, value interface{}) {
			if _, err = tmpFile.WriteString(value.(*Menu).GetValue() + "\n"); err != nil {
				err = eris.Wrap(err, "Unable to write the temporary file for rofi menus")

				return
			}
		})

		cmdToExec := fmt.Sprintf("rofi -i -dmenu -selected-row 0 %v -p %s -input %s",
			shellescape.QuoteCommand(menu.GetOptions()), menu.GetValue(), tmpFile.Name())

		var captured string

		captured, cmdErr = menu.runCommand(cmdToExec)
		if cmdErr != nil {
			err = eris.Wrap(cmdErr.(error), "Unable to execute rofi")

			return
		}

		captured = strings.Trim(captured, "\n")
		_, selecValue := menu.items.Find(func(index int, value interface{}) bool {
			selecMI := value.(*Menu)

			return selecMI.GetValue() == captured
		})

		if selecValue != nil {
			if res, err = selecValue.(*Menu).GetSelection(); err == nil {
				res.Insert(0, menu)
			}
		} else {
			err = eris.Errorf("Unable to find the selected menu item")

			return
		}
	} else {
		res = doublylinkedlist.New()
		res.Add(menu)
	}

	return res, err //nolint:wrapcheck // Has been wrapped by eris
}

// AddItem Adds an item to the menu. This shall be another menu.
func (menu *Menu) AddItem(item Menu) (err error) {
	itemIdx, _ := menu.items.Find(func(index int, value interface{}) bool {
		selecMI := value.(*Menu)

		return selecMI.GetValue() == item.GetValue()
	})

	if itemIdx != -1 {
		err = eris.Errorf("Menu already contains item '%s'", item.GetValue())

		return
	}

	menu.items.Add(&item)

	return
}

// RemoveItem Removes the item from the menu.
func (menu *Menu) RemoveItem(item *Menu) (err error) {
	itemIdx, _ := menu.items.Find(func(index int, value interface{}) bool {
		selecMI := value.(*Menu)

		return selecMI.GetValue() == item.GetValue()
	})

	if itemIdx == -1 {
		err = eris.Errorf("Menu does not contains item '%s'", item.GetValue())

		return
	}

	menu.items.Remove(itemIdx)

	return
}
