package main

import "fmt"

func main() {
	var c chan int
	c = make(chan int)
	sync := make(chan bool)
	go func() {
		A(c)
		sync <- true
	}()
	go func() {
		B(c)
		sync <- true
	}()
	<-sync
	<-sync
	fmt.Println("fin du main")
}

func A(out chan int) {
	out <- 3
}

func B(in chan int) {
	var x int
	x = <-in
	fmt.Println(x)
}

/*
Lors de l'exécution du programme initial :

Le programme va se bloquer (deadlock)
Cela se produit car l'envoi dans A (out <- 3) est bloquant en attendant que quelqu'un reçoive
L'exécution est séquentielle, donc B n'est jamais atteint
Aucun affichage ne sera effectué

Pour résoudre ce problème avec une seule goroutine :


Il faut mettre le go devant l'appel à A(c), donc :

go A(c)
B(c)

Cela permet à A de s'exécuter en parallèle pendant que le programme principal continue vers B

Avec cette modification (un seul go), la trace sera :
3
fin du main
Car :

A s'exécute en parallèle et envoie 3
B reçoit le 3 et l'affiche
Le programme principal termine avec "fin du main"

Avec une goroutine pour chaque appel :

go A(c)
go B(c)
fmt.Println("fin du main")
La trace probable sera :
fin du main
Car :
Le programme principal termine avant que les goroutines aient eu le temps de s'exécuter
Les goroutines n'ont pas le temps de s'exécuter complètement

Pour obtenir l'exécution souhaitée, il faut :

Ajouter un mécanisme de synchronisation comme WaitGroup ou une syncro
Voici le code corrigé :

go func main() {
    var c chan int
    c = make(chan int)
    var wg sync.WaitGroup
    wg.Add(2)  // Pour attendre 2 goroutines

    go func() {
        A(c)
        wg.Done()
    }()

    go func() {
        B(c)
        wg.Done()
    }()

    wg.Wait()  // Attendre que les goroutines terminent
    fmt.Println("fin du main")
}
Avec cette version, on s'assure que :

Les goroutines ont le temps de s'exécuter
Le programme attend leur terminaison avant de finir
La trace sera maintenant :

3
fin du main

*/
