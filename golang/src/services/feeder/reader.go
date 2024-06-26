package feederservice

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Emad-am/feeder/tools"
)

func (feeder *Feeder) startReader() {
	user := conf.DDE.UserName
	password := conf.DDE.Password

	c := *feeder.C
	conn := *feeder.Conn
	reader := bufio.NewReader(conn)

	initialResponse, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading initial response: %v\n", err)
		return
	}
	fmt.Printf("Server: %s", initialResponse)

	_, err = fmt.Fprintf(conn, user+"\n")

	secondResponse, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading second response: %v\n", err)
		return
	}
	fmt.Printf("Server: %s", secondResponse)

	loginPrompt, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading third response: %v\n", err)
		return
	}
	fmt.Printf("Server: %s", loginPrompt)

	_, err = fmt.Fprintf(conn, user+"\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending username: %v\n", err)
		return
	}
	fmt.Printf("Sent username: %s\n", user)

	passwordPrompt, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password prompt: %v\n", err)
		return
	}
	fmt.Printf("Server: %s", passwordPrompt)

	_, err = fmt.Fprintf(conn, password+"\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending password: %v\n", err)
		return
	}
	fmt.Printf("Sent password: %s\n", password)

	authResponse, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading authentication response: %v\n", err)
		return
	}

	fmt.Printf("Authentication response: %s", authResponse)
	// waitChannel := make(chan struct{}, 1)
	go func() {
		// defer func() {
		// 	waitChannel <- struct{}{}
		// }()

		for {
			res, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading authentication response: %v\n", err)
				return
			}
			if res == "" {
				fmt.Println("NIL")
				continue
			}
			str := tools.Trim(res, []string{"\n", "\r", "."})
			splitted := strings.Split(str, " ")
			tick := Tick{
				Symbol: splitted[0],
				Ask:    tools.ToInt(splitted[1]),
				Bid:    tools.ToInt(splitted[2]),
				Time:   time.Now().UnixMilli(),
			}

			go func(data Tick) {
				c <- data
			}(tick)
		}
	}()

	// <-waitChannel
}
