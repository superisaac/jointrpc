package misc

import (
	"github.com/mitchellh/mapstructure"
)

// Decode takes an input structure and uses reflection to translate it to
// the output structure. output must be a pointer to a map or struct.
func DecodeStruct(input interface{}, output interface{}) error {
        config := &mapstructure.DecoderConfig{
                Metadata: nil,
		TagName:  "json",
                Result:   output,
        }

        decoder, err := mapstructure.NewDecoder(config)
        if err != nil {
                return err
        }

        return decoder.Decode(input)
}
