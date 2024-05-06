package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Добро пожаловать в Калькулятор \n" +
		"Правила использования программы:\n" +
		"1. Введите выражение в одну строку.\n" +
		"2. Калькулятор принимает на вход числа от 1 до 10 включительно, не более.\n" +
		"3. Калькулятор умеет работать только с целыми числами.\n" +
		"4. Калькулятор умеет работать только с арабскими (1, 2, 3, 4, 5…) или с римскими (I, II, III, IV, V…) числами одновременно.\n" +
		"5. Калькулятор умеет выполнять операции сложения, вычитания, умножения и деления с двумя числами.\n" +
		"\n\n" +
		"Введите математическое выражение:")
	intType, first, second, sign, err := readLine()
	if err != nil {
		fmt.Println("Ошибка ввода данных:\n", err)
		return
	}
	if intType == "arab" {
		firstNum, err1 := strconv.Atoi(first)
		if err1 != nil {
			fmt.Println("Ошибка перевода строки в число:\n", err1)
			return
		}
		secondNum, err2 := strconv.Atoi(second)
		if err2 != nil {
			fmt.Println("Ошибка перевода строки в число:\n", err2)
			return
		}
		res, err3 := calculator(firstNum, secondNum, sign)
		if err3 != nil {
			fmt.Println("Ошибка работы программы:\n", err3)
			return
		} else {
			fmt.Println("Результат: ", res)
		}
	} else {
		firstNum := RomanInt(first)
		secondNum := RomanInt(second)
		res, err1 := calculator(firstNum, secondNum, sign)
		if err1 != nil {
			fmt.Println("Ошибка работы программы:\n", err1)
			return
		} else {
			final, err2 := IntRoman(res)
			if err2 != nil {
				fmt.Println("Ошибка работы программы:", err2)
				return
			}
			fmt.Println("Результат: ", final)
		}
	}
}

func calculator(first int, second int, sign string) (int, error) {
	if first > 10 || second > 10 {
		return 8, errorHandler(8)
	}
	switch {
	case sign == "+":
		return first + second, nil
	case sign == "-":
		return first - second, nil
	case sign == "*":
		return first * second, nil
	case sign == "/" && second != 0:
		return first / second, nil
	case sign == "/" && second == 0:
		return 4, errorHandler(4)
	default:
		return 5, errorHandler(5)
	}
}
func readLine() (string, string, string, string, error) {
	stdin := bufio.NewReader(os.Stdin)
	usInput, _ := stdin.ReadString('\n')
	usInput = strings.TrimSpace(usInput)
	intType, first, second, sign, err := checkInput(usInput)
	if err != nil {
		return "", "", "", "", err
	}
	return intType, first, second, sign, err
}

func checkInput(input string) (string, string, string, string, error) { //Проверка ввода данных
	r := regexp.MustCompile("\\s+")
	replace := r.ReplaceAllString(input, "")
	arr := strings.Split(replace, "")
	var intType, first, second, sign string
	for index, value := range arr {
		isN := Number(value)
		isS := Sign(value)
		isR := RomanNumber(value)
		if !isN && !isS && !isR {
			return "", "", "", "", errorHandler(1)
		}
		if isS {
			if sign != "" {
				return "", "", "", "", errorHandler(6)
			} else {
				sign = arr[index]
			}
		}
		if (isN && intType != "roman") || (isR && intType != "arab") {
			if intType == "" {
				if isN {
					intType = "arab"
				} else {
					intType = "roman"
				}
			}
			if first == "" && !(index+1 == len(arr)) && Sign(arr[index+1]) {
				slice := arr[0:(index + 1)]
				first = strings.Join(slice, "")
			} else if index+1 == len(arr) && first != "" {
				slice := arr[(len(first) + 1):]
				second = strings.Join(slice, "")
			}
		} else if (intType == "arab" && isR) || (intType == "roman" && isN) {
			return "", "", "", "", errorHandler(2)
		}
	}
	if second == "" || first == "" || sign == "" {
		return "", "", "", "", errorHandler(3)
	}
	return intType, first, second, sign, nil
}

func Number(c string) bool { //Проверка числа
	if c >= "0" && c <= "9" {
		return true
	} else {
		return false
	}
}

func Sign(c string) bool { //Проверка знака
	if c == "+" || c == "-" || c == "/" || c == "*" {
		return true
	} else {
		return false
	}
}
func RomanNumber(c string) bool { //Проверка числа Римское
	_, ok := dict[c]
	if ok {
		return true
	} else {
		return false
	}
}

func errorHandler(code int) error {
	return errors.New(errorDict[code])
}

var errorDict = map[int]string{
	1: "Неверные символы. Пожалуйста, используйте только арабские или римские цифры одновременно, а так же только математические операторы '+', '-', '/', '*' !",
	2: "Неверный ввод. Пожалуйста, используйте только арабские или только римские цифры! ",
	3: "Неверное количество аргументов. Для работы калькулятора необходимо два числа и математический оператор!",
	4: "Деление на 0!",
	5: "Ошибка вычисления!",
	6: "Неверное количество аргументов. Пожалуйста, введите только два числа и один математический оператор!",
	7: "Ошибка отображения ответа, так как в римской системе нет отрицательных чисел!",
	8: "Введите только числа от 0 до 10 включительно!",
}

var dict = map[string]int{
	"M":   1000,
	"CM":  900,
	"D":   500,
	"CD":  400,
	"C":   100,
	"XC":  90,
	"L":   50,
	"XL":  40,
	"X":   10,
	"IX":  9,
	"V":   5,
	"IV":  4,
	"III": 3,
	"II":  2,
	"I":   1,
}

func RomanInt(roman string) int { //Перевод в INT
	var res int
	arr := strings.Split(roman, "")
	for index, value := range arr {
		if index+1 != len(arr) && dict[value] < dict[arr[index+1]] {
			res -= dict[value]
		} else {
			res += dict[value]
		}
	}
	return res
}

func IntRoman(number int) (string, error) { //Перевод из INT
	if number <= 0 {
		return "", errorHandler(7)
	}
	arr1 := [15]int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 3, 2, 1}
	arr2 := [15]string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "III", "II", "I"}
	var str string
	for number > 0 {
		for i := 0; i < 15; i++ {
			if arr1[i] <= number {
				str += arr2[i]
				number -= arr1[i]
				break
			}
		}
	}
	return str, nil
}
