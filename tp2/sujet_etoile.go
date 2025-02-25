package main

import (
	"fmt"
	"sync"
)

const (
	nbNodes  = 5
	nbRounds = 3
	null     = -1
)

type Message struct {
	Source      int
	Destination int
	Data        int
	TTL         int
}

func node(id int, ins []chan Message, outs []chan Message, isCenter bool) {
	// Pour chaque tour
	for round := 0; round < nbRounds; round++ {
		// Phase d'envoi
		msg := Message{
			Source:      id,
			Destination: -1,
			Data:        id, // Envoi de son propre ID comme donnée
			TTL:         3,
		}

		// Envoi aux voisins
		for i := range outs {
			select {
			case outs[i] <- msg:
				fmt.Printf("Node %d sent message in round %d\n", id, round)
			default:
				fmt.Printf("Node %d couldn't send in round %d\n", id, round)
			}
		}

		// Phase de réception
		for i := range ins {
			select {
			case received := <-ins[i]:
				fmt.Printf("Node %d received message from %d in round %d\n",
					id, received.Source, round)
			default:
				fmt.Printf("Node %d no message to receive in round %d\n",
					id, round)
			}
		}
	}
	fmt.Printf("Node %d finished all rounds\n", id)
}

func main() {
	var wg sync.WaitGroup

	// Création des canaux avec buffer
	inToCenter := make([]chan Message, nbNodes-1)
	outFromCenter := make([]chan Message, nbNodes-1)

	// Initialisation des canaux avec buffer
	for i := 0; i < nbNodes-1; i++ {
		inToCenter[i] = make(chan Message, 1)
		outFromCenter[i] = make(chan Message, 1)
	}

	wg.Add(nbNodes)

	// Lancement du nœud central
	go func() {
		defer wg.Done()
		node(0, inToCenter, outFromCenter, true)
	}()

	// Lancement des nœuds périphériques
	for i := 1; i < nbNodes; i++ {
		go func(id int) {
			defer wg.Done()
			nodeIndex := id - 1
			node(id,
				[]chan Message{outFromCenter[nodeIndex]},
				[]chan Message{inToCenter[nodeIndex]},
				false)
		}(i)
	}

	wg.Wait()
}
func receive(ins []chan Message) []Message {
	messages := make([]Message, len(ins))
	var wg sync.WaitGroup

	wg.Add(len(ins))
	for i := range ins {
		go func(idx int) {
			defer wg.Done()
			messages[idx] = <-ins[idx]
		}(i)
	}
	wg.Wait()
	return messages
}

func send(messages []Message, outs []chan Message) {
	var wg sync.WaitGroup

	wg.Add(len(outs))
	for i := range outs {
		go func(idx int) {
			defer wg.Done()
			outs[idx] <- messages[idx]
		}(i)
	}
	wg.Wait()
}
