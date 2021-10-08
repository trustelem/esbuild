// This API exposes esbuild's two main operations: building and transforming.
// It's intended for integrating esbuild into other tools as a library.
//
// If you are just trying to run esbuild from Go without the overhead of
// creating a child process, there is also an API for the command-line
// interface itself: https://godoc.org/github.com/evanw/esbuild/pkg/cli.
//
// Build API
//
// This function runs an end-to-end build operation. It takes an array of file
// paths as entry points, parses them and all of their dependencies, and
// returns the output files to write to the file system. The available options
// roughly correspond to esbuild's command-line flags.
//
// Example usage:
//
//     package main
//
//     import (
//         "os"
//
//         "github.com/trustelem/esbuild/pkg/api"
//     )
//
//     func main() {
//         result := api.Build(api.BuildOptions{
//             EntryPoints: []string{"input.js"},
//             Outfile:     "output.js",
//             Bundle:      true,
//             Write:       true,
//             LogLevel:    api.LogLevelInfo,
//         })
//
//         if len(result.Errors) > 0 {
//             os.Exit(1)
//         }
//     }
//
// Transform API
//
// This function transforms a string of source code into JavaScript. It can be
// used to minify JavaScript, convert TypeScript/JSX to JavaScript, or convert
// newer JavaScript to older JavaScript. The available options roughly
// correspond to esbuild's command-line flags.
//
// Example usage:
//
//     package main
//
//     import (
//         "fmt"
//         "os"
//
//         "github.com/trustelem/esbuild/pkg/api"
//     )
//
//     func main() {
//         jsx := `
//             import * as React from 'react'
//             import * as ReactDOM from 'react-dom'
//
//             ReactDOM.render(
//                 <h1>Hello, world!</h1>,
//                 document.getElementById('root')
//             );
//         `
//
//         result := api.Transform(jsx, api.TransformOptions{
//             Loader: api.LoaderJSX,
//         })
//
//         fmt.Printf("%d errors and %d warnings\n",
//             len(result.Errors), len(result.Warnings))
//
//         os.Stdout.Write(result.Code)
//     }
//
package api

import "github.com/trustelem/esbuild/internal/fs"

type SourceMap uint8

const (
	SourceMapNone SourceMap = iota
	SourceMapInline
	SourceMapLinked
	SourceMapExternal
	SourceMapInlineAndExternal
)

type SourcesContent uint8

const (
	SourcesContentInclude SourcesContent = iota
	SourcesContentExclude
)

type LegalComments uint8

const (
	LegalCommentsDefault LegalComments = iota
	LegalCommentsNone
	LegalCommentsInline
	LegalCommentsEndOfFile
	LegalCommentsLinked
	LegalCommentsExternal
)

type JSXMode uint8

const (
	JSXModeTransform JSXMode = iota
	JSXModePreserve
)

type Target uint8

const (
	DefaultTarget Target = iota
	ESNext
	ES5
	ES2015
	ES2016
	ES2017
	ES2018
	ES2019
	ES2020
	ES2021
)

type Loader uint8

const (
	LoaderNone Loader = iota
	LoaderJS
	LoaderJSX
	LoaderTS
	LoaderTSX
	LoaderJSON
	LoaderText
	LoaderBase64
	LoaderDataURL
	LoaderFile
	LoaderBinary
	LoaderCSS
	LoaderDefault
)

type Platform uint8

const (
	PlatformBrowser Platform = iota
	PlatformNode
	PlatformNeutral
)

type Format uint8

const (
	FormatDefault Format = iota
	FormatIIFE
	FormatCommonJS
	FormatESModule
)

type EngineName uint8

const (
	EngineChrome EngineName = iota
	EngineEdge
	EngineFirefox
	EngineIOS
	EngineNode
	EngineSafari
)

type Engine struct {
	Name    EngineName
	Version string
}

type Location struct {
	File       string
	Namespace  string
	Line       int // 1-based
	Column     int // 0-based, in bytes
	Length     int // in bytes
	LineText   string
	Suggestion string
}

type Message struct {
	PluginName string
	Text       string
	Location   *Location
	Notes      []Note

	// Optional user-specified data that is passed through unmodified. You can
	// use this to stash the original error, for example.
	Detail interface{}
}

type Note struct {
	Text     string
	Location *Location
}

type StderrColor uint8

const (
	ColorIfTerminal StderrColor = iota
	ColorNever
	ColorAlways
)

type LogLevel uint8

const (
	LogLevelSilent LogLevel = iota
	LogLevelVerbose
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

type Charset uint8

const (
	CharsetDefault Charset = iota
	CharsetASCII
	CharsetUTF8
)

type TreeShaking uint8

const (
	TreeShakingDefault TreeShaking = iota
	TreeShakingFalse
	TreeShakingTrue
)

////////////////////////////////////////////////////////////////////////////////
// Build API

type FsLike = fs.FsLike

type BuildOptions struct {
	Color    StderrColor `json:",omitempty"`
	LogLimit int         `json:",omitempty"`
	LogLevel LogLevel    `json:",omitempty"`

	Sourcemap      SourceMap      `json:",omitempty"`
	SourceRoot     string         `json:",omitempty"`
	SourcesContent SourcesContent `json:",omitempty"`

	Target  Target   `json:",omitempty"`
	Engines []Engine `json:",omitempty"`

	MinifyWhitespace  bool          `json:",omitempty"`
	MinifyIdentifiers bool          `json:",omitempty"`
	MinifySyntax      bool          `json:",omitempty"`
	Charset           Charset       `json:",omitempty"`
	TreeShaking       TreeShaking   `json:",omitempty"`
	IgnoreAnnotations bool          `json:",omitempty"`
	LegalComments     LegalComments `json:",omitempty"`

	JSXMode     JSXMode `json:",omitempty"`
	JSXFactory  string  `json:",omitempty"`
	JSXFragment string  `json:",omitempty"`

	Define    map[string]string `json:",omitempty"`
	Pure      []string          `json:",omitempty"`
	KeepNames bool              `json:",omitempty"`

	GlobalName        string            `json:",omitempty"`
	Bundle            bool              `json:",omitempty"`
	PreserveSymlinks  bool              `json:",omitempty"`
	Splitting         bool              `json:",omitempty"`
	Outfile           string            `json:",omitempty"`
	Metafile          bool              `json:",omitempty"`
	Outdir            string            `json:",omitempty"`
	Outbase           string            `json:",omitempty"`
	AbsWorkingDir     string            `json:",omitempty"`
	Platform          Platform          `json:",omitempty"`
	Format            Format            `json:",omitempty"`
	External          []string          `json:",omitempty"`
	MainFields        []string          `json:",omitempty"`
	Conditions        []string          `json:",omitempty"` // For the "exports" field in "package.json"
	Loader            map[string]Loader `json:",omitempty"`
	ResolveExtensions []string          `json:",omitempty"`
	Tsconfig          string            `json:",omitempty"`
	OutExtensions     map[string]string `json:",omitempty"`
	PublicPath        string            `json:",omitempty"`
	Inject            []string          `json:",omitempty"`
	Banner            map[string]string `json:",omitempty"`
	Footer            map[string]string `json:",omitempty"`
	NodePaths         []string          `json:",omitempty"` // The "NODE_PATH" variable from Node.js

	EntryNames string `json:",omitempty"`
	ChunkNames string `json:",omitempty"`
	AssetNames string `json:",omitempty"`

	EntryPoints         []string     `json:",omitempty"`
	EntryPointsAdvanced []EntryPoint `json:",omitempty"`

	Stdin          *StdinOptions `json:",omitempty"`
	FS             FsLike        `json:",omitempty"`
	Write          bool          `json:",omitempty"`
	AllowOverwrite bool          `json:",omitempty"`
	Incremental    bool          `json:",omitempty"`
	Plugins        []Plugin      `json:",omitempty"`

	Watch *WatchMode `json:",omitempty"`
}

type EntryPoint struct {
	InputPath  string
	OutputPath string
}

type WatchMode struct {
	OnRebuild func(BuildResult)
}

type StdinOptions struct {
	Contents   string
	ResolveDir string
	Sourcefile string
	Loader     Loader
}

type BuildResult struct {
	Errors   []Message
	Warnings []Message

	OutputFiles []OutputFile
	Metafile    string

	Rebuild func() BuildResult // Only when "Incremental: true"
	Stop    func()             // Only when "Watch: true"
}

type OutputFile struct {
	Path     string
	Contents []byte
}

func Build(options BuildOptions) BuildResult {
	return buildImpl(options).result
}

////////////////////////////////////////////////////////////////////////////////
// Transform API

type TransformOptions struct {
	Color    StderrColor `json:",omitempty"`
	LogLimit int         `json:",omitempty"`
	LogLevel LogLevel    `json:",omitempty"`

	Sourcemap      SourceMap      `json:",omitempty"`
	SourceRoot     string         `json:",omitempty"`
	SourcesContent SourcesContent `json:",omitempty"`

	Target     Target   `json:",omitempty"`
	Format     Format   `json:",omitempty"`
	GlobalName string   `json:",omitempty"`
	Engines    []Engine `json:",omitempty"`

	MinifyWhitespace  bool          `json:",omitempty"`
	MinifyIdentifiers bool          `json:",omitempty"`
	MinifySyntax      bool          `json:",omitempty"`
	Charset           Charset       `json:",omitempty"`
	TreeShaking       TreeShaking   `json:",omitempty"`
	IgnoreAnnotations bool          `json:",omitempty"`
	LegalComments     LegalComments `json:",omitempty"`

	JSXMode     JSXMode `json:",omitempty"`
	JSXFactory  string  `json:",omitempty"`
	JSXFragment string  `json:",omitempty"`

	TsconfigRaw string `json:",omitempty"`
	Footer      string `json:",omitempty"`
	Banner      string `json:",omitempty"`

	Define    map[string]string `json:",omitempty"`
	Pure      []string          `json:",omitempty"`
	KeepNames bool              `json:",omitempty"`

	Sourcefile string `json:",omitempty"`
	Loader     Loader `json:",omitempty"`
}

type TransformResult struct {
	Errors   []Message
	Warnings []Message

	Code []byte
	Map  []byte
}

func Transform(input string, options TransformOptions) TransformResult {
	return transformImpl(input, options)
}

////////////////////////////////////////////////////////////////////////////////
// Serve API

type ServeOptions struct {
	Port      uint16
	Host      string
	Servedir  string
	OnRequest func(ServeOnRequestArgs)
}

type ServeOnRequestArgs struct {
	RemoteAddress string
	Method        string
	Path          string
	Status        int
	TimeInMS      int // The time to generate the response, not to send it
}

type ServeResult struct {
	Port uint16
	Host string
	Wait func() error
	Stop func()
}

func Serve(serveOptions ServeOptions, buildOptions BuildOptions) (ServeResult, error) {
	return serveImpl(serveOptions, buildOptions)
}

////////////////////////////////////////////////////////////////////////////////
// Plugin API

type SideEffects uint8

const (
	SideEffectsTrue SideEffects = iota
	SideEffectsFalse
)

type Plugin struct {
	Name  string
	Setup func(PluginBuild)
}

type PluginBuild struct {
	InitialOptions *BuildOptions
	OnStart        func(callback func() (OnStartResult, error))
	OnEnd          func(callback func(result *BuildResult))
	OnResolve      func(options OnResolveOptions, callback func(OnResolveArgs) (OnResolveResult, error))
	OnLoad         func(options OnLoadOptions, callback func(OnLoadArgs) (OnLoadResult, error))
}

type OnStartResult struct {
	Errors   []Message
	Warnings []Message
}

type OnResolveOptions struct {
	Filter    string
	Namespace string
}

type OnResolveArgs struct {
	Path       string
	Importer   string
	Namespace  string
	ResolveDir string
	Kind       ResolveKind
	PluginData interface{}
}

type OnResolveResult struct {
	PluginName string

	Errors   []Message
	Warnings []Message

	Path        string
	External    bool
	SideEffects SideEffects
	Namespace   string
	PluginData  interface{}

	WatchFiles []string
	WatchDirs  []string
}

type OnLoadOptions struct {
	Filter    string
	Namespace string
}

type OnLoadArgs struct {
	Path       string
	Namespace  string
	PluginData interface{}
}

type OnLoadResult struct {
	PluginName string

	Errors   []Message
	Warnings []Message

	Contents   *string
	ResolveDir string
	Loader     Loader
	PluginData interface{}

	WatchFiles []string
	WatchDirs  []string
}

type ResolveKind uint8

const (
	ResolveEntryPoint ResolveKind = iota
	ResolveJSImportStatement
	ResolveJSRequireCall
	ResolveJSDynamicImport
	ResolveJSRequireResolve
	ResolveCSSImportRule
	ResolveCSSURLToken
)

////////////////////////////////////////////////////////////////////////////////
// FormatMessages API

type MessageKind uint8

const (
	ErrorMessage MessageKind = iota
	WarningMessage
)

type FormatMessagesOptions struct {
	TerminalWidth int
	Kind          MessageKind
	Color         bool
}

func FormatMessages(msgs []Message, opts FormatMessagesOptions) []string {
	return formatMsgsImpl(msgs, opts)
}

////////////////////////////////////////////////////////////////////////////////
// AnalyzeMetafile API

type AnalyzeMetafileOptions struct {
	Color   bool
	Verbose bool
}

func AnalyzeMetafile(metafile string, opts AnalyzeMetafileOptions) string {
	return analyzeMetafileImpl(metafile, opts)
}
