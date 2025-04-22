package service

import (
	"bufio"
	"os"
)

type FileProducer struct {
	filePath string
}

// NewFileProducer конструкттор для FileProducer
func NewFileProducer(filePath string) *FileProducer {
	return &FileProducer{filePath: filePath}
}

// чтение строки из файла
func (fp *FileProducer) Produce() ([]string, error) {
	file, err := os.Open(fp.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	// построчное чтение данных из файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // читает следующую строку из файла
		lines = append(lines, scanner.Text()) // возвращает строку, которую только что прочитал сканер
	}

	// если сканер ничего не прочитал
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
