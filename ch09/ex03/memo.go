package memo

import "context"

type request struct {
	key      string
	response chan<- result
	ctx      context.Context
}

type entry struct {
	res   result
	ready chan struct{}
}

type Memo struct {
	requests chan request
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

type Func func(ctx context.Context, key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func IsCancelled(err error) bool {
	type cancelled interface {
		Cancelled() bool
	}
	c, ok := err.(cancelled)
	return ok && c.Cancelled()
}

func (memo *Memo) Get(ctx context.Context, key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response, ctx}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e != nil {
			select {
			case <-e.ready:
				if IsCancelled(e.res.err) {
					delete(cache, req.key)
					e = nil
				}
			default:
			}
		}

		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(req.ctx, f, req.key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(ctx context.Context, f Func, key string) {
	e.res.value, e.res.err = f(ctx, key)
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}
