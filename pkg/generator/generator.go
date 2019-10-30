package generator

import (
	"github.com/jmoiron/sqlx"
	"github.com/lumina-tech/gooq/pkg/generator/metadata"
	"github.com/lumina-tech/gooq/pkg/generator/plugin"
	"github.com/lumina-tech/gooq/pkg/generator/postgres"
)

type Generator struct {
	plugins []plugin.Plugin
}

func NewGenerator(
	plugins ...plugin.Plugin,
) *Generator {
	return &Generator{plugins: plugins}
}

func (gen *Generator) Run(
	db *sqlx.DB,
) error {
	loader := postgres.NewPostgresLoader()
	data, err := metadata.NewData(db, loader)
	if err != nil {
		return err
	}
	for _, plugin := range gen.plugins {
		if err := plugin.GenerateCode(data); err != nil {
			return err
		}
	}
	return nil
}
