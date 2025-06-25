////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"errors"
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

func TestCheckErr_DefaultBehavior(t *testing.T) {
	// override exitFunc to capture the code
	origExit := exitFunc
	var code int
	exitFunc = func(c int) { code = c }
	defer func() { exitFunc = origExit }()

	// capture output
	buf := &bytes.Buffer{}

	// invoke CheckErr
	CheckErr(
		errors.New("boom"),
		WithWriter(buf),
	)

	// verify exit code
	if code != 1 {
		t.Errorf("exit code = %d; want 1", code)
	}

	out := stripANSI(buf.String())
	for _, want := range []string{
		"boom",          // original error
		"check error",   // default op
		"runtime_error", // default category
		"severity",      // default detail key
		"critical",      // default detail value
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q missing %q", out, want)
		}
	}
}

func TestCheckErr_WithOverrides(t *testing.T) {
	// override exitFunc
	origExit := exitFunc
	var code int
	exitFunc = func(c int) { code = c }
	defer func() { exitFunc = origExit }()

	// capture output
	buf := &bytes.Buffer{}

	// call with every override
	CheckErr(
		errors.New("fail"),
		WithWriter(buf),
		WithExitCode(42),
		WithOp("load-conf"),
		WithCategory("cfg_err"),
		WithMessage("couldn't load config"),
		WithDetails(map[string]any{"path": "/etc/app.cfg"}),
	)

	// verify exit code
	if code != 42 {
		t.Errorf("exit code = %d; want 42", code)
	}

	out := stripANSI(buf.String())
	wantLines := []string{
		"Op       load-conf,",
		"Message  couldn't load config,",
		"Err      fail,",
		"path     /etc/app.cfg,",
		"Category cfg_err,",
	}
	for _, want := range wantLines {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got:\n%s", want, out)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
