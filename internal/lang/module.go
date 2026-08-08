package lang

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type projectConfig struct {
	Name  string            `json:"name"`
	Root  string            `json:"root"`
	Paths map[string]string `json:"paths"`
}

type projectContext struct {
	ConfigDir   string
	ProjectRoot string
	Aliases     map[string]string
}

type moduleUnit struct {
	File        string // absolute path
	Program     *Program
	ImportKeys  []string            // absolute dependency files, in source order
	ImportPos   map[string]Position // dependency abs path -> import path token position
	ImportNames map[string]string   // dependency abs path -> binding name (module namespace)
}

// moduleNamespace is the value bound by `import "…"` (SF-3.1). Exports are
// accessed with `.` (e.g. `json.parse`).
type moduleNamespace struct {
	name    string
	exports map[string]any
}

// RunFile discovers the project, loads the import graph for entryPath, and executes
// dependency modules before the entry module.
func RunFile(entryPath string, output io.Writer, options ExecuteOptions) *Diagnostic {
	absEntry, err := filepath.Abs(entryPath)
	if err != nil {
		return &Diagnostic{Code: "S031", Message: "Cannot resolve path `" + entryPath + "`.", File: entryPath}
	}
	ctx, d := discoverProject(absEntry)
	if d != nil {
		return d
	}
	units, order, d := loadModuleGraph(absEntry, ctx)
	if d != nil {
		return d
	}
	exports := map[string]map[string]any{}
	importSymbols := map[string]map[string]symbol{}
	for _, file := range order {
		if isStdModulePath(file) {
			mod, d := loadStdModule(stdImportFromPath(file))
			if d != nil {
				d.File = file
				return d
			}
			exports[file] = mod.exports
			importSymbols[file] = stdSymbols(mod)
			continue
		}
		unit := units[file]
		importedValues := map[string]any{}
		importedSyms := map[string]symbol{}
		for _, dep := range unit.ImportKeys {
			binding := unit.ImportNames[dep]
			if binding == "" {
				binding = moduleBindingNameFromPath(dep)
			}
			if _, exists := importedValues[binding]; exists {
				return &Diagnostic{
					Code:      "S034",
					Message:   "Imported module `" + binding + "` collides with another import.",
					File:      unit.File,
					Pos:       unit.ImportPos[dep],
					Fix:       "Import only one module that binds `" + binding + "`, or rename a project module path.",
					Expected:  "unique module binding `" + binding + "`",
					Actual:    "duplicate import binding",
					Operation: "import",
				}
			}
			importedValues[binding] = &moduleNamespace{name: binding, exports: exports[dep]}
			importedSyms[binding] = symbol{kind: "module", arity: -1}
		}
		if d := resolveModule(unit.Program, importedSyms); d != nil {
			return d
		}
		env, d := executeModule(unit.Program, output, options, importedValues)
		if d != nil {
			return d
		}
		moduleExports := map[string]any{}
		moduleSyms := map[string]symbol{}
		for _, name := range unit.Program.Exports {
			value, ok := env.get(name.Lexeme)
			if !ok {
				return &Diagnostic{
					Code:    "S036",
					Message: "Export `" + name.Lexeme + "` was not defined.",
					File:    unit.File,
					Pos:     name.Pos,
				}
			}
			moduleExports[name.Lexeme] = value
			moduleSyms[name.Lexeme] = symbolFromValue(value)
		}
		exports[file] = moduleExports
		importSymbols[file] = moduleSyms
	}
	return nil
}

func discoverProject(entryAbs string) (projectContext, *Diagnostic) {
	dir := filepath.Dir(entryAbs)
	for {
		configPath := filepath.Join(dir, "vol.config.json")
		if info, err := os.Stat(configPath); err == nil && !info.IsDir() {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return projectContext{}, &Diagnostic{Code: "S035", Message: "Cannot read `" + configPath + "`.", File: entryAbs}
			}
			var cfg projectConfig
			if err := json.Unmarshal(data, &cfg); err != nil {
				return projectContext{}, &Diagnostic{Code: "S035", Message: "Invalid `vol.config.json`: " + err.Error() + ".", File: configPath}
			}
			root := cfg.Root
			if root == "" {
				root = "."
			}
			projectRoot := filepath.Clean(filepath.Join(dir, root))
			aliases := map[string]string{}
			for key, value := range cfg.Paths {
				aliases[key] = value
			}
			return projectContext{ConfigDir: dir, ProjectRoot: projectRoot, Aliases: aliases}, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return projectContext{
		ConfigDir:   filepath.Dir(entryAbs),
		ProjectRoot: filepath.Dir(entryAbs),
		Aliases:     map[string]string{},
	}, nil
}

func loadModuleGraph(entryAbs string, ctx projectContext) (map[string]*moduleUnit, []string, *Diagnostic) {
	units := map[string]*moduleUnit{}
	var visiting []string
	visitingSet := map[string]bool{}
	var order []string

	var visit func(fileAbs string) *Diagnostic
	visit = func(fileAbs string) *Diagnostic {
		if units[fileAbs] != nil && containsString(order, fileAbs) {
			return nil
		}
		if isStdModulePath(fileAbs) && containsString(order, fileAbs) {
			return nil
		}
		if visitingSet[fileAbs] {
			cycle := append(append([]string{}, visiting...), fileAbs)
			return &Diagnostic{
				Code:    "S033",
				Message: "Import cycle: " + strings.Join(displayPaths(cycle, ctx.ProjectRoot), " -> ") + ".",
				File:    fileAbs,
			}
		}
		if isStdModulePath(fileAbs) {
			if _, d := loadStdModule(stdImportFromPath(fileAbs)); d != nil {
				d.File = fileAbs
				return d
			}
			if !containsString(order, fileAbs) {
				order = append(order, fileAbs)
			}
			units[fileAbs] = &moduleUnit{File: fileAbs, Program: &Program{File: fileAbs}, ImportPos: map[string]Position{}}
			return nil
		}
		visitingSet[fileAbs] = true
		visiting = append(visiting, fileAbs)
		source, err := os.ReadFile(fileAbs)
		if err != nil {
			return &Diagnostic{Code: "S031", Message: "Cannot read module `" + fileAbs + "`.", File: fileAbs}
		}
		program, d := Parse(fileAbs, string(source))
		if d != nil {
			return d
		}
		unit := &moduleUnit{
			File:        fileAbs,
			Program:     program,
			ImportPos:   map[string]Position{},
			ImportNames: map[string]string{},
		}
		seenDep := map[string]bool{}
		for _, statement := range program.Statements {
			imp, ok := statement.(*ImportStatement)
			if !ok {
				continue
			}
			resolved, d := resolveImportPath(imp.Path.Lexeme, imp.Path.Pos, fileAbs, ctx)
			if d != nil {
				return d
			}
			binding := moduleBindingName(imp.Path.Lexeme, resolved)
			if !seenDep[resolved] {
				seenDep[resolved] = true
				unit.ImportKeys = append(unit.ImportKeys, resolved)
				unit.ImportPos[resolved] = imp.Path.Pos
				unit.ImportNames[resolved] = binding
			}
			if d := visit(resolved); d != nil {
				return d
			}
		}
		units[fileAbs] = unit
		visiting = visiting[:len(visiting)-1]
		delete(visitingSet, fileAbs)
		if !containsString(order, fileAbs) {
			order = append(order, fileAbs)
		}
		return nil
	}
	if d := visit(entryAbs); d != nil {
		return nil, nil, d
	}
	return units, order, nil
}

func resolveImportPath(importPath string, pos Position, fromFile string, ctx projectContext) (string, *Diagnostic) {
	if isReservedStdImport(importPath) {
		if _, d := loadStdModule(importPath); d != nil {
			d.File = fromFile
			d.Pos = pos
			return "", d
		}
		return stdModulePath(importPath), nil
	}
	expanded := expandAlias(importPath, ctx.Aliases)
	if expanded == "" {
		return "", &Diagnostic{
			Code:    "S035",
			Message: "Unknown import alias in `" + importPath + "`.",
			File:    fromFile,
			Pos:     pos,
			Fix:     "Define the alias under `paths` in `vol.config.json`.",
		}
	}
	if strings.Contains(expanded, "..") {
		return "", &Diagnostic{
			Code:    "S032",
			Message: "Import path `" + importPath + "` must not contain `..`.",
			File:    fromFile,
			Pos:     pos,
		}
	}
	candidateFile := filepath.Clean(filepath.Join(ctx.ProjectRoot, expanded+".vol"))
	candidateMod := filepath.Clean(filepath.Join(ctx.ProjectRoot, expanded, "mod.vol"))
	var chosen string
	if info, err := os.Stat(candidateFile); err == nil && !info.IsDir() {
		chosen = candidateFile
	} else if info, err := os.Stat(candidateMod); err == nil && !info.IsDir() {
		chosen = candidateMod
	} else {
		return "", &Diagnostic{
			Code:    "S031",
			Message: "Module `" + importPath + "` not found (tried `" + relToRoot(candidateFile, ctx.ProjectRoot) + "` and `" + relToRoot(candidateMod, ctx.ProjectRoot) + "`).",
			File:    fromFile,
			Pos:     pos,
		}
	}
	if !underRoot(chosen, ctx.ProjectRoot) {
		return "", &Diagnostic{
			Code:    "S032",
			Message: "Import `" + importPath + "` escapes the project root.",
			File:    fromFile,
			Pos:     pos,
		}
	}
	return chosen, nil
}

func expandAlias(importPath string, aliases map[string]string) string {
	if !strings.HasPrefix(importPath, "@") {
		return importPath
	}
	keys := make([]string, 0, len(aliases))
	for key := range aliases {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })
	for _, key := range keys {
		if importPath == key {
			return aliases[key]
		}
		if strings.HasPrefix(importPath, key+"/") {
			return filepath.ToSlash(filepath.Join(aliases[key], strings.TrimPrefix(importPath, key+"/")))
		}
	}
	return ""
}

func underRoot(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func relToRoot(path, root string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

func displayPaths(paths []string, root string) []string {
	out := make([]string, len(paths))
	for i, path := range paths {
		out[i] = relToRoot(path, root)
	}
	return out
}

func containsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func symbolFromValue(value any) symbol {
	switch v := value.(type) {
	case *function:
		return symbol{kind: "function", arity: len(v.declaration.Parameters)}
	case *structType:
		return symbol{kind: "struct", arity: -1}
	case *builtinFunction:
		return symbol{kind: "builtin", arity: -1}
	case *moduleNamespace:
		return symbol{kind: "module", arity: -1}
	case *dictValue, *nativeHandle, *httpReplyValue:
		return symbol{kind: "value", arity: -1}
	default:
		return symbol{kind: "value", arity: -1}
	}
}

// moduleBindingName is the identifier installed by an import (path basename).
func moduleBindingName(importPath, resolvedAbs string) string {
	if strings.HasPrefix(importPath, "@std/") {
		return strings.TrimPrefix(importPath, "@std/")
	}
	if filepath.Base(resolvedAbs) == "mod.vol" {
		return filepath.Base(filepath.Dir(resolvedAbs))
	}
	name := importPath
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	if strings.HasPrefix(name, "@") {
		name = strings.TrimPrefix(name, "@")
	}
	return name
}

func moduleBindingNameFromPath(resolvedAbs string) string {
	if isStdModulePath(resolvedAbs) {
		return strings.TrimPrefix(stdImportFromPath(resolvedAbs), "@std/")
	}
	if filepath.Base(resolvedAbs) == "mod.vol" {
		return filepath.Base(filepath.Dir(resolvedAbs))
	}
	base := filepath.Base(resolvedAbs)
	return strings.TrimSuffix(base, ".vol")
}
