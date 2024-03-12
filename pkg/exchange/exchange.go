package exchange


import (
    "fmt"
	"Go_net/pkg/common"
)



///////////////Variables///////////////////////
/*
var or = common.Order{O_type:"ask", Price:10, Volume:15}
var or1 = common.Order{O_type:"ask",Price: 10,Volume: 1}
var or2 = common.Order{O_type:"ask",Price: 10, Volume:5}
var bid1 = common.Order{O_type:"bid",Price: 9,Volume: 9}
var bid2 = common.Order{O_type:"bid",Price: 9, Volume:1}
var bid3 = common.Order{O_type:"bid",Price: 7,Volume: 100}
var bid4 = common.Order{O_type:"bid",Price: 9,Volume: 2}
*/

/*
var filled_quotes []common.Order_Filled

var type_ob bool

var ask_l [common.ORDER_BOOK_LENGTH]*common.Queue
var bid_l [common.ORDER_BOOK_LENGTH]*common.Queue

var incoming_q []common.Order
*/
/////////////////////////////////////////


// It puts the single Order in the Order Book
func Inserter(l_in_quo *[]common.Order, order_bo common.Order_Book) {

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



//////////// Matching Engine /////////////

// Resolve the fact that every quoted is matced
// if the order filled is empty do not considered in the
// loop
func Match(order_b common.Order_Book)(fill common.Order_Filled){

	var matches [2]common.Order
	lb,_,la,_:=Find_best(order_b)

	if order_b.Ask[la].Is_Empty() || order_b.Bid[lb].Is_Empty(){

		fmt.Printf("\n-----No match-----\nEmpy list\n\nAsk list: %v\nBid list: %v",order_b.Ask[la], order_b.Bid[lb])

		// if no order is filled the price will be -1
		fill=common.Order_Filled{
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

			fill=common.Order_Filled{
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

			fill=common.Order_Filled{
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
		fill=common.Order_Filled{
			Price: -1,
		}

	} else{
		fmt.Printf("\n-----No match-----\n")

		// if no order is filled the price will be -1
		fill=common.Order_Filled{
			Price: -1,
		}
	}

	//Order_Book_print()

	return 	fill

}




// if there is no best quote it will rerurn {0,[]}
func Find_best (ob common.Order_Book) (level_b int, best_b []common.Order,level_a int, best_a []common.Order) {
	
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


func Engine(list_incoming *[]common.Order,Or_Bo *common.Order_Book,list_filled *[]common.Order_Filled){

    for{

        if len(*list_incoming)!=0{

            Inserter(list_incoming,*Or_Bo)
            Filled_Or:= Match(*Or_Bo)
            common.Order_Book_print(*Or_Bo,common.ORDER_BOOK_LENGTH,false)
            if Filled_Or.Price==-1{
                continue
            }
            *list_filled = append(*list_filled,Filled_Or)
            common.Order_Book_print(*Or_Bo,common.ORDER_BOOK_LENGTH,false)

            for { // the -1 means that there is no match

                Filled_Or:= Match(*Or_Bo)

                if Filled_Or.Price==-1{
                    break
                }

                *list_filled = append(*list_filled,Filled_Or)
                common.Order_Book_print(*Or_Bo,common.ORDER_BOOK_LENGTH,false)

            }
            fmt.Printf("End main engine loop: %v",*list_filled)
            common.Order_Book_print(*Or_Bo,common.ORDER_BOOK_LENGTH,false)
            
        }
    }
}



///////////////////////////////////////////////////////////////////////////////
/*
// StartExchange starts the Exchange
func StartExchange() {
    

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

	Inserter(&incoming_q, OB)
	
	common.Order_Book_print(OB,ORDER_BOOK_LENGTH,false) 

}
*/

