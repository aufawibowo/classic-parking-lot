package parking

import (
	"container/heap"
	"strconv"
)

type ParkingLot struct {
	capacity  int
	freeSlots *IntMinHeap // heap of unused slot numbers
	slots     []string    // index i holds car number or ""
	carToSlot map[string]int
}

func NewParkingLot(capacity int) *ParkingLot {
	h := &IntMinHeap{}
	for i := 1; i <= capacity; i++ {
		heap.Push(h, i)
	}
	return &ParkingLot{
		capacity:  capacity,
		freeSlots: h,
		slots:     make([]string, capacity), // all ""
		carToSlot: make(map[string]int),
	}
}

func (p *ParkingLot) Park(car string) (int, bool) {
	if p.freeSlots.Len() == 0 {
		return 0, false
	}
	slot := heap.Pop(p.freeSlots).(int)
	p.slots[slot-1] = car
	p.carToSlot[car] = slot
	return slot, true
}

func (p *ParkingLot) Leave(car string, hours int) (int, int, bool) {
	slot, ok := p.carToSlot[car]
	if !ok {
		return 0, 0, false
	}

	// fee
	fee := 10
	if hours > 2 {
		fee += (hours - 2) * 10
	}

	// cleanup
	delete(p.carToSlot, car)
	p.slots[slot-1] = ""
	heap.Push(p.freeSlots, slot)

	return slot, fee, true
}

func (p *ParkingLot) Status() [][2]string {
	out := make([][2]string, 0, p.capacity)
	for i, car := range p.slots {
		if car != "" {
			out = append(out, [2]string{strconv.Itoa(i + 1), car})
		}
	}
	return out
}
