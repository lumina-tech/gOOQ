package plugin

import "github.com/lumina-tech/gooq/pkg/generator/metadata"

type Plugin interface {
	Name() string
}

type CodeGenerator interface {
	GenerateCode(data *metadata.Data) error
}
