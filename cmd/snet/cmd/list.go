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
	"fmt"
	"io"

	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trusch/snet/pkg/snet"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list messages",
	Long:  `list messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn := getConn(cmd)
		cli := snet.NewCoreClient(conn)
		ctx := context.Background()

		feed, _ := cmd.Flags().GetString("feed")
		dataType, _ := cmd.Flags().GetString("type")
		stream, _ := cmd.Flags().GetBool("stream")

		iter, err := cli.List(ctx, &snet.ListRequest{
			Feed:   feed,
			Type:   dataType,
			Stream: stream,
		})
		if err != nil {
			logrus.Fatal(err)
		}

		marshaler := &jsonpb.Marshaler{
			Indent: "  ",
		}
		for {
			item, err := iter.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				logrus.Fatal(err)
			}
			bs, _ := marshaler.MarshalToString(item)
			fmt.Printf("%v\n", bs)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("addr", "a", "localhost:3001", "address to connect to")
	listCmd.Flags().String("cert", "", "client cert")
	listCmd.Flags().String("key", "", "client key")
	listCmd.Flags().String("ca", "", "ca certificate")
	listCmd.Flags().Bool("insecure", false, "skip tls verify")
	listCmd.Flags().Bool("plaintext", false, "no tls at all")
	listCmd.Flags().Duration("timeout", 0, "operation timeout, example: 2s")

	listCmd.Flags().String("feed", "", "feed to list from")
	listCmd.Flags().String("type", "", "data type to include in list")
	listCmd.Flags().Int("from-seq", 0, "first sequence number")
	listCmd.Flags().Int("to-seq", 0, "last sequence number")
	listCmd.Flags().Bool("stream", false, "stream new messages")
}
