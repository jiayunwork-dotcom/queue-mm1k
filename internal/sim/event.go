package sim

import "container/heap"

// event is one scheduled transition of the event table: an arrival or a
// departure at a future time.
type event struct {
	time  float64
	kind  eventKind
}

type eventKind int

const (
	arrival eventKind = iota
	departure
)

// eventHeap is a min-heap of events ordered by time, the engine's event
// table. heap.Interface guarantees the next event is always at index 0.
type eventHeap []event

func (h eventHeap) Len() int            { return len(h) }
func (h eventHeap) Less(i, j int) bool  { return h[i].time < h[j].time }
func (h eventHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }

func (h *eventHeap) Push(x any) {
	*h = append(*h, x.(event))
}

func (h *eventHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

// schedule inserts an event into the table.
func (h *eventHeap) schedule(t float64, k eventKind) {
	heap.Push(h, event{time: t, kind: k})
}

// next removes and returns the earliest event.
func (h *eventHeap) next() (event, bool) {
	if h.Len() == 0 {
		return event{}, false
	}
	return heap.Pop(h).(event), true
}
