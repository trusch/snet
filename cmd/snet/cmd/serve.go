// Copyright © 2019 Tino Rusch <tino.rusch@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	grpcserver "github.com/contiamo/goserver/grpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trusch/snet/pkg/events"
	"github.com/trusch/snet/pkg/fetcher"
	"github.com/trusch/snet/pkg/server"
	"github.com/trusch/snet/pkg/snet"
	"golang.org/x/crypto/ed25519"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run the snet server",
	Long:  `run the snet server.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		opts := loadServeOpts(cmd)
		keys := make([]ed25519.PrivateKey, len(opts.FeedKeys))
		for idx, keyFile := range opts.FeedKeys {
			bs, err := ioutil.ReadFile(keyFile)
			if err != nil {
				logrus.Fatal(err)
			}
			key := ed25519.NewKeyFromSeed(bs)
			keys[idx] = key
			fmt.Printf("serving feed %x\n", key.Public())
		}

		db, err := sql.Open("postgres", "postgres://localhost:5432?user=postgres&dbname=postgres&sslmode=disable")
		if err != nil {
			logrus.Fatal(err)
		}

		eventManager := events.NewManager(ctx, 32)

		coreServer, err := server.NewCoreHandler(db, keys)
		if err != nil {
			logrus.Fatal(err)
		}

		peersServer, err := server.NewPeersHandler(db, eventManager)
		if err != nil {
			logrus.Fatal(err)
		}

		// setup grpc server with options
		grpcServer, err := grpcserver.New(&grpcserver.Config{
			Options: []grpcserver.Option{
				grpcserver.WithCredentials(opts.TLSCert, opts.TLSKey, opts.TLSCACert),
				grpcserver.WithLogging("snet"),
				grpcserver.WithMetrics(),
				grpcserver.WithRecovery(),
				grpcserver.WithReflection(),
			},
			Register: func(srv *grpc.Server) {
				snet.RegisterCoreServer(srv, coreServer)
				snet.RegisterPeersServer(srv, peersServer)
			},
		})
		if err != nil {
			logrus.Fatal(err)
		}

		// start server
		go grpcserver.ListenAndServe(ctx, opts.Addr, grpcServer)

		fetchWorker := fetcher.New(db, eventManager)
		err = fetchWorker.Work(ctx)
		if err != nil {
			logrus.Error(err)
		}

		<-c
		cancel()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("addr", ":3001", "listening address")
	serveCmd.Flags().String("tls-cert", "", "TLS certificate")
	serveCmd.Flags().String("tls-key", "", "TLS key certificate")
	serveCmd.Flags().String("tls-ca-cert", "", "TLS ca certificate to verify clients")
	serveCmd.Flags().StringSlice("feed-key", []string{}, "feed key file (can be repeated)")
}

type serveOpts struct {
	Addr                       string
	TLSCert, TLSKey, TLSCACert string
	FeedKeys                   []string
}

func loadServeOpts(cmd *cobra.Command) serveOpts {
	cert, _ := cmd.Flags().GetString("tls-cert")
	key, _ := cmd.Flags().GetString("tls-key")
	ca, _ := cmd.Flags().GetString("tls-ca-cert")
	addr, _ := cmd.Flags().GetString("addr")
	feedKeys, _ := cmd.Flags().GetStringSlice("feed-key")
	return serveOpts{addr, cert, key, ca, feedKeys}
}
