/*
canal envoiId (int)
Processus 1 -------------------------→

	Contrôleur

Processus 2 -------------------------→    ↓

Processus 3 -------------------------→

	canal receptionElu (int)

Processus 1 ←--------------------------

Processus 2 ←--------------------------

Processus 3 ←--------------------------
*/

// 2
package main

import "fmt"

func main() {
	// Déclaration des tableaux de canaux
	var envoiId [3]chan int      // Canal pour envoyer les IDs
	var receptionElu [3]chan int // Canal pour recevoir l'élu

	// Allocation des canaux
	for i := range envoiId {
		envoiId[i] = make(chan int)
		receptionElu[i] = make(chan int)
	}

	// Lancement des processus
	for i := 0; i < 3; i++ {
		go processus(i, envoiId[i], receptionElu[i])
	}

	// Lancement du contrôleur
	controleur(envoiId, receptionElu)
}

//3
// Processus utilisateur
func processus(id int, envoi chan int, reception chan int) {
	for {
		// Envoi de l'ID au contrôleur
		envoi <- id

		// Attente du résultat
		elu := <-reception

		// Affichage du résultat
		if elu == id {
			fmt.Printf("Processus %d : Je suis élu\n", id)
		} else {
			fmt.Printf("Processus %d : Le processus %d est élu\n", id, elu)
		}
	}
}

// Processus contrôleur
func controleur(envoi [3]chan int, reception [3]chan int) {
	for {
		maxId := -1

		// Réception des IDs
		for i := 0; i < 3; i++ {
			id := <-envoi[i]
			if id > maxId {
				maxId = id
			}
		}

		// Envoi du résultat à tous
		for i := 0; i < 3; i++ {
			reception[i] <- maxId
		}
	}
}
