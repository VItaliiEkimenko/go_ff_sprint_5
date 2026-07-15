package actioninfo

import (
	"fmt"
	"log"
)

type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Printf("ошибка парсинга строки %q: %v", data, err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			log.Printf("ошибка формирования информации для строки %q: %v", data, err)
			continue
		}

		fmt.Println(info)
	}
}
