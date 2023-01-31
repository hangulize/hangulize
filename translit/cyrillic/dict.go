//go:generate python -c "import json; exec(open('cyrtranslit-mapping.v1.1.py').read()); json.dump(TRANSLIT_DICT, open('dict.json', 'w'))"

package cyrillic

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type mapping struct {
	ToLatin    map[string]string `json:"tolatin"`
	ToCyrillic map[string]string `json:"tocyrillic"`
}

//go:embed dict.json
var dictJSON []byte

var dict map[string]mapping

func loadDict() map[string]mapping {
	if dict != nil {
		return dict
	}

	if err := json.Unmarshal(dictJSON, &dict); err != nil {
		panic(fmt.Errorf("cyrillic mapping.json: %w", err))
	}
	return dict
}
