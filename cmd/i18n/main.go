package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/itskreisler/i18n-go/generator"
	"github.com/itskreisler/i18n-go/validator"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "generate":
		generate(os.Args[2:])
	case "validate":
		validate(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func generate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	dir := fs.String("dir", "example/locales", "locale directory")
	pkg := fs.String("package", "locales", "Go package name")
	out := fs.String("out", "example/locales/keys_gen.go", "generated output")
	fs.Parse(args)

	if err := generator.Generate(*dir, *pkg, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func validate(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	dir := fs.String("dir", "example/locales", "locale directory")
	defaultLocale := fs.String("default", "en", "default locale")
	fs.Parse(args)

	report, err := validator.Validate(*dir, *defaultLocale)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !report.Valid() {
		for locale, keys := range report.Missing {
			for _, key := range keys {
				fmt.Printf("missing: %s -> %s\n", locale, key)
			}
		}
		for locale, keys := range report.Extra {
			for _, key := range keys {
				fmt.Printf("extra: %s -> %s\n", locale, key)
			}
		}
		os.Exit(1)
	}

	fmt.Println("i18n validation passed")
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: i18n <generate|validate> [flags]")
}
