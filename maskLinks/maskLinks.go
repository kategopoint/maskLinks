package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// maskLinks принимает строку и маскирует ссылки, начинающиеся с http://
func maskLinks(input string) string {
	buffer := make([]byte, 0, len(input))

	// Перебираем байты входной строки
	i := 0
	for i < len(input) {
		// Проверяем, начинается ли текущая позиция с "http://"
		if i+7 <= len(input) && input[i:i+7] == "http://" {

			// Добавляем "http://" в буфер
			buffer = append(buffer, input[i:i+7]...)
			i += 7 // Пропускаем "http://"

			// Маскируем ссылку и добавляем * в buffer
			for (i < len(input)) && (input[i] != ' ') {
				buffer = append(buffer, '*')
				i++
			}

		} else {
			// Если не ссылка, просто добавляем байт в buffer
			buffer = append(buffer, input[i])
			i++
		}
	}

	return string(buffer)
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// inputmaskLinks := "Here's my spammy page: http://hehefouls.netHAHAHA see you."
	// output := maskLinks(inputmaskLinks)

	output := maskLinks(strings.TrimSpace(input)) // Удаляем лишние пробелы в начале и конце строки

	fmt.Println(output)
}
