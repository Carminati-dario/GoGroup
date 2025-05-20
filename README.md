# GoGroup
Chat multiutente in Go basata su TCP per comunicazioni in tempo reale. Supporta utenti attivi e visualizzatori passivi (screen), permettendo la creazione di gruppi di discussione personalizzati. L'architettura client-server gestisce efficacemente connessioni simultanee, offrendo un sistema di messaggistica leggero e funzionale.
## Struttura del Progetto
Il progetto è suddiviso in tre componenti principali:

### Server: 
Gestisce le connessioni, i gruppi e il routing dei messaggi

### Client: 
Interfaccia utente per inviare messaggi nei gruppi

### Screen: 
Visualizzatore passivo per monitorare i messaggi di uno specifico gruppo

## Funzionalità

- Comunicazione TCP in tempo reale
- Supporto per gruppi di chat personalizzati
- Client per utenti standard e modalità "screen" di sola visualizzazione
- Gestione di username unici
- Architettura scalabile fino a 10 connessioni simultanee

## Istruzioni per l'uso
1. Avvio del server
2. Connessione di client su terminale differente
3. Connessione come server su terminale differente

## Comandi disponibili 
- ```esci```: Termina la connessione
- ``` GROUP: [nome-gruppo]```: Crea o si unisce a un gruppo
