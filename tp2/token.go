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

	// Lancement des processus
	for i := 0; i < nbNodes; i++ {
		go func(id int) {
			// id est une copie de i au moment de la création

			node(id, tabCan[id], tabCan[(id+1)%nbNodes], id == 1)
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

// node est une fonction qui gère le comportement d'un nœud dans l'anneau
// Paramètres :
// - id : identifiant unique du nœud
// - in : canal d'entrée pour recevoir les messages du nœud précédent
// - out : canal de sortie pour envoyer des messages au nœud suivant
// - aLeJeton : booléen indiquant si ce nœud commence avec le jeton
func node(id int, in <-chan int, out chan<- int, aLeJeton bool) {

	// Initialisation de la valeur à envoyer (null par défaut)
	tokenAenvoyé := null // Au départ, tous les nœuds préparent null à envoyer

	// Si ce nœud doit commencer avec le jeton (cas du nœud 1)
	if aLeJeton {
		tokenAenvoyé = jeton // On prépare le jeton à envoyer
	}

	// Boucle principale qui s'exécute pour chaque tour
	for i := 0; i < nbTours; i++ {
		// PHASE 1: RÉCEPTION
		msgRecu := <-in // On attend et on reçoit un message du nœud précédent

		// PHASE 2: ENVOI
		out <- tokenAenvoyé // On envoie ce qu'on avait préparé (jeton ou null)

		// PHASE 3: MISE À JOUR DE L'ÉTAT
		if msgRecu == jeton { // Si on vient de recevoir le jeton
			tokenAenvoyé = jeton // On le prépare pour l'envoi au prochain tour
			// Affichage spécial pour noter la réception du jeton
			fmt.Printf("Node %d a reçu le jeton au tour %d\n", id, i)
		} else { // Si on n'a pas reçu le jeton
			tokenAenvoyé = null // On prépare null pour l'envoi au prochain tour
		}

		// PHASE 4: AFFICHAGE DE L'ÉTAT
		// Affiche l'identifiant du nœud, le numéro du tour et la valeur reçue
		fmt.Printf("id : %d, tour de boucle %d, valeur %d \n",
			id, i, msgRecu)
	}
}

/*
Réponses :
a) Oui, le processus 0 est toujours en avance sur les autres.
b) Oui, P0 a le temps d'attente le plus petit (0 secondes car idid = 00 = 0)
c) Non, P0 n'arrive pas à avoir plus d'un tour d'avance sur les autres. On peut voir qu'il commence toujours un nouveau tour, mais ne peut pas prendre plus d'avance à cause du modèle synchrone : il doit attendre de recevoir un message avant de pouvoir avancer, ce qui le synchronise avec les autres processus.
La communication synchrone agit comme un régulateur qui empêche un processus d'aller trop vite par rapport aux autres, même s'il a un temps de traitement plus court.
*/
