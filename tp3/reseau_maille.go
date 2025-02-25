package main

import (
	"fmt"
	"sync"
)

const (
	nbNodes = 7  // 7 nœuds dans le réseau
	nbChan  = 14 // Nombre total de canaux pour un réseau maillé avec au moins 2 voisins par nœud
)

// Pour une topologie avec 7 nœuds où chaque nœud a au moins 2 voisins
// Par exemple:
// 0 connecté à 1,2
// 1 connecté à 0,2,3
// 2 connecté à 0,1,4
// 3 connecté à 1,4,5
// 4 connecté à 2,3,6
// 5 connecté à 3,6
// 6 connecté à 4,5
func receive(messIn []int, in []chan int, wg *sync.WaitGroup) {
	wg.Add(len(in))
	for i := 0; i < len(in); i++ {
		go func(i int) {
			defer wg.Done()
			messIn[i] = <-in[i]
		}(i)
	}
}

func send(messOut []int, out []chan int, wg *sync.WaitGroup) {
	wg.Add(len(out))
	for i := 0; i < len(out); i++ {
		go func(i int) {
			defer wg.Done()
			out[i] <- messOut[i]
		}(i)
	}
}

func communication(in []chan int, messIn []int, out []chan int, messOut []int) {
	var wg sync.WaitGroup

	send(messOut, out, &wg)
	receive(messIn, in, &wg)

	wg.Wait()
}
func node(id int, in []chan int, out []chan int) {
	// Préparation des tableaux de messages
	messIn := make([]int, len(in))
	messOut := make([]int, len(out))

	// Préparation des messages à envoyer (son propre ID)
	for i := range messOut {
		messOut[i] = id
	}

	// Communication
	communication(in, messIn, out, messOut)

	// Affichage des voisins découverts
	fmt.Printf("Node %d - Voisins découverts : ", id)
	for port, neighborId := range messIn {
		fmt.Printf("port %d -> noeud %d, ", port, neighborId)
	}
	fmt.Println()
}
func main() {
	// Création du tableau global de canaux
	var tabChan [nbChan]chan int
	for i := range tabChan {
		tabChan[i] = make(chan int)
	}

	// Création des tableaux de communication
	var comIn [nbNodes][]chan int
	var comOut [nbNodes][]chan int

	// Définition de la topologie
	// Node 0 : connecté à 1,2
	comIn[0] = []chan int{tabChan[0], tabChan[2]}
	comOut[0] = []chan int{tabChan[1], tabChan[3]}

	// Node 1 : connecté à 0,2,3
	comIn[1] = []chan int{tabChan[1], tabChan[4], tabChan[6]}
	comOut[1] = []chan int{tabChan[0], tabChan[5], tabChan[2]}

	// Node 2 : connecté à 0,1,4
	comIn[2] = []chan int{tabChan[3], tabChan[5], tabChan[0]}
	comOut[2] = []chan int{tabChan[2], tabChan[4], tabChan[1]}

	// Node 3 : connecté à 1,5
	comIn[3] = []chan int{tabChan[2], tabChan[4]}
	comOut[3] = []chan int{tabChan[3], tabChan[5]}

	// Node 4 : connecté à 2,6
	comIn[4] = []chan int{tabChan[1], tabChan[6]}
	comOut[4] = []chan int{tabChan[0], tabChan[2]}

	// Node 5 : connecté à 3,6
	comIn[5] = []chan int{tabChan[5], tabChan[3]}
	comOut[5] = []chan int{tabChan[4], tabChan[6]}

	// Node 6 : connecté à 4,5
	comIn[6] = []chan int{tabChan[2], tabChan[6]}
	comOut[6] = []chan int{tabChan[1], tabChan[5]}
	// ... définir les autres connexions selon votre schéma

	// Lancement des goroutines
	var wg sync.WaitGroup
	wg.Add(nbNodes)

	for i := 0; i < nbNodes; i++ {
		go func(j int) {
			defer wg.Done()
			node(j, comIn[j], comOut[j])
		}(i)
	}

	wg.Wait()
}
