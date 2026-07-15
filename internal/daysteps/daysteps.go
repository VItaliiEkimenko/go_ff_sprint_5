package daysteps

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type DaySteps struct {
	Steps      int
	Duration   time.Duration
	personaldata.Personal
}

func (ds *DaySteps) Parse(datastring string) (err error) {
	// 1. Сразу отвергаем любые пробельные символы во всей строке
	if strings.ContainsAny(datastring, " \t\r\n") {
		return fmt.Errorf("в строке не должно быть пробелов или других пробельных символов")
	}

	if datastring == "" {
		return fmt.Errorf("пустая строка")
	}

	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("неверный формат строки: ожидается 2 поля (шаги, длительность), получено %d", len(parts))
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	// 2. Валидация шагов: разрешаем опциональный '+' в начале, дальше только цифры
	if stepsStr == "" {
		return fmt.Errorf("отсутствует значение количества шагов")
	}

	i := 0
	// Если первый символ '+' — пропускаем его
	if len(stepsStr) > 0 && stepsStr[0] == '+' {
		i = 1
	}

	// После возможного '+' должны идти только цифры
	for j := i; j < len(stepsStr); j++ {
		r := stepsStr[j]
		if r < '0' || r > '9' {
			return fmt.Errorf("количество шагов должно содержать только цифры")
		}
	}

	// Запрет ведущего нуля: "0123" — ошибка. "0" технически допустим синтаксически,
	// но будет отвергнут проверкой steps <= 0 ниже.
	if i == 0 && len(stepsStr) > 1 && stepsStr[0] == '0' {
		return fmt.Errorf("недопустимый формат количества шагов: ведущий ноль")
	}
	if i == 1 && len(stepsStr) > 2 && stepsStr[1] == '0' {
		// Случай "+0123" тоже считаем ошибкой (ведущий ноль после плюса)
		return fmt.Errorf("недопустимый формат количества шагов: ведущий ноль")
	}

	var steps int
	_, err = fmt.Sscanf(stepsStr, "%d", &steps)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге количества шагов: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть положительным")
	}
	ds.Steps = steps

	// 3. Длительность: НЕ делаем TrimSpace — пробелы должны быть ошибкой
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге длительности: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной")
	}
	ds.Duration = duration

	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	if ds.Steps <= 0 || ds.Duration <= 0 {
		return "", fmt.Errorf("недопустимые данные тренировки")
	}

	distance := spentenergy.Distance(ds.Steps, float64(ds.Height))

	calories, err := spentenergy.WalkingSpentCalories(
		ds.Steps,
		ds.Weight,
		ds.Height,
		ds.Duration,
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps,
		distance,
		calories,
	), nil
}
