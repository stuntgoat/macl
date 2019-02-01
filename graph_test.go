package main

import (
	"testing"
)

func newCoinMap(coins ...CoinKey) map[CoinKey]bool {
	m := map[CoinKey]bool{}
	for _, t := range coins {
		m[t] = true
	}
	return m
}

func Test_FindConsecutive(t *testing.T) {
	/*

	   a _ _ a
	   _ _ a _
	   _ a a _
	   a _ _ _

	*/
	m := newCoinMap(
		CoinKey{0, 0},
		CoinKey{0, 3},
		CoinKey{1, 2},
		CoinKey{2, 1},
		CoinKey{2, 2},
		CoinKey{3, 0})

	pg := PlayerGraph{m}
	if !pg.FindConsecutive(0, 3, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(1, 2, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(2, 1, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(3, 0, 4) {
		t.Error("expecting to find 4")
	}

	if !pg.FindConsecutive(0, 3, 3) {
		t.Error("expecting to find 3")
	}

	if pg.FindConsecutive(0, 3, 5) {
		t.Error("expecting NOT to find 5 consecutive")
	}

	/*

	   a _ _ a
	   _ _ a _
	   _ _ a _
	   a a a a

	*/
	m = newCoinMap(
		CoinKey{0, 0},
		CoinKey{0, 3},

		CoinKey{1, 2},
		CoinKey{2, 2},

		CoinKey{3, 0},
		CoinKey{3, 1},
		CoinKey{3, 2},
		CoinKey{3, 3})

	pg = PlayerGraph{m}
	if !pg.FindConsecutive(3, 0, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(3, 1, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(3, 2, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(3, 3, 4) {
		t.Error("expecting to find 4")
	}

	if !pg.FindConsecutive(3, 3, 3) {
		t.Error("expecting to find 3")
	}

	if pg.FindConsecutive(3, 2, 5) {
		t.Error("expecting NOT to find 5 consecutive")
	}

	/*

	   a _ _ a
	   _ a a _
	   _ _ a _
	   a a _ a

	*/
	m = newCoinMap(
		CoinKey{0, 0},
		CoinKey{0, 3},

		CoinKey{1, 1},
		CoinKey{1, 2},

		CoinKey{2, 2},

		CoinKey{3, 0},
		CoinKey{3, 1},
		CoinKey{3, 3})

	pg = PlayerGraph{m}
	if !pg.FindConsecutive(0, 0, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(1, 1, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(2, 2, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(3, 3, 4) {
		t.Error("expecting to find 4")
	}

	if !pg.FindConsecutive(0, 0, 3) {
		t.Error("expecting to find 3")
	}

	if pg.FindConsecutive(0, 0, 5) {
		t.Error("expecting NOT to find 5 consecutive")
	}

	/*

	   a _ _ a
	   _ _ a a
	   _ _ a a
	   a _ _ a

	*/
	m = newCoinMap(
		CoinKey{0, 0},
		CoinKey{0, 3},

		CoinKey{1, 2},
		CoinKey{1, 3},

		CoinKey{2, 2},
		CoinKey{2, 3},

		CoinKey{3, 0},
		CoinKey{3, 3})

	pg = PlayerGraph{m}
	if !pg.FindConsecutive(0, 3, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(1, 3, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(2, 3, 4) {
		t.Error("expecting to find 4")
	}
	if !pg.FindConsecutive(3, 3, 4) {
		t.Error("expecting to find 4")
	}

	if !pg.FindConsecutive(0, 3, 3) {
		t.Error("expecting to find 3")
	}

	if pg.FindConsecutive(0, 0, 5) {
		t.Error("expecting NOT to find 5 consecutive")
	}

}
