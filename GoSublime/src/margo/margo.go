package margo

import (
	"time"

	"disposa.blue/margo/golang"
	"disposa.blue/margo/mg"
)

// Margo is the entry-point to margo
func Margo(ma mg.Args) {
	ma.Store.Use(

		// use gocode for autocompletion
		&golang.Gocode{
			// autocompete packages that are not yet imported
			// this goes well with GoImports
			UnimportedPackages: true,

			// show the function parameters. this can take up a lot of space
			ShowFuncParams: true,
		},

		&golang.SyntaxCheck{},
		golang.GoImports,
		golang.GoInstall("-i"),
		// golang.GoTest("-race"),
		golang.GoTest(),

		// run `golint` on save
		&golang.Linter{Label: "Go/Lint", Name: "golint"},
		&golang.Linter{Label: "Go/GoConst", Name: "goconst", Args: []string{"."}},
		&golang.Linter{Label: "Go/UsedExports", Name: "usedexports", Args: []string{"."}},
		&golang.Linter{Label: "Go/IneffAssign", Name: "ineffassign", Args: []string{"-n", "."}},
		// &golang.Linter{Label: "Go/Cyclo", Name: "cyclo", Args: []string{"--max-complexity", "15", "."}},
		// &golang.Linter{Label: "Go/Interfacer", Name: "interfacer", Args: []string{"./..."}},
		// &golang.Linter{Label: "Go/ErrorCheck", Name: "errcheck", Args: []string{"-ignoretests", "."}},
		// &golang.Linter{Label: "Go/Unconver", Name: "unconvert", Args: []string{"."}},

		golang.Snippets,
		MySnippets,
	)
}

// DayTimeStatus adds the current day and time to the status bar
var DayTimeStatus = mg.Reduce(func(mx *mg.Ctx) *mg.State {
	if _, ok := mx.Action.(mg.Started); ok {
		dispatch := mx.Store.Dispatch
		// kick off the ticker when we start
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			for range ticker.C {
				dispatch(mg.Render)
			}
		}()
	}

	// we always want to render the time
	// otherwise it will sometimes disappear from the status bar
	now := time.Now()
	format := "Mon, 15:04"
	if now.Second()%2 == 0 {
		format = "Mon, 15 04"
	}
	return mx.AddStatus(now.Format(format))
})

// MySnippets is a slice of functions returning our own snippets
var MySnippets = golang.SnippetFuncs{
	func(cx *golang.CompletionCtx) []mg.Completion {
		// if we're not in a block (i.e. function), do nothing
		if !cx.Scope.Is(golang.BlockScope) {
			return nil
		}

		return []mg.Completion{
			{
				Query: "if err",
				Title: "err != nil { return }",
				Src:   "if ${1:err} != nil {\n\treturn $0\n}",
			},
		}
		//
	},
	func(cx *golang.CompletionCtx) []mg.Completion {
		// if we're not in a block (i.e. function), do nothing
		if !cx.Scope.Is(golang.BlockScope) {
			return nil
		}

		return []mg.Completion{
			{
				Query: "ppp",
				Title: "pprint",
				Src:   "pp.Println($0)",
			},
			{
				Query: "terrorf",
				Title: "t.Errorf() condition",
				Src: `if $1 {
	t.Errorf("$2: $3 = ($4); want ($5)", $6)
}
`,
			},
			{
				Query: "terror",
				Title: "t.Error() condition",
				Src: `if $1 {
	t.Error("$2: $3 = ($4); want ($5)")
}
`,
			},
		}
		//
	},
}
