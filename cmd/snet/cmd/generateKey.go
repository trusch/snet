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
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"
)

// generateKeyCmd represents the generateKey command
var generateKeyCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate a key",
	Long:  `generate a ed25519 key.`,
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		_, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			logrus.Fatal(err)
		}
		seed := fmt.Sprintf("%x", priv.Seed())
		feed := fmt.Sprintf("%x", priv.Public())
		if out != "" {
			err := ioutil.WriteFile(out, priv.Seed(), 0400)
			if err != nil {
				logrus.Fatal(err)
			}
		}
		fmt.Println("feed:", feed)
		fmt.Println("seed:", seed)
	},
}

func init() {
	rootCmd.AddCommand(generateKeyCmd)
	generateKeyCmd.Flags().StringP("out", "o", "", "output file for the key")
}
