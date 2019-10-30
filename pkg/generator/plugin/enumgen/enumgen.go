package enumgen

import (
	"sort"
	"strings"
	"time"

	"github.com/lumina-tech/gooq/pkg/generator/utils"

	"github.com/lumina-tech/gooq/pkg/generator/metadata"
)

type EnumGenerator struct {
	outputFile string
}

func NewEnumGenerator(
	outputFile string,
) *EnumGenerator {
	return &EnumGenerator{
		outputFile: outputFile,
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
	enumTemplate := utils.GetTemplate(enumTemplate)
	return utils.RenderToFile(enumTemplate, gen.outputFile, args)
}
