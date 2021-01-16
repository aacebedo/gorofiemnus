package gorofimenus

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
func (menu *Menu) Options() (res Options) {
	return menu.options
}

// SetOptions Sets the options of the menu to be passed to rofi.
func (menu *Menu) SetOptions(options Options) {
	menu.options = options
}

// GetValue Gets the value of the item when the menu is displayed in rofi.
func (menu *Menu) Value() (res string) {
	return menu.value
}

func (menu *Menu) runCommand(cmdToExec string) (res string, err error) {
	MainLogger.Debugf("Running command '%s'", cmdToExec)
	c := cmd.NewCommand(cmdToExec, cmd.WithStandardStreams, cmd.WithInheritedEnvironment(cmd.EnvVars{}))

	err = c.Execute()
	if err != nil {
		VerboseLogger.Errorf("Original error when running command: %s", cmdToExec)
		err = eris.Wrapf(InternalError, "Error while executing command '%s'", cmdToExec)

		return
	}

	if c.ExitCode() != 0 {
		MainLogger.Errorf("Unable to execute command '%s': '%s'", cmdToExec, c.Stderr())
		err = eris.Wrapf(InternalError, "Unable to execute command '%s', exit code is '%d'", cmdToExec, c.ExitCode())

		return
	}

	MainLogger.Debugf("Command '%s' successfully executed", cmdToExec)

	res = c.Stdout()

	return
}

// GetSelection Triggers rofi to request the selection to the user.
//nolint:funlen // One line more than the limit, cannot be solved
func (menu *Menu) GetSelection(keepInputFile bool) (res *doublylinkedlist.List, err error) {
	MainLogger.Debug("Getting selection")

	if menu.items.Size() != 0 { //nolint:nestif // Detect 5 level of nesting however it is not
		if _, err = menu.runCommand("rofi -v"); err != nil {
			err = eris.Wrap(InternalError, "Unable to find rofi")

			return
		}

		MainLogger.Debug("Creating temporary file to feed the input of rofi")

		var tmpFile *os.File

		tmpFile, err = ioutil.TempFile(os.TempDir(), "gorofimenu-")
		if err != nil {
			VerboseLogger.Errorf("An error happened when writing in the the temporary rofi input file: %s", err)
			err = eris.Wrap(InternalError, "Unable to create the temporary file for rofi menus")

			return
		}
		defer tmpFile.Close()
		defer func() {
			if !keepInputFile {
				MainLogger.Debug("Removing temporary rofi input file")
				os.Remove(tmpFile.Name())
			}
		}()

		MainLogger.Debug("Temporary rofi input file successfully created")
		MainLogger.Debug("Filling the rofi input file")
		menu.items.Each(func(index int, value interface{}) {
			if _, err = tmpFile.WriteString(value.(*Menu).Value() + "\n"); err != nil {
				VerboseLogger.Errorf("An error happened when writing in the the temporary rofi input file: %s", err)
				err = eris.Wrap(InternalError, "Unable to write the temporary file for rofi menus")

				return
			}
			MainLogger.Debugf("Added '%s'", value.(*Menu).Value())
		})

		MainLogger.Debug("Rofi input file successfully filled")

		cmdToExec := fmt.Sprintf("rofi -i -dmenu -selected-row 0 -p %s -input %s %v",
			shellescape.Quote(shellescape.StripUnsafe(menu.Value())), tmpFile.Name(), shellescape.QuoteCommand(menu.Options()))

		var captured string

		MainLogger.Debug("Displaying the rofi window to get the user's selection")

		captured, err = menu.runCommand(cmdToExec)
		if err != nil {
			err = eris.Wrap(err, "Unable to obtain the selection")

			return
		}

		captured = strings.Trim(captured, "\n")

		MainLogger.Debugf("Selection is successful, selected value is '%s'", captured)

		_, selectedValue := menu.items.Find(func(index int, value interface{}) bool {
			selectedMenu := value.(*Menu)

			return selectedMenu.Value() == captured
		})

		if selectedValue != nil {
			MainLogger.Debug("Requesting the selection for submenu")

			if res, err = selectedValue.(*Menu).GetSelection(keepInputFile); err == nil {
				res.Insert(0, menu)
			}

			MainLogger.Debugf("Submenu selection is successful, submenu selection is '%s'", menu.Value())
		} else {
			err = eris.Wrap(InternalError, "Unable to find the selected menu item")

			return
		}
	} else {
		res = doublylinkedlist.New()
		res.Add(menu)
	}

	return res, err //nolint:wrapcheck // Has been wrapped by eris
}

// AddSubMenu Adds a submenu to the menu.
func (menu *Menu) AddSubMenu(subMenu Menu) (err error) {
	MainLogger.Debugf("Adding submenu '%s' to menu '%s'", subMenu.Value(), menu.Value())
	itemIdx, _ := menu.items.Find(func(index int, value interface{}) bool {
		selectedMenu := value.(*Menu)

		return selectedMenu.Value() == subMenu.Value()
	})

	if itemIdx != -1 {
		err = eris.Wrapf(AlreadyExistsError, "Menu '%s' already contains submenu '%s'", menu.Value(), subMenu.Value())

		return
	}

	menu.items.Add(&subMenu)
	MainLogger.Debugf("Successfully added submenu '%s' to menu '%s'", subMenu.Value(), menu.Value())

	return
}

// RemoveSubMenu Removes the submenu from the menu.
func (menu *Menu) RemoveSubMenu(subMenu *Menu) (err error) {
	itemIdx, _ := menu.items.Find(func(index int, value interface{}) bool {
		selectedMenu := value.(*Menu)

		return selectedMenu.Value() == subMenu.Value()
	})

	if itemIdx == -1 {
		err = eris.Wrapf(NotExistsError, "Menu '%s' does not contains submenu '%s'", menu.Value(), subMenu.Value())

		return
	}

	menu.items.Remove(itemIdx)

	return
}
