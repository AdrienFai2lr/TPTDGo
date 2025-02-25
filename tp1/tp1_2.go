package main

import "fmt"

func main() {
	sourceOut := make(chan int)
	seuillageOut := make(chan int)
	sync := make(chan bool)
	go func() {
		source(sourceOut)
		sync <- true
	}()
	go func() {
		seuillage(100, sourceOut, seuillageOut)
		sync <- true
	}()
	go func() {
		affichage(seuillageOut)
		sync <- true
	}()
	for i := 0; i < 3; i++ {
		<-sync
	}
}
func source(s chan int) {
	var valeur int = 5
	s <- valeur
}
func seuillage(seuil int, s1 chan int, s2 chan int) {
	pixelVal := <-s1
	if pixelVal < seuil {
		s2 <- 0
	} else {
		s2 <- 1
	}
}
func affichage(s3 chan int) {
	fmt.Println(<-s3)
}
