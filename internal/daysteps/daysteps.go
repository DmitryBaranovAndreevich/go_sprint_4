package daysteps

import (
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	steps := 0
	duration := time.Duration(0)

	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return steps, duration, fmt.Errorf("to many data: %d , data: %s", len(parts), data)
	}

	parseSteps, err := strconv.Atoi(parts[0])
	if err != nil {
		return steps, duration, fmt.Errorf("invalid steps: %s, data: %s", parts[0], data)
	}

	if parseSteps <= 0 {
		return steps, duration, fmt.Errorf("steps must more than zero")
	}

	parseDuration, err := time.ParseDuration(parts[1])
	if err != nil || parseDuration <= 0 {
		return steps, duration, fmt.Errorf("invalid duration: %s, data: %s", parts[1], data)
	}

	return parseSteps, parseDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)

	if err != nil {
		fmt.Println("error: ", err)
		return ""
	}

	distance := float64(steps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("error: ", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
