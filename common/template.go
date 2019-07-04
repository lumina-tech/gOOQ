package common

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/knq/snaker"
	"golang.org/x/tools/imports"
)

var (
	templateFuncs = map[string]interface{}{
		"capitalize":     capitalize,
		"dict":           dictionary,
		"snakeToCamelID": snaker.SnakeToCamelIdentifier,
		"toLower":        strings.ToLower,
		"toUpper":        strings.ToUpper,
	}
)

func init() {
	snaker.AddInitialisms("OS")
}

func GetComponentTemplates(basePath string) ([]string, error) {
	return filepath.Glob(filepath.Join(basePath, "resources/template/components/*.tmpl"))
}

func GetTemplateFuncs() map[string]interface{} {
	return templateFuncs
}

func RenderToFile(tpl *template.Template, filename string, data interface{}) error {
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return err
	}
	if err := write(filename, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func gofmt(filename string, b []byte) ([]byte, error) {
	out, err := imports.Process(filename, b, nil)
	if err != nil {
		return b, errors.New("unable to gofmt")
	}
	return out, nil
}

func write(filename string, b []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return errors.New("failed to create directory")
	}

	formatted := b
	if strings.HasSuffix(filename, ".go") {
		formatted, err = gofmt(filename, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gofmt failed: %s\n", err.Error())
			formatted = b
		}
	}

	err = ioutil.WriteFile(filename, formatted, 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s", filename)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Template Functions
///////////////////////////////////////////////////////////////////////////////

func capitalize(value string) string {
	if len(value) == 0 {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func dictionary(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dictionary call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dictionary keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
