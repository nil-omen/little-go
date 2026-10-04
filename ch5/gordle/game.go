package gordle

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// errInvalidWordLength is returned when the guess has the wrong number of characters.
var errInvalidWordLength = fmt.Errorf("invalid guess, word doesn't have the same number of characters as the solution")

// Game holds all the information we need to play a game of gordle.
type Game struct {
	reader      *bufio.Reader
	solution    []rune
	maxAttempts int
}

// New returns a Game, which can be used to Play!
func New(playerInput io.Reader, corpus []string, maxAttempts int) (*Game, error) {
	if len(corpus) == 0 {
		return nil, ErrCorpusIsEmpty
	}
	g := &Game{
		reader:      bufio.NewReader(playerInput),
		solution:    splitToUppercaseCharacters(pickWord(corpus)),
		maxAttempts: maxAttempts,
	}

	return g, nil
}

// Play runs the game.
func (g *Game) Play() {
	fmt.Println("Welcome to Gordle!")

	// break condition: we've reached the maximum number of attempts
	for currentAttempt := 1; currentAttempt <= g.maxAttempts; currentAttempt++ {
		guess := g.ask()

		fb := computeFeedback(guess, g.solution)

		fmt.Println(fb.String())

		if slices.Equal(guess, g.solution) {
			fmt.Printf("🎉 You won! You found it in %d guess(es)! The word was: %s.\n", currentAttempt, string(g.solution))
			return
		}
	}
	// we've exhausted the number of allowed attempts
	fmt.Printf("😞 You've lost! The solution was: %s. \n", string(g.solution))

}

// ask reads input until a valid suggestion is made (and returned).
func (g *Game) ask() []rune {
	fmt.Printf("Enter a %d-charcter guess:\n", len(g.solution))
	for {
		playerInput, _, err := g.reader.ReadLine()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "Gordle failed to read your guess: %s\n", err.Error())
			continue
		}

		guess := splitToUppercaseCharacters(string(playerInput))

		// Verifying the guess length
		err = g.validateGuess(guess)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Your attempt is invalid with Gordle's solution: %s.\n", err.Error())
		} else {
			return guess
		}
	}
}

// validateGuess ensures the guess is valid enough.
func (g *Game) validateGuess(guess []rune) error {
	if len(guess) != len(g.solution) {
		return fmt.Errorf("expected %d, got %d, %w", len(g.solution), len(guess), errInvalidWordLength)
	}

	return nil
}

// splitToUppercaseCharacters is a naive implementation to turn a string into a list of characters.
func splitToUppercaseCharacters(input string) []rune {
	return []rune(strings.ToUpper(input))
}

// computeFeedback verifies every character of the guess against the solution.
func computeFeedback(guess, solution []rune) feedback {
	result := make(feedback, len(guess))
	used := make([]bool, len(solution))

	if len(guess) != len(solution) {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error! Guess and solution have different lengths: %d vs %d", len(guess), len(solution))
		return result
	}

	// check for correct letters
	for posInGuess, character := range guess {
		if character == solution[posInGuess] {
			result[posInGuess] = correctPosition
			used[posInGuess] = true
		}
	}

	// look for letters in wrong position
	for posInGuess, character := range guess {
		if result[posInGuess] != absentCharacter {
			// The character has already been marked, ignore it.
			continue
		}

		for posInSolution, target := range solution {
			if used[posInSolution] {
				// The letter of the solution is already assigned to a letter of the guess.
				// Skip to the next letter of the solution.
				continue
			}

			if character == target {
				result[posInGuess] = wrongPosition
				used[posInSolution] = true // Skip to the next letter of the guess.
				break

			}
		}
	}

	return result
}
