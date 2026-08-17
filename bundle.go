package i18n

import (
    "io/fs"
    "path/filepath"
    "strings"

    gi18n "github.com/nicksnyder/go-i18n/v2/i18n"
    "github.com/BurntSushi/toml"
    "golang.org/x/text/language"
)

type Bundle struct {
    inner *gi18n.Bundle
}

func NewBundle(defaultLanguage string) *Bundle {
    tag := language.Make(defaultLanguage)
    b := gi18n.NewBundle(tag)
    b.RegisterUnmarshalFunc("toml", toml.Unmarshal)
    return &Bundle{inner: b}
}

func (b *Bundle) LoadFS(fsys fs.FS, files ...string) error {
    for _, name := range files {
        if _, err := b.inner.LoadMessageFileFS(fsys, name); err != nil {
            return err
        }
    }
    return nil
}

func (b *Bundle) LoadDir(fsys fs.FS, dir string) error {
    entries, err := fs.ReadDir(fsys, dir)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
            continue
        }
        path := filepath.Join(dir, entry.Name())
        if _, err := b.inner.LoadMessageFileFS(fsys, path); err != nil {
            return err
        }
    }
    return nil
}

func (b *Bundle) Inner() *gi18n.Bundle {
    return b.inner
}
