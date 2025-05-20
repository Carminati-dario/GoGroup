package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	connHost     = "127.0.0.1" // host locale (server)
	portaBase    = 8080        // porta (solo per l'apertura delle connessioni)
	connType     = "tcp"
	numMaxClient = 10 // numero massimo di client accettati
)

// ConnessioneSafe implementa il tipo connessione in maniera concorrente (
type ConnessioneSafe struct {
	mux sync.Mutex
	c   map[int]net.Conn
}

type clientType struct {
	client     net.Conn
	port       int
	typeScreen bool
	idGruppo   []int
	screenId   int
}

var (
	l               [numMaxClient + 1]net.Listener
	isPortaOccupata [numMaxClient + 1]bool
	connessioni     [numMaxClient + 1]net.Conn
	screenGroup     int                          //tiene traccia dell'ultimo screen connesso
	gruppi          []string                     //array dei gruppi
	usernames       []string                     //array degli username
	client          [numMaxClient + 1]clientType //numero di client che si possono connettere
)

func main() {
	screenGroup = 0
	fmt.Println("Avvio server multichat " + connType + " su " + connHost + ":" + strconv.Itoa(portaBase) + ".")
	var err error
	for i := 0; i < numMaxClient+1; i++ { // Preparazione connessione per tutte le porte
		l[i], err = net.Listen(connType, ":"+strconv.Itoa(i+portaBase))
		if err != nil {
			fmt.Println("Errore apertura listener:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Connessione su porta " + strconv.Itoa(portaBase+i) + " in funzione.")
		defer l[i].Close()
		isPortaOccupata[0] = true // Occupo la porta per le connessioni iniziali
	}

	for { // Ripete all'infinito per accettare più connessioni
		fmt.Println("Server multichat pronto e in attesa di connessioni")
		conn, err := l[0].Accept()
		if err != nil {
			fmt.Println("Errore di connessione:", err.Error())
			return
		}
		fmt.Println("Ricevuta connessione con " + conn.RemoteAddr().String() + ". Verifico presenza slot liberi...")

		// Genera il nuovo socket
		nuovaPorta := portaBase // Intero di partenza per la ricerca di una porta libera
		portaNonTrovata := true
		for portaNonTrovata && nuovaPorta-portaBase < numMaxClient+1 {
			if !isPortaOccupata[nuovaPorta-portaBase] {
				isPortaOccupata[nuovaPorta-portaBase] = true
				portaNonTrovata = false
			} else {
				nuovaPorta++
			}
		}
		if nuovaPorta-portaBase == numMaxClient+1 {
			fmt.Println("Connessione rifiutata per mancanza slot liberi.")
			fmt.Fprintf(conn, "Chat al completo!\n")
			continue
		}
		// In questa esercitazione il numero di porta successivo è progressivo e non casuale
		// Inoltre non è stato implementato alcun meccanismo di autenticazione. Questa scelta quali problemi potrebbe comportare?
		fmt.Println("Nuova porta individuata:", strconv.Itoa(nuovaPorta))
		go handleConnessione(nuovaPorta)

		fmt.Println(conn.LocalAddr().String() + ": invio a " + conn.RemoteAddr().String() + " la nuova porta " + strconv.Itoa(nuovaPorta))
		fmt.Fprintf(conn, strconv.Itoa(nuovaPorta)+"\n") // Invia al client la porta da usare per la comunicazione
	}
}

func username_acceptable(username string) bool {
	for i := 0; i < len(usernames); i++ {
		if usernames[i] == username {
			return false
		}
	}
	return true
}

func contains(element string) bool {
	for _, i := range gruppi {
		if i == element {
			return true
		}
	}
	return false
}

func findI(element string) int {
	for i := 0; i < len(gruppi); i++ {
		if gruppi[i] == element {
			return i
		}
	}
	return -1
}

// func handleConnessione gestisce le singole connessioni dei client ricevendo e inviando messaggi
func handleConnessione(porta int) {
	fmt.Println("Attesa di connessione su porta " + strconv.Itoa(porta))

	p := porta - portaBase
	var err error
	connessioni[p], err = l[p].Accept()
	if err != nil {
		fmt.Println("Errore di connessione:", err.Error())
		return
	}

	fmt.Println("Client " + connessioni[p].RemoteAddr().String() + " connesso su " + strconv.Itoa(porta))
	client[p].client = connessioni[p]
	client[p].port = porta
	continua := true
	var username string

	//Inserimento username da parte del client e conseguente identificazione degli screen
	for {
		netData, err := bufio.NewReader(connessioni[p]).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		username = strings.TrimSpace(netData)
		if username != "screen" {
			if username_acceptable(username) {
				usernames = append(usernames, username)
				break
			} else {
				fmt.Fprintf(connessioni[p], "Username già esistente, inserisci un altro:\n")
			}
		} else {
			client[p].typeScreen = true
			client[p].screenId = screenGroup
			screenGroup++
			broadcastMessaggio(p, "Gruppo di appartenenza "+string(client[p].screenId), client[p].screenId)
			fmt.Println("Gruppo di appartenenza " + strconv.Itoa(client[p].screenId))
			break
		}
	}

	if !client[p].typeScreen {
		for continua {
			netData, err := bufio.NewReader(connessioni[p]).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			message := strings.TrimSpace(string(netData))

			mUpper := strings.ToUpper(message)
			if strings.HasPrefix(mUpper, "GROUP: ") {
				nameGroup := message[7:]
				if !contains(nameGroup) {
					gruppi = append(gruppi, nameGroup)
				}
				client[p].idGruppo = append(client[p].idGruppo, findI(nameGroup))
			}

			// Controllo sulla richiesta di terminazione della connessione
			if message == "esci" {
				continua = false
			} else {
				for _, i := range client[p].idGruppo {
					message = username + " - > " + message + "\n"
					//fmt.Println("MESSAGGIO VALE: " + message)
					broadcastMessaggio(p, message, i)
				}
			}

			fmt.Println(connessioni[p].LocalAddr().String() + ": " + connessioni[p].RemoteAddr().String() + " -> " + message + " -> " + strconv.Itoa(client[p].idGruppo[0]))
		}

		fmt.Println(connessioni[p].LocalAddr().String() + ": client " + connessioni[p].RemoteAddr().String() + " disconnesso.")

		// Libera le informazioni occupate
		client[p].client = nil
		client[p].idGruppo = nil
		client[p].port = 0
		client[p].typeScreen = false
		connessioni[p] = nil
		isPortaOccupata[p] = false

	}
}

// func broadcastMessaggio si occupa di inviare il messaggio a tutti i client connessi
func broadcastMessaggio(p int, netData string, id int) {
	for i, conn := range connessioni {
		if i != 0 && i != p && connessioni[i] != nil && client[i].typeScreen && client[i].screenId == id { // Evita di inviare il messaggio a se stesso e a connessioni che non esistono
			fmt.Fprintf(conn, netData)
		}
	}
}
