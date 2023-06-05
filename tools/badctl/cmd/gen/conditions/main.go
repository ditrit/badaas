package conditions

import (
	"errors"
	"fmt"
	"go/types"
	"os"
	"strings"

	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
	"github.com/ditrit/badaas/tools/badctl/cmd/log"
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

// GenConditionsCmd Run func
func generateConditions(cmd *cobra.Command, args []string) {
	log.SetLevel()
	// Inspect package and use type checker to infer imported types
	pkgs := loadPackages(args)

	// Get the package of the file with go:generate comment or in command params
	destPkg := os.Getenv("GOPACKAGE")
	if destPkg == "" {
		destPkg = viper.GetString(DestPackageKey)
		if destPkg == "" {
			cmderrors.FailErr(errors.New("config --dest_package or use go generate"))
		}
	}

	// Generate conditions for each package
	for _, pkg := range pkgs {
		generateConditionsForPkg(destPkg, pkg)
	}
}

// Generates a file with conditions for each BaDORM model in the package
func generateConditionsForPkg(destPkg string, pkg *packages.Package) {
	log.Logger.Infof("Generating conditions for types in package %q", pkg.Types.Name())

	for _, name := range pkg.Types.Scope().Names() {
		object := getObject(pkg, name)
		if object != nil {
			file := NewConditionsFile(
				destPkg,
				strings.ToLower(object.Name())+"_conditions.go",
			)

			err := file.AddConditionsFor(object)
			if err != nil {
				// object is not a BaDORM model, do not generate conditions
				continue
			}

			err = file.Save()
			if err != nil {
				cmderrors.FailErr(err)
			}
		}
	}
}

// Load package information from paths
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

// Get object by name in the package
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
