package web

type HandleFuncV1 func(ctx *Context) (next bool)

// type ChainV1 []HandleFuncV1

type Middleware func(next HandleFunc) HandleFunc

type Chain []HandleFunc

type ChainV1 struct {
	handlers []HandleFuncV1
}

func (c ChainV1) Run(ctx *Context) {
	for _, h := range c.handlers {
		next := h(ctx)
		if !next {
			return
		}

	}
}
