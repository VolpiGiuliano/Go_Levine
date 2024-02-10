package main

import (
	"fmt"
)


const(
	ORDER_BOOK_LENGTH int8 = 20
)


// Rappresents the single order.
//
// o_type : string {"bid","ask"}
//
// price : float32
//
// volume : int32
type Order struct {
	o_type string
	price  float32
	volume int32
}

/////////////Queue////////////////

// Queue that the orders need to follow in any given price.
// The order is first-in first-out (time priority) FIFO
type Queue struct {
	items []Order
}

func (q *Queue) Enqueue(i Order) {
	q.items = append(q.items, i)
}

func (q *Queue) Dequeue() Order {
	to_remove := q.items[0]
	q.items = q.items[1:]

	return to_remove
}

////////////////Order Book//////////////////

type Order_Book struct {
	ask [ORDER_BOOK_LENGTH]*Queue
	bid [ORDER_BOOK_LENGTH]*Queue
}

////////////////////////////////
// Made up orders for testing

var or = Order{"ask", 10, 15}
var or1 = Order{"ask", 10, 1}
var or2 = Order{"ask", 10, 5}
var bid1 = Order{"bid", 9, 9}
var bid2 = Order{"bid", 9, 1}
var bid3 = Order{"bid", 7, 100}
var bid4 = Order{"bid", 9, 2}
var bid5 = Order{"bid", 5, 5}
var bido1 = Order{"bid", 2, 1}
var bido2 = Order{"bid", 2, 2}


var ask_l [ORDER_BOOK_LENGTH]*Queue
var bid_l [ORDER_BOOK_LENGTH]*Queue

var incoming_q []Order

/////////////////////////////////////////



// It puts the single Order in
func inserter(l_in_quo *[]Order, l_ask [ORDER_BOOK_LENGTH]*Queue, l_bid [ORDER_BOOK_LENGTH]*Queue) {

	fmt.Println("                  §")
	fmt.Printf("++++++++++++++++++++++++++++++++++++\nStart insertion\nINFUNC Incoming quote: %v     Pointer: %p\n", l_in_quo, l_in_quo)

	switch {
	case len(*l_in_quo) == 1 && (*l_in_quo)[0].o_type == "ask":

		l_ask[int((*l_in_quo)[0].price)].Enqueue((*l_in_quo)[0])
		//return

	case len(*l_in_quo) == 1 && (*l_in_quo)[0].o_type == "bid":

		l_bid[int((*l_in_quo)[0].price)].Enqueue((*l_in_quo)[0])
		//return

	case len(*l_in_quo) > 1:

		for _, v := range *l_in_quo {
			if v.o_type == "ask" {
				l_ask[int(v.price)].Enqueue(v)
			} else if v.o_type == "bid" {
				l_bid[int(v.price)].Enqueue(v)
			}
		}

	case len(*l_in_quo) == 0:
		fmt.Printf("No incoming quotes\n++++++++++++++++++++++++++++++++++++\n")
		fmt.Println("                  §")
		return

	}

	*l_in_quo = nil
	fmt.Printf("End insertion\nINFUNC Incoming quote: %v     Pointer: %p\n++++++++++++++++++++++++++++++++++++\n", l_in_quo, l_in_quo)
	fmt.Println("                  §")

}

func ticker(list [ORDER_BOOK_LENGTH]*Queue, name string) {

	fmt.Printf("\nList %v\n", name)

	for index, pp := range list {

		if len(pp.items) != 0 {
			fmt.Printf("Value: %v   %v  | Point:%v     %p\n", index, pp, *&pp.items[0], &pp.items[0])
		} else {
			fmt.Printf("Value: %v   %v  | Point:%v     %p\n", index, pp, *&pp.items, &pp.items)
		}

	}
}

func Order_Book_print(OB Order_Book, lenght_OB int8) {
	//fmt.Printf("Order Book: %v  \n %v \n", OB, OB.ask)
	fmt.Println("                  °")
	fmt.Println("----------- Order Book --------------")
	fmt.Println("            Bid        Ask ")
	for  i := 0; i < int(lenght_OB); i++ {

		fmt.Printf("Level: %v | %v   -   %v  |\n", i, OB.bid[i], OB.ask[i])
	}

	fmt.Println("-------------------------------------")
	fmt.Println("                  °")
}



func main() {

	for price_i := range ask_l {
		p_que := Queue{}
		ask_l[price_i] = &p_que
	}

	for price_i := range bid_l {
		p_que := Queue{}
		bid_l[price_i] = &p_que
	}
	

	incoming_q = append(incoming_q, or, or1, or2, bid1, bid2, bid3, bid4)

	OB := Order_Book{ask_l, bid_l}

	inserter(&incoming_q, ask_l, bid_l)


/*
	fmt.Println("----------- INITIAL LIST -----------")
	fmt.Println("ASK list")
	for ind, pp := range ask_l {

		if len(pp.items) != 0 {
			fmt.Printf("Value: %v   %v  | Point:%v     %p\n", ind, *pp, *&pp.items[0], &pp.items[0])
		} else {
			fmt.Printf("Value: %v   %v  | Point:%v     %p\n", ind, *pp, *&pp.items, &pp.items)
		}
	}

	fmt.Println("BID list")
	for ind, pp := range bid_l {

		if len(pp.items) != 0 {
			fmt.Printf("Value: %v   %v  | Point:%v     %p\n", ind, *pp, *&pp.items[0], &pp.items[0])
		} else {
			fmt.Printf("Value: %v   %v  | Point:%v     %p\n", ind, *pp, *&pp.items, &pp.items)
		}
	}

	fmt.Println("------------------------------------")

	fmt.Printf("Incoming quote: %v     Pointer: %p\n", incoming_q, &incoming_q)
*/
	//ticker(bid_l, "Bids")
	//ticker(ask_l, "Ask")

	

	Order_Book_print(OB,ORDER_BOOK_LENGTH)

	// test ob insert

	incoming_q = append(incoming_q, bid5) 
	inserter(&incoming_q, OB.ask, OB.bid)
	Order_Book_print(OB,ORDER_BOOK_LENGTH)

	fmt.Printf("Incoming after test of 5: %v \n",incoming_q)
	// test que order

	incoming_q = append(incoming_q, bido1,bido2) 
	inserter(&incoming_q, OB.ask, OB.bid)
	Order_Book_print(OB,ORDER_BOOK_LENGTH)
	fmt.Printf("Taking out: %v \n",OB.bid[2].Dequeue())
	Order_Book_print(OB,ORDER_BOOK_LENGTH)

	fmt.Printf("Incoming after test of order: %v \n",incoming_q)

	// test inseter reset


}
