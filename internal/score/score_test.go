package score

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// CA-4: of two wins with 50 and 41 shots, the best win is 41, whatever the
// order; losses only add to the losses of their level.
func TestBoard_AddWin(t *testing.T) {
	for _, shots := range [][]int{{50, 41}, {41, 50}} {
		b := Board{}
		for _, n := range shots {
			b.AddWin("dificil", n)
		}
		if want := (Record{Wins: 2, BestWin: 41}); b["dificil"] != want {
			t.Errorf("after wins with %v shots, dificil = %+v, want %+v", shots, b["dificil"], want)
		}
	}

	b := Board{}
	b.AddLoss("facil")
	b.AddLoss("facil")
	b.AddWin("mestre", 60)
	want := Board{"facil": {Losses: 2}, "mestre": {Wins: 1, BestWin: 60}}
	if !reflect.DeepEqual(b, want) {
		t.Errorf("board = %+v, want %+v", b, want)
	}
}

// A saved board loads back the same, in the JSON format of the plan, and a
// missing file is an empty board.
func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "battleship", "placar.json")
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load(missing) error = %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("Load(missing) = %#v, want an empty board", got)
	}

	b := Board{"dificil": {Wins: 3, Losses: 2, BestWin: 41}}
	if err := Save(path, b); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if want := `{"dificil":{"wins":3,"losses":2,"best_win":41}}`; string(data) != want {
		t.Errorf("file = %s, want %s", data, want)
	}
	if got, err = Load(path); err != nil || !reflect.DeepEqual(got, b) {
		t.Errorf("Load() = %+v, %v, want %+v", got, err, b)
	}

	// Saving again replaces the file and leaves no temporary file behind.
	if err := Save(path, Board{}); err != nil {
		t.Fatalf("Save(empty) error = %v", err)
	}
	if got, err = Load(path); err != nil || len(got) != 0 {
		t.Errorf("Load() after saving an empty board = %+v, %v", got, err)
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("folder has %d files, want only placar.json", len(entries))
	}
}

// CA-7: a corrupt file, or one that is not a file, is an error, and the file
// is kept as it was.
func TestLoad_Corrupt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "placar.json")
	if err := os.WriteFile(path, []byte("{nope"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Error("Load(corrupt) error = nil, want an error")
	}
	if data, _ := os.ReadFile(path); string(data) != "{nope" {
		t.Errorf("file after Load = %q, want it unchanged", data)
	}
	if _, err := Load(dir); err == nil {
		t.Error("Load(folder) error = nil, want an error")
	}
	if _, err := Load(""); err == nil {
		t.Error(`Load("") error = nil, want an error`)
	}
}

// CA-7: saving in a folder without write permission, under a file, or with
// no path is an error.
func TestSave_Unwritable(t *testing.T) {
	readOnly := t.TempDir()
	if err := os.Chmod(readOnly, 0o500); err != nil {
		t.Fatalf("Chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(readOnly, 0o700) })
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	for _, path := range []string{
		filepath.Join(readOnly, "placar.json"),
		filepath.Join(file, "placar.json"),
		"",
	} {
		if err := Save(path, Board{"facil": {Wins: 1, BestWin: 30}}); err == nil {
			t.Errorf("Save(%q) error = nil, want an error", path)
		}
	}
}

// The default path is placar.json in the battleship folder of the user's
// configuration directory.
func TestDefaultPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("AppData", filepath.Join(home, "AppData"))
	dir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir: %v", err)
	}
	got, err := DefaultPath()
	if want := filepath.Join(dir, "battleship", "placar.json"); err != nil || got != want {
		t.Errorf("DefaultPath() = %q, %v, want %q", got, err, want)
	}
}
