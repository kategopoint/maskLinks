package presenter

import "os"

type FilePresenter struct {
	filePath string
}

func NewFilePresenter(filePath string) *FilePresenter {
	return &FilePresenter{filePath: filePath}
}

func (fp *FilePresenter) Present(lines []string) error {

	// os.O_RDWR	Открыть файл для чтения и записи (Read + Write)
	// os.O_CREATE	Создать файл, если он не существует
	// os.O_TRUNC	Очистить содержимое файла (обрезка до нулевой длины), если он существует
	// права доступа (0644 — владелец читает/пишет, остальные только читают)
	file, err := os.OpenFile(fp.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}
