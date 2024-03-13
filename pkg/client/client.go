package client

import (
	"Go_Levine/pkg/common"
	"encoding/gob"
	
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	//"time"
	)




func Type_Order() common.Order {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your position (bid/ask): ")
		bid_ask_in, _ := reader.ReadString('\n')
		bid_ask := strings.TrimSpace(bid_ask_in)

		fmt.Print("Enter your price: ")
		price, _ := reader.ReadString('\n')
		price = strings.TrimSpace(price)
		pric, _ := strconv.Atoi(price)

		fmt.Print("Enter your volume: ")
		volu, _ := reader.ReadString('\n')
		volu = strings.TrimSpace(volu)
		volum, _ := strconv.Atoi(volu)

		inOrd:= common.Order{
			O_type:bid_ask,
			Price:int32(pric),
			Volume:int32(volum),
			}
        return inOrd
}



// StartClient starts the client
func StartClient() {
    // Client implementation here
	//conn, err := net.Dial("tcp", "192.168.1.133:12345")
    conn, err := net.Dial("tcp", "localhost:12345")
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
	fmt.Printf("\n\nNetwork: %v\nString: %v\n",conn.LocalAddr().Network(),conn.LocalAddr().String())
    defer conn.Close()

    reader := bufio.NewReader(os.Stdin)

	// name
	fmt.Printf("\nEnter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	User_info := common.User{Name:name}

	err = common.SendMessage(conn, 2, User_info)
	if err != nil {
		fmt.Println("Error sending User_info:", err)
		return
	}
	fmt.Printf("\nInfo sent: %v\n",User_info)
	
 

    // loop until there are no more matches
	var exi string="n"
	for{
		
		fmt.Print("\nWanna: \n (q) Quit\n (t) Trade\n (v) See OrderBook\n\n")
		inp, _ := reader.ReadString('\n')
		exi = strings.TrimSpace(inp)
		
		
		if exi== "q"{                            
			fmt.Print("\n\n$$$$$$$$$$$$$$$$$$$$\n$       BYE!       $\n$$$$$$$$$$$$$$$$$$$$\n\n")
			break

		} else if exi== "t" {
			inOrd := Type_Order()

			fmt.Printf("Address of inord= %v:\t%p\n", inOrd, &inOrd)
			
			err = common.SendMessage(conn, 1, inOrd)
			if err != nil {
				fmt.Println("Error sending MessageOne:", err)
				return
			}
			fmt.Printf("\nOrder sent: %v\n",inOrd)

		} else if exi== "v" {

			encoder := gob.NewEncoder(conn)

			// Encode the message type
			err := encoder.Encode(3)
			if err != nil {
				fmt.Errorf("error encoding message type: %v", err)
				continue //fmt.Errorf("error encoding message type: %v", err)
			}
			
			// Encode and send the actual data
			err = encoder.Encode(User_info) //Why no :????????????
			if err != nil {
				fmt.Errorf("error encoding data: %v", err)
				continue//fmt.Errorf("error encoding data: %v", err)
			}
			//var code int
			var OB common.Order_Book

			//time.Sleep(2 * time.Second)
		/*	
			err = common.ReciveMessage(conn,4,&OB)
			if err != nil {
				fmt.Println("Error Reciving Data request:", err)
				return
			}
*/
		
			decoder := gob.NewDecoder(conn)
			/*
			err__f := decoder.Decode(&code)
			fmt.Printf("Code: %v \n",code)
			if err__f != nil {
				fmt.Println("Error decoding order book type message:", err__f)
				continue
			}
			*/
			err_m := decoder.Decode(&OB)
			if err_m != nil {
				fmt.Println("Error decoding order book data:", err_m)
				continue
			}
		/*	*/
			fmt.Printf("Received Order Book: %+v\n", OB)
		
			// Print the received order book
			fmt.Println("Received Order Book:")
			common.Order_Book_print(OB, common.ORDER_BOOK_LENGTH, false)
		}

	}

}

