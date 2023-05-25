package conditions

import (
	"errors"
	"fmt"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/ditrit/verdeter"
	"github.com/spf13/cobra"

	"golang.org/x/tools/go/packages"
)

var GenConditionsCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "conditions",
	Short: "Generate conditions to query your objects using BaDORM",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateConditions,
	Args:  cobra.MinimumNArgs(1),
})

var destPkg string

func generateConditions(cmd *cobra.Command, args []string) {
	// Inspect package and use type checker to infer imported types
	pkgs := loadPackages(args)

	// Get the package of the file with go:generate comment
	destPkg = os.Getenv("GOPACKAGE")
	if destPkg == "" {
		// TODO que tambien se pueda usar solo
		failErr(errors.New("this command should be called using go generate"))
	}

	for _, pkg := range pkgs {
		log.Println(pkg.Types.Path())
		log.Println(pkg.Types.Name())

		for _, name := range pkg.Types.Scope().Names() {
			object := getObject(pkg, name)
			if object != nil {
				log.Println(name)

				file := NewConditionsFile(
					destPkg,
					strings.ToLower(object.Name())+"_conditions.go",
				)

				err := file.AddConditionsFor(object)
				if err != nil {
					continue
				}

				err = file.Save()
				if err != nil {
					failErr(err)
				}
			}
		}
	}
}

func loadPackages(paths []string) []*packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes}
	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}

	// print compilation errors of source packages
	packages.PrintErrors(pkgs)

	return pkgs
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getObject(pkg *packages.Package, name string) types.Object {
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			name, pkg))
	}

	// Generate only if it is a declared type
	object, ok := obj.(*types.TypeName)
	if !ok {
		return nil
	}

	return object
}

// TODO add logs
