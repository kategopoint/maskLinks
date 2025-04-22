package service

import "fmt"

type Producer interface {
	Produce() ([]string, error)
}

type Presenter interface {
	Present([]string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

// NewService конструктор для создания нового сервиса
func NewService(prod Producer, pres Presenter) *Service {
	return &Service{prod: prod, pres: pres}
}

// главный метод запуска
func (s *Service) Run() error {
	lines, err := s.prod.Produce()
	if err != nil {
		return fmt.Errorf("error producing data: %w", err)
	}

	maskedLines := s.maskLinksSlice(lines)

	if err := s.pres.Present(maskedLines); err != nil {
		return fmt.Errorf("error presenting data: %w", err)
	}

	return nil
}

// maskLinksFunc маскирует ссылки каждой строки []string
func (s *Service) maskLinksSlice(input []string) []string {
	var maskedLines []string
	for _, line := range input {
		maskedLines = append(maskedLines, maskLinks(line))
	}
	return maskedLines
}

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

// func (s *Service) Produce() ([]string, error) {

// 	data, err := os.ReadFile("linesForFunc.txt")

// 	if err != nil {
// 		crtFile, err := os.Create("linesForFunc.txt")
// 		if err != nil {
// 			panic(err)
// 		}
// 		crtFile.WriteString("Here's my spammy page: http://hehefouls.netHAHAHA see you.\nNo links here!")
// 		defer crtFile.Close()
// 	}

// 	lines := bytes.Split(data, []byte{'\n'})

// 	var dataLines []string
// 	for _, line := range lines {
// 		dataLines = append(dataLines, string(line))
// 	}

// 	return dataLines, err
// }
