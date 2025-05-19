package peer

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// PeerInfo representa a un nodo en la red
type PeerInfo struct {
	ID   int
	IP   string
	Port string
}

// Peer representa al nodo actual
type Peer struct {
	ID    int
	Port  string
	Peers []PeerInfo
}

// NewPeer crea un nuevo nodo Peer
func NewPeer(id int, port string, peers []PeerInfo) *Peer {
	return &Peer{
		ID:    id,
		Port:  port,
		Peers: peers,
	}
}

// StartListener inicia la escucha para recibir archivos
func (p *Peer) StartListener() {
	ln, err := net.Listen("tcp", ":"+p.Port)
	if err != nil {
		fmt.Println("Error al iniciar listener:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Nodo", p.ID, "escuchando en puerto", p.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:", err)
			continue
		}
		go p.handleConnection(conn)
	}
}

// handleConnection recibe archivos entrantes con nombre original
func (p *Peer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Leer nombre del archivo
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer nombre del archivo:", err)
		return
	}
	filename = strings.TrimSpace(filename)

	// Crear archivo en carpeta shared con el nombre original
	file, err := os.Create("shared/" + filename)
	if err != nil {
		fmt.Println("Error al crear archivo:", err)
		return
	}
	defer file.Close()

	// Leer y guardar contenido
	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println("Error al recibir archivo:", err)
		return
	}

	fmt.Println("📥 Archivo recibido correctamente como:", filename)
}

// SendFile envía un archivo a un nodo destino
func (p *Peer) SendFile(filePath, addr string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Enviar nombre del archivo
	filename := filepath.Base(filePath)
	_, err = fmt.Fprintf(conn, filename+"\n")
	if err != nil {
		return err
	}

	// Enviar contenido del archivo
	_, err = io.Copy(conn, file)
	return err
}

// GetLocalIP retorna la IP local de la máquina
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
