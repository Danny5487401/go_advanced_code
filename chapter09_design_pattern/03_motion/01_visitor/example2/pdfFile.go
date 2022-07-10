package visitor

// PdfFile PdfFile
type PdfFile struct {
	path string
}

// Accept Accept
func (f *PdfFile) Accept(visitor Visitor) error {
	return visitor.Visit(f)
}
