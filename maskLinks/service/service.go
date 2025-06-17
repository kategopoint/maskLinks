package service

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

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

	// Достаем []string из txt файла
	lines, err := s.prod.Produce()
	if err != nil {
		return fmt.Errorf("error producing data: %w", err)
	}

	// группа горутин с обработкой ошибок и лимитом в 10 одновременных задач
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(10)

	// bufferPool для эффективного использования памяти
	bufferPool := &sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}

	// буф.канальчик для результатов (емкость = кол-во строк)
	results := make(chan struct {
		index int
		line  string
	}, len(lines))

	// обрабатываем каждую строку в отдельной горутине
	for i, line := range lines {

		i, line := i, line // фиксируем значения для каждой горутины

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				buf := bufferPool.Get().(*bytes.Buffer) // берем буфер из пула
				buf.Reset()                             // очищаем буфер
				defer bufferPool.Put(buf)               // возвращаем буфер в пул через defer

				masked := s.MaskLinksWithBuffer(line, buf)

				// отправка результата в канал
				results <- struct {
					index int
					line  string
				}{i, masked}

				return nil
			}

		})
	}

	// ждем завершения и закрываем канал
	go func() {
		g.Wait()
		close(results)
	}()

	// собираем результаты из канала в слайс
	maskedLines := make([]string, len(lines))
	for res := range results {
		maskedLines[res.index] = res.line
	}

	// проверяем ошибки выполнения горутин
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error in worker: %w", err)
	}

	// Записываем в txt файл
	if err := s.pres.Present(maskedLines); err != nil {
		return fmt.Errorf("error presenting data: %w", err)
	}

	return nil
}

// maskLinksFunc маскирует ссылки каждой строки []string
// func (s *Service) MaskLinksSlice(input []string) []string {
// 	var maskedLines []string
// 	for _, line := range input {
// 		maskedLines = append(maskedLines, MaskLinks(line))
// 	}
// 	return maskedLines
// }

// maskLinks принимает строку и маскирует ссылки, начинающиеся с http://
func (s *Service) MaskLinksWithBuffer(input string, buf *bytes.Buffer) string {

	// buffer := make([]byte, 0, len(input))
	buf.Grow(len(input))

	i := 0
	for i < len(input) {
		// Проверяем, начинается ли текущая позиция с "http://"
		if i+7 <= len(input) && input[i:i+7] == "http://" {

			// Добавляем "http://" в буфер
			// buffer = append(buffer, input[i:i+7]...)
			buf.WriteString("http://")
			i += 7 // Пропускаем "http://"

			// Маскируем ссылку и добавляем '*' в буфер
			for (i < len(input)) && (input[i] != ' ') {
				// buffer = append(buffer, '*')
				buf.WriteByte('*')
				i++
			}

		} else {
			// Если не ссылка, просто добавляем байт в buffer
			// buffer = append(buffer, input[i])
			buf.WriteByte(input[i])
			i++
		}
	}

	// return string(buffer)
	return buf.String()
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
