package core

import (
	"fmt"
	"strings"

	"github.com/aacebedo/gorofimenus/core/loggers"
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
	res, retCode = cmd.CaptureStandardOutput(func() (retCode interface{}) {
		loggers.MainLogger.Debugf("Command to exec is '%s'", cmdToExec)
		c := cmd.NewCommand(cmdToExec, cmd.WithStandardStreams, cmd.WithInheritedEnvironment(cmd.EnvVars{}))

		if cmdErr := c.Execute(); cmdErr != nil {
			retCode = eris.Wrap(cmdErr, "Unable to execute command")
		}
		if c.ExitCode() != 0 {
			return eris.Errorf("Unable to execute command: '%s'", c.Stderr())
		}

		return
	})

	return
}

// GetSelection Triggers rofi to request the selection to the user.
func (menu *Menu) GetSelection() (res string, err error) {
	if menu.items.Size() != 0 { //nolint:nestif // Detect 5 level of nesting however it is not
		if _, cmdErr := menu.runCommand("rofi -v"); cmdErr != nil {
			err = eris.Wrap(cmdErr.(error), "Unable to find rofi")

			return
		}

		items := []string{}
		menu.items.Each(func(index int, value interface{}) {
			items = append(items, shellescape.Quote(value.(*Menu).GetValue()))
		})

		itemsStr := strings.Join(items, "\n")
		itemsStr = strings.TrimSuffix(itemsStr, "\n")

		cmdToExec := fmt.Sprintf("echo '%s' | rofi -i -dmenu -selected-row 0 %v", itemsStr,
			shellescape.QuoteCommand(menu.GetOptions()))
		captured, cmdErr := menu.runCommand(cmdToExec)

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
			res, err = selecValue.(*Menu).GetSelection()
		} else {
			err = eris.Errorf("Unable to find the selected menu item")

			return
		}
	} else {
		res = menu.GetValue()
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
