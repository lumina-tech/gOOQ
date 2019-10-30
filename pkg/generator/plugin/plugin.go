package plugin

import "github.com/lumina-tech/gooq/pkg/generator/metadata"

type Plugin interface {
	GenerateCode(data *metadata.Data) error
}
