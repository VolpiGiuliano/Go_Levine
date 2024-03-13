package server

import (
	"Go_Levine/pkg/common"
	"Go_Levine/pkg/exchange"
	"encoding/gob"
	"fmt"
	"net"
	//"time"
)



var or = common.Order{O_type:"ask", Price:10, Volume:15}
var or1 = common.Order{O_type:"ask",Price: 10,Volume: 1}
var or2 = common.Order{O_type:"ask",Price: 10, Volume:5}
var bid1 = common.Order{O_type:"bid",Price: 9,Volume: 9}
var bid2 = common.Order{O_type:"bid",Price: 9, Volume:1}
var bid3 = common.Order{O_type:"bid",Price: 7,Volume: 100}
var bid4 = common.Order{O_type:"bid",Price: 9,Volume: 2}


var ask_l [common.ORDER_BOOK_LENGTH]*common.Queue
var bid_l [common.ORDER_BOOK_LENGTH]*common.Queue

var incoming_q []common.Order
var filled_quotes []common.Order_Filled

var type_ob bool



// StartServer starts the server
func StartServer() {
    ///////////////////////////////////////////////////////////////////////////////

	for price_i := range ask_l {
		p_que := common.Queue{}
		ask_l[price_i] = &p_que
	}

	for price_i := range bid_l {
		p_que := common.Queue{}
		bid_l[price_i] = &p_que
	}
	
	OB := common.Order_Book{Ask:ask_l,Bid: bid_l}

	incoming_q = append(incoming_q, or, or2,or1, bid1, bid2, bid3, bid4)

	exchange.Inserter(&incoming_q, OB)
	
	common.Order_Book_print(OB,common.ORDER_BOOK_LENGTH,false)

    ///////////////////////////////////////////////////////////////////////////////

    // Server implementation here
    listener, err := net.Listen("tcp", ":12345")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server is listening on port 12345")

    go exchange.Engine(&incoming_q,&OB,&filled_quotes)


    for {

        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        go handleConnection(conn,&incoming_q,&OB)

        if len(incoming_q)!=0{
            exchange.Inserter(&incoming_q, OB)
            common.Order_Book_print(OB,common.ORDER_BOOK_LENGTH,false)
        }
        common.Order_Book_print(OB,common.ORDER_BOOK_LENGTH,false)

    }


}



// handleConnection handles a single client connection
func handleConnection(conn net.Conn,list *[]common.Order,OB *common.Order_Book) {
	fmt.Printf("Collegato! %v\n",conn.RemoteAddr())
    defer conn.Close()


    for {
        fmt.Println("Return to the start of the loop")
        decoder := gob.NewDecoder(conn)
		// Read the message type
		var messageType int
		err := decoder.Decode(&messageType)
		if err != nil {
			fmt.Printf("\nError decoding message type: %v\n", err)
			return
		}

        switch messageType {
            case 1: // Recived an Order
                // Decode and handle MessageOne

                var order_reci common.Order
                err_m := decoder.Decode(&order_reci)
                if err_m != nil {
                    fmt.Println("Error decoding MessageOne:", err_m)
                    return
                }
                fmt.Printf("Received MessageOrder: %+v\n", order_reci)

                *list = append(*list, order_reci)
                // Respond to MessageOne
                response := "Response to MessageOne"
                fmt.Println("Sending response:", response)
                conn.Write([]byte(response))
                
            
            case 2:// User info
                // Decode and handle MessageTwo
                var user_info common.User
                err := decoder.Decode(&user_info)
                if err != nil {
                    fmt.Println("Error decoding User:", err)
                    return
                }
                fmt.Printf("Received User info: %+v\n", user_info)

                // Respond to MessageTwo
                response := fmt.Sprintf("Response to User %v",user_info.Name)
                fmt.Println("Sending response:", response)
                conn.Write([]byte(response))

            case 3:// Oeder book
                

                var user_info common.User
                err := decoder.Decode(&user_info)
                if err != nil {
                    fmt.Println("Error decoding User for OB:", err)
                    return
                }
                fmt.Printf("Received User info for OB: %+v\n", user_info)
             

                fmt.Printf("Conn: %v\n",conn)

                
                // Sending the OB

                encoder := gob.NewEncoder(conn)
                err_mm := encoder.Encode(*OB)
                if err_mm != nil {
                    fmt.Println("Error encoding order book:", err_mm)
                    //return
                }

                ////////////////////////////////////////////
                fmt.Println("Order book sent to client")
                //time.Sleep(5 * time.Second)
  
            default:
                fmt.Println("Unknown message type:", messageType)
                return
        }

        

    }

    
}


