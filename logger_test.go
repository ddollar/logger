package logger

import (
	"bytes"
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
)

var (
	buffer bytes.Buffer
)

func NewLogger() *Logger {
	buffer.Truncate(0)
	return NewWriter("ns=test", &buffer)
}

func TestNew(t *testing.T) {
	log := New("ns=test")
	assertEquals(t, log.namespace, "ns=test")
}

func TestAt(t *testing.T) {
	log := NewLogger()
	log.At("target").Logf("foo=bar")
	assertLine(t, buffer.String(), `ns=test at=target foo=bar`)
}

func TestAtOverrides(t *testing.T) {
	log := NewLogger()
	log.At("target1").At("target2").Logf("foo=bar")
	assertLine(t, buffer.String(), `ns=test at=target2 foo=bar`)
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	log := NewLogger()
	FromContext(log.Append("with=context").WithContext(ctx)).Logf("foo=bar")
	assertLine(t, buffer.String(), `ns=test with=context foo=bar`)
}
func TestError(t *testing.T) {
	log := NewLogger()
	log.Error(errors.New("broken"))
	assertLine(t, buffer.String(), `ns=test error="broken"`)
}

func TestLog(t *testing.T) {
	log := NewLogger()
	log.Logf("string=%q int=%d float=%0.2f", "foo", 42, 3.14159)
	assertLine(t, buffer.String(), `ns=test string="foo" int=42 float=3.14`)
}

func TestNamespace(t *testing.T) {
	log := NewLogger()
	log.Namespace("foo=bar").Namespace("baz=qux").Logf("fred=barney")
	assertLine(t, buffer.String(), `ns=test foo=bar baz=qux fred=barney`)
}

func TestReplace(t *testing.T) {
	log := NewLogger()
	log.Namespace("baz=qux1").Replace("baz", "qux2").Logf("foo=bar")
	assertLine(t, buffer.String(), `ns=test baz=qux2 foo=bar`)
}

func TestReplaceExisting(t *testing.T) {
	log := NewLogger()
	log.Namespace("foo=bar").Namespace("baz=qux").Replace("baz", "zux").Logf("thud=grunt")
	assertLine(t, buffer.String(), `ns=test foo=bar baz=zux thud=grunt`)
}

func TestStart(t *testing.T) {
	log := NewLogger()
	log.Start().Successf("num=%d", 42)
	assertContains(t, buffer.String(), "elapsed=")
}

func TestStep(t *testing.T) {
	log := NewLogger()
	log.Step("target").Logf("foo=bar")
	assertLine(t, buffer.String(), `ns=test step=target foo=bar`)
}

func TestStepOverrides(t *testing.T) {
	log := NewLogger()
	log.Step("target1").Step("target2").Logf("foo=bar")
	assertLine(t, buffer.String(), `ns=test step=target2 foo=bar`)
}

func TestSuccess(t *testing.T) {
	log := NewLogger()
	log.Successf("num=%d", 42)
	assertLine(t, buffer.String(), `ns=test state=success num=42`)
}

func assertContains(t *testing.T, got, search string) {
	if strings.Index(got, search) == -1 {
		t.Errorf("\n   expected: %q\n to contain: %q", got, search)
	}
}

func assertEquals(t *testing.T, got, search string) {
	if got != search {
		t.Errorf("\n   expected: %q\n to equal: %q", got, search)
	}
}

func assertLine(t *testing.T, got, search string) {
	search = search + "\n"
	if search != got {
		t.Errorf("\n   expected: %q\n to be: %q", got, search)
	}
}

func assertMatch(t *testing.T, got, search string) {
	r := regexp.MustCompile(search)

	if !r.MatchString(got) {
		t.Errorf("\n   expected: %q\n   to match: %q", got, search)
	}
}
