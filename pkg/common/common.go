package common

import (
	"fmt"

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

var l_order_filled []Order_Filled

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

////////////////////////////////

//////////// Matching Engine /////////////

// Resolve the fact that every quoted is matced
// if the order filled is empty do not considered in the
// loop
func Match(order_b Order_Book)(fill Order_Filled){

	var matches [2]Order
	lb,_,la,_:=Find_best(order_b)

	if order_b.Ask[la].Is_Empty() || order_b.Bid[lb].Is_Empty(){

		fmt.Printf("\n-----No match-----\nEmpy list\n\nAsk list: %v\nBid list: %v",order_b.Ask[la], order_b.Bid[lb])

		// if no order is filled the price will be -1
		fill=Order_Filled{
			Price: -1,
		}
		return
	}


	matches[0]=order_b.Ask[la].Observe()
	matches[1]=order_b.Bid[lb].Observe()

	if lb==la {       // Normal match
		
		if matches[0].Volume == matches[1].Volume { // Same volume -> so we take out both ordes from the Queue
			
			ask_f :=order_b.Ask[la].Dequeue()
			bid_f :=order_b.Bid[lb].Dequeue()

			fill=Order_Filled{
				Og_bid: bid_f,
				Og_ask: ask_f,
				Price: ask_f.Price,
				Vol_filled: ask_f.Volume,
	
			}
			
		}else if matches[0].Volume > matches[1].Volume{ // Ask>Bid Volume
			bid_f :=order_b.Bid[lb].Dequeue()

			fill=order_b.Ask[la].Items[0].Partial_fill(bid_f) //modify ask
			fmt.Printf("\nFill: %v (ask volume >bid vol)-------\n",fill)


		}else if matches[0].Volume < matches[1].Volume{ // Ask<Bid Volume
			ask_f :=order_b.Ask[la].Dequeue()
			fill=order_b.Bid[lb].Items[0].Partial_fill(ask_f) //modify ask
			fmt.Printf("\nFill: %v (ask volume <bid vol)------\n",fill)
		}	
		
	// overshoot
	} else if lb>la { // Bid price is larger than the best Ask
				

		if matches[0].Volume == matches[1].Volume { // Same volume -> so we take out both ordes from the Queue
			
			ask_f :=order_b.Ask[la].Dequeue()
			bid_f :=order_b.Bid[lb].Dequeue()

			fill=Order_Filled{
				Og_bid: bid_f,
				Og_ask: ask_f,
				Price: ask_f.Price,
				Vol_filled: ask_f.Volume,
	
			}
			
		}else if matches[0].Volume > matches[1].Volume{ // Ask>Bid Volume
			bid_f :=order_b.Bid[lb].Dequeue()

			fill=order_b.Ask[la].Items[0].Partial_fill(bid_f) //modify ask
			fmt.Printf("\nFill: %v (ask volume >bid vol)-------\n",fill)


		}else if matches[0].Volume < matches[1].Volume{ // Ask<Bid Volume
			ask_f :=order_b.Ask[la].Dequeue()
			fill=order_b.Bid[lb].Items[0].Partial_fill(ask_f) //modify ask
			//fill.price=ask_f.price
			fmt.Printf("\nFill: %v (ask volume <bid vol)------\n",fill)
		}	


	} else if  len(matches[1].O_type)==0 || len(matches[0].O_type)==0 {
		fmt.Printf("\n-----No match-----\nEmpy list\n\n%v %v",matches[0],matches[1])

		// if no order is filled the price will be -1
		fill=Order_Filled{
			Price: -1,
		}

	} else{
		fmt.Printf("\n-----No match-----\n")

		// if no order is filled the price will be -1
		fill=Order_Filled{
			Price: -1,
		}
	}


	return 	fill

}


//////////////////////////////////////

// Made up orders for testing
/*
var or = Order{"ask", 10, 15}
var or1 = Order{"ask", 10, 1}
var or2 = Order{"ask", 10, 5}
var bid1 = Order{"bid", 9, 9}
var bid2 = Order{"bid", 9, 1}
var bid3 = Order{"bid", 7, 100}
var bid4 = Order{"bid", 9, 2}

var bid5 = Order{"bid", 10, 15}
var bido1 = Order{"bid", 2, 1}
var bido2 = Order{"bid", 2, 2}
// test 134
var bido138 = Order{"bid", 10, 2}
// test bid>ask
var askbigger = Order{"ask", 9, 2} // it will interact with bid1

*/

var ask_l [ORDER_BOOK_LENGTH]*Queue
var bid_l [ORDER_BOOK_LENGTH]*Queue

var incoming_q []Order

/////////////////////////////////////////



// It puts the single Order in the Order Book
func Inserter(l_in_quo *[]Order, order_bo Order_Book) {

	l_ask:= order_bo.Ask
	l_bid := order_bo.Bid

	fmt.Println("                  §")
	fmt.Printf("++++++++++++++++++++++++++++++++++++\nStart insertion\nINFUNC Incoming quote: %v     Pointer: %p\n", l_in_quo, l_in_quo)

	switch {
	case len(*l_in_quo) == 1 && (*l_in_quo)[0].O_type == "ask":

		l_ask[int((*l_in_quo)[0].Price)].Enqueue((*l_in_quo)[0])
		//return

	case len(*l_in_quo) == 1 && (*l_in_quo)[0].O_type == "bid":

		l_bid[int((*l_in_quo)[0].Price)].Enqueue((*l_in_quo)[0])
		//return

	case len(*l_in_quo) > 1:

		for _, v := range *l_in_quo {
			if v.O_type == "ask" {
				l_ask[int(v.Price)].Enqueue(v)
			} else if v.O_type == "bid" {
				l_bid[int(v.Price)].Enqueue(v)
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
