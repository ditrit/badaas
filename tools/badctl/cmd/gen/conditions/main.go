package conditions

import (
	"errors"
	"fmt"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
	"github.com/ditrit/verdeter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"golang.org/x/tools/go/packages"
)

var GenConditionsCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "conditions",
	Short: "Generate conditions to query your objects using BaDORM",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateConditions,
	Args:  cobra.MinimumNArgs(1),
})

const DestPackageKey = "dest_package"

func init() {
	err := GenConditionsCmd.LKey(
		DestPackageKey, verdeter.IsStr, "d",
		"Destination package (not used if ran with go generate)",
	)
	if err != nil {
		cmderrors.FailErr(err)
	}
}

func generateConditions(cmd *cobra.Command, args []string) {
	// Inspect package and use type checker to infer imported types
	pkgs := loadPackages(args)

	// Get the package of the file with go:generate comment
	destPkg := os.Getenv("GOPACKAGE")
	if destPkg == "" {
		destPkg = viper.GetString(DestPackageKey)
		if destPkg == "" {
			cmderrors.FailErr(errors.New("config --dest_package or use go generate"))
		}
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
					cmderrors.FailErr(err)
				}
			}
		}
	}
}

func loadPackages(paths []string) []*packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes}
	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		cmderrors.FailErr(fmt.Errorf("loading packages for inspection: %v", err))
	}

	// print compilation errors of source packages
	packages.PrintErrors(pkgs)

	return pkgs
}

func getObject(pkg *packages.Package, name string) types.Object {
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		cmderrors.FailErr(fmt.Errorf("%s not found in declared types of %s",
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
