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


type Order_Filled struct {
	//bid_ID int
	//ask_ID int
	og_bid Order // orgiginal quote
	og_ask Order // 
	price  float32
	vol_filled_ask int32
	vol_filled_bid int32
}


var l_order_filled []Order_Filled

/////////////Queue////////////////

// Queue that the orders need to follow in any given price.
// The order is first-in first-out (time priority) FIFO.
// To enter and to 
// Observe gives the first in line Order without taking it out of the Queue
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

func (q *Queue) Observe() Order {
	to_see := q.items[0]

	return to_see
}


////////////////Order Book//////////////////

type Order_Book struct {
	ask [ORDER_BOOK_LENGTH]*Queue
	bid [ORDER_BOOK_LENGTH]*Queue
}

////////////////////////////////

//////////// Matching Engine /////////////

func match(order_b Order_Book)(fill Order_Filled){

	var matches [2]Order
	lb,_,la,_:=find_best(order_b)
	matches[0]=order_b.ask[la].Observe()
	matches[1]=order_b.bid[lb].Observe()


	if lb==la {       // Normal match
		if matches[0].volume == matches[1].volume{ // Same volume -> so we take out both ordes from the Queue
			ask_f :=order_b.ask[la].Dequeue()
			bid_f :=order_b.bid[lb].Dequeue()
			fill=Order_Filled{
				og_bid: bid_f,
				og_ask: ask_f,
				price: ask_f.price,
				vol_filled_ask: ask_f.volume,
				vol_filled_bid: bid_f.volume,
			}
			
		}
		
	} else if lb>la { // Bid price is larger than the best Ask

	} else if lb<la { // Ask price is larger than the best Bid

	} else{
		fmt.Printf("\n-----No match-----\n")
	}
	// var hitted []Order_Filled
	//fill=append(fill,order_b.ask[la].Observe(),order_b.bid[lb].Observe())

	return 	

}


//////////////////////////////////////

// Made up orders for testing

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


var ask_l [ORDER_BOOK_LENGTH]*Queue
var bid_l [ORDER_BOOK_LENGTH]*Queue

var incoming_q []Order

/////////////////////////////////////////



// It puts the single Order in the Order Book
func inserter(l_in_quo *[]Order, order_bo Order_Book) {

	l_ask:= order_bo.ask
	l_bid := order_bo.bid

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



func Size_Level(level_list []Order)(size int){

	for i:=0; i< len(level_list);i++{
		size= size+int(level_list[i].volume)
	}

	return
}



func Order_Book_print(OB Order_Book, lenght_OB int8,size_only bool) {
	//fmt.Printf("Order Book: %v  \n %v \n", OB, OB.ask)
	fmt.Println("                  °")
	fmt.Println("----------- Order Book --------------")
	fmt.Println("            Bid        Ask ")
	
	for  i := 0; i < int(lenght_OB); i++ {

		if size_only{
			fmt.Printf("| %v | - %v - | %v |\n",Size_Level(OB.bid[i].items),i,Size_Level(OB.ask[i].items))
		}else{
			fmt.Printf("Level: %v | %v   -   %v  |\n", i, OB.bid[i], OB.ask[i])
		}
		//Size_Level(OB.bid[i].items) //i, OB.bid[i].items, OB.ask[i].items)
	}

	fmt.Println("-------------------------------------")
	fmt.Println("                  °")
}


// if there is no best quote it will rerurn {0,[]}
func find_best (ob Order_Book) (level_b int, best_b []Order,level_a int, best_a []Order) {
	
	fmt.Printf("##################################\n\nSearching for the BEST ASK\n")
	for index := len(ob.bid)-1;index >= 0; index-- {
		fmt.Printf("Index: %v | Items: %v  ->len:%v\n",index, ob.bid[index].items,len(ob.bid[index].items))

		if len(ob.bid[index].items)!=0{
			level_b= index
			best_b=ob.bid[index].items
			break
		}

	}

	fmt.Printf("\n===================================\n\nSearching for the BEST ASK\n")
	for index, value := range ob.ask{
		fmt.Printf("Index: %v | Items: %v  ->len:%v\n",index, value.items,len(value.items))

		if len(value.items)!=0{
			level_a= index
			best_a=value.items
			break
		}

	}
	fmt.Printf("\n##################################\n\n")
	return
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

	inserter(&incoming_q, OB)

	Order_Book_print(OB,ORDER_BOOK_LENGTH,false)

	lb,vb,la,va:=find_best(OB)
	fmt.Printf("\nTest best bid: %v %v \nTest best ask: %v %v \n",lb,vb,la,va)
	fmt.Printf("\n\nSize of:%v  ----->   %v  \n\n",la,Size_Level(va))

	//bid5
	incoming_q = append(incoming_q, bid5)
	inserter(&incoming_q, OB)

	Order_Book_print(OB,ORDER_BOOK_LENGTH,false)
	a:= match(OB)
	l_order_filled= append(l_order_filled, a)
	fmt.Printf("\n\nFilled orders : %v\n\n",a)
	fmt.Printf("\n\nList of Filled orders : %v\n\n",l_order_filled)
	Order_Book_print(OB,ORDER_BOOK_LENGTH,false)
}

