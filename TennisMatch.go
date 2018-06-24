package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func (pl *Player) makePlay(c chan bool) {
	// Generating seed by time
	rand.Seed(time.Now().UTC().UnixNano())

	// The player misses about 30% of times
	missed := rand.Intn(10) < 3
	c <- missed

}

//The player representation with their score total and position for the round
type Player struct {
	number   int
	points   int
	games    int
	sets     int
	position string
}

var (
	p1          = &Player{number: 1, games: 0, sets: 0, points: 0}
	p2          = &Player{number: 2, games: 0, sets: 0, points: 0}
	server      = p1
	receiver    = p2
	totalSets   = 3
	totalGames  = 6
	totalPontos = 4
)

func swapPosition() {
	if server == p1 {
		receiver = p1
		server = p2
	} else {
		server = p1
		receiver = p2
	}
}

func GotPoint(num int) {
	fmt.Printf("\n~~~~~~~~~~~~~\n")
	fmt.Printf("O jogador nº%d fez um ponto! \n", num)
	fmt.Printf("Jogador 1 - Pontos: %d \n", p1.points)
	fmt.Printf("Jogador 2 - Pontos: %d \n", p2.points)
	fmt.Printf("~~~~~~~~~~~~~\n")
}

func WonGame(num int) {
	fmt.Printf("\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("O jogador nº%d venceu o game! \n", num)
	fmt.Printf("Jogador 1 - Pontos: %d \n", p1.points)
	fmt.Printf("Jogador 2 - Pontos: %d \n", p2.points)
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
}

func WonSet(num int) {
	fmt.Printf("\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("O jogador nº%d venceu o set! \n", num)
	fmt.Printf("Jogador 1 - Games: %d \n", p1.games)
	fmt.Printf("Jogador 2 - Games: %d \n", p2.games)
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
}

func WonMatch(num int) {
	fmt.Printf("\n\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("O jogador nº%d venceu a partida! \n", num)
	fmt.Printf("Jogador 1 - Sets: %d \n", p1.sets)
	fmt.Printf("Jogador 2 - Sets: %d \n", p2.sets)
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
}

func main() {
	var mutx = &sync.Mutex{}
	ch := make(chan bool)
	dur := 500 * time.Millisecond

	// While no one won enough sets to win the match
	for p1.sets < totalSets && p2.sets < totalSets {

		// While no one won enough games to win a set
		for (p1.games < totalGames || p1.games-p2.games < 2) &&
			(p2.games < totalGames || p2.games-p1.games < 2) {

			// While no one got enough points to win a game
			for (p1.points < totalPontos || p1.points-p2.points < 2) &&
				(p2.points < totalPontos || p2.points-p1.points < 2) {

				// Go back and forth until someone scores
				for true {

					// Server makes first play
					mutx.Lock()
					go server.makePlay(ch)
					mutx.Unlock()

					// Server missed
					if <-ch {
						receiver.points++
						GotPoint(receiver.number)
						swapPosition()
						break
					}

					// Receiver makes second play
					mutx.Lock()
					go receiver.makePlay(ch)
					mutx.Unlock()

					// Receiver misses
					if <-ch {
						server.points++
						GotPoint(server.number)
						swapPosition()
						break
					}
					time.Sleep(dur)
				}
			}

			// Checks if p1 won a game
			if p1.points >= totalPontos && p1.points-p2.points >= 2 {
				p1.games++
				WonGame(p1.number)
			} else {
				p2.games++
				WonGame(p2.number)
			}
			// Resets points
			p1.points = 0
			p2.points = 0

		}

		// Checks if p1 won a set
		if p1.games >= totalGames && p1.games-p2.games >= 2 {
			p1.sets++
			WonSet(p1.number)
		} else {
			p1.sets++
			WonSet(p2.number)
		}
		// Resets games
		p1.games = 0
		p2.games = 0
	}

	// Checks if p1 won the match
	if p1.sets >= totalSets {
		WonMatch(p1.number)
	} else {
		WonMatch(p2.number)
	}

}
