package main

import (
	"fmt"
	"log"

	"github.com/OrlandoHdz/orlylibs/transfers"
)

func main() {

	fmt.Println("Iniciando prueba UploadFiles")

	td := transfers.Transfers{
		DestinationUser:   "ftp4772",
		DestinationPass:   "bRN0cJRK",
		DestinationHost:   "webconnect.seresnet.com",
		DestinationPort:   "22",
		DestinationFolder: "/recepcion/delfor_d04a",
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
