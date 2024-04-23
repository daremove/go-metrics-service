// Пакет utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import "sync"

// Parallelize выполняет несколько функций параллельно.
// Эта функция принимает произвольное количество аргументов функций, которые не принимают параметров и не возвращают значения.
// Она использует sync.WaitGroup для ожидания завершения всех переданных функций.
func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}
