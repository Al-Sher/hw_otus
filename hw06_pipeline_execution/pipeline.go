package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(fanIn(done, in))
	}

	return fanIn(done, in)
}

func fanIn(done In, in In) Out {
	ch := make(Bi)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				ch <- v
			}
		}
	}()

	return ch
}
