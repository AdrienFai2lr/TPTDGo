package main

import (
	"fmt"
	"time"
)

const nbNodes int = 4
const nbTours int = 4
const null int = -1

func main() {
	// Création des canaux de communication
	var tabCan [nbNodes + 1]chan int
	sync := make(chan bool)
	for i := range tabCan {
		tabCan[i] = make(chan int)
	}
	now := time.Now()
	// Lancement des processus
	for i := 0; i < nbNodes; i++ {
		go func(id int) {
			// id est une copie de i au moment de la création

			node(id, tabCan[i], tabCan[(i+1)%nbNodes], now)
			sync <- true
		}(i)
	}
	// Injection d'une première valeur pour démarrer le système
	tabCan[0] <- null
	for i := 1; i < nbNodes; i++ {
		<-sync
	}
	fmt.Println("fin du main")
}

func node(id int, in <-chan int, out chan<- int, now time.Time) {
	for i := 0; i < nbTours; i++ {
		// Phase d'échange des messages
		msgRecu := <-in // Réception
		out <- null     // Envoi du message null
		// Délai proportionnel au carré de l'id
		time.Sleep(time.Duration(id*id) * time.Second)
		// Phase de changement d'état
		fmt.Printf("id : %d, tour de boucle %d, valeur %d temps écoulé %v\n", id, i, msgRecu, time.Since(now))
	}
}

/*
Réponses :
a) Oui, le processus 0 est toujours en avance sur les autres.
b) Oui, P0 a le temps d'attente le plus petit (0 secondes car idid = 00 = 0)
c) Non, P0 n'arrive pas à avoir plus d'un tour d'avance sur les autres. On peut voir qu'il commence toujours un nouveau tour, mais ne peut pas prendre plus d'avance à cause du modèle synchrone : il doit attendre de recevoir un message avant de pouvoir avancer, ce qui le synchronise avec les autres processus.
La communication synchrone agit comme un régulateur qui empêche un processus d'aller trop vite par rapport aux autres, même s'il a un temps de traitement plus court.
*/
