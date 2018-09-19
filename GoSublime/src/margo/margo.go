package margo

import (
	"time"

	"margo.sh/golang"
	"margo.sh/mg"
)

// Margo is the entry-point to margo
func Margo(m mg.Args) {
	m.Use(
		&mg.MOTD{
			// Interval, if set, specifies how often to automatically fetch messages from Endpoint
			// Interval: 3600e9, // automatically fetch updates every hour
		},
		&golang.Gocode{
			Source:          false,
			ShowFuncParams:  true,
			ProposeTests:    false,
			ProposeBuiltins: true,
		},

		mg.NewReducer(func(mx *mg.Ctx) *mg.State {
			return mx.SetConfig(mx.Config.EnabledForLangs(
				mg.AllLangs,
			))
		}),
		&golang.SyntaxCheck{},
		golang.GoImports,
		golang.GoInstallDiscardBinaries("-i"),
		// golang.GoTest("-race"),
		golang.GoTest(),

		// run gometalinter on save
		// &golang.Linter{Name: "gometalinter", Args: []string{
		// 	"--disable=gas",
		// 	"--fast",
		// }},

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
		&golang.Guru{},
		&DayTimeStatus{},
		&golang.GoCmd{},
		&golang.GocodeCalltips{
			Source: false,
		},

		// Add user commands for running tests and benchmarks
		// gs: this adds support for the tests command palette `ctrl+.`,`ctrl+t` or `cmd+.`,`cmd+t`
		&golang.TestCmds{
			// additional args to add to the command when running tests and examples
			TestArgs: []string{},

			// additional args to add to the command when running benchmarks
			BenchArgs: []string{"-benchmem"},
		},

		// golang.GoVet(),
	)
}

// DayTimeStatus adds the current day and time to the status bar
type DayTimeStatus struct {
	mg.ReducerType
}

func (dts DayTimeStatus) ReducerMount(mx *mg.Ctx) {
	// kick off the ticker when we start
	dispatch := mx.Store.Dispatch
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			dispatch(mg.Render)
		}
	}()
}

func (dts DayTimeStatus) Reduce(mx *mg.Ctx) *mg.State {
	// we always want to render the time
	// otherwise it will sometimes disappear from the status bar
	now := time.Now()
	format := "Mon, 15:04"
	if now.Second()%2 == 0 {
		format = "Mon, 15 04"
	}
	return mx.AddStatus(now.Format(format))
}

// MySnippets is a slice of functions returning our own snippets
var MySnippets = golang.SnippetFuncs(
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
			{
				Query: "tcases",
				Title: "test cases",
				Src: `tcs := []struct {
    name string
    $1
}{}
for _, tc := range tcs {
    t.Run(tc.name, func(t *testing.T) {
    })
}`,
			},
		}
		//
	},
)

var qtSnippets = golang.SnippetFuncs(
	func(cx *golang.CompletionCtx) []mg.Completion {
		return []mg.Completion{
			{
				Query: "QtqmlResource",
				Title: "Load a qml resource",
				Src: `
    core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
    core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)

    gui.NewQGuiApplication(len(os.Args), os.Args)
    if quickcontrols2.QQuickStyle_Name() == "" {
        quickcontrols2.QQuickStyle_SetStyle("Material")
    }

    var engine = qml.NewQQmlApplicationEngine(nil)
    engine.Load(core.NewQUrl3("qrc:/${0:filepath}.qml", 0))
    gui.QGuiApplication_Exec()
`,
			},
			{
				Query: "QtuiResource",
				Title: "Load a ui resource",
				Src: `
    core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)
    core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

    widgets.NewQApplication(len(os.Args), os.Args)
    window := widgets.NewQMainWindow(nil, 0)
    dialog, err := qtlib.LoadResource(window, "./qml/${0:filepath}.ui")
    if err != nil {
        ${1:return err}
    }
    window.SetupUi(dialog)

    dialog.Show()
    widgets.QApplication_Exec()
                `,
			},
			{
				Query: "QtuiFindElement",
				Title: "Find and create an element from a widget",
				Src: `
    ${1:name} := widgets.NewQ${2:SomethingPointer}(
        ${3:parent}.FindChild("${1:name}", core.Qt__FindChildrenRecursively).Pointer(),
    )
    $4
                `,
			},
		}
	},
)
