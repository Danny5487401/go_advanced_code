package visitor

// PPTFile PPTFile
type PPTFile struct {
	path string
}

// Accept Accept
func (f *PPTFile) Accept(visitor Visitor) error {
	return visitor.Visit(f)
}
