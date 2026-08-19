package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"help"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "queue-mm1k") {
		t.Error("help output missing program name")
	}
}

func TestRunAnalyticSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analytic", "--lambda", "0.8", "--mu", "1.0", "--k", "10"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "pi[") {
		t.Error("analytic output missing distribution")
	}
}

func TestRunSimulateSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"simulate", "--lambda", "0.8", "--mu", "1.0", "--k", "10",
		"--seed", "42", "--arrivals", "1000"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "arrivals") {
		t.Error("simulate output missing arrivals")
	}
}

func TestRunCompareSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"compare", "--lambda", "0.8", "--mu", "1.0", "--k", "10",
		"--seed", "42", "--arrivals", "50000"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "cross-rules passed") {
		t.Errorf("compare output missing cross-rules; got: %s", stdout.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"bogus"}, &stdout, &stderr)
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
}
