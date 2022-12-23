package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/rek7/chargen-go/pkg/chargen"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chargen",
	Short: "chargen - Wrapper for github.com/rek7/chargen-go library to stand up a server or run a client.",
}

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"srv"},
	Short:   "Chargen TCP/UDP client",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Panicf("issue parsing host " + err.Error())
		}
		proto, err := cmd.Flags().GetString("protocol")
		if err != nil {
			log.Panicf("issue parsing proto " + err.Error())
		}

		i := strings.LastIndex(host, ":")
		if i == -1 {
			log.Panicln("target needs to be ip:port")
		}

		ip := host[:i]
		port, err := strconv.Atoi(host[i+1:])
		if err != nil {
			log.Panicf("issue converting port %v", err)
		}

		server := chargen.NewServer()
		if proto == "tcp" {
			ln, err := net.Listen(proto, host)
			if err != nil {
				log.Panicf("issue listening " + err.Error())
			}

			log.Printf("listening on tcp: %v\n", ln.Addr())
			if err := server.ServeTCP(ln); err != nil {
				log.Panicf("issue serving " + err.Error())
			}
		} else {
			ln, err := net.ListenUDP(proto, &net.UDPAddr{
				Port: port,
				IP:   net.ParseIP(ip),
			})
			if err != nil {
				log.Panicf("issue listening " + err.Error())
			}

			log.Printf("listening on udp: %v\n", ln.LocalAddr())
			if err := server.ServeUDP(ln); err != nil {
				log.Panicf("issue serving " + err.Error())
			}
		}

	},
}

var clientCmd = &cobra.Command{
	Use:     "client",
	Aliases: []string{"cli"},
	Short:   "Chargen UDP/TCP client",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Panicf("issue parsing host " + err.Error())
		}

		proto, err := cmd.Flags().GetString("protocol")
		if err != nil {
			log.Panicf("issue parsing proto " + err.Error())
		}

		src, err := cmd.Flags().GetString("src")
		if err != nil {
			log.Panicf("issue parsing src " + err.Error())
		}

		mode, err := cmd.Flags().GetBool("mode")
		if err != nil {
			log.Panicf("issue parsing mode " + err.Error())
		}

		cli, err := chargen.NewClient(proto, host)
		if err != nil {
			log.Panicf("issue creating client " + err.Error())
		}

		defer cli.Close()
		if src != "" {
			if src == "t" {
				src = ""
			}
			if err := cli.UpdateSrcIP(net.IP(src)); err != nil {
				log.Panicf("issue uypdating ip " + err.Error())
			}
		}

		if mode {
			for {
				line, err := cli.Read()
				if err != nil {
					log.Panicf("issue reading " + err.Error())
				}
				fmt.Println(string(line))
			}
		} else {
			cli.Write(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(clientCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	clientCmd.Flags().StringP("host", "i", "127.0.0.1:19", "Host, including port, default 127.0.0.1:19")
	clientCmd.Flags().StringP("protocol", "p", "tcp", "Protocol, default udp. Only accepts tcp/udp")
	clientCmd.Flags().StringP("src", "s", "", "Spoof src ip, 't' for random ip if you want to specify one pass it as an arg")
	clientCmd.Flags().BoolP("mode", "m", false, "false (default) to write, true to read from server")

	serverCmd.Flags().StringP("host", "i", ":19", "Host, including port, default :19")
	serverCmd.Flags().StringP("protocol", "p", "tcp", "Protocol, default udp. Only accepts tcp/udp")
	Execute()
}
