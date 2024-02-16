package client

import (
	"Go_net/pkg/common"

	"encoding/gob"
	"fmt"
	"net"
    "os"
	"strconv"
	"strings"
    "bufio"
)



func sendMessage(conn net.Conn, messageType int, data interface{}) error {
	// Create a new encoder for writing messages
	encoder := gob.NewEncoder(conn)

	// Encode the message type
	err := encoder.Encode(messageType)
	if err != nil {
		return fmt.Errorf("error encoding message type: %v", err)
	}

	// Encode and send the actual data
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding data: %v", err)
	}

	return nil
}





// StartClient starts the client
func StartClient() {
    // Client implementation here
    conn, err := net.Dial("tcp", "localhost:12345")
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
	
    defer conn.Close()

    reader := bufio.NewReader(os.Stdin)
    // loop until there are no more matches
	var exi string="n"
	for{
		
		fmt.Print("Wanna: \n (q) Quit\n (t) Trade")
		inp, _ := reader.ReadString('\n')
		exi = strings.TrimSpace(inp)
		
		
		if exi== "q"{                            
			fmt.Print("\n\n$$$$$$$$$$$$$$$$$$$$\n$       BYE!       $\n$$$$$$$$$$$$$$$$$$$$\n\n")
			break
		} else if exi== "t" {
			inOrd := Type_Order()

			fmt.Printf("Address of inord=%d:\t%p\n", inOrd, &inOrd)
			
			err = sendMessage(conn, 1, inOrd)
			if err != nil {
				fmt.Println("Error sending MessageOne:", err)
				return
			}
			fmt.Printf("\nOrder sent: %v\n",inOrd)
		}

		/*

		fmt.Print("Enter your position (bid/ask): ")
		bid_ask_in, _ := reader.ReadString('\n')
		bid_ask := strings.TrimSpace(bid_ask_in)

		fmt.Print("Enter your volume: ")
		volu, _ := reader.ReadString('\n')
		volu = strings.TrimSpace(volu)
		volum, _ := strconv.Atoi(volu)

		fmt.Print("Enter your price: ")
		price, _ := reader.ReadString('\n')
		price = strings.TrimSpace(price)
		pric, _ := strconv.Atoi(price)

		inOrd:= common.Order{
			O_type:bid_ask,
			Price:int32(pric),
			Volume:int32(volum),
			}
		*/


	}
    // Client implementation for interacting with the server
    // ...
}


func Type_Order() common.Order {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your position (bid/ask): ")
		bid_ask_in, _ := reader.ReadString('\n')
		bid_ask := strings.TrimSpace(bid_ask_in)

		fmt.Print("Enter your volume: ")
		volu, _ := reader.ReadString('\n')
		volu = strings.TrimSpace(volu)
		volum, _ := strconv.Atoi(volu)

		fmt.Print("Enter your price: ")
		price, _ := reader.ReadString('\n')
		price = strings.TrimSpace(price)
		pric, _ := strconv.Atoi(price)

		inOrd:= common.Order{
			O_type:bid_ask,
			Price:int32(pric),
			Volume:int32(volum),
			}
        return inOrd
}