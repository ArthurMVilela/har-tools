package encoding

import (
	"encoding/json"
	"os"

	"github.com/ArthurMVilela/har-tools/pkg/model"
)

func LoadHARFromFile(file string) (*model.HAR, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ParseFromJSON(data)
}

func ParseFromJSON(data []byte) (*model.HAR, error) {
	har := new(model.HAR)

	err := json.Unmarshal(data, har)
	if err != nil {
		return nil, err
	}

	return har, nil
}

func EncodeToJSON(src any, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(src, "", "	")
	}

	return json.Marshal(src)
}
