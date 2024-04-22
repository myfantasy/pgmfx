package jobber

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/myfantasy/mfctx"
)

var ErrCheckOwnershipError = fmt.Errorf("fail to check ownership")
var ErrCheckOwnershipFalse = fmt.Errorf("ownership was lost")

type CheckIsOwner func(ctx context.Context, worker string) (isOwner bool, err error)

type Ownership struct {
	ctxBase *mfctx.Crumps

	mx sync.Mutex

	isActive         bool
	ctxActive        *mfctx.Crumps
	activeCancelFunc context.CancelCauseFunc

	worker string

	checkIsOwner CheckIsOwner

	errorCounter int
}

type OwnershipChecker struct {
	mx  sync.Mutex
	ows map[string]*Ownership

	ctxBase      *mfctx.Crumps
	checkIsOwner CheckIsOwner
}

func MakeOwnership(ctxBase *mfctx.Crumps, worker string, checkIsOwner CheckIsOwner) *Ownership {
	return &Ownership{
		ctxBase:      ctxBase,
		worker:       worker,
		checkIsOwner: checkIsOwner,
	}
}

func MakeOwnershipChecker(ctxBase *mfctx.Crumps, checkIsOwner CheckIsOwner) *OwnershipChecker {
	return &OwnershipChecker{
		ctxBase:      ctxBase,
		checkIsOwner: checkIsOwner,
		ows:          make(map[string]*Ownership),
	}
}

func (o *Ownership) GetActiveContext() (ctx context.Context, isActive bool) {
	o.mx.Lock()
	defer o.mx.Unlock()

	if o.isActive {
		return o.ctxActive, true
	}

	return nil, false
}

func (o *Ownership) CheckActive() {
	isOwner, err := o.checkIsOwner(o.ctxBase, o.worker)

	o.mx.Lock()
	defer o.mx.Unlock()
	if err != nil {
		o.errorCounter--
	}
	if o.errorCounter <= 0 && o.isActive {
		o.activeCancelFunc(ErrCheckOwnershipError)
		o.isActive = false
	}
	if o.errorCounter < 0 {
		o.errorCounter = 0
	}

	if err != nil {
		return
	}

	o.errorCounter++
	if o.errorCounter > 3 {
		o.errorCounter = 3 // a kind of magic
	}

	if !isOwner && o.isActive {
		o.activeCancelFunc(ErrCheckOwnershipFalse)
		o.isActive = false
	}

	if isOwner && !o.isActive {
		o.ctxActive, o.activeCancelFunc = o.ctxBase.WithCancelCause()
		o.isActive = true
	}
}

// UpdateJob do CheckActive() every 15 sec
func (o *Ownership) UpdateJob() {
	for o.ctxBase.Err() != nil {
		o.CheckActive()

		select {
		case <-time.After(15 * time.Second):
		case <-o.ctxActive.Done():
			return
		}
	}
}

func (owc *OwnershipChecker) InitWorker(worker string) {
	owc.mx.Lock()
	defer owc.mx.Unlock()

	_, ok := owc.ows[worker]
	if !ok {
		owc.ows[worker] = MakeOwnership(owc.ctxBase, worker, owc.checkIsOwner)

		go owc.ows[worker].UpdateJob()
	}
}

func (owc *OwnershipChecker) GetActiveContext(worker string) (ctx context.Context, isActive bool) {
	owc.mx.Lock()
	defer owc.mx.Unlock()

	ow, ok := owc.ows[worker]

	if !ok {
		return nil, false
	}

	return ow.GetActiveContext()
}
