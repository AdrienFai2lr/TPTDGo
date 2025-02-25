package main

import (
	"fmt"
)

const finLigne = -1
const finImage = -2
const k = 5
const ordre = 3
const nbX = 10

func main() {
	var tabCan [ordre + 1]chan [3]float64
	for i := range tabCan {
		tabCan[i] = make(chan [3]float64)
	}

	go source(tabCan[0])
	for i := 1; i <= ordre; i++ { // Changé ici
		go traitement(float64(i), tabCan[i-1], tabCan[i])
	}
	affichage(tabCan[ordre])
}
func source(out chan<- [3]float64) {
	for i := 0; i < nbX; i++ {
		x := float64(i)
		out <- [3]float64{x, 1, 0}
	}
	close(out)
}
func traitement(k float64, s1 <-chan [3]float64, s2 chan<- [3]float64) {
	for values := range s1 {
		x, PUk_1, Sk_1 := values[0], values[1], values[2]
		PUk := PUk_1 * x
		Sk := Sk_1 + PUk/k
		s2 <- [3]float64{x, PUk, Sk}
	}
	close(s2)
}
func affichage(s3 chan [3]float64) {
	for values := range s3 {
		x, y, Sn := values[0], values[1], values[2]
		fmt.Println("valeur de y :", y)
		fmt.Printf("P%d(%.2f) = ,%.6f\n", ordre, x, Sn)
	}
}
