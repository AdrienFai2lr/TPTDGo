package main

import (
	"fmt"
)

const nbNodes int = 4
const nbTours int = 4
const null int = -1
const jeton int = -2

func main() {
	// Création des canaux de communication
	var tabCan [nbNodes + 1]chan int
	sync := make(chan bool)
	for i := range tabCan {
		tabCan[i] = make(chan int)
	}
	valeurLead := [4]int{15, 2, 15, 25}
	// Lancement des processus
	for i := 0; i < nbNodes; i++ {
		go func(id int) {
			// id est une copie de i au moment de la création

			node(id, tabCan[id], tabCan[(id+1)%nbNodes], valeurLead[id])
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

// node gère le comportement d'un nœud dans l'algorithme d'élection du leader
// Paramètres:
// - id: identifiant unique du nœud
// - in: canal d'entrée pour recevoir les messages du nœud précédent
// - out: canal de sortie pour envoyer des messages au nœud suivant
// - valeurLeader: valeur utilisée pour l'élection du leader
func node(id int, in <-chan int, out chan<- int, valeurLeader int) {
	// Variable indiquant si ce nœud est le leader (initialement false)
	estLeader := false
	// Nouvelle variable pour stocker l'id du leader
	idLeader := -1 // Initialement inconnu

	// Premier message à envoyer : sa propre valeur d'élection
	msgEnvoie := valeurLeader

	// Affichage du message de démarrage avec l'ID et la valeur d'élection
	fmt.Printf("node %d demarre avec la valeur : %d\n", id, valeurLeader)

	// Boucle principale : exécute le nombre de tours défini
	for i := 0; i < nbTours; i++ {
		// PHASE 1 : COMMUNICATION
		msgRecu := <-in  // Attend et reçoit un message du voisin précédent
		out <- msgEnvoie // Envoie le message préparé au voisin suivant

		// PHASE 2 : DÉCISION
		// Analyse le message reçu et décide de l'action à prendre
		switch {
		case msgRecu == null:
			// Si on reçoit null, on enverra null au prochain tour
			msgEnvoie = null

		case msgRecu == valeurLeader:
			// Si on reçoit sa propre valeur, on devient leader
			estLeader = true
			msgEnvoie = null // Le leader envoie null
			fmt.Printf("Node %d devient LEADER au tour %d\n", id, i)

		case msgRecu < valeurLeader:
			// Si la valeur reçue est inférieure, on l'ignore (envoie null)
			msgEnvoie = null

		case msgRecu > valeurLeader:
			// Si la valeur reçue est supérieure, on la propage
			msgEnvoie = msgRecu
		}

		// PHASE 3 : AFFICHAGE
		// Prépare le statut à afficher
		status := "non-leader"
		if estLeader {
			status = "LEADER"
		}

		// Affiche l'état complet du nœud pour ce tour
		fmt.Printf("id: %d, tour: %d, reçu: %d, statut: %s\n",
			id, i, msgRecu, status)
	}
	// Phase 2 : Diffusion de l'identifiant du leader
	msgEnvoie = null // Réinitialisation du message
	if estLeader {
		msgEnvoie = id // Le leader envoie son identifiant
	}

	// Nouvelle boucle pour la diffusion (même nombre de tours que la phase 1)
	for i := 0; i < nbTours; i++ {
		msgRecu := <-in
		out <- msgEnvoie

		// Si on reçoit un id valide (différent de null)
		if msgRecu != null {
			idLeader = msgRecu  // On mémorise l'id du leader
			msgEnvoie = msgRecu // On propage l'information
			fmt.Printf("Node %d a appris que le leader est %d au tour %d\n",
				id, idLeader, i)
		} else {
			msgEnvoie = null // Sinon on propage null
		}

		// Affichage avec l'id du leader connu
		status := "non-leader"
		if estLeader {
			status = "LEADER"
		}
		fmt.Printf("Phase 2 - id: %d, tour: %d, reçu: %d, statut: %s, leader connu: %d\n",
			id, i, msgRecu, status, idLeader)
	}
}
