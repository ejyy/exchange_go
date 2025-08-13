package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/ejyy/exchange_go/exchange"
)

// Server handles TCP connections and bridges to the exchange
type Server struct {
	exchange *exchange.Exchange
	actions  chan *exchange.Action
	clients  map[net.Conn]bool
	mutex    sync.RWMutex
}

// Init creates a new TCP server instance
func (srv *Server) Init(ex *exchange.Exchange, actions chan *exchange.Action) {
	srv.exchange = ex
	srv.actions = actions
	srv.clients = make(map[net.Conn]bool)

	if err := srv.Start(exchange.TCPPort); err != nil {
		fmt.Println("Server failed: ", err)
	}
}

// Start begins listening on the specified port and processes actions
func (s *Server) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("TCP server listening on port %d\n", port)

	// Start action processor goroutine
	go s.processActions()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		s.addClient(conn)
		go s.handleClient(conn)
	}
}

// addClient registers a new client connection
func (s *Server) addClient(conn net.Conn) {
	s.mutex.Lock()
	s.clients[conn] = true
	s.mutex.Unlock()
}

// removeClient unregisters a client connection
func (s *Server) removeClient(conn net.Conn) {
	s.mutex.Lock()
	delete(s.clients, conn)
	s.mutex.Unlock()
	conn.Close()
}

// broadcast sends a message to all connected clients
func (s *Server) broadcast(message string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for conn := range s.clients {
		if _, err := fmt.Fprintln(conn, message); err != nil {
			// Client disconnected, remove it
			go s.removeClient(conn)
		}
	}
}

// handleClient processes commands from a single client
func (s *Server) handleClient(conn net.Conn) {
	defer s.removeClient(conn)
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		s.processCommand(conn, line)
	}
}

// processCommand parses and executes client commands
func (s *Server) processCommand(conn net.Conn, cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch strings.ToUpper(parts[0]) {
	case "LIMIT":
		s.handleLimit(conn, parts[1:])
	case "CANCEL":
		s.handleCancel(conn, parts[1:])
	default:
		fmt.Fprintln(conn, "ERROR: Unknown command")
	}
}

// handleLimit processes LIMIT orders
// Format: LIMIT <symbol> <side> <price> <size> <trader>
func (s *Server) handleLimit(conn net.Conn, args []string) {
	if len(args) != 5 {
		fmt.Fprintln(conn, "ERROR: LIMIT requires 5 arguments: symbol side price size trader")
		return
	}

	symbol := args[0]

	var side exchange.Side
	switch strings.ToUpper(args[1]) {
	case "BID", "BUY":
		side = exchange.Bid
	case "ASK", "SELL":
		side = exchange.Ask
	default:
		fmt.Fprintln(conn, "ERROR: Invalid side, use BID/BUY or ASK/SELL")
		return
	}

	price, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		fmt.Fprintln(conn, "ERROR: Invalid price")
		return
	}

	size, err := strconv.ParseUint(args[3], 10, 32)
	if err != nil {
		fmt.Fprintln(conn, "ERROR: Invalid size")
		return
	}

	trader, err := strconv.ParseUint(args[4], 10, 16)
	if err != nil {
		fmt.Fprintln(conn, "ERROR: Invalid trader ID")
		return
	}

	s.exchange.Limit(symbol, exchange.Price(price), exchange.Size(size), side, exchange.TraderID(trader))
	fmt.Fprintln(conn, "OK")
}

// handleCancel processes CANCEL orders
// Format: CANCEL <orderID>
func (s *Server) handleCancel(conn net.Conn, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(conn, "ERROR: CANCEL requires 1 argument: orderID")
		return
	}

	orderID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(conn, "ERROR: Invalid order ID")
		return
	}

	s.exchange.Cancel(exchange.OrderID(orderID))
	fmt.Fprintln(conn, "OK")
}

// processActions reads from the actions channel and broadcasts to clients
func (s *Server) processActions() {
	for action := range s.actions {
		s.broadcast(action.String())
	}
}
