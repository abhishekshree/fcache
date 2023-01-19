package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/abhishekshree/fcache/cache"
	"github.com/abhishekshree/fcache/client"
	"github.com/abhishekshree/fcache/protocol"
)

type ServerConfig struct {
	ListenAddr string
	IsLeader   bool // by some consensus algorithm
}

type Server struct {
	ServerConfig
	members map[*client.Client]struct{}
	cache   cache.Cacher
}

func NewServer(config ServerConfig, c cache.Cacher) *Server {
	return &Server{
		ServerConfig: config,
		cache:        c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("server listening on %s", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept: %v", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			log.Printf("failed to read: %v", err)
			return
		}
		go s.handleCommand(conn, cmd)
	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandSet:
		s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
		s.handleGetCommand(conn, v)
	case *protocol.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) error {
	fmt.Println("member just joined the cluster:", conn.RemoteAddr())

	s.members[client.NewFromConn(conn)] = struct{}{}

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) error {
	// log.Printf("GET %s", cmd.Key)

	resp := protocol.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = protocol.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)

	go func() {
		for member := range s.members {
			err := member.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TTL)
			if err != nil {
				log.Println("forward to member error:", err)
			}
		}
	}()

	resp := protocol.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = protocol.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = protocol.StatusOK
	_, err := conn.Write(resp.Bytes())

	return err
}
