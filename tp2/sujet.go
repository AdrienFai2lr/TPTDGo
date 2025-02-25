package main

import (
	"fmt"
	"sync"
)

const nbNodes int = 5

func node(id int, in, out []chan int, decLocale bool) {

	var messIn, messOut []int
	messIn = make([]int, len(in))
	messOut = make([]int, len(out))

	var decision bool

	//le premier mesg envoyé par un noeud à son voisin est égal a 1 avis favorable,-1 sinon
	if decLocale {
		messIn[0] = 1
	} else {
		messIn[0] = -1
	}
	//boucle
	for i := 0; i < nbNodes; i++ {
		communication(in, messIn, out, messOut)
		if decLocale {
			messOut[0] = messIn[0] + 1
		} else {
			messOut[0] = messIn[0] - 1
		}
	}
	decision = messIn[0] >= 0

	//ne pas modifier cette ligne finale
	fmt.Println("node:", id, "decision:", decision)
}

func main() {
	//valeurs de vote des noeuds : MODIFIABLE
	tabval := [nbNodes]bool{true, true, false, true, false}

	//topologie 5 noeuds en anneau unidirectionnel
	var tabCan [nbNodes][]chan int

	for i := range tabCan {
		tabCan[i] = make([]chan int, 1)
	}

	for i := range tabCan {
		tabCan[i][0] = make(chan int)
	}

	var wg sync.WaitGroup
	wg.Add(nbNodes)
	for i := 0; i < nbNodes; i++ {
		go func(i int) {
			node(i, tabCan[i], tabCan[(i+1)%nbNodes], tabval[i])
			wg.Done()
		}(i)

	}
	wg.Wait()

}

func communication(in []chan int, messIn []int, out []chan int, messOut []int) {
	var wg sync.WaitGroup
	send(messOut, out, &wg)
	receive(messIn, in, &wg)
	wg.Wait()
}

func receive(messIn []int, in []chan int, wg *sync.WaitGroup) {
	wg.Add(len(in))
	for i := 0; i < len(in); i++ {
		go func(i int) {
			messIn[i] = <-in[i]
			wg.Done()
		}(i)
	}
}

func send(messOut []int, out []chan int, wg *sync.WaitGroup) {
	wg.Add(len(out))
	for i := 0; i < len(out); i++ {
		go func(i int) {
			out[i] <- messOut[i]
			wg.Done()
		}(i)
	}
}
