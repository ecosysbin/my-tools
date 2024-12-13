//func main() {
// 1. function
//	wg := &sync.WaitGroup{}
//	wg.Add(10)
//	for i := 0; i < 10; i++{
//		go goroutin(wg)
//	}
//	wg.Wait()
//	fmt.Println("all goroutin finish")
//}
//
//func goroutin(w *sync.WaitGroup) {
//	defer w.Done()
//	fmt.Println("a")
//}
// 等待所有协程执行完成
func main() {
	//chann := make(chan int, 10)
	//for i := 0; i < 10; i++ {
	//	go func(in chan int){
	//		defer func(){chann <- 1}()
	//		fmt.Println("this channel!")
	//	}(chann)
	//}

	// 2. function
	//for {
	//	if len(chann) == 10 {
	//		close(chann)
	//		break
	//	}
	//}
	// 3. function
	//for i := 0; i < 10; i++ {
	//	<- chann
	//}
	//fmt.Println("all channel finish!")
	fmt.Println(test1())
}

func test1() int {
	t := 5
	defer func() int{
		t = t + 5
		return t
	}()

	return t
}
