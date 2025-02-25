package main

import "fmt"

const finLigne = -1
const finImage = -2

func main() {
	sourceOut := make(chan int)
	seuillageOut := make(chan int)

	go source(sourceOut)

	go seuillage(100, sourceOut, seuillageOut)

	affichage(seuillageOut)

}
func source(s chan int) {
	image := [12]int{100, 200,
		150, finLigne, 32, 250, 18, finLigne, 47, 242, 99, finImage}

	for _, e := range image {
		s <- e
	}

}
func seuillage(seuil int, s1 chan int, s2 chan int) {
	for {
		pixelVal := <-s1
		switch pixelVal {
		case finLigne, finImage:
			s2 <- pixelVal
		default:
			if pixelVal < seuil {
				s2 <- 0
			} else {
				s2 <- 1
			}
		}
	}
}
func affichage(s3 chan int) {
	for {
		val := <-s3
		switch val {
		case finLigne:
			fmt.Println()
		case finImage:
			fmt.Println()
			fmt.Println()
			return
		default:
			fmt.Print(val)
		}

	}
}
