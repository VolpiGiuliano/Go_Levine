package common

import (
	"fmt"
	"encoding/gob"
	"net"
)


const(
	ORDER_BOOK_LENGTH int = 20
)

type User struct{
	Name string
}

// Rappresents the single order.
//
// o_type : string {"bid","ask"}
//
// price : float32
//
// volume : int32
type Order struct {
	O_type string
	Price  int32
	Volume int32
	
}


type Order_Filled struct {
	Og_bid Order // original quote
	Og_ask Order // 
	Price  int32
	Vol_filled int32
}


func (bigger_order *Order) Partial_fill(other_order Order) Order_Filled {
	
	//bigger_order.volume = bigger_order.volume - other_order.volume

	var ord_filled Order_Filled

	if bigger_order.O_type=="bid" {

		ord_filled:=Order_Filled{
			Og_bid: *bigger_order,
			Og_ask: other_order,
			Price: bigger_order.Price,
			Vol_filled: other_order.Volume,
		}
		bigger_order.Volume = bigger_order.Volume - other_order.Volume
		return ord_filled
		//fmt.Printf("%v",ord_filled)

	} else if bigger_order.O_type=="ask"{

		ord_filled:=Order_Filled{
			Og_bid: other_order,
			Og_ask: *bigger_order,
			Price: bigger_order.Price,
			Vol_filled: other_order.Volume,
		}
		bigger_order.Volume = bigger_order.Volume - other_order.Volume
		
		return ord_filled
		//fmt.Printf("%v",ord_filled)
	} else{
		bigger_order.Volume = bigger_order.Volume - other_order.Volume
		return ord_filled
	}
		
	
}


/////////////Queue////////////////

// Queue that the orders need to follow in any given price.
// The order is first-in first-out (time priority) FIFO.
// To enter and to 
// Observe gives the first in line Order without taking it out of the Queue

type Queue struct {
	Items []Order
}

func (q *Queue) Enqueue(i Order) {
	q.Items = append(q.Items, i)
}

func (q *Queue) Dequeue() Order {
	to_remove := q.Items[0]
	q.Items = q.Items[1:]

	return to_remove
}

func (q *Queue) Observe() Order {
	to_see := q.Items[0]

	return to_see
}


func (q *Queue) Is_Empty() bool {
	var is bool
	if len(q.Items)==0{
		is= true
	} else{
		is= false
	}
	return is
}

////////////////Order Book//////////////////

type Order_Book struct {
	Ask [ORDER_BOOK_LENGTH]*Queue
	Bid [ORDER_BOOK_LENGTH]*Queue
}


/////////////////////////////////////////


// Fuction usefull for the Order_Book_print()
func Size_Level(level_list []Order)(size int){

	for i:=0; i< len(level_list);i++{
		size= size+int(level_list[i].Volume)
	}

	return
}


func Order_Book_print(OB Order_Book, lenght_OB int,size_only bool) {
	//fmt.Printf("Order Book: %v  \n %v \n", OB, OB.ask)
	fmt.Println("                  °")
	fmt.Println("----------- Order Book --------------")
	fmt.Println("            Bid        Ask ")
	
	for  i := 0; i < int(lenght_OB); i++ {

		if size_only{
			fmt.Printf("| %v | - %v - | %v |\n",Size_Level(OB.Bid[i].Items),i,Size_Level(OB.Ask[i].Items))
		}else{
			fmt.Printf("Level: %v | %v   -   %v  |\n", i, OB.Bid[i], OB.Ask[i])
		}
		//Size_Level(OB.bid[i].items) //i, OB.bid[i].items, OB.ask[i].items)
	}

	fmt.Println("-------------------------------------")
	fmt.Println("                  °")
}


// if there is no best quote it will rerurn {0,[]}
func Find_best (ob Order_Book) (level_b int, best_b []Order,level_a int, best_a []Order) {
	
	fmt.Printf("##################################\n\nSearching for the BEST BID\n")
	for index := len(ob.Bid)-1;index >= 0; index-- {
		fmt.Printf("Index: %v | Items: %v  ->len:%v\n",index, ob.Bid[index].Items,len(ob.Bid[index].Items))

		if len(ob.Bid[index].Items)!=0{
			level_b= index
			best_b=ob.Bid[index].Items
			break
		}

	}

	fmt.Printf("\n===================================\n\nSearching for the BEST ASK\n")
	for index, value := range ob.Ask{
		fmt.Printf("Index: %v | Items: %v  ->len:%v\n",index, value.Items,len(value.Items))

		if len(value.Items)!=0{
			level_a= index
			best_a=value.Items
			break
		}

	}
	fmt.Printf("\n##################################\n\n")
	return
}




// connections

func SendMessage(conn net.Conn, messageType int, data interface{}) error {
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

func ReciveMessage(conn net.Conn, messageType int, data interface{}) error {
	
	decoder := gob.NewDecoder(conn)
	/*
	err__f := decoder.Decode(&code)
	if err__f != nil {
		fmt.Println("Error decoding order book type message:", err__f)
		continue
	}
	*/
	err_m := decoder.Decode(&data)
	if err_m != nil {
		fmt.Println("Error decoding order book data:", data)
		return fmt.Errorf("error decoding order book data: %v", err_m)
		
	}
	return nil
}


