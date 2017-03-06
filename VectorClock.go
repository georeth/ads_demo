package ads_demo

import "fmt"

type VectorClock struct {
	B []int
	R int
}

func NewVectorClock(numServer int) *VectorClock {
	var clock VectorClock
	clock.B = make([]int, numServer)
	clock.R = 0
	return &clock
}

func (clock *VectorClock) Ready(now *VectorClock) bool {
	for i := 0; i < len(now.B); i++ {
		if clock.B[i] > now.B[i] {
			return false
		}
	}
	if clock.R > now.R {
		return false
	}
	return true
}

// return previous
func (clock *VectorClock) Red() int {
	return clock.R
}
func (clock *VectorClock) Copy() VectorClock {
	oldb := make([]int, len(clock.B))
	copy(oldb, clock.B)
	return VectorClock{oldb, clock.R}
}
func (clock *VectorClock) Tick(id int, color COLOR) VectorClock {
	old := clock.Copy()
	clock.B[id]++
	if color == RED {
		clock.R++
	}

	return old
}

func (clock *VectorClock) Print(id int) {
	fmt.Printf("#%d [", id)
	for i := 0; i < len(clock.B); i++ {
		fmt.Printf(" %d", clock.B[i])
	}
	fmt.Printf(" ; %d ]\n", clock.R)
}
