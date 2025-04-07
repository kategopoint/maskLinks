package main

import (
	"fmt"
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
			for (i < len(input)) && (input[i] != ' ') { // || (input[i] != '\n')
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
	inputmaskLinks := "Here's my spammy page: http://hehefouls.netHAHAHA see you."
	output := maskLinks(inputmaskLinks)
	fmt.Println(output) // Ожидаемый вывод: Here's my spammy page: http://******************* see you.
}
