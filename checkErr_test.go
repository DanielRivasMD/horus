////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestRegisterErrorAndGetErrorRegistry(t *testing.T) {
	// reset the global registry
	errorTypeRegistry = make(map[string]int)

	// nil error should be ignored
	RegisterError(nil)
	if len(errorTypeRegistry) != 0 {
		t.Fatalf("expected empty registry after RegisterError(nil), got %v", errorTypeRegistry)
	}

	// plain error → "unknown"
	RegisterError(errors.New("plain"))
	reg := GetErrorRegistry()
	if got, want := reg["unknown"], 1; got != want {
		t.Errorf("unknown count = %d; want %d", got, want)
	}

	// categorized Herror → its Category
	he := NewCategorizedHerror("op", "mycat", "msg", errors.New("inner"), nil)
	RegisterError(he)
	reg = GetErrorRegistry()
	if got, want := reg["mycat"], 1; got != want {
		t.Errorf("mycat count = %d; want %d", got, want)
	}

	// ensure GetErrorRegistry returns a copy
	reg["mycat"] = 42
	reg2 := GetErrorRegistry()
	if got, want := reg2["mycat"], 1; got != want {
		t.Errorf("registry was mutated externally; got %d, want %d", got, want)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// captureOutput temporarily redirects stdout and returns what was printed.
func captureOutput_checkErr(f func()) (string, int) {
	// override exitFunc to capture code
	var code int
	exitFunc = func(c int) { code = c }

	// capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// run
	f()

	// restore stdout
	w.Close()
	os.Stdout = old

	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	return buf.String(), code
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestCheckErr_DefaultBehavior(t *testing.T) {
	out, code := captureOutput_checkErr(func() {
		CheckErr(errors.New("boom"))
	})

	if code != 1 {
		t.Errorf("exit code = %d; want 1", code)
	}
	if !strings.Contains(out, "boom") {
		t.Errorf("output %q does not mention the original error", out)
	}
	if !strings.Contains(out, "check error") {
		t.Errorf("output %q does not mention default op name", out)
	}
	if !strings.Contains(out, "runtime_error") {
		t.Errorf("output %q does not mention default category", out)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestCheckErr_WithOverrides(t *testing.T) {
	out, code := captureOutput_checkErr(func() {
		CheckErr(
			errors.New("fail"),
			WithOp("load-conf"),
			WithCategory("cfg_err"),
			WithMessage("couldn't load config"),
			WithDetails(map[string]any{"path": "/etc/app.cfg"}),
		)
	})

	if code != 1 {
		t.Errorf("exit code = %d; want 1", code)
	}
	for _, want := range []string{
		"load-conf",
		"cfg_err",
		"couldn't load config",
		"path",
		"/etc/app.cfg",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("override %q missing in output: %q", want, out)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
