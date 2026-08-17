package i18n

import (
	"io/fs"
	"sort"
	"strings"

	gi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type Key string

type Localizer struct {
	bundle *Bundle
	inner  *gi18n.Localizer
}

// New builds a Localizer from the locale files embedded in fsys.
//
// langs[0] is the preferred language. Any remaining languages act as a
// fallback chain: a message missing from the preferred language is resolved
// from the last language (the bundle's default language, which is the one
// go-i18n falls back to) instead of panicking. A key missing from every
// loaded language still panics at lookup time.
func New(fsys fs.FS, langs ...string) (*Localizer, error) {
	if len(langs) == 0 {
		return nil, ErrNoLanguage
	}

	files, err := listLocaleFiles(fsys)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrNoLocaleFiles
	}

	// The last listed language is the bundle's default language, i.e. the
	// fallback target when a message is missing from the preferred language.
	bundle := NewBundle(langs[len(langs)-1])
	if err := bundle.LoadFS(fsys, files...); err != nil {
		return nil, err
	}

	return &Localizer{
		bundle: bundle,
		inner:  gi18n.NewLocalizer(bundle.inner, langs...),
	}, nil
}

func (l *Localizer) T(key Key) string {
	return l.mustLocalize(&gi18n.LocalizeConfig{
		MessageID: string(key),
	})
}

func (l *Localizer) TD(key Key, data map[string]any) string {
	return l.mustLocalize(&gi18n.LocalizeConfig{
		MessageID:    string(key),
		TemplateData: data,
	})
}

func (l *Localizer) TP(key Key, count any, data map[string]any) string {
	return l.mustLocalize(&gi18n.LocalizeConfig{
		MessageID:    string(key),
		PluralCount:  count,
		TemplateData: data,
	})
}

// mustLocalize localizes cfg, panicking only when nothing could be resolved.
//
// go-i18n returns a non-nil MessageNotFoundErr both when a key is missing
// entirely AND when it was resolved through the default-language fallback
// (in which case the message is set). MustLocalize would panic in both
// cases, so we call Localize and inspect the message instead.
func (l *Localizer) mustLocalize(cfg *gi18n.LocalizeConfig) string {
	msg, err := l.inner.Localize(cfg)
	if err != nil && msg == "" {
		panic(err)
	}
	return msg
}

func (l *Localizer) Localize(cfg *gi18n.LocalizeConfig) (string, error) {
	return l.inner.Localize(cfg)
}

func (l *Localizer) Bundle() *Bundle {
	return l.bundle
}

func listLocaleFiles(fsys fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	// Only active.<lang>.toml files are locales, matching the validator and
	// the //go:embed active.*.toml convention used across the project.
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "active.") || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	return files, nil
}
