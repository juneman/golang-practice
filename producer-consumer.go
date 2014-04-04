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
	fmt.Printf("Good id: %d, state:%s\n", this.id, this.state)
}

func product(capcity int, speed float64) chan Good {
	line := make(chan Good, capcity)

	go func() {
		for i := 1; ; i++ {
			good := new(Good)
			good.setId(i)
			good.setState("Product")

			good.detail()
			line <- *good
			time.Sleep((time.Duration) ((time.Duration)(1000.0/speed) * time.Millisecond))
		}
	}()

	return line
}

func transport(in chan Good, capcity int, speed float64) chan Good {
	out := make(chan Good, capcity)

	go func() {
		for {
			good := <-in
			good.setState("Transporting")
			good.detail()

			out <- good
			time.Sleep((time.Duration) ((time.Duration)(1000.0/speed) * time.Millisecond))
		}
	}()

	return out
}

func consummer(in chan Good, speed float64) {
		for {
			good := <-in

			good.setState("Consuming")
			good.detail()
			time.Sleep((time.Duration) ((time.Duration)(1000.0/speed) * time.Millisecond))
		}
}

func main() {
    capcity := 10
    product_speed := 15.0
    trans_speed := 8.0
    consuming_speed := 12.0

	consummer(transport(product(capcity, product_speed), capcity, trans_speed), consuming_speed)
}
