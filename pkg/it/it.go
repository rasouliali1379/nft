package it

import (
	"fmt"
	"log"
)

func Must(err error) {
	if err != nil {
		panic(fmt.Errorf("must failed: %w", err))
	}
}

func MustR[T any](t T, err error) T {
	Must(err)
	return t
}

func Should(err error) {
	if err != nil {
		log.Printf("should failed: %v", err)
	}
}

func ShouldR[T any](t T, err error) T {
	Should(err)
	return t
}
