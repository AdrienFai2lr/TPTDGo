package main
import (
	"fmt"
)
const finLigne = -1
const finImage = -2
const k = 5

func main() {
	const nbProc=3
	var tabCan[nbProc+1] chan int 
	for i :=range tabCan{
		tabCan[i]=make(chan int)
	}
	tabk:=[3]int{2,3,4}
	
	go source(tabCan[0])
	for i:=0;i<nbProc;i++{
		go traitement(tabk[i], tabCan[i], tabCan[i+1])
	}
	affichage(tabCan[nbProc])

}
func source(s chan int) {
	tabEntier := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, e := range tabEntier {
		s <- e
	}
	close(s)
}
func traitement(val int, s1 <-chan int, s2 chan<- int) {
	for i := range s1 {
		res := i + val
		s2 <- res
	}
	close(s2)
}
func affichage(s3 chan int) {
	for i := range s3{
		fmt.Println(i)
	}
}
