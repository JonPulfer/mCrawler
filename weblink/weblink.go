package weblink

// Definition of the usage type of the link.
const (
	Hyperlink uint = iota
	Image
	Javascript
)

// Resource details describing the link.
type Resource struct {
	TargetURI string
	Type      uint
}
