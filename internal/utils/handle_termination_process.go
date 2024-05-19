// Package utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// HandleTerminationProcess устанавливает обработчик для сигналов прерывания и завершения работы,
// вызывая переданную функцию cleanup при получении этих сигналов.
func HandleTerminationProcess(cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}
