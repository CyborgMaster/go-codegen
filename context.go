package codegen

import (
	"go/token"
	"go/types"
	"log"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/jinzhu/inflection"
)

// Context represents the context in which a code generation operation is run.
type Context struct {
	Dir         string
	SearchPaths []string
	Fset        *token.FileSet
	Templates   map[string]*template.Template
	PackageName string
	Imports     map[string]bool
	Results     map[string]string
}

// NewContext initializes a new code generation context.
func NewContext(dir string, searchPaths []string) (*Context, error) {
	result := &Context{
		Dir:         dir,
		SearchPaths: searchPaths,
		Fset:        token.NewFileSet(),
		PackageName: "main", // default to main
		Results:     map[string]string{},
		Templates:   map[string]*template.Template{},
		Imports:     map[string]bool{},
	}
	return result, result.Populate()
}

// Populate fills in the rest of the context based upon the context's
// config.
func (ctx *Context) Populate() error {
	for _, dir := range ctx.SearchPaths {
		err := ctx.searchDir(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

var templateFunctions template.FuncMap

func init() {
	templateFunctions = sprig.TxtFuncMap()
	templateFunctions["singular"] = inflection.Singular
	templateFunctions["plural"] = inflection.Plural
	templateFunctions["typeName"] = types.ExprString
}

func (ctx *Context) searchDir(dir string) error {
	// search directory for every template in the package
	pat := path.Join(dir, "*.tmpl")
	paths, err := filepath.Glob(pat)
	dirName := path.Base(dir)

	if err != nil {
		return err
	}

	for _, p := range paths {
		base := path.Base(p)
		name := base[:len(base)-len(".tmpl")]

		t, err := template.New(base).Funcs(templateFunctions).ParseFiles(p)
		if err != nil {
			return err
		}

		// Add the template with and without the package name
		ctx.Templates[name] = t
		ctx.Templates[dirName+"."+name] = t
	}

	log.Printf("found %d templates in %s", len(paths), dir)

	return nil
}