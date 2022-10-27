package dev07

func or(channels ...<-chan any) <-chan any {
	exit := make(chan any)
	for _, v := range channels {
		go func(c <-chan any) {
			select {
			case <-c:
				exit <- "end"
				close(exit)
			}
		}(v)
	}

	<-exit
	return exit
}
