// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"path"
	"testing"

	"github.com/limetext/backend/keys"
	"github.com/limetext/backend/parser"
	"github.com/limetext/backend/render"
	qp "github.com/quarnster/parser"
)

func TestGetEditor(t *testing.T) {
	editor := GetEditor()
	if editor == nil {
		t.Error("Expected an editor, but got nil")
	}
}

func TestLoadKeyBindings(t *testing.T) {
	ed := GetEditor()

	if ed.defaultKB.KeyBindings().Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}
}

func TestLoadSettings(t *testing.T) {
	ed := GetEditor()
	switch ed.Platform() {
	case "windows":
		if res := ed.Settings().Get("font_face", ""); res != "Consolas" {
			t.Errorf("Expected windows font_face be Consolas, but is %s", res)
		}
	case "darwin":
		if res := ed.Settings().Get("font_face", ""); res != "Menlo" {
			t.Errorf("Expected OSX font_face be Menlo, but is %s", res)
		}
	default:
		if res := ed.Settings().Get("font_face", ""); res != "Monospace" {
			t.Errorf("Expected Linux font_face be Monospace, but is %s", res)
		}
	}
}

func TestNewWindow(t *testing.T) {
	editor := GetEditor()
	l := len(editor.Windows())

	w := editor.NewWindow()
	defer w.Close()

	if len(editor.Windows()) != l+1 {
		t.Errorf("Expected 1 window, but got %d", len(editor.Windows()))
	}
}

func TestRemoveWindow(t *testing.T) {
	editor := GetEditor()
	l := len(editor.Windows())

	w0 := editor.NewWindow()
	defer w0.Close()

	editor.remove(w0)

	if len(editor.Windows()) != l {
		t.Errorf("Expected the window to be removed, but %d still remain", len(editor.Windows()))
	}

	w1 := editor.NewWindow()
	defer w1.Close()

	w2 := editor.NewWindow()
	defer w2.Close()

	editor.remove(w1)

	if len(editor.Windows()) != l+1 {
		t.Errorf("Expected the window to be removed, but %d still remain", len(editor.Windows()))
	}
}

func TestSetActiveWindow(t *testing.T) {
	editor := GetEditor()

	w1 := editor.NewWindow()
	defer w1.Close()

	w2 := editor.NewWindow()
	defer w2.Close()

	if editor.ActiveWindow() != w2 {
		t.Error("Expected the newest window to be active, but it wasn't")
	}

	editor.SetActiveWindow(w1)

	if editor.ActiveWindow() != w1 {
		t.Error("Expected the first window to be active, but it wasn't")
	}
}

func TestSetFrontend(t *testing.T) {
	f := DummyFrontend{}

	editor := GetEditor()
	editor.SetFrontend(&f)

	if editor.Frontend() != &f {
		t.Errorf("Expected a DummyFrontend to be set, but got %T", editor.Frontend())
	}
}

func TestClipboard(t *testing.T) {
	editor := GetEditor()

	// Put back whatever was already there.
	clip := editor.GetClipboard()
	defer editor.SetClipboard(clip)

	want := "test0"

	editor.SetClipboard(want)

	if got := editor.GetClipboard(); got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}

	want = "test1"

	editor.SetClipboard(want)

	if got := editor.GetClipboard(); got != want {
		t.Errorf("Expected %q to be on the clipboard, but got %q", want, got)
	}
}

func TestHandleInput(t *testing.T) {
	// FIXME: This test causes a panic.
	t.Skip("Avoiding pointer issues causing a panic.")

	editor := GetEditor()
	kp := keys.KeyPress{Key: 'i'}

	editor.HandleInput(kp)

	if ki := <-editor.keyInput; ki != kp {
		t.Errorf("Expected %s to be on the input buffer, but got %s", kp, ki)
	}
}

type dummyColorSc struct {
	name string
}

func (d *dummyColorSc) Name() string {
	return d.name
}

func (d *dummyColorSc) Spice(*render.ViewRegions) render.Flavour { return render.Flavour{} }

func TestAddColorScheme(t *testing.T) {
	cs := new(dummyColorSc)
	ed := GetEditor()

	ed.AddColorScheme("test/path", cs)
	if ret := ed.colorSchemes["test/path"]; ret != cs {
		t.Errorf("Expected 'test/path' color scheme %v, but got %v", cs, ret)
	}
}

func TestGetColorScheme(t *testing.T) {

}

func TestColorSchemes(t *testing.T) {

}

type dummySyntax struct {
	name      string
	filetypes []string
	data      string
}

func (d *dummySyntax) Name() string {
	return d.name
}

func (d *dummySyntax) FileTypes() []string {
	return d.filetypes
}

func (d *dummySyntax) Parser(data string) (parser.Parser, error) {
	d.data = data
	return d, nil
}

func (d *dummySyntax) Parse() (*qp.Node, error) { return nil, nil }

func TestAddSyntax(t *testing.T) {
	syn := new(dummySyntax)
	ed := GetEditor()

	ed.AddSyntax("test/path", syn)
	if ret := ed.syntaxes["test/path"]; ret != syn {
		t.Errorf("Expected 'test/path' syntax %v, but got %v", syn, ret)
	}
}

func TestGetSyntax(t *testing.T) {

}

func TestSyntaxes(t *testing.T) {

}

func init() {
	ed := GetEditor()
	ed.Init()
	ed.AddPackagesPath("shipped", path.Join("testdata", "shipped"))
	ed.AddPackagesPath("default", path.Join("testdata", "default"))
	ed.AddPackagesPath("user", path.Join("testdata", "user"))
}