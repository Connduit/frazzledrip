package network

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/utils"
)

var (
	listenerMutex   sync.Mutex
	currentListener net.Listener
)

// StartReverseTCPListener binds synchronously so the caller receives bind errors.
func StartReverseTCPListener(port int) error {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()
	if currentListener != nil {
		return fmt.Errorf("listener is already running")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	currentListener = listener
	go acceptClients(listener)
	log.Printf("TCP listener started on %s", listener.Addr())
	return nil
}

func acceptClients(listener net.Listener) {
	defer func() {
		listener.Close()
		listenerMutex.Lock()
		if currentListener == listener {
			currentListener = nil
		}
		listenerMutex.Unlock()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("TCP listener stopped: %v", err)
			return
		}
		listenerMutex.Lock()
		if currentListener != listener {
			listenerMutex.Unlock()
			conn.Close()
			return
		}
		addr := conn.RemoteAddr().String()
		client := &models.Client{ID: utils.GenerateClientID(), Conn: conn, LastSeen: time.Now(), IPAddress: addr, IsActive: true, Type: "TCP", Listener: conn.LocalAddr().String(), Messages: make([]models.ClientMessage, 0)}
		AddClient(addr, client)
		listenerMutex.Unlock()
		log.Printf("Client connected: %s from %s", client.ID, addr)
		go handleClientConnection(client)
	}
}

// Stopping accepts does not close existing session connections.
func StopReverseTCPListener() error {
	listenerMutex.Lock()
	listener := currentListener
	currentListener = nil
	listenerMutex.Unlock()
	if listener == nil {
		return nil
	}
	err := listener.Close()
	if errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}

func ListenerStatus() (bool, int) {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()
	if currentListener == nil {
		return false, 0
	}
	return true, currentListener.Addr().(*net.TCPAddr).Port
}

func IsListening() bool { running, _ := ListenerStatus(); return running }
