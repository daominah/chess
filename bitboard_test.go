package chess

import (
	"fmt"
	"strconv"
	"testing"
)

type bitboardTestPair struct {
	initial  uint64
	reversed uint64
}

var (
	tests = []bitboardTestPair{
		{
			uint64(1),
			uint64(9223372036854775808),
		},
		{
			uint64(18446744073709551615),
			uint64(18446744073709551615),
		},
		{
			uint64(0),
			uint64(0),
		},
	}
)

func TestBitboardReverse(t *testing.T) {
	for _, p := range tests {
		r := uint64(bitboard(p.initial).Reverse())
		if r != p.reversed {
			t.Fatalf("bitboard reverse of %s expected %s but got %s", intStr(p.initial), intStr(p.reversed), intStr(r))
		}
	}
}

func TestBitboardOccupied(t *testing.T) {
	m := map[Square]bool{
		B3: true,
	}
	bb := newBitboard(m)
	if bb.Occupied(B3) != true {
		t.Fatalf("bitboard occupied of %s expected %t but got %t", bb, true, false)
	}
}

func BenchmarkBitboardReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := uint64(9223372036854775807)
		bitboard(u).Reverse()
	}
}

func intStr(i uint64) string {
	return bitboard(i).String()
}

func TestBitboard_Draw(t *testing.T) {
	if true {
		// following code calculates rank attack ray (similar to example in
		// https://www.chessprogramming.org/Hyperbola_Quintessence)
		// example rank with the Rook and some occupied squares: "1 1 0 R 0 0 1 1",
		// result rank attack targets:                           "0 1 1 0 1 1 1 0"
		tmp, _ := strconv.ParseUint(`0000000011010011000000000000000000000000000000000000000000000000`, 2, 64)
		occupancy := bitboard(tmp)
		fmt.Println("occupancy rank 2", occupancy.Draw())
		rook := bitboard(uint64(1) << (63 - uint(D2)))
		fmt.Println("slider (the Rook on D2)", rook.Draw())
		fmt.Println("o-r", (occupancy - rook).Draw()) // clear the slider
		oSub2r := occupancy - 2*rook
		fmt.Println(`borrow 1 from the first blocker on the left: o-2r`, oSub2r.Draw())
		OSub2R := (occupancy.Reverse() - 2*rook.Reverse()).Reverse()
		fmt.Println(`borrow 1 from the first blocker on the right: (o'-2r')'`, OSub2R.Draw())
		attackSet := oSub2r ^ OSub2R
		fmt.Println("attackSet", attackSet.Draw())
	}
	if false {
		// https://www.chessprogramming.org/Hyperbola_Quintessence,
		// file attack example
		tmp, _ := strconv.ParseUint("0000000000010000000000000000000000000000000000000001000000010000", 2, 64)
		occupancy := bitboard(tmp)
		t.Log("occupancy", occupancy.Draw())
		rook := bbForSquare(D2)
		t.Log("slider", rook.Draw())
		t.Log("o-r", (occupancy - rook).Draw()) // clear the slider
		oSub2r := (occupancy - 2*rook) & bbFileD
		t.Log("o-2r", oSub2r.Draw()) // borrows "one" from next blocker (blocker value to 0, attack ray direction to blocker to 1)
		OSub2R := (occupancy.Reverse() - 2*rook.Reverse()).Reverse() & bbFileD
		t.Log("o'-2r'", OSub2R.Draw()) // other direction
		attackSet := oSub2r ^ OSub2R
		t.Log("attackSet", attackSet.Draw())
	}
}
