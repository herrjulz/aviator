package processor

import "github.com/JulzDiverse/aviator/validator"

type SpruceProcessor struct {
	spruce *validator.Spruce
}

func Process(cfg validator.Spruce) ([]byte, error) {
	//processor := SpruceProcessor{&cfg}
	return []byte{}, nil
}
