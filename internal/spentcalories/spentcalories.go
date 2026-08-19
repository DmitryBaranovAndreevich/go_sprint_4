package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	steps := 0
	activity := ""
	duration := time.Duration(0)

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return steps, activity, duration, fmt.Errorf("to many data: %d, data: %s", len(parts), data)
	}

	parseSteps, err := strconv.Atoi(parts[0])
	if err != nil || parseSteps <= 0 {
		return steps, activity, duration, fmt.Errorf("invalid steps: %s, data: %s", parts[0], data)
	}

	parseDuration, err := time.ParseDuration(parts[2])
	if err != nil || parseDuration <= 0 {
		return steps, activity, duration, fmt.Errorf("invalid duration: %s, data: %s", parts[1], data)
	}

	return parseSteps, parts[1], parseDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * height
	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	return distance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	switch activity {
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distance, speed, calories), nil

	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distance, speed, calories), nil

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid data: steps: %d, weight: %.2f, height: %.2f, duration: %s", steps, weight, height, duration)
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	spentCalories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return spentCalories * walkingCaloriesCoefficient, nil

}
