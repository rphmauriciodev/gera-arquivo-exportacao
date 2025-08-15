package service

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/config"
	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/repository"
)

func ExportProcedures(cfg *config.Config, tasks <-chan string, wg *sync.WaitGroup) {

	defer wg.Done()

	outDir := cfg.File.Out

	db, err := config.ConnectDb(cfg)

	if err != nil {
		log.Fatal(err)
	}

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

func ProceduresFromConfig(cfg *config.Config) ([]string, error) {

	var procs []string
	var err error

	if cfg.File.Procs != "" {
		for _, p := range strings.Split(cfg.File.Procs, ",") {
			procs = append(procs, strings.TrimSpace(p))
		}
	} else if cfg.File.File != "" {
		procs, err = readProceduresFromFile(cfg.File.File)
		if err != nil {
			log.Fatal("Erro ao ler arquivo:", err)
		}
	} else {
		err = errors.New("nenhuma procedure informada (use -procs ou -file no config.json)")
		log.Fatal(err)
	}
	return procs, err
}

func readProceduresFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var procs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			procs = append(procs, line)
		}
	}
	return procs, scanner.Err()
}
