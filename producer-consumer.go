package main

import (
	"fmt"
    "time"
)

type Good struct {
	id    int
	state string
}

func (this *Good) setId(id int) {
	this.id = id
}

func (this *Good) setState(s string) {
	this.state = s
}

func (this *Good) detail() {
	fmt.Printf("Good id: %d, state:%s\n\n", this.id, this.state)
}

func product() chan Good {
	line := make(chan Good)

	go func() {
		for i := 1; ; i++ {
			good := new(Good)
			good.setId(i)
			good.setState("Product")

			good.detail()
			line <- *good
			time.Sleep(1 * time.Second)
		}
	}()

	return line
}

func transport(in chan Good) chan Good {
	out := make(chan Good)

	go func() {
		for {
			good := <-in
			good.setState("Transporting")
			good.detail()

			out <- good
		}
	}()

	return out
}

func consummer(in chan Good) {
		for {
			good := <-in

			good.setState("Consuming")
			good.detail()
		}
}

func main() {
	consummer(transport(product()))
}
