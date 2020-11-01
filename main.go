// Este programa é um Esboço
// um simples Rascunho de um Sync
// Ele irá sincronizar arquivos com
// nuvem, vamos iniciar com Buckets
// S3 da AWS e DigitalOcean
// Mas pode estender para qualquer
// Serviços de Buckets Ex:
// Google Driver, Mega, Azure etc..

package main

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jeffotoni/gcolor"
	"github.com/jeffotoni/gconcat"
)

type SendFile struct {
	Path      string
	TypeEvent string
}

type IntRange struct {
	min, max int
}

var (
	WORKES = 10
	m      sync.Map
)

func main() {

	// Rebendo path como parametro
	// Path será onde o sync irá
	// mapear para que tudo que
	// ocorra neste diretório
	// possa ser capturado
	// e enviado para Nuvem/Bucket etc...
	pathFile := flag.String("path", "", "nome do arquivo ou diretorio a ser enviado")
	workerFlag := flag.Int("worker", WORKES, "worker concorrentes")
	flag.Parse()

	// Path é obrigatório
	if len(*pathFile) == 0 {
		flag.PrintDefaults()
		return
	}

	if *workerFlag > 0 {
		WORKES = *workerFlag
	}

	/// mapeamento dos diretorios para notify
	// este mapeamento é simplesmente uma
	//  das formas de escutar atraves
	// de eventos do kernel do S.O
	// para que possamos caputrar
	// e saber o que ocorreu nos diretórios
	go func() {
		gcolor.PrintCyan("Inicio mapeamento...")
		for {
			i := 0
			if err := filepath.Walk(*pathFile,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						gcolor.PrintRed("Error walk...")
						return err
					}

					// Armazenando o path
					// somente uma simulacao
					// do mapeamento que
					// iremos fazer para notify
					if isDir(path) {
						m.Store(i, path)
						i++
					}
					return nil

				}); err != nil {
				gcolor.PrintRed("Error dir", err.Error())
			}
			gcolor.PrintCyan("Fim mapeamento\n")
			<-time.After(time.Hour) // in hour
		}
	}()

	// aguardando o mapeamento
	time.Sleep(time.Second * 5)

	////////////////////////////////////////
	// aqui temos declarações de channels
	var a int = 1
	var b int = 93
	var done = make(chan bool, 1)
	jobs := make(chan SendFile)
	var watcherEvent = make(chan string)
	var eventList = []string{"CREATE", "MODIFY"}

	// notify será implementado here
	// Aqui é uma simulação o notify
	// em 5 em 5 segundos ele lança um evento
	// Eventos são:
	//  CREATE -> arquivos em diretorios
	//  MODIFY -> atulizacao do arquivo
	go func() {
		i := 0
		for {
			i++
			go func() {
				////////////////////////////////////////
				////////////////////////////////////////
				////////////////////////////////////////
				// simulacao de
				// geracao event
				// Ele pega os paths mapeados e captura
				// quando ouver alterações
				// lembrando é uma simulação do notify
				event := "MODIFY"
				pathEvent := "/dir1/dir2/"
				rand.Seed(time.Now().UnixNano())
				n := a + rand.Intn(b-a+1)
				if n%2 == 0 {
					event = eventList[0]
				}

				dirlist, ok := m.Load(n)
				if ok {
					pathEvent = dirlist.(string)
				}
				////////////////////////////////////////

				// colocando o evento no CHANNEL
				// apos receber no channel o mesmo
				// sera executando por um Worker
				watcherEvent <- gconcat.Build(event, ":", pathEvent, "file_", n, ".pdf")
			}()

			if i == WORKES {
				i = 0
				time.Sleep(time.Second * 5)
				println("........................................................................................")
			}
		}
	}()

	// Aqui esta subindo diversos
	// workers para executar os
	// jobs
	for w := 1; w <= WORKES; w++ {
		go worker(jobs)
	}

	// escutando eventos notify
	go func() {
		for {
			select {
			case event := <-watcherEvent:
				vet := strings.Split(event, ":")
				typeEvent := vet[0]
				pathfile := vet[1]
				jobs <- SendFile{
					Path:      pathfile,
					TypeEvent: typeEvent,
				}
			}
		}
	}()
	<-done
}

// Aqui onde irá descarregar
// os channels
// o que entrar na fila
// será executado aqui
// Este Worker está aguardando
// que a fila seja preenchida
// o worker foi instanciado
// no inicio, foi feito
// dinamicamente
func worker(jobs <-chan SendFile) {
	for {
		select {
		case j := <-jobs:
			SendFileFile(j)
		}
	}
}

// Este metodo é responsavel
// por enviar o arquivo para
// nuvem, seja bucket aws s3
// seja space do digitalOcean
// ou qualquer outro que desejar
// utilizar.
func SendFileFile(job SendFile) {
	gcolor.PrintYellow("SendFile:", job.Path, " -> ", job.TypeEvent)
	time.Sleep(time.Millisecond * 700)
}

// get Random
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

// validando se é diretorio
func isDir(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

// validando se é um arquivo
func fileExist(name string) bool {
	if stat, err := os.Stat(name); err == nil && !stat.IsDir() {
		return true
	}
	return false
}
