package i18n_test

import (
	"testing"

	i18n "github.com/itskreisler/i18n-go"
	"github.com/itskreisler/i18n-go/example/locales"
)

func TestLocalizer(t *testing.T) {
	loc, err := i18n.New(locales.FS, "es")
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
