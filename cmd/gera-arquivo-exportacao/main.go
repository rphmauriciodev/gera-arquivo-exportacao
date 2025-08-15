package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/config"
	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/service"
)

func main() {
	configPath := flag.String("config", "../../config.json", "Arquivo de configuração JSON")
	server := flag.String("server", "", "Servidor SQL Server")
	user := flag.String("user", "", "Usuário do banco")
	pass := flag.String("pass", "", "Senha do banco")
	dbname := flag.String("db", "", "Nome do banco de dados")
	file := flag.String("file", "", "Arquivo com lista de procedures")
	outDir := flag.String("out", "", "Diretório de saída")
	procsFlag := flag.String("procs", "", "Lista de procedures separadas por vírgula")
	workersNum := flag.Int("workers", 5, "Número de workers concorrentes")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		cfg = &config.Config{}
	}

	config.OverrideConfig(cfg, server, user, pass, dbname, file, outDir, procsFlag, workersNum)

	procs, err := service.ProceduresFromConfig(cfg)

	if err != nil {
		log.Fatal(err)
	}

	tasks := make(chan string, len(procs))

	var wg sync.WaitGroup
	for i := 0; i < cfg.File.NumWorkers; i++ {
		wg.Add(1)
		go service.ExportProcedures(cfg, tasks, &wg)
	}

	for _, proc := range procs {
		tasks <- proc
	}
	close(tasks)

	wg.Wait()
	fmt.Println("Todas as procedures foram processadas!")
}
