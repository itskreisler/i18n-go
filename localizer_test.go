package i18n

import (
    "testing"

    "github.com/kreisler/i18n/example/locales"
)

func TestLocalizer(t *testing.T) {
    loc, err := New(locales.FS, "es")
    if err != nil {
        t.Fatal(err)
    }

    if got := loc.T(locales.Hello); got != "¡Hola!" {
        t.Fatalf("got %q", got)
    }

    if got := loc.TD(locales.WelcomeUser, map[string]any{"Name": "Kreisler"}); got != "¡Bienvenido, Kreisler!" {
        t.Fatalf("got %q", got)
    }

    if got := loc.TP(locales.Items, 3, map[string]any{"Count": 3}); got != "Tienes 3 elementos." {
        t.Fatalf("got %q", got)
    }
}
