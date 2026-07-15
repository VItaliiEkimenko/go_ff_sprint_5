package trainings

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(dataString string) error {
	parts := strings.Split(dataString, ",")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат строки: ожидается 3 поля (шаги, тип, длительность), получено %d", len(parts))
	}

	var steps int
	_, err := fmt.Sscan(strings.TrimSpace(parts[0]), &steps)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге количества шагов: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть положительным")
	}
	t.Steps = steps

	t.TrainingType = strings.TrimSpace(parts[1])

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return fmt.Errorf("ошибка при парсинге длительности: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной")
	}
	t.Duration = duration

	return nil
}

func (t Training) ActionInfo() (string, error) {
		if t.TrainingType != "Ходьба" && t.TrainingType != "Бег" {
		return "", fmt.Errorf("неизвестный тип тренировки: %s", t.TrainingType)
	}

	if t.Steps <= 0 || t.Duration <= 0 {
		return "", fmt.Errorf("недопустимые данные тренировки")
	}

	durationHours := t.Duration.Hours()
	if durationHours == 0 {
		return "", fmt.Errorf("нулевая длительность тренировки")
	}

	distance := spentenergy.Distance(t.Steps, t.Height)

	var calories float64
	var err error

	switch t.TrainingType {
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", fmt.Errorf("неподдерживаемый тип тренировки: %s", t.TrainingType)
	}
	if err != nil {
		return "", err
	}

	speed := distance / durationHours

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType,
		durationHours,
		distance,
		speed,
		calories,
	), nil
}
