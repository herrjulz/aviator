package modifier

import (
	"errors"

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
	if mod.Delete != "" {
		if yml, err := m.goml.Delete(file, mod.Delete); err == nil {
			return yml, nil
		}
	} else if mod.Set != "" {
		if yml, err := m.goml.Set(file, mod.Set, mod.Value); err == nil {
			return yml, nil
		}
	} else if mod.Update != "" {
		if yml, err := m.goml.Update(file, mod.Update, mod.Value); err == nil {
			return yml, nil
		}
	} else {
		return nil, errors.New("modification path not provided")
	}
	return nil, err
}
