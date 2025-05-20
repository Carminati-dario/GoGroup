// Client multichat
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	connHost = "127.0.0.1" // host server
	connPort = "8080"      // porta
	connType = "tcp"       // tipo connessione
)

func main() {
	fmt.Println("Connetto in " + connType + " al server multichat " + connHost + ":" + connPort + " (scrivere \"esci\" per terminare)")

	connIniziale, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Errore di connessione:", err.Error())
		os.Exit(1)
	}

	fmt.Println(connIniziale.LocalAddr().String() + ": connesso al server " + connIniziale.RemoteAddr().String())

	// Attende il messaggio di riposta del server, il quale conterrà nuova porta da utilizzare. In tale modo la porta connPort potrà essere liberata per altre connessioni
	nuovaPorta, err := bufio.NewReader(connIniziale).ReadString('\n')
	if err != nil {
		fmt.Println("Errore nel recupero della nuova porta:", err.Error())
		os.Exit(1)
	}

	if strings.TrimSpace(nuovaPorta) == "Chat al completo!" { // Non ci sono slot liberi, quindi termino
		// Notare che viene inviato l'errore come un comune messaggio di testo
		fmt.Println("Chat al completo! Riprovare in seguito.")
		os.Exit(1)
	}

	fmt.Println(connIniziale.LocalAddr().String() + ": nuovo socket comunicato dal server: " + connHost + ":" + nuovaPorta)
	connIniziale.Close()

	nuovoSocket := net.JoinHostPort(connHost, strings.TrimSpace(nuovaPorta))
	conn, err := net.Dial(connType, nuovoSocket)
	if err != nil {
		fmt.Println("Errore di connessione 2:", err.Error())
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Errore in creazione buffer di invio:", err.Error())
		os.Exit(1)
	}
	fmt.Fprintf(conn, "screen\n")

	for {

		fmt.Print(">> ")
		// Ricezione del messaggio dal server
		messaggio, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Println("Errore in creazione buffer di ricezione:", err.Error())
			os.Exit(1)
		}

		fmt.Print(messaggio)

	}
}
