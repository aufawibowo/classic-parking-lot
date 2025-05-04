package parking

import (
	"testing"
)

func mustPark(t *testing.T, lot *ParkingLot, car string, wantSlot int) {
	t.Helper()
	gotSlot, ok := lot.Park(car)
	if !ok {
		t.Fatalf("expected park(%s) to succeed", car)
	}
	if gotSlot != wantSlot {
		t.Fatalf("park(%s): got slot %d, want %d", car, gotSlot, wantSlot)
	}
}

func TestParkingLot_BasicFlow(t *testing.T) {
	const capacity = 3
	lot := NewParkingLot(capacity)

	// 1. newly-created lot should expose all free slots through the heap.
	if lot.freeSlots.Len() != capacity {
		t.Fatalf("got %d freeSlots, want %d", lot.freeSlots.Len(), capacity)
	}

	// 2. Park up to full capacity – must allocate lowest-numbered slots.
	mustPark(t, lot, "B​1", 1)
	mustPark(t, lot, "B​2", 2)
	mustPark(t, lot, "B​3", 3)

	// 3. Further parking must fail.
	if slot, ok := lot.Park("OVERFLOW"); ok || slot != 0 {
		t.Fatalf("expected parking to be full, got slot=%d ok=%v", slot, ok)
	}

	// 4. Leave with >2 hours – fee = 10 + (hours-2)*10.
	slot, fee, ok := lot.Leave("B​2", 4) // 4 h → 10 + 20 = 30
	if !ok || slot != 2 || fee != 30 {
		t.Fatalf("leave returned slot=%d fee=%d ok=%v (want 2,30,true)", slot, fee, ok)
	}

	// 5. Freed slot should become the next allocation target.
	mustPark(t, lot, "C​1", 2)

	// 6. Status must list occupied slots in ascending order.
	want := [][2]string{
		{"1", "B​1"},
		{"2", "C​1"},
		{"3", "B​3"},
	}
	got := lot.Status()
	if len(got) != len(want) {
		t.Fatalf("status len=%d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("status[%d]=%v want %v", i, got[i], want[i])
		}
	}
}

func TestParkingLot_Errors(t *testing.T) {
	lot := NewParkingLot(1)
	mustPark(t, lot, "X", 1)

	// Leaving an unknown car.
	if _, _, ok := lot.Leave("UNKNOWN", 1); ok {
		t.Fatalf("expected leave to fail for unknown car")
	}

	// Parking the same registration twice should just give “lot full” because slot occupied.
	if _, ok := lot.Park("X"); ok {
		t.Fatalf("expected duplicate car to fail (lot full)")
	}
}
