package main

import "fmt"

func main() {
	c := make(chan int)
	d := make(chan int)
	s1 := make(chan bool)
	s2 := make(chan bool)
	go A(c, d, 6, s1)
	go A(d, c, 7, s2)
	<-s1
	<-s2
	fmt.Println("fin du main")
}

func A(c1, c2 chan int, init int, synch chan bool) {
	var x1, x2 int
	x2 = init
	go func() {
		c2 <- x2 // L'envoi se fait dans une goroutine
	}()
	x1 = <-c1 // Pendant que la réception continue
	fmt.Println(x1)
	synch <- true
}

/*

Schéma de l'organisation :

Main Program
     |
     |-- Channel c (int)
     |-- Channel d (int)
     |-- Channel s1 (bool)
     |-- Channel s2 (bool)
     |
     |-----> Goroutine A1 (init=6)
     |       - Reçoit de c
     |       - Envoie sur d
     |       - Envoie sur s1
     |
     |-----> Goroutine A2 (init=7)
             - Reçoit de d
             - Envoie sur c
             - Envoie sur s2

À l'exécution :

Le programme va se bloquer (deadlock)
Les deux goroutines A sont bloquées en attente de réception (x1 = <-c1)
Aucune ne peut progresser car chacune attend que l'autre envoie en premier
Aucun affichage ne sera effectué

Version modifiée pour résoudre le blocage :

go func main() {
    c := make(chan int)
    d := make(chan int)
    s1 := make(chan bool)
    s2 := make(chan bool)

    go func() {
        go A(c, d, 6, s1)
        c <- 0  // Envoi initial pour débloquer
    }()

    go func() {
        go A(d, c, 7, s2)
        d <- 0  // Envoi initial pour débloquer
    }()

    <-s1
    <-s2
    fmt.Println("fin du main")
}
La trace sera :
0
0
fin du main

Version réduite avec une seule goroutine :

go func main() {
    c := make(chan int)
    d := make(chan int)
    s1 := make(chan bool)
    s2 := make(chan bool)

    go func() {
        c <- 0  // Envoi initial
        d <- 0  // Envoi initial
        A(c, d, 6, s1)
        A(d, c, 7, s2)
    }()

    <-s1
    <-s2
    fmt.Println("fin du main")
}
Cette version :

Utilise une seule goroutine qui gère les deux appels à A
Fait les envois initiaux avant d'appeler les fonctions A
Produit la même trace que la version précédente
Est plus simple à comprendre et à maintenir
Évite la création de goroutines supplémentaires

Dans les deux cas (3 et 4), le blocage est résolu en initialisant les canaux avec des valeurs, permettant ainsi aux fonctions A de commencer leur exécution.*/
