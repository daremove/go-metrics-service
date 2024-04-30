// Пакет utils предоставляет утилитные функции и структуры, используемые во всем приложении.
package utils

import "time"

// Backoff предоставляет механизм для реализации экспоненциальной стратегии отката (backoff).
type Backoff struct {
	step     int             // Текущий шаг в стратегии отката.
	strategy []time.Duration // Список интервалов времени для каждого шага отката.
}

// NewBackoff создает новый экземпляр Backoff с предопределенной стратегией.
func NewBackoff() Backoff {
	return Backoff{
		step:     -1,
		strategy: []time.Duration{time.Second, 2 * time.Second, 5 * time.Second},
	}
}

// Duration возвращает следующую длительность отката и флаг, указывающий, можно ли продолжить.
func (b *Backoff) Duration() (time.Duration, bool) {
	b.step++

	if b.step > len(b.strategy)-1 {
		return 0, false
	}

	return b.strategy[b.step], true
}

// Reset сбрасывает счетчик шагов отката к начальному состоянию.
func (b *Backoff) Reset() {
	b.step = -1
}
