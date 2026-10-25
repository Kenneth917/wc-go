package features

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

// bin is the wc binary under test. Set WC_BIN to test a prebuilt binary, or
// WC_DOCKER_IMAGE to run the same scenarios against a container image.
var bin string

type world struct {
	dir    string
	output string
	code   int
}

func (w *world) aFileWithContent(name string, content *godog.DocString) error {
	return os.WriteFile(filepath.Join(w.dir, name), []byte(content.Content+"\n"), 0o644)
}

func (w *world) iRun(command string) error { return w.run(command, "") }

func (w *world) iRunWithStdin(command string, stdin *godog.DocString) error {
	return w.run(command, stdin.Content+"\n")
}

func (w *world) run(command, stdin string) error {
	args := strings.Fields(command)
	if len(args) == 0 || args[0] != "wc" {
		return fmt.Errorf("command must start with wc, got %q", command)
	}
	var cmd *exec.Cmd
	if img := os.Getenv("WC_DOCKER_IMAGE"); img != "" {
		dockerArgs := append([]string{"run", "-i", "--rm", "-v", w.dir + ":/work", "-w", "/work", img}, args[1:]...)
		cmd = exec.Command("docker", dockerArgs...)
	} else {
		cmd = exec.Command(bin, args[1:]...)
		cmd.Dir = w.dir
	}
	cmd.Stdin = strings.NewReader(stdin)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	w.output = out.String()
	w.code = 0
	if ee, ok := err.(*exec.ExitError); ok {
		w.code = ee.ExitCode()
	} else if err != nil {
		return err
	}
	return nil
}

// normalize collapses runs of spaces so scenarios don't depend on column padding.
func normalize(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, l := range lines {
		lines[i] = strings.Join(strings.Fields(l), " ")
	}
	return strings.Join(lines, "\n")
}

func (w *world) theOutputShouldBe(want *godog.DocString) error {
	if got := normalize(w.output); got != normalize(want.Content) {
		return fmt.Errorf("output mismatch\n got: %q\nwant: %q", got, normalize(want.Content))
	}
	return nil
}

func (w *world) theExitCodeShouldBe(code int) error {
	if w.code != code {
		return fmt.Errorf("exit code %d, want %d (output: %q)", w.code, code, w.output)
	}
	return nil
}

func initScenario(sc *godog.ScenarioContext) {
	w := &world{}
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		dir, err := os.MkdirTemp("", "wc-feature-")
		*w = world{dir: dir}
		return ctx, err
	})
	sc.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return ctx, os.RemoveAll(w.dir)
	})
	sc.Step(`^a file "([^"]*)" with content:$`, w.aFileWithContent)
	sc.Step(`^I run "([^"]*)"$`, w.iRun)
	sc.Step(`^I run "([^"]*)" with stdin:$`, w.iRunWithStdin)
	sc.Step(`^the output should be:$`, w.theOutputShouldBe)
	sc.Step(`^the exit code should be (\d+)$`, w.theExitCodeShouldBe)
}

func TestFeatures(t *testing.T) {
	if os.Getenv("WC_DOCKER_IMAGE") == "" {
		bin = os.Getenv("WC_BIN")
		if bin == "" {
			bin = filepath.Join(t.TempDir(), "wc")
			if out, err := exec.Command("go", "build", "-o", bin, "../cmd/wc").CombinedOutput(); err != nil {
				t.Fatalf("build failed: %v\n%s", err, out)
			}
		}
	}
	suite := godog.TestSuite{
		ScenarioInitializer: initScenario,
		Options:             &godog.Options{Format: "pretty", Paths: []string{"."}, TestingT: t},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
