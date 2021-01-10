package gorofimenus

import (
	"os"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

// MenuLoader Utility struct to load the menus from YAML based sources.
type MenuLoader struct{}

// MenuDescriptor Deserialized YAML representation of the menu.
type MenuDescriptor struct {
	Value    string           `yaml:"value,omitempty"`
	Options  []string         `yaml:"options,omitempty"`
	Submenus []MenuDescriptor `yaml:"submenus,omitempty"`
}

// NewMenuLoader Create a new MenuLoader.
func NewMenuLoader() (res *MenuLoader) {
	res = &MenuLoader{}

	return
}

// LoadMenusFromFile Loads menus from a YAML file.
func (ml *MenuLoader) LoadMenusFromFile(filepath string) (res *Menu, err error) {
	menuDesc := &MenuDescriptor{}

	MainLogger.Debugf("Loading menu from file '%s'", filepath)

	file, err := os.Open(filepath)
	if err != nil {
		VerboseLogger.Errorf("Original menu description loading error: '%s'", err)
		err = eris.Wrap(ConfigurationLoadingError, "Unable to open menu description file")

		return
	}
	defer file.Close()

	MainLogger.Debugf("Creating YAML decoder to decode menu")

	d := yaml.NewDecoder(file)
	if err = d.Decode(&menuDesc); err != nil {
		VerboseLogger.Errorf("Original menu description file parsing error: '%s'", err)
		err = eris.Wrap(ConfigurationLoadingError, "Unable to parse YAML string")

		return
	}

	res, err = ml.loadMenu(*menuDesc)

	return
}

// LoadMenusFromString Loads menus from a YAML string.
func (ml *MenuLoader) LoadMenusFromString(strToLoad string) (res *Menu, err error) {
	menuDesc := &MenuDescriptor{}

	MainLogger.Debugf("Loading menu from string '%s'", strToLoad)

	if err = yaml.Unmarshal([]byte(strToLoad), &menuDesc); err != nil {
		VerboseLogger.Errorf("Original menu description string parsing error: '%s'", err)
		err = eris.Wrap(ConfigurationLoadingError, "Unable to parse YAML string")

		return
	}

	res, err = ml.loadMenu(*menuDesc)

	MainLogger.Debug("Menu successfully loaded")

	return
}

func (ml *MenuLoader) loadMenu(menuDesc MenuDescriptor) (res *Menu, err error) {
	res = NewMenu(menuDesc.Value)
	MainLogger.Debugf("Recursion for loading menu '%s'", menuDesc.Value)

	if len(menuDesc.Options) != 0 {
		MainLogger.Debugf("Options of menu '%s' are '%s'", menuDesc.Value, menuDesc.Options)
		res.SetOptions(menuDesc.Options)
	}

	for _, val := range menuDesc.Submenus {
		var subMenu *Menu

		subMenu, err = ml.loadMenu(val)
		if err != nil {
			err = eris.Wrap(err, "Unable to load menu")

			return
		}

		err = res.AddSubMenu(*subMenu)
	}

	return
}
