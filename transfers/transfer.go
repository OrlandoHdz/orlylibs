// Transfiere archivos entre servidores ftp o sftp

package transfers

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Transfers struct {
	DestinationUser    string
	DestinationPass    string
	DestinationHost    string
	DestinationPort    string
	DestinationFolder  string
	destinationHostKey ssh.PublicKey
}

// NewTransfers crea una nueva  Transfer que contiene los valoes necesarios
func NewTransfers(trans Transfers) Transfers {

	hostKey := getsHostKey(trans.DestinationHost)

	t := Transfers{
		DestinationUser:    trans.DestinationUser,
		DestinationPass:    trans.DestinationPass,
		DestinationHost:    trans.DestinationHost,
		DestinationPort:    trans.DestinationPort,
		DestinationFolder:  trans.DestinationFolder,
		destinationHostKey: hostKey,
	}

	return t

}

// UploadFiles sube archivos al servidor destino
func (t *Transfers) UploadFiles(files []string) error {

	// Abre conexion al servidor seres
	server, client, err := t.getServerConnection()
	if err != nil {
		log.Println("Error al abrir conexion al servidor")
		return err
	}
	defer client.Close()
	defer server.Close()

	for _, file := range files {
		sFile := strings.Split(file, "/")
		fileDestination := t.DestinationFolder + "/" + sFile[len(sFile)-1]
		// Crea el archivo fisico en el servior
		aFile, err := server.Create(fileDestination)
		if err != nil {
			log.Println("Error al crear archivo fisico:", fileDestination)
			return err
		}
		defer aFile.Close()
		// Abre el archivo origen
		fileOrigin, err := os.Open(file)
		if err != nil {
			log.Println("Error al abrir archivo origen:", file)
			return err
		}
		defer fileOrigin.Close()
		// Copia archivo origen a destino
		_, err = io.Copy(aFile, fileOrigin)
		if err != nil {
			log.Println("Error al copiar el archivo a seres:")
			return err
		}
	}

	return nil
}

// getServerConnection  retorna la conexion del servidor destino sftp
func (t *Transfers) getServerConnection() (*sftp.Client, *ssh.Client, error) {

	config := ssh.ClientConfig{
		User: t.DestinationUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(t.DestinationPass),
		},
		HostKeyCallback: ssh.FixedHostKey(t.destinationHostKey),
	}

	// connect
	conn, err := ssh.Dial("tcp", t.DestinationHost+":"+t.DestinationPort, &config)
	if err != nil {
		return nil, nil, err
	}

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, nil, err
	}

	return client, conn, err
}

// getsHostKey obtiene el hostkey del servidor remoto, esto es para los SFTP
func getsHostKey(host string) ssh.PublicKey {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("Error obteniendo llave publica %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("No se encontro la llave publica para: %s", host)
	}

	return hostKey
}
