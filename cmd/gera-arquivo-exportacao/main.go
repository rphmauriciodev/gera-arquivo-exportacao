package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/config"
	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/db"
	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/repository"
	"github.com/rphmauriciodev/gera-arquivo-exportacao.git/internal/worker"
)

func main() {
	configPath := flag.String("config", "config.json", "Arquivo de configuração JSON")
	server := flag.String("server", "", "Servidor SQL Server")
	user := flag.String("user", "", "Usuário do banco")
	pass := flag.String("pass", "", "Senha do banco")
	dbname := flag.String("db", "", "Nome do banco de dados")
	file := flag.String("file", "", "Arquivo com lista de procedures")
	outDir := flag.String("out", "", "Diretório de saída")
	procsFlag := flag.String("procs", "", "Lista de procedures separadas por vírgula")
	workersNum := flag.Int("workers", 5, "Número de workers concorrentes")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		cfg = &config.Config{}
	}

	overrideConfig(cfg, server, user, pass, dbname, file, outDir, procsFlag, workersNum)

	if cfg.Server.Server == "" || cfg.Server.User == "" || cfg.Server.Password == "" || cfg.Server.Database == "" {
		log.Fatal("Erro: servidor, usuário, senha e banco são obrigatórios.")
	}

	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}
	defer conn.Close()

	var procs []string
	if cfg.File.Procs != "" {
		for _, p := range strings.Split(cfg.File.Procs, ",") {
			procs = append(procs, strings.TrimSpace(p))
		}
	} else if cfg.File.File != "" {
		procs, err = repository.ReadProceduresFromFile(cfg.File.File)
		if err != nil {
			log.Fatal("Erro ao ler arquivo:", err)
		}
	} else {
		log.Fatal("Nenhuma procedure informada (use -procs ou -file no config.json).")
	}

	os.MkdirAll(cfg.File.Out, os.ModePerm)

	tasks := make(chan string, len(procs))

	var wg sync.WaitGroup
	for i := 0; i < cfg.File.NumWorkers; i++ {
		wg.Add(1)
		go worker.Start(conn, cfg.File.Out, tasks, &wg)
	}

	for _, proc := range procs {
		tasks <- proc
	}
	close(tasks)

	wg.Wait()
	fmt.Println("Todas as procedures foram processadas!")
}

func overrideConfig(cfg *config.Config, server, user, pass, dbname, file, outDir, procsFlag *string, workersNum *int) {
	if *server != "" {
		cfg.Server.Server = *server
	}
	if *user != "" {
		cfg.Server.User = *user
	}
	if *pass != "" {
		cfg.Server.Password = *pass
	}
	if *dbname != "" {
		cfg.Server.Database = *dbname
	}
	if *file != "" {
		cfg.File.File = *file
	}
	if *outDir != "" {
		cfg.File.Out = *outDir
	}
	if *procsFlag != "" {
		cfg.File.Procs = *procsFlag
	}
	if *workersNum > 0 {
		cfg.File.NumWorkers = *workersNum
	}
	if cfg.File.NumWorkers == 0 {
		cfg.File.NumWorkers = 5
	}
}
