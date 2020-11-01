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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeffotoni/gcolor"
)

type SendFile struct {
	Path      string
	TypeEvent string
}

var (
	WORKES = 50
)

func main() {

	// Rebendo path como parametro
	// Path será onde o sync irá
	// mapear para que tudo que
	// ocorra neste diretório
	// possa ser capturado
	// e enviado para Nuvem/Bucket etc...
	pathFile := flag.String("path", "", "nome do arquivo ou diretorio a ser enviado")
	flag.Parse()

	// Path é obrigatório
	if len(*pathFile) == 0 {
		flag.PrintDefaults()
		return
	}

	////////////////////////////////////////
	// aqui temos declarações de channels
	var done = make(chan bool, 1)
	var watcherEvent = make(chan string)
	var eventList = []string{"CREATE", "MODIFY"}
	// notify será implementado here
	// Aqui é uma simulação o notify
	// em 5 em 5 segundos ele lança um evento
	// Eventos são:
	//  CREATE -> arquivos em diretorios
	//  MODIFY -> atulizacao do arquivo
	go func() {
		for {
			watcherEvent <- "CREATE:/dir1/dir2/dir3/file1.pdf"
			<-time.After(time.Second * 5)
		}
	}()

	/// mapeamento dos diretorios para notify
	go func() {
		gcolor.PrintCyan("Inicio de mapeamento...")
		for {
			if err := filepath.Walk(*pathFile,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						gcolor.PrintRed("Error walk...")
						return err
					}
					// mappper here...
					println("mapper here:", path)
					return nil

				}); err != nil {
				gcolor.PrintRed("Error dir", err.Error())
			}
			gcolor.PrintCyan("Fim de mapeamento")
			<-time.After(time.Hour) // in hour
		}
	}()

	jobs := make(chan SendFile)

	// subindo workers
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

func worker(jobs <-chan SendFile) {
	for {
		select {
		case j := <-jobs:
			SendFileFile(j)
		}
	}
}

func SendFileFile(job SendFile) {
	gcolor.PrintYellow("enviando SendFilefile....:", job.Path, " -> ", job.TypeEvent)
	time.Sleep(time.Millisecond * 700)
}
