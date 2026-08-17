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

func New(fsys fs.FS, langs ...string) (*Localizer, error) {
    if len(langs) == 0 {
        return nil, ErrNoLanguage
    }

    bundle := NewBundle(langs[0])
    files, err := listLocaleFiles(fsys)
    if err != nil {
        return nil, err
    }

    if err := bundle.LoadFS(fsys, files...); err != nil {
        return nil, err
    }

    return &Localizer{
        bundle: bundle,
        inner:  gi18n.NewLocalizer(bundle.inner, langs...),
    }, nil
}

func (l *Localizer) T(key Key) string {
    return l.inner.MustLocalize(&gi18n.LocalizeConfig{
        MessageID: string(key),
    })
}

func (l *Localizer) TD(key Key, data map[string]any) string {
    return l.inner.MustLocalize(&gi18n.LocalizeConfig{
        MessageID:    string(key),
        TemplateData: data,
    })
}

func (l *Localizer) TP(key Key, count any, data map[string]any) string {
    return l.inner.MustLocalize(&gi18n.LocalizeConfig{
        MessageID:    string(key),
        PluralCount:  count,
        TemplateData: data,
    })
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

    files := make([]string, 0, len(entries))
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
            continue
        }
        files = append(files, entry.Name())
    }
    sort.Strings(files)
    return files, nil
}
