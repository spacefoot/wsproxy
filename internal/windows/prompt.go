//go:build windows

package windows

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/huh"
)

type Action int

const (
	INSTALL Action = iota
	REMOVE
	START
	STOP
	RESTART
	EXIT
)

func Prompt() error {
	var input Action
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Action]().
				Title("Action").
				Options(
					huh.NewOption("Démarrer", START),
					huh.NewOption("Arrêter", STOP),
					huh.NewOption("Redémarrer", RESTART),
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
		fmt.Println("Service arrêté")
	case RESTART:
		if err := Stop(); err != nil {
			return err
		}
		fmt.Println("Service arrêté")
		if err := Start(); err != nil {
			return err
		}
		fmt.Println("Service démarré")
	case EXIT:
		return nil
	default:
		return errors.New("Invalid input")
	}

	time.Sleep(2 * time.Second)

	return nil
}
