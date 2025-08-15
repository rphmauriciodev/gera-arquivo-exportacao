package worker

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/repository"
)

func Start(db *sql.DB, outDir string, tasks <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for proc := range tasks {
		code, err := repository.GetProcedureCode(db, proc)
		if err != nil {
			log.Printf("Erro ao buscar procedure %s: %v\n", proc, err)
			continue
		}

		filePath := fmt.Sprintf("%s/%s.sql", outDir, proc)
		err = os.WriteFile(filePath, []byte(code), 0644)
		if err != nil {
			log.Printf("Erro ao salvar arquivo %s: %v\n", filePath, err)
			continue
		}

		fmt.Printf("Procedure %s exportada para %s\n", proc, filePath)
	}
}
