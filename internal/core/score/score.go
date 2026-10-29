// Package score keeps the wins and losses against the machine across games,
// per machine level, and stores them in a JSON file.
package score

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrNoPath means there is no file to read or write the score from, such as
// when the user has no configuration directory.
var ErrNoPath = errors.New("no score file path")

// Record is the score of one level.
type Record struct {
	Wins    int `json:"wins"`
	Losses  int `json:"losses"`
	BestWin int `json:"best_win"` // fewest player shots in a win, 0 with no win
}

// Board is the score of every level, by level name.
type Board map[string]Record

// AddWin records a win at level in shots, which becomes the best win if it
// took fewer shots than every earlier win.
func (b Board) AddWin(level string, shots int) {
	r := b[level]
	r.Wins++
	if r.BestWin == 0 || shots < r.BestWin {
		r.BestWin = shots
	}
	b[level] = r
}

// AddLoss records a loss at level.
func (b Board) AddLoss(level string) {
	r := b[level]
	r.Losses++
	b[level] = r
}

// Load reads the board at path. A missing file is an empty board; a file
// that cannot be read or is not a valid board is an error.
func Load(path string) (Board, error) {
	if path == "" {
		return nil, ErrNoPath
	}
	data, err := os.ReadFile(path) //nolint:gosec // the score file of the user, chosen by the program
	if errors.Is(err, fs.ErrNotExist) {
		return Board{}, nil
	}
	if err != nil {
		return nil, err
	}
	var b Board
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("read score %s: %w", path, err)
	}
	if b == nil {
		b = Board{}
	}
	return b, nil
}

// Save writes b to path, creating its directory, through a temporary file
// renamed into place so a failed write never leaves half a file.
func Save(path string, b Board) error {
	if path == "" {
		return ErrNoPath
	}
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) //nolint:errcheck // gone after the rename; only cleans up a failed save
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// DefaultPath is placar.json in the battleship folder of the user's
// configuration directory (os.UserConfigDir).
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "battleship", "placar.json"), nil
}
