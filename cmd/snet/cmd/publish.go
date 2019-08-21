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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trusch/snet/pkg/snet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a message to a feed",
	Long:  `publish a message to a feed.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn := getConn(cmd)
		cli := snet.NewCoreClient(conn)
		ctx := context.Background()

		feed, _ := cmd.Flags().GetString("feed")
		data, _ := cmd.Flags().GetString("data")
		dataType, _ := cmd.Flags().GetString("type")

		msg, err := cli.Create(ctx, &snet.CreateRequest{
			Feed:    feed,
			Content: []byte(data),
			Type:    dataType,
		})

		if err != nil {
			logrus.Fatal(err)
		}

		marshaler := &jsonpb.Marshaler{
			Indent: "  ",
		}
		bs, _ := marshaler.MarshalToString(msg)
		fmt.Printf("%v\n", bs)
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringP("addr", "a", "localhost:3001", "address to connect to")
	publishCmd.Flags().String("cert", "", "client cert")
	publishCmd.Flags().String("key", "", "client key")
	publishCmd.Flags().String("ca", "", "ca certificate")
	publishCmd.Flags().Bool("insecure", false, "skip tls verify")
	publishCmd.Flags().Bool("plaintext", false, "no tls at all")
	publishCmd.Flags().Duration("timeout", 0, "operation timeout, example: 2s")

	publishCmd.Flags().String("feed", "", "feed to publish to")
	publishCmd.Flags().String("data", "", "data to publish")
	publishCmd.Flags().String("type", "", "data type to publish")
}

func getConn(cmd *cobra.Command) *grpc.ClientConn {
	creds, err := clientCredentials(cmd)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "unable to create client credentials"))
	}
	addr, _ := cmd.Flags().GetString("addr")
	conn, err := grpc.Dial(addr, creds)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "unable to dial grpc address"))
	}
	return conn
}

func clientCredentials(cmd *cobra.Command) (grpc.DialOption, error) {
	caFile, _ := cmd.Flags().GetString("ca")
	certFile, _ := cmd.Flags().GetString("cert")
	keyFile, _ := cmd.Flags().GetString("key")
	plaintext, _ := cmd.Flags().GetBool("plaintext")
	skipVerify, _ := cmd.Flags().GetBool("insecure")
	if plaintext {
		logrus.Debug("using insecure")
		return grpc.WithInsecure(), nil
	}
	cfg := &tls.Config{}
	if skipVerify {
		logrus.Debug("credentials skip certificate verification")
		cfg.InsecureSkipVerify = true
	}
	if certFile != "" {
		certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, errors.Wrap(err, "could not load server key pair")
		}
		cfg.Certificates = append(cfg.Certificates, certificate)
	}
	if caFile != "" {
		logrus.Debugf("credentials CA \"%s\"\n", caFile)
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, errors.Wrap(err, "could not read ca certificate")
		}
		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append ca cert")
		}
		cfg.RootCAs = certPool
	}
	logrus.Debug("credentials CA empty")
	creds := grpc.WithTransportCredentials(credentials.NewTLS(cfg))
	return creds, nil
}
