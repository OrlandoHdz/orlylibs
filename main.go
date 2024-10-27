package main

import (
	"fmt"
	"log"

	"github.com/OrlandoHdz/orlylibs/transfers"
)

func main() {

	fmt.Println("Iniciando prueba UploadFiles")

	td := transfers.Transfers{
		DestinationUser:   "uuuuuuuu",
		DestinationPass:   "ppppppp",
		DestinationHost:   "serverxxxxx.com",
		DestinationPort:   "22",
		DestinationFolder: "/recepcion/datos",
	}

	tt := transfers.NewTransfers(td)

	ff := []string{"/Users/orlando/Downloads/I.830.830_1729763293-3077914.EDI", "/Users/orlando/Downloads/DOC_ALT_1729763293-3077914_2.edi"}

	err := tt.UploadFiles(ff)

	if err != nil {
		log.Fatalf("Ocurrio un error: %v", err.Error())
	} else {
		fmt.Println("Proceso ejecutado")
	}

}
