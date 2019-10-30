package enumgen

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lumina-tech/gooq/pkg/generator/utils"

	"github.com/lumina-tech/gooq/pkg/generator/metadata"
)

type EnumGenerator struct {
	modelTemplateString string
	modelOutputFile     string
}

func NewEnumGenerator(
	modelTemplatePath, modelOutputPath, dbname string,
) *EnumGenerator {
	modelOutputFile := fmt.Sprintf("%s/%s_enum.generated.go",
		modelOutputPath, dbname)
	return &EnumGenerator{
		modelTemplateString: modelTemplatePath,
		modelOutputFile:     modelOutputFile,
	}
}

func (gen *EnumGenerator) GenerateCode(
	data *metadata.Data,
) error {
	enums := append(data.Enums, data.ReferenceTableEnums...)
	sort.SliceStable(enums, func(i, j int) bool {
		return strings.Compare(enums[i].Name, enums[j].Name) < 0
	})
	args := templateArgs{
		Package:   "model",
		Timestamp: time.Now().Format(time.RFC3339),
		Enums:     enums,
	}
	enumTemplate := utils.GetTemplate(gen.modelTemplateString)
	return utils.RenderToFile(enumTemplate, gen.modelOutputFile, args)
}
