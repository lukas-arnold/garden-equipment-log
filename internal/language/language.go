package language

import (
	"embed"
	"encoding/json"
	"strings"
)

//go:embed *.json
var languageFiles embed.FS

type Translations map[string]string

var languages = make(map[string]Translations)

func LoadLanguages() error {
	entries, err := languageFiles.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			data, err := languageFiles.ReadFile(entry.Name())
			if err != nil {
				return err
			}

			var t Translations
			if err := json.Unmarshal(data, &t); err != nil {
				return err
			}

			lang := strings.TrimSuffix(entry.Name(), ".json")
			languages[lang] = t
		}
	}

	return nil
}

func T(lang, key string) string {
	if l, ok := languages[lang]; ok {
		if v, ok := l[key]; ok {
			return v
		}
	}
	return key
}
