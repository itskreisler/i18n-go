package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	report, err := Validate("../example/locales", "en")
	if err != nil {
		t.Fatal(err)
	}
	if !report.Valid() {
		t.Fatalf("expected valid report: %#v", report)
	}
}

// TestValidatePluralNamedNamespace is a regression test: namespaces whose
// sub-keys are plural categories (e.g. [ram.zero]) used to make the whole
// namespace count as one key, silently hiding every real missing key inside it.
func TestValidatePluralNamedNamespace(t *testing.T) {
	dir := t.TempDir()
	writeLocale(t, dir, "active.en.toml", `[hello]
other = "Hello!"

[ram.zero]
other = "0 MB"

[ram.all_off]
other = "all off"

[menu.prompt]
other = "Press enter"
`)
	// es is missing [menu.prompt].
	writeLocale(t, dir, "active.es.toml", `[hello]
other = "¡Hola!"

[ram.zero]
other = "0 MB"

[ram.all_off]
other = "todo apagado"
`)

	report, err := Validate(dir, "en")
	if err != nil {
		t.Fatal(err)
	}

	if report.Valid() {
		t.Fatal("expected invalid report: es is missing menu.prompt")
	}
	if got := report.Missing["es"]; len(got) != 1 || got[0] != "menu.prompt" {
		t.Errorf("Missing[es] = %v, want [menu.prompt]", got)
	}
	if len(report.Extra) != 0 {
		t.Errorf("Extra = %v, want none", report.Extra)
	}
}

func writeLocale(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
