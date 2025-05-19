package peer

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Peer representa la estructura del archivo peers.json
type Peer struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// GetLocalIP obtiene la IP local (no loopback)
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// LoadPeers carga los peers desde un archivo JSON
func LoadPeers(path string) ([]Peer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var peers []Peer
	if err := json.Unmarshal(data, &peers); err != nil {
		return nil, err
	}

	return peers, nil
}

// GetLocalID compara la IP local con las del JSON y retorna el ID correspondiente
func GetLocalID(peers []Peer) (int, error) {
	localIP, err := GetLocalIP()
	if err != nil {
		return 0, fmt.Errorf("no se pudo obtener IP local: %w", err)
	}

	for _, p := range peers {
		if p.IP == localIP {
			return p.ID, nil
		}
	}

	return 0, fmt.Errorf("la IP local %s no se encuentra en peers.json", localIP)
}
