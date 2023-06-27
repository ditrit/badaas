package conditions

import (
	"github.com/dave/jennifer/jen"

	"github.com/ditrit/badaas/tools/badctl/cmd/version"
)

type File struct {
	destPkg string
	jenFile *jen.File
	name    string
}

func NewFile(destPkg, name string) *File {
	// Start a new file in destination package
	f := jen.NewFile(destPkg)

	// Add a package comment, so IDEs detect files as generated
	f.PackageComment("Code generated by badctl v" + version.Version + ", DO NOT EDIT.")

	return &File{
		destPkg: destPkg,
		name:    name,
		jenFile: f,
	}
}

func (file File) Add(codes ...jen.Code) {
	for _, code := range codes {
		file.jenFile.Add(code)
	}
}

// Write generated file
func (file File) Save() error {
	return file.jenFile.Save(file.name)
}
