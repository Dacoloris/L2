package dev10

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

var timeout string

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "time limit to establish connection")
	flag.Parse()
}

func main() {
	re := regexp.MustCompile(`\ds`)
	if !re.MatchString(timeout) {
		log.Fatal("invalid timeout format: " + timeout)
	}

	if len(flag.Args()) < 2 {
		log.Fatal("usage: --timeout=10s host port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	toInt, _ := strconv.Atoi(timeout[:len(timeout)-1])
	to := time.Duration(toInt) * time.Second

	var conn net.Conn
	var err error

	start := time.Now()
	for time.Since(start) < to {
		conn, err = net.Dial("tcp", host+":"+port)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("unable to establish connection after timeout: %v", to)
	}
	defer conn.Close()
	log.Printf("connected to %s:%s", host, port)

	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Print("Message from server: " + msg)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := scanner.Text()
		_, err := fmt.Fprintf(conn, in+"\n")
		if err != nil {
			log.Fatal("connection closed")
		}
	}
}
