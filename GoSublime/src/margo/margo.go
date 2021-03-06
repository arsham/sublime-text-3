package margo

import (
	"time"

	"margo.sh/golang"
	"margo.sh/mg"
)

// Margo is the entry-point to margo
func Margo(m mg.Args) {
	m.Use(
		&mg.MOTD{},
		&golang.Gocode{
			ShowFuncParams:      true,
			ShowFuncResultNames: true,
			ProposeBuiltins:     true,
		},
		&golang.MarGocodeCtl{
			ImporterMode: golang.SrcImporterWithFallback,
		},
		&golang.GocodeCalltips{},
		mg.NewReducer(func(mx *mg.Ctx) *mg.State {
			return mx.SetConfig(mx.Config.EnabledForLangs(
				mg.AllLangs,
			))
		}),

		&golang.SyntaxCheck{},
		golang.GoImports,
		// golang.GoInstallDiscardBinaries("-i"),
		// golang.GoTest("-race"),
		golang.GoTest("-short"),

		&golang.Linter{Name: "golangci-lint", Label: "golangci", Args: []string{
			"run",
			"--fast",
			"--enable=prealloc",
			"--enable=gosimple",
			"--enable=staticcheck",
			"--enable=unused",
			"--enable=gocritic",
			"--enable=unparam",
			"--enable=interfacer",
			"--enable=unconvert",
			"--enable=nakedret",
			"--enable=maligned",
			"--skip-dirs=vendor",
			"--tests=false",
			// "--new-from-rev=HEAD~1",
		}},
		&golang.Linter{Name: "golint", Label: "Go/Lint"},

		golang.Snippets,
		MySnippets,
		&DayTimeStatus{},
		&golang.GoCmd{},
		&golang.Guru{},

		// Add user commands for running tests and benchmarks
		&golang.TestCmds{
			// additional args to add to the command when running tests and examples
			TestArgs: []string{},

			// additional args to add to the command when running benchmarks
			BenchArgs: []string{"-benchmem"},
		},
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
		if !cx.Scope.Is(golang.BlockScope) || !cx.IsTestFile {
			return nil
		}
		return []mg.Completion{
			{
				Query: "test error",
				Title: "t.Error() condition",
				Src: `
	if $1 {
		t.Error("$2: $3 = ($4); want ($5)")
	}`,
			},
			{
				Query: "test errorf",
				Title: "t.Errorf() condition",
				Src: `
	if $1 {
		t.Errorf("$2: $3 = ($4); want ($5)", $6)
	}`,
			},
			{
				Query: "test cases",
				Title: "Test Cases",
				Src: `
	tcs := map[string]struct {
		$1
	}{}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			$2
		})
	}`,
			},
			{
				Query: "patch method",
				Title: "Monkey path a method",
				Src: `
	monkey.PatchInstanceMethod(reflect.TypeOf(${1:instance}), "${2:method name}", func(${3:receiver}, ${4:args}) ${5:return values} {
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(${1:instance}), "${2:method name}")
	`,
			},
			{
				Query: "patch func",
				Title: "Monkey path a function",
				Src: `
	monkey.Patch(${1:time.Now}, ${2:func() time.Time} {
	})
	defer monkey.Unpatch(${1:time.Now})
	`,
			},
			{
				Query: "db mock",
				Title: "Get db mock",
				Src: `
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}()
	$1`,
			},
		}
	},
	func(cx *golang.CompletionCtx) []mg.Completion {
		if !cx.IsTestFile {
			return nil
		}
		return []mg.Completion{
			{
				Query: "random generator",
				Title: "random generator",
				Src: `
var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(count int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, count)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
`,
			},
		}
	},

	func(cx *golang.CompletionCtx) []mg.Completion {
		if !cx.Scope.Is(golang.BlockScope) {
			return nil
		}
		return []mg.Completion{
			{
				Query: "ppp",
				Title: "pprint",
				Src:   "pp.Println($0)",
			},
		}
	},
	func(cx *golang.CompletionCtx) []mg.Completion {
		return []mg.Completion{
			{
				Query: "fuzz function",
				Title: "Fuzz Function",
				Src: `
func Fuzz${1:funcName}(data []byte) int {
    if !${2:notInteresting} {
        return -1
    }
    if _, err := ${1:funcName}(p); err != nil {
        return 0
    }
    return 1
}
`,
			},
		}
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
