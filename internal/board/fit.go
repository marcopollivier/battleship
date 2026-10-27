package board

// maxFitAttempts is how many placements Fits tries before it gives up and
// answers true. The default fleet never comes close to it.
const maxFitAttempts = 1_000_000

// Fits reports whether the ships of the fleet not placed yet can still all
// be placed on b, following the placing rules. It searches every way to
// place them, without randomness, and leaves b unchanged.
//
// If the search passes maxFitAttempts, which only a large custom fleet can
// do, Fits answers true: it only answers false when it is sure.
func Fits(b *Board) bool {
	return fitsWithin(b, maxFitAttempts)
}

// fitsWithin is Fits with limit in place of maxFitAttempts.
func fitsWithin(b *Board, limit int) bool {
	fits, attempts := fitSearch(b, limit)
	return fits || attempts > limit
}

// fitSearch looks for a way to place the missing ships of b on a copy of it.
// It reports whether it found one and how many placements it tried, and
// stops once that passes limit.
func fitSearch(b *Board, limit int) (fits bool, attempts int) {
	s := &fitter{board: b.clone(), ships: b.Remaining(), limit: limit}
	return s.place(0, 0), s.attempts
}

// fitter is the state of one fitSearch: the board copy, the ships to place,
// largest first, the placements tried so far and the most it may try.
type fitter struct {
	board    *Board
	ships    []ShipType
	attempts int
	limit    int
}

// place tries every spot, from spot from on, for ship i, and the ships after
// it on each spot that works, undoing it when they do not fit. A spot is a
// cell, in rows and then columns, and one of the ship's orientations.
//
// Ships of the same type are interchangeable, so the next ship of the same
// type only tries spots after this one: any other order is the same
// arrangement again.
func (s *fitter) place(i, from int) bool {
	if i == len(s.ships) {
		return true
	}
	t := s.ships[i]
	orientations := ValidOrientations(t)
	if len(orientations) == 0 {
		return false
	}
	for spot := from; spot < Size*Size*len(orientations); spot++ {
		if s.attempts > s.limit {
			return false
		}
		s.attempts++
		cell := spot / len(orientations)
		ref := Coord{Col: cell % Size, Row: cell / Size}
		if s.board.Place(t, ref, orientations[spot%len(orientations)]) != nil {
			continue
		}
		next := 0
		if i+1 < len(s.ships) && s.ships[i+1] == t {
			next = spot + 1
		}
		if s.place(i+1, next) {
			return true
		}
		s.board.RemoveLast()
	}
	return false
}
