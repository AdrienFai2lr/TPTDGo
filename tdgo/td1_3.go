package main

func main() {
	s1 := make(chan bool)
	s2 := make(chan bool)
	s3 := make(chan bool)
	s4 := make(chan bool)
	go func() {
		A()
		s1 <- true
	}()
	go func() {
		B()
		s2 <- true
	}()
	go func() {
		<-s1
		<-s2
		C()
		s3 <- true
	}()
	go func() {
		D()
		s4 <- true
	}()
	<-s3
	<-s4
}

// 2eme
func main() {
	s1 := make(chan bool)
	s2 := make(chan bool)
	s3 := make(chan bool)
	s4 := make(chan bool)
	go func() {
		A()
		s1 <- true
	}()
	go func() {
		B()
		s2 <- true
	}()
	go func() {
		<-s1
		<-s2
		C()
		s3 <- true
	}()
	<-s3
	go func() {
		D()
		s4 <- true
	}()
	<-s4
}

/*
Programme 2 (droite).

Analyse des ordres d'exécution :

Programme 1 :
Copy                Main
                 |
    +------------+------------+
    |            |           |
    A            B        Wait(s1,s2)
    |            |           |
    s1           s2          C
                             |
                             s3
                             |
                             D
                             |
                             s4
Programme 2 :
Copy                Main
                 |
    +------------+------------+
    |            |           |
    A            B        Wait(s1)
    |            |        Wait(s2)
    s1           s2          C
                             |
                             s3
                             |
                             D
                             |
                             s4
Différences clés :

Dans P1, la fonction C attend simultanément s1 et s2
Dans P2, la fonction C attend s1 puis s2 séquentiellement
Le reste des synchronisations est identique


Calcul des temps d'exécution :


A : 5 unités
B : 1 unité
C : 2 unités
D : 3 unités

Programme 1 :
Temps des fonctions :
A (5) et B (1) démarrent en parallèle à t=0
B finit à t=1
A finit à t=5
C (2) démarre après avoir reçu s1 ET s2, donc à t=5
C finit à t=7
D (3) démarre en parallèle avec C à t=0 (!)
D finit à t=3
Le programme attend s3 puis s4, donc finit à t=7

A et B s'exécutent en parallèle
Temps jusqu'à C = max(5, 1) = 5 unités
C commence après max(A,B) = 5 unités
C prend 2 unités
D commence après C = 7 unités
D prend 3 unités
Total = 10 unités

Programme 2 :
A (5) et B (1) démarrent en parallèle à t=0
B finit à t=1
A finit à t=5
C (2) démarre après avoir reçu s1 ET s2, donc à t=5
C finit à t=7
D (3) ne démarre qu'après réception de s3, donc à t=7
D finit à t=10

A et B s'exécutent en parallèle
C attend séquentiellement s1 puis s2 = max(5, 1) = 5 unités
C prend 2 unités
D commence après C = 7 unités
D prend 3 unités
Total = 10 unités

Conclusion :

Bien que les programmes aient des structures de synchronisation légèrement différentes
Dans ce cas précis avec ces temps d'exécution, ils prennent le même temps total
Cependant, dans le cas général, ils pourraient avoir des comportements différents selon les temps d'exécution des fonctions*/
