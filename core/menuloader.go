package core

import (
	"os"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

// MenuLoader Utility struct to load the menus from YAML based sources.
type MenuLoader struct {
}

// MenuDescriptor Deserialized YAML representation of the menu.
type MenuDescriptor struct {
	Value    *string           `yaml:"value,omitempty"`
	Options  *[]string         `yaml:"options,omitempty"`
	Submenus *[]MenuDescriptor `yaml:"submenus,omitempty"`
}

// NewMenuLoader Create a new MenuLoader.
func NewMenuLoader() (res *MenuLoader) {
	res = &MenuLoader{}

	return
}

// LoadMenusFromFile Loads menus from a YAML file.
func (ml *MenuLoader) LoadMenusFromFile(filepath string) (res *Menu, err error) {
	mdesc := &MenuDescriptor{}

	file, err := os.Open(filepath)
	if err != nil {
		err = eris.Wrap(err, "Unable to open menu description file")

		return
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err = d.Decode(&mdesc); err != nil {
		err = eris.Wrap(err, "Unable to parse YAML string")

		return
	}

	res, err = ml.loadMenu(*mdesc)

	return
}

// LoadMenusFromString Loads menus from a YAML string.
func (ml *MenuLoader) LoadMenusFromString(strToLoad string) (res *Menu, err error) {
	mdesc := &MenuDescriptor{}

	if err = yaml.Unmarshal([]byte(strToLoad), &mdesc); err != nil {
		err = eris.Wrap(err, "Unable to parse YAML string")

		return
	}

	res, err = ml.loadMenu(*mdesc)

	return
}

func (ml *MenuLoader) loadMenu(menuDesc MenuDescriptor) (res *Menu, err error) {
	if menuDesc.Value == nil {
		err = eris.New("Missing value in menu")

		return
	}

	res = NewMenu(*menuDesc.Value)

	if menuDesc.Options != nil {
		res.SetOptions(*menuDesc.Options)
	}

	if menuDesc.Submenus != nil {
		for _, val := range *menuDesc.Submenus {
			subMenu, err2 := ml.loadMenu(val)
			if err2 != nil {
				err = eris.Wrapf(err2, "Unable to load menu")

				return
			}

			err = res.AddItem(*subMenu)
		}
	}

	return
}
