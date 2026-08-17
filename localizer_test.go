package i18n_test

import (
	"errors"
	"testing"
	"testing/fstest"

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

func TestNewWithoutLanguage(t *testing.T) {
	if _, err := i18n.New(fstest.MapFS{}); !errors.Is(err, i18n.ErrNoLanguage) {
		t.Fatalf("New without languages = %v, want ErrNoLanguage", err)
	}
}

func TestNewWithoutLocaleFiles(t *testing.T) {
	if _, err := i18n.New(fstest.MapFS{}, "en"); !errors.Is(err, i18n.ErrNoLocaleFiles) {
		t.Fatalf("New on empty FS = %v, want ErrNoLocaleFiles", err)
	}
}

func TestNewIgnoresNonActiveToml(t *testing.T) {
	fsys := fstest.MapFS{
		"active.en.toml": &fstest.MapFile{Data: []byte("[hello]\nother = \"Hello!\"\n")},
		"config.toml":    &fstest.MapFile{Data: []byte("not a locale\n")},
	}

	loc, err := i18n.New(fsys, "en")
	if err != nil {
		t.Fatal(err)
	}
	if got := loc.T(i18n.Key("hello")); got != "Hello!" {
		t.Fatalf("got %q, want Hello!", got)
	}
}

// TestFallbackToNextLanguage covers the headline fallback feature: a key
// missing in the preferred language is resolved from the last language
// (the bundle default) instead of panicking.
func TestFallbackToNextLanguage(t *testing.T) {
	fsys := fstest.MapFS{
		"active.en.toml": &fstest.MapFile{Data: []byte("[hello]\nother = \"Hello!\"\n\n[bye]\nother = \"Bye!\"\n")},
		"active.es.toml": &fstest.MapFile{Data: []byte("[hello]\nother = \"¡Hola!\"\n")},
	}

	loc, err := i18n.New(fsys, "es", "en")
	if err != nil {
		t.Fatal(err)
	}

	if got := loc.T(i18n.Key("hello")); got != "¡Hola!" {
		t.Errorf("hello (es) = %q, want ¡Hola!", got)
	}
	if got := loc.T(i18n.Key("bye")); got != "Bye!" {
		t.Errorf("bye (missing in es) = %q, want fallback Bye!", got)
	}
}

// TestFallbackTemplate verifies fallback works for templated messages too.
func TestFallbackTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"active.en.toml": &fstest.MapFile{Data: []byte("[welcome_user]\nother = \"Welcome, {{.Name}}!\"\n")},
		"active.es.toml": &fstest.MapFile{Data: []byte("[hello]\nother = \"¡Hola!\"\n")},
	}

	loc, err := i18n.New(fsys, "es", "en")
	if err != nil {
		t.Fatal(err)
	}

	if got := loc.TD(i18n.Key("welcome_user"), map[string]any{"Name": "Kreisler"}); got != "Welcome, Kreisler!" {
		t.Errorf("welcome_user (fallback) = %q, want Welcome, Kreisler!", got)
	}
}

// TestMissingKeyPanics verifies the typed-key contract: a key missing from
// every loaded language still panics at lookup time.
func TestMissingKeyPanics(t *testing.T) {
	fsys := fstest.MapFS{
		"active.en.toml": &fstest.MapFile{Data: []byte("[hello]\nother = \"Hello!\"\n")},
		"active.es.toml": &fstest.MapFile{Data: []byte("[hello]\nother = \"¡Hola!\"\n")},
	}

	loc, err := i18n.New(fsys, "es", "en")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for a key missing in every language")
		}
	}()
	_ = loc.T(i18n.Key("no_such_key"))
}
