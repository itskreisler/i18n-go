package validator

import (
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
