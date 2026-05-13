package main

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzeFilesReportsFunctionAndFileMetrics(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Plain() {}

func Risky(a, b bool, xs []int) int {
	if a && b {
		return 1
	}
	for _, x := range xs {
		if x > 0 {
			return x
		}
	}
	return 0
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	report, err := analyzeFiles([]string{path})
	if err != nil {
		t.Fatalf("analyze files: %v", err)
	}
	if len(report.functions) != 2 {
		t.Fatalf("functions = %d", len(report.functions))
	}
	risky := report.functions[1]
	if risky.name != "Risky" {
		t.Fatalf("function name = %q", risky.name)
	}
	if risky.cyclo != 5 {
		t.Fatalf("cyclo = %d", risky.cyclo)
	}
	if risky.cognitive < 5 {
		t.Fatalf("cognitive = %d, want nested control flow to be counted", risky.cognitive)
	}
	if risky.maintainabilityIndex <= 0 || risky.lines == 0 || risky.halsteadVolume <= 0 {
		t.Fatalf("function metrics = %+v", risky)
	}
	if len(report.files) != 1 || report.files[0].maintainabilityIndex <= 0 {
		t.Fatalf("file metrics = %+v", report.files)
	}
}

func TestHalsteadTokenHelpers(t *testing.T) {
	if !skipHalsteadToken(token.COMMENT) {
		t.Fatalf("comment token should be skipped")
	}
	if !skipHalsteadToken(token.SEMICOLON) {
		t.Fatalf("semicolon token should be skipped")
	}
	if skipHalsteadToken(token.IDENT) {
		t.Fatalf("identifier token should be counted")
	}
	if key := halsteadTokenKey(token.IDENT, "name"); key != "name" {
		t.Fatalf("literal key = %q, want name", key)
	}
	if key := halsteadTokenKey(token.ADD, ""); key != "+" {
		t.Fatalf("operator key = %q, want +", key)
	}
}

func TestAnalyzePathsSkipsTestsAndFormatsGocycloRows(t *testing.T) {
	dir := t.TempDir()
	production := filepath.Join(dir, "sample.go")
	testFile := filepath.Join(dir, "sample_test.go")
	if err := os.WriteFile(production, []byte("package sample\nfunc Build() {}\n"), 0o600); err != nil {
		t.Fatalf("write production: %v", err)
	}
	if err := os.WriteFile(testFile, []byte("package sample\nfunc TestBuild() {}\n"), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
	report, err := analyzePaths([]string{dir})
	if err != nil {
		t.Fatalf("analyze paths: %v", err)
	}
	if len(report.functions) != 1 || report.functions[0].name != "Build" {
		t.Fatalf("functions = %+v", report.functions)
	}

	var out bytes.Buffer
	failed := printReport(&out, report, options{gocyclo: true})
	if failed {
		t.Fatalf("gocyclo report failed unexpectedly")
	}
	output := out.String()
	if !strings.Contains(output, "1 sample Build ") || !strings.Contains(output, "sample.go:2:1") {
		t.Fatalf("gocyclo output = %q", output)
	}
}

func TestCollectGoFileHandlesWalkErrorsAndSkippedDirs(t *testing.T) {
	walkErr := errors.New("walk failed")
	if err := collectGoFile("broken", nil, walkErr, nil); !errors.Is(err, walkErr) {
		t.Fatalf("walk error = %v, want %v", err, walkErr)
	}

	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.Mkdir(gitDir, 0o700); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	pkgDir := filepath.Join(dir, "pkg")
	if err := os.Mkdir(pkgDir, 0o700); err != nil {
		t.Fatalf("mkdir pkg: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	for _, entry := range entries {
		var files []string
		err := collectGoFile(filepath.Join(dir, entry.Name()), entry, nil, &files)
		if entry.Name() == ".git" {
			if !errors.Is(err, filepath.SkipDir) {
				t.Fatalf("skip dir error = %v, want filepath.SkipDir", err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("collect non-skipped dir: %v", err)
		}
		if len(files) != 0 {
			t.Fatalf("files = %+v, want none for directory", files)
		}
	}
}

func TestRunHelpListsMaintainabilityBaselineFlags(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-h"}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("exit = %d, want help usage exit; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	help := stderr.String()
	for _, want := range []string{
		"-function-mi-baseline",
		"-write-function-mi-baseline",
		"-mi-baseline",
		"-write-mi-baseline",
		"-fail-only",
	} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q: %s", want, help)
		}
	}
}

func TestCognitiveComplexityDoesNotNestElseIfChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Classify(value int) string {
	if value == 0 {
		return "zero"
	} else if value == 1 {
		return "one"
	} else if value == 2 {
		return "two"
	}
	return "many"
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	report, err := analyzeFiles([]string{path})
	if err != nil {
		t.Fatalf("analyze files: %v", err)
	}
	if got := report.functions[0].cognitive; got != 3 {
		t.Fatalf("cognitive = %d, want flat else-if chain score 3", got)
	}
}

func TestCognitiveComplexityCoversSwitchSelectAndBlocks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Route(value any, ready <-chan struct{}) int {
	total := 0
	switch typed := value.(type) {
	case int:
		total = typed
	default:
		total = 1
	}
	select {
	case <-ready:
		total++
	default:
		total--
	}
	{
		total = total + 1
	}
	return total
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	report, err := analyzeFiles([]string{path})
	if err != nil {
		t.Fatalf("analyze files: %v", err)
	}
	if got := report.functions[0].cognitive; got < 4 {
		t.Fatalf("cognitive = %d, want switch/select/block paths counted", got)
	}
}

func TestBranchScoreCountsFlowBreakingBranches(t *testing.T) {
	for _, tok := range []token.Token{token.GOTO, token.BREAK, token.CONTINUE} {
		if got := branchScore(&ast.BranchStmt{Tok: tok}); got != 1 {
			t.Fatalf("branchScore(%s) = %d, want 1", tok, got)
		}
	}
	if got := branchScore(&ast.BranchStmt{Tok: token.FALLTHROUGH}); got != 0 {
		t.Fatalf("branchScore(FALLTHROUGH) = %d, want 0", got)
	}
}

func TestCommentLinesInRangeOnlyCountsContainedGroups(t *testing.T) {
	src := `package sample

// outside
func Build() {
	// inside
	println("ok")
}
`
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, "sample.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse sample: %v", err)
	}
	fn := parsed.Decls[0].(*ast.FuncDecl)
	if got := commentLinesInRange(parsed.Comments, fn.Pos(), fn.End()); got != 1 {
		t.Fatalf("commentLinesInRange = %d, want 1", got)
	}
}

func TestRunReportsThresholdFailures(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Risky(a, b bool) int {
	if a && b {
		return 1
	}
	return 0
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-cyclo-over", "1", "-cognitive-over", "1", "-function-mi-under", "99", "-mi-under", "99", path}, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("exit = %d, want threshold failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	errText := stderr.String()
	for _, want := range []string{"cyclomatic threshold 1 exceeded", "cognitive threshold 1 exceeded", "function maintainability index", "maintainability index", "lines=", "cyclo=", "halstead_volume="} {
		if !strings.Contains(errText, want) {
			t.Fatalf("stderr missing %q: %s", want, errText)
		}
	}
}

func TestRunFunctionMaintainabilityGatePassesSmallFunctions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Build(name string) string {
	if name == "" {
		return "default"
	}
	return name
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-function-mi-under", "65", path}, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("exit = %d, want pass; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "maintainability_index=") {
		t.Fatalf("stdout missing function maintainability index: %s", stdout.String())
	}
}

func TestRunFailOnlySuppressesPassingMetricRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Build(name string) string {
	if name == "" {
		return "default"
	}
	return name
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-fail-only", "-function-mi-under", "65", "-mi-under", "65", path}, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("exit = %d, want pass; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want quiet passing gate", stdout.String())
	}
}

func TestFunctionMaintainabilityBaselineAllowsHistoricalLowMI(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	baselinePath := filepath.Join(dir, "mi-baseline.json")
	src := `package sample

func Build(name string) string {
	if name == "" {
		return "default"
	}
	return name
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-function-mi-under", "80", "-write-function-mi-baseline", baselinePath, path}, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("write baseline exit = %d; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	exit = run([]string{"-function-mi-under", "80", "-function-mi-baseline", baselinePath, path}, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("baseline exit = %d, want pass; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}

	withNewLowMI := src + `
func NewRisk(flag bool) string {
	for i := 0; i < 3; i++ {
		if flag && i > 1 {
			return "flag"
		}
	}
	return "clear"
}
`
	if err := os.WriteFile(path, []byte(withNewLowMI), 0o600); err != nil {
		t.Fatalf("write changed sample: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	exit = run([]string{"-function-mi-under", "80", "-function-mi-baseline", baselinePath, path}, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("new low MI exit = %d, want failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "function MI baseline missing") {
		t.Fatalf("stderr missing baseline miss: %s", stderr.String())
	}
}

func TestFileMaintainabilityBaselineAllowsHistoricalLowMI(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	baselinePath := filepath.Join(dir, "file-mi-baseline.json")
	src := `package sample

func Build(name string) string {
	if name == "" {
		return "default"
	}
	return name
}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-mi-under", "90", "-write-mi-baseline", baselinePath, path}, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("write file baseline exit = %d; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	exit = run([]string{"-mi-under", "90", "-mi-baseline", baselinePath, path}, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("file baseline exit = %d, want pass; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}

	changedPath := filepath.Join(dir, "new_low.go")
	if err := os.WriteFile(changedPath, []byte(src), 0o600); err != nil {
		t.Fatalf("write new low file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	exit = run([]string{"-mi-under", "90", "-mi-baseline", baselinePath, dir}, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("new low file exit = %d, want failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "file MI baseline missing") {
		t.Fatalf("stderr missing file baseline miss: %s", stderr.String())
	}
}

func TestMaintainabilityBaselineComparisonUsesRoundedMetrics(t *testing.T) {
	var stderr bytes.Buffer
	opts := options{err: &stderr, functionMIUnder: 70, miUnder: 70}

	fn := functionMetric{
		file:                 "sample.go",
		name:                 "Build",
		maintainabilityIndex: 69.86,
	}
	functionBaseline := map[string]functionMIBaselineRecord{
		functionKey(fn): {MaintainabilityIndex: 69.9},
	}
	if functionMIThresholdFails(fn, opts, functionBaseline, true) {
		t.Fatalf("function baseline unexpectedly failed: %s", stderr.String())
	}

	file := fileMetric{
		file:                 "sample.go",
		maintainabilityIndex: 69.86,
	}
	fileBaseline := map[string]fileMIBaselineRecord{
		fileKey(file): {MaintainabilityIndex: 69.9},
	}
	if fileFails(file, opts, fileBaseline, true) {
		t.Fatalf("file baseline unexpectedly failed: %s", stderr.String())
	}
}

func TestMissingFunctionMaintainabilityBaselineFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Build() {}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-function-mi-under", "70", "-function-mi-baseline", filepath.Join(dir, "missing.json"), path}, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("exit = %d, want missing baseline failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "read function MI baseline") {
		t.Fatalf("stderr missing baseline read error: %s", stderr.String())
	}
}

func TestMissingFileMaintainabilityBaselineFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

func Build() {}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-mi-under", "70", "-mi-baseline", filepath.Join(dir, "missing.json"), path}, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("exit = %d, want missing baseline failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "read file MI baseline") {
		t.Fatalf("stderr missing file baseline read error: %s", stderr.String())
	}
}

func TestUnsupportedFileMaintainabilityBaselineSchemaFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	baselinePath := filepath.Join(dir, "file-mi-baseline.json")
	if err := os.WriteFile(path, []byte("package sample\nfunc Build() {}\n"), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	if err := os.WriteFile(baselinePath, []byte(`{"schema_version":"wrong","threshold":70,"files":[]}`), 0o600); err != nil {
		t.Fatalf("write baseline: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-mi-under", "70", "-mi-baseline", baselinePath, path}, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("exit = %d, want unsupported baseline failure; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "unsupported file MI baseline schema") {
		t.Fatalf("stderr missing unsupported schema error: %s", stderr.String())
	}
}

func TestRunReturnsUsageForBadFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exit := run([]string{"-unknown"}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("exit = %d, want usage error; stdout=%s stderr=%s", exit, stdout.String(), stderr.String())
	}
}

func TestReceiverNameVariants(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package sample

type Box[T any] struct{}
type Pair[A, B any] struct{}

func (b Box[T]) Value() {}
func (b *Box[T]) Pointer() {}
func (p Pair[A, B]) Both() {}
`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write sample: %v", err)
	}
	report, err := analyzeFiles([]string{path})
	if err != nil {
		t.Fatalf("analyze files: %v", err)
	}
	names := []string{report.functions[0].name, report.functions[1].name, report.functions[2].name}
	for _, want := range []string{"(Box).Value", "(Box).Pointer", "(Pair).Both"} {
		if !containsString(names, want) {
			t.Fatalf("function names = %v, want %s", names, want)
		}
	}
}

func TestIndexListReceiverInnerExpr(t *testing.T) {
	expr := &ast.IndexListExpr{X: ast.NewIdent("Pair")}
	inner, ok := indexListReceiverInnerExpr(expr)
	if !ok {
		t.Fatalf("index list receiver was not recognized")
	}
	if ident, ok := inner.(*ast.Ident); !ok || ident.Name != "Pair" {
		t.Fatalf("inner = %#v", inner)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
