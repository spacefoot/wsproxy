//go:build windows

package windows

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/huh"
)

const (
	INSTALL int = iota
	REMOVE
	START
	STOP
	EXIT
)

func Prompt() error {
	var input int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Action").
				Options(
					huh.NewOption("Démarrer", START),
					huh.NewOption("Arrêter", STOP),
					huh.NewOption("Installer", INSTALL),
					huh.NewOption("Supprimer", REMOVE),
					huh.NewOption("Exit", EXIT),
				).
				Value(&input),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch input {
	case INSTALL:
		if err := Install(); err != nil {
			return err
		}
		fmt.Println("Service installé")
		if err := Start(); err != nil {
			return err
		}
		fmt.Println("Service démarré")
	case REMOVE:
		_ = Stop()
		if err := Remove(); err != nil {
			return err
		}
		fmt.Println("Service marquer pour suppression. Effectif au prochain redémarrage du système.")
	case START:
		if err := Start(); err != nil {
			return err
		}
		fmt.Println("Service démarré")
	case STOP:
		if err := Stop(); err != nil {
			return err
		}
		fmt.Println("Service Arrêté")
	case EXIT:
		return nil
	default:
		return errors.New("Invalid input")
	}

	time.Sleep(2 * time.Second)

	return nil
}
