package plugin

import "github.com/lumina-tech/gooq/pkg/generator/data"

type Plugin interface {
	Name() string
}

type CodeGenerator interface {
	GenerateCode(data *data.Data) error
}
