package modifier

import (
	"github.com/JulzDiverse/aviator"
	"github.com/JulzDiverse/aviator/gomlclient"
)

type Modifier struct {
	goml aviator.GomlClient
}

func NewModifier(goml aviator.GomlClient) *Modifier {
	return &Modifier{
		goml: goml,
	}
}

func New() *Modifier {
	return &Modifier{
		goml: gomlclient.New(),
	}
}

func (m *Modifier) Modify(file []byte, mod aviator.Modify) ([]byte, error) {
	var err error
	modified := file
	if len(mod.Delete) > 0 {
		for _, v := range mod.Delete {
			if yml, err := m.goml.Delete(modified, v); err == nil {
				modified = yml
			}
		}
	}
	if len(mod.Set) > 0 {
		for _, set := range mod.Set {
			if yml, err := m.goml.Set(modified, set.Path, set.Value); err == nil {
				modified = yml
			}
		}
	}
	if len(mod.Update) > 0 {
		for _, update := range mod.Update {
			if yml, err := m.goml.Update(modified, update.Path, update.Value); err == nil {
				modified = yml
			}
		}
	}
	return modified, err
}
