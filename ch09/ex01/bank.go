package bank

type withdraw struct {
	amount   int
	resultCh chan bool
}

var deposits = make(chan int)
var balances = make(chan int)
var withdraws = make(chan withdraw)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	r := make(chan bool)
	withdraws <- withdraw{amount: amount, resultCh: r}
	return <-r
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case w := <-withdraws:
			if balance < w.amount {
				w.resultCh <- false
			} else {
				balance -= w.amount
				w.resultCh <- true
			}
		}
	}
}

func init() {
	go teller()
}
