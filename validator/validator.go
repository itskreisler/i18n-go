package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Report struct {
	Missing map[string][]string
	Extra   map[string][]string
}

func Validate(dir, defaultLocale string) (*Report, error) {
	locales, err := readLocales(dir)
	if err != nil {
		return nil, err
	}

	base, ok := locales[defaultLocale]
	if !ok {
		return nil, fmt.Errorf("default locale %q not found", defaultLocale)
	}

	report := &Report{
		Missing: map[string][]string{},
		Extra:   map[string][]string{},
	}

	for locale, keys := range locales {
		if locale == defaultLocale {
			continue
		}

		for key := range base {
			if !keys[key] {
				report.Missing[locale] = append(report.Missing[locale], key)
			}
		}

		for key := range keys {
			if !base[key] {
				report.Extra[locale] = append(report.Extra[locale], key)
			}
		}

		sort.Strings(report.Missing[locale])
		sort.Strings(report.Extra[locale])
	}

	return report, nil
}

func (r *Report) Valid() bool {
	return len(r.Missing) == 0 && len(r.Extra) == 0
}

func readLocales(dir string) (map[string]map[string]bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]bool{}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "active.") || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		locale := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "active."), ".toml")
		var tree map[string]any

		if _, err := toml.DecodeFile(filepath.Join(dir, entry.Name()), &tree); err != nil {
			return nil, err
		}

		keys := map[string]bool{}
		flatten("", tree, keys)
		result[locale] = keys
	}

	return result, nil
}

func flatten(prefix string, value map[string]any, out map[string]bool) {
	if isMessageTable(value) {
		out[prefix] = true
		return
	}

	for key, child := range value {
		table, ok := child.(map[string]any)
		if !ok {
			id := key
			if prefix != "" {
				id = prefix + "." + key
			}
			out[id] = true
			continue
		}

		id := key
		if prefix != "" {
			id = prefix + "." + key
		}
		flatten(id, table, out)
	}
}

func isMessageTable(table map[string]any) bool {
	// A message table's keys are only plural categories (plus an optional
	// description). Requiring ALL keys to be categories - not just one -
	// stops a namespace like [ram.zero] from being mistaken for a message
	// just because one of its sub-keys is named "zero".
	for key := range table {
		switch key {
		case "description", "zero", "one", "two", "few", "many", "other":
		default:
			return false
		}
	}
	return len(table) > 0
}
