package exec

type Middleware func(next Handler) Handler
