package presenter

import "os"

type FilePresenter struct {
	filePath string
}

func NewFilePresenter(filePath string) *FilePresenter {
	return &FilePresenter{filePath: filePath}
}

func (fp *FilePresenter) Present(lines []string) error {
	file, err := os.Create(fp.filePath)
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
