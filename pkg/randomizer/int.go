package randomizer

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func GenerateRandomNumber(min, max int64) (int64, error) {
	if min >= max || min < 0 {
		return 0, errors.New("некоректний діапазон")
	}

	// Обчислюємо довжину діапазону.
	diff := big.NewInt(max - min)

	// Генеруємо випадкове число.
	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0, err
	}

	// Додаємо мінімальне значення до отриманого числа, щоб знаходитись в заданому діапазоні.
	randomNum := n.Int64() + min
	return randomNum, nil
}
