package enumgen

import "github.com/lumina-tech/gooq/pkg/generator/metadata"

type templateArgs struct {
	Timestamp string
	Package   string
	Schema    string
	Enums     []metadata.Enum
}
