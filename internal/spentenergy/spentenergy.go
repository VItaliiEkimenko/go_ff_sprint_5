package spentenergy

import (
	"fmt"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm = 1000 // количество метров в километре.
	minInH = 60   // количество минут в часе.
	stepLengthCoefficient = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

func Distance(steps int, height float64) float64 {
	stepsFloat := float64(steps)
	stepLength := height * stepLengthCoefficient
	distanceInMeters := stepsFloat * stepLength
	return distanceInMeters / mInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceKm := Distance(steps, height)
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distanceKm / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов не может быть отрицательным: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным: %f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным: %f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной: %v", duration)
	}

	meanSpeedVal := MeanSpeed(steps, height, duration)
	if meanSpeedVal == 0 {
		return 0, nil
	}

	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeedVal * durationInMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверки входных данных — важно, чтобы на нулевых шагах была ошибка
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным: %f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным: %f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной: %v", duration)
	}

	distanceKm := Distance(steps, height)
	calories := weight * distanceKm * walkingCaloriesCoefficient
	return calories, nil
}

