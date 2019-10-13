package gooq

// Serializer interface
type Renderable interface {
	Render(builder *Builder)
}
