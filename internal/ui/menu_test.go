package ui

import (
	"regexp"
	"strings"
	"testing"
)

// escapes captures the truecolor sequences lipgloss emits, so a row's colouring
// can be inspected without depending on where lipgloss puts resets.
var escapes = regexp.MustCompile(`\x1b\[38;2;(\d+);(\d+);(\d+)m`)

// coloursIn returns the distinct foreground colours used in a rendered row.
func coloursIn(row string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range escapes.FindAllStringSubmatch(row, -1) {
		key := m[1] + "," + m[2] + "," + m[3]
		if !seen[key] {
			seen[key] = true
			out = append(out, key)
		}
	}
	return out
}

// menuRows renders the menu block on its own.
//
// Taking rows out of a full Render() would include the frame's │ on each side and
// count the border's colour as one of the row's — which is exactly the false
// positive this helper exists to avoid.
func menuRows(t *testing.T, p Panel) []string {
	t.Helper()

	rows := strings.Split(darkTheme().menu(p), "\n")
	if len(rows) != len(p.Menu) {
		t.Fatalf("rendered %d menu rows, want %d", len(rows), len(p.Menu))
	}
	return rows
}

// One row, one colour. An earlier version split each row across a cursor colour,
// a label colour and a description colour, which made the cursor compete with the
// label instead of leading the eye to it.
func TestEachMenuRowUsesASingleColour(t *testing.T) {
	forceTrueColor(t)

	for i, row := range menuRows(t, demoPanel()) {
		if got := coloursIn(row); len(got) != 1 {
			t.Errorf("menu row %d uses %d colours (%v), want 1", i, len(got), got)
		}
	}
}

func TestSelectedRowIsGoldEndToEnd(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel() // status is selected

	theme := darkTheme()
	wantGold := rgbOf(theme.Gold)
	wantPlain := rgbOf(theme.Steel)

	for i, row := range menuRows(t, p) {
		got := coloursIn(row)[0]
		want := wantPlain
		if i == p.Selected {
			want = wantGold
		}
		if got != want {
			t.Errorf("menu row %d is rgb(%s), want rgb(%s)", i, got, want)
		}
	}
}

// The whole row turns gold, not just the label: the cursor, the bullet and the
// description come with it.
func TestSelectedRowColoursTheDescriptionToo(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()

	row := menuRows(t, p)[p.Selected]
	plain := strip(row)

	for _, want := range []string{"❯", "▸", p.Menu[p.Selected].Label, p.Menu[p.Selected].Desc} {
		if !strings.Contains(plain, want) {
			t.Errorf("the selected row is missing %q", want)
		}
	}
	if got := coloursIn(row); len(got) != 1 || got[0] != rgbOf(darkTheme().Gold) {
		t.Errorf("the selected row uses %v, want a single gold sweep", got)
	}
}

// Disabled rows look identical to enabled ones now, by explicit request: colour
// carries selection, nothing else. The refusal is behavioural, and this test pins
// that it stays behavioural rather than creeping back into the palette.
func TestDisabledRowsAreNotDimmed(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()

	rows := menuRows(t, p)
	var enabled, disabled string
	for i, item := range p.Menu {
		if i == p.Selected {
			continue
		}
		if item.Enabled {
			enabled = coloursIn(rows[i])[0]
		} else {
			disabled = coloursIn(rows[i])[0]
		}
	}
	if disabled == "" {
		t.Skip("the demo panel has no disabled rows")
	}
	if enabled != "" && enabled != disabled {
		t.Errorf("enabled rows are rgb(%s) and disabled rgb(%s); they must match", enabled, disabled)
	}
	if disabled != rgbOf(darkTheme().Steel) {
		t.Errorf("disabled rows are rgb(%s), want the ordinary text colour", disabled)
	}
}

// rgbOf reports the "r,g,b" a colour actually renders as.
//
// It asks lipgloss rather than converting the hex directly: lipgloss round-trips
// colours through a float representation, so #E4E4EC is emitted as rgb(227,227,236),
// not rgb(228,228,236). Probing through the same path both sides makes the
// assertion immune to that rounding.
func rgbOf(hex string) string {
	found := coloursIn(Fg(hex).Render("x"))
	if len(found) != 1 {
		return "unrenderable:" + hex
	}
	return found[0]
}

// Both movement commands are offered, always. `upgrade` fetches a published release and
// `update` pulls a checkout; which one applies depends on the machine, and hiding the other
// would make the menu change shape between them.
func TestBothUpgradeAndUpdateAreOffered(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.Menu = []MenuItem{
		{Label: "upgrade", Desc: "fetch the newest release · relink", Enabled: true},
		{Label: "update", Desc: "pull this checkout · rebuild · relink", Enabled: false},
	}

	out := strip(darkTheme().Render(p))
	for _, label := range []string{"upgrade", "update"} {
		if !strings.Contains(out, label) {
			t.Errorf("the menu does not offer %q:\n%s", label, out)
		}
	}
}

// Disabled, not dimmed and not absent. Colour carries selection and nothing else, so a
// disabled row keeps full contrast and reads as available-but-inert — which is why absence is
// the only way it could disappear, and it must not.
func TestTheInapplicableActionIsDisabledNotHidden(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.Menu = []MenuItem{
		{Label: "upgrade", Desc: "fetch the newest release · relink", Enabled: false},
		{Label: "status", Desc: "32 linked", Enabled: true},
	}

	out := strip(darkTheme().Render(p))
	if !strings.Contains(out, "upgrade") {
		t.Errorf("a disabled action was hidden:\n%s", out)
	}
	if !strings.Contains(out, "fetch the newest release") {
		t.Errorf("a disabled action lost its description:\n%s", out)
	}
}
