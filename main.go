package main

import "sync"

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go StartEcho()
	wg.Add(1)
	go StartBot()
	wg.Wait()

}
