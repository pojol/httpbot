package sizewg

import (
	"context"
	"fmt"
	"math"
	"sync"
)

type SizeWaitGroup struct {
	Size    int
	wg      sync.WaitGroup
	blockch chan struct{}
}

func New(size int) SizeWaitGroup {
	if size <= 0 || size > math.MaxInt16 {
		panic(fmt.Errorf("not allow size %v", size))
	}

	return SizeWaitGroup{
		Size:    size,
		wg:      sync.WaitGroup{},
		blockch: make(chan struct{}, size),
	}
}

func (swg *SizeWaitGroup) Enter() {
	swg.EnterWithContext(context.Background())
}

func (swg *SizeWaitGroup) EnterWithContext(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case swg.blockch <- struct{}{}:
		break
	}

	swg.wg.Add(1)
}

func (swg *SizeWaitGroup) Leave() {
	<-swg.blockch
	swg.wg.Done()
}

func (swg *SizeWaitGroup) Wait() {
	swg.wg.Wait()
}
