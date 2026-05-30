package language

import (
	"embed"
	"encoding/json"
)

//go:embed *.json
var languageFiles embed.FS

type Translations map[string]string

var languages = map[string]Translations{}

func LoadLanguages() error {
	files := map[string]string{
		"en": "en.json",
		"de": "de.json",
	}

	for lang, path := range files {
		data, err := languageFiles.ReadFile(path)
		if err != nil {
			return err
		}

		var t Translations
		err = json.Unmarshal(data, &t)
		if err != nil {
			return err
		}

		languages[lang] = t
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
