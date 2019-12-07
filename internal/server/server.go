// Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"time"

	"github.com/grab/async"
	"github.com/grab/talaria/internal/config"
	"github.com/grab/talaria/internal/monitor"
	"github.com/grab/talaria/internal/presto"
	"github.com/grab/talaria/internal/table"
	talaria "github.com/grab/talaria/proto"
	"google.golang.org/grpc"
)

const (
	ctxTag  = "server"
	errTag  = "error"
	funcTag = "func"
)

// Membership represents a contract required for recovering cluster information.
type Membership interface {
	Members() []string
}

// Storage represents an eventlog storage contract.
type Storage interface {
	io.Closer
	Append(key, value []byte, ttl time.Duration) error
	Range(seek, until []byte, f func(key, value []byte) bool) error
}

// ------------------------------------------------------------------------------------------------------------

// New creates a new talaria server.
func New(prestoCfg *config.Presto, monitor monitor.Client, tables ...table.Table) *Server {
	const maxMessageSize = 32 * 1024 * 1024 // 32 MB
	server := &Server{
		server:    grpc.NewServer(grpc.MaxRecvMsgSize(maxMessageSize)),
		prestoCfg: prestoCfg,
		monitor:   monitor,
		tables:    make(map[string]table.Table),
	}

	// Register the gRPC server
	talaria.RegisterIngressServer(server.server, server)

	// Build a registry of tables
	for _, table := range tables {
		server.tables[table.Name()] = table
	}
	return server
}

// Server represents the talaria server which should implement presto thrift interface.
type Server struct {
	server    *grpc.Server           // The underlying gRPC server
	prestoCfg *config.Presto         // The presto configuration
	monitor   monitor.Client         // The monitoring layer
	cancel    context.CancelFunc     // The cancellation function for the server
	tables    map[string]table.Table // The list of tables
}

// Listen starts listening on presto RPC & gRPC.
func (s *Server) Listen(ctx context.Context, prestoPort, grpcPort int32) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// Asynchronously start the gRPC listener
	async.Invoke(ctx, func(ctx context.Context) (interface{}, error) {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		defer lis.Close()
		if err != nil {
			return nil, err
		}

		if err := s.server.Serve(lis); err != nil {
			return nil, err
		}
		return nil, nil
	})

	// Serve presto and block
	return presto.Serve(ctx, int32(prestoPort), s)
}

// Close closes the server and related resources.
func (s *Server) Close() {
	s.server.GracefulStop()
	s.cancel()

	// Close all the open tables
	for _, t := range s.tables {
		if err := t.Close(); err != nil {
			s.monitor.Errorf(err.Error())
		}
	}
}

// handlePanic handles the panic and logs it out.
func (s *Server) handlePanic() {
	if r := recover(); r != nil {
		s.monitor.Errorf("panic recovered: %ss \n %s", r, debug.Stack())
	}
}
