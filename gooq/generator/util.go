package generator

import (
	"io/ioutil"
	"text/template"

	"github.com/lumina-tech/lumina/apps/server/pkg/common"
)

func getTemplate(
	templatePath string,
) *template.Template {
	templateFuncMap := common.GetTemplateFuncs()
	schemaTmplateFile, err := ioutil.ReadFile(templatePath)
	check(err)
	return template.Must(template.New(ModelTemplateFilename).
		Funcs(templateFuncMap).Parse(string(schemaTmplateFile)))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
