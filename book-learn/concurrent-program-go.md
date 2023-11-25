# 1.golang 创建一个协程

## 1.1 go语句和goroutine

```go
go func() {
  fmt.println("Go! Goroutine!")
}
```

Go运行时系统对go语句中的函数的执行是并发。go语句执行的时候，其中的go函数会被单独放入一个goroutine中，这之后该go函数的执行独立于当前goroutine运行。

## 1.2 time.Sleep()干预go执行

在Go中有很多方法干预多个G的执行顺序

```go
package main

import "time"

func main() {
	go println("Go! Goroutine!")
	time.Sleep(time.Millisecond)
}
```

time.Sleep的作用是让调用它的goroutine暂停一段时间，runtime.GoSched函数是暂停当前的G，这里也有效。

## 1.3 go函数添加函数声明

```go
import (
	"fmt"
	"time"
)

func main() {
	names := []string{"Eric", "Harry", "Robert", "Mark"}
	for _, name := range names {
		go func(name string) {
			fmt.Printf("Hello %s\n", names)
		}(name)
	}
	time.Sleep(time.Millisecond)
}

```



# 2.golang停止一个协程



在 Golang 中，协程的停止可以通过几种方式来实现：

## 1、使用 `context.Context` 进行协程的取消

可以使用 `context.Context` 来控制协程的生命周期，从而达到停止协程的目的。

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func printNum(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("协程被取消")
			return
		default:
			fmt.Printf("子协程的数字：%d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go printNum(ctx, 10)

	//3秒后取消协程
	time.Sleep(3 * time.Second)
	cancel()

	//等待协程结束
	time.Sleep(1 * time.Second)
	fmt.Println("出协程退出")
}
```

在上述代码中，我们通过 `context.Context` 来实现了协程的取消。在 `printNum` 函数中，我们使用 `select` 语句来监听 `ctx.Done()` 是否被关闭，如果关闭了，则协程被取消。在 `main` 函数中，我们通过 `context.WithCancel` 函数创建了一个带有取消功能的上下文对象 `ctx`，并将其传递给 `printNum` 函数。在主协程中，我们等待 3 秒钟后调用 `cancel` 函数来取消协程的执行。最后，我们等待协程结束，并输出 `主协程退出`。

## 2、使用 `channel` 进行协程的关闭

可以通过关闭 `channel` 来实现协程的停止。

```go
package main

import (
	"fmt"
	"time"
)

func printNum(stop <-chan struct{}, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-stop:
			fmt.Println("协程被关闭")
			return
		default:
			fmt.Printf("子协程的数字：%d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	stopChan := make(chan struct{})
	go printNum(stopChan, 10)
	time.Sleep(3 * time.Second)
	close(stopChan)

	time.Sleep(time.Second * 1)
	fmt.Println("主协程退出")
}
```

在上述代码中，我们使用 `channel` 来实现了协程的停止。在 `printNum` 函数中，我们通过监听 `stopCh` 的关闭来实现协程的停止。在 `main` 函数中，我们创建了一个 `stopCh` 的无缓冲通道，并将其传递给 `printNum` 函数。在主协程中，我们等待 3 秒钟后关闭 `stopCh` 通道，从而停止协程的执行。

## 3、使用布尔变量

使用一个 boolean 变量来控制协程的执行。在协程中，检查这个变量的值，如果为 true，则退出协程的执行。

```go
package main

import (
	"fmt"
	"time"
)

func worker(stopChan chan bool) {
	for {
		select {
		case <-stopChan:
			fmt.Println("stop...")
			return
		default:
			fmt.Println("working")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
func main() {
	stopChan := make(chan bool)
	go worker(stopChan)

	time.Sleep(2 * time.Second)
	stopChan <- false
	time.Sleep(time.Second)
}
```

# 3.Golang协程状态

在Golang中，协程状态分为以下几种：

1. 新建状态（New）：协程被创建后，但是还没有开始执行任务。
2. 运行状态（Runnable）：协程已经被调度，并且正在执行任务。
3. 阻塞状态（Blocked）：协程因为等待某些条件而被暂停执行，比如等待IO操作完成、等待锁释放、等待channel接收数据等。
4. 死亡状态（Dead）：协程已经执行完成或者因为某些异常而结束。

在Golang中，协程的状态是由调度器来控制的，程序员无法直接控制协程的状态。当一个协程处于阻塞状态时，调度器会把该协程从可运行队列中移除，并且把其它可运行的协程加入到可运行队列中等待调度。当阻塞状态的协程恢复执行时，调度器会把该协程重新放回到可运行队列中，等待下一次调度执行。

```go
package main

import (
	"fmt"
	"time"
)

func A() {
	fmt.Println("New: 协程已被创建但还未开始执行任务")
}

func B() {
	fmt.Println("Runnable: 协程正在执行任务")
	time.Sleep(time.Second)
}

func main() {
	//新建
	go A()

	//运行
	go B()

	ch := make(chan bool)
	go func() {
		fmt.Println("Blocked: 协程因为等待channel接收数据而被暂停执行")
		<-ch
	}()

	// 死亡状态
	go func() {
		fmt.Println("Dead: 协程执行完成")
	}()

	time.Sleep(time.Second * 2)
}
```

在这个例子中，我们创建了四个协程并打印了协程的状态，每个协程的状态分别为：

1. `foo()`函数所在协程的状态为新建状态（New）。
2. `bar()`函数所在协程的状态为运行状态（Runnable）。
3. 包含`<-ch`语句的协程的状态为阻塞状态（Blocked）。
4. 包含`fmt.Println("Dead: 协程执行完成")`语句的协程的状态为死亡状态（Dead）。

> 注意，在主协程中，我们使用了`time.Sleep()`函数来延迟程序的执行，以确保所有协程都有足够的时间来执行并展示其状态。

# 4.golang协程安全

## 什么是协程安全

协程安全（goroutine safety）是指在多个协程（goroutines）并发执行时，对共享变量和资源的访问不会导致数据竞争和不一致性的问题。

在Golang中，协程是由Go运行时（Go runtime）调度的轻量级线程。由于协程之间并发访问共享变量的问题可能导致数据竞争和不一致性，因此需要采取一些措施来确保协程安全。

## 协程安全的解决办法

Golang提供了几种方法来确保协程安全，包括：

1. 使用互斥锁（Mutex）：互斥锁是一种常用的同步机制，用于保护共享变量的访问。在任何时候只有一个协程可以持有互斥锁，并访问共享变量。其他协程需要等待锁被释放才能访问共享变量。
2. 使用读写锁（RWMutex）：读写锁是一种特殊的互斥锁，可以同时允许多个协程对共享变量进行读操作，但在写操作时需要排他地持有锁。使用读写锁可以提高并发访问共享变量的效率。
3. 使用通道（Channel）：通道是一种用于协程之间通信的机制，通道可以确保数据的同步和安全访问。使用通道可以避免协程之间访问共享变量的问题。
4. 避免共享状态：如果可能的话，应该尽量避免使用共享变量，而是将状态封装在单个协程内部，然后使用通道或函数参数进行通信。

> 需要注意的是，并不是所有的协程都需要保证并发安全，只有在协程之间共享变量或资源时才需要考虑并发安全的问题。在设计并发程序时，应该根据实际需求和情况来选择合适的并发安全策略。

使用mutex :

```go
package main

import (
	"fmt"
	"sync"
)

var (
	num int
	m   sync.Mutex
)

func increase(wg *sync.WaitGroup, caseNum int) {
	for i := 0; i < 500; i++ {
		m.Lock()
		num += 1
		fmt.Println("第", caseNum, "个case", num)
		m.Unlock()
	}
	wg.Done()
}

func main() {
	num = 0
	wg := sync.WaitGroup{}

	wg.Add(2)
	go increase(&wg, 1)
	go increase(&wg, 2)

	wg.Wait()
	fmt.Println(num)
}
```

# 5.golang共享变量和临界区

在并发编程中，共享变量是指多个线程或协程都可以访问和修改的同一个变量，临界区则是指在程序执行过程中，访问共享变量的代码块。当多个线程或协程同时访问共享变量时，就可能发生数据竞争，导致程序出现错误或不确定的行为。

为了避免数据竞争，需要在访问共享变量的代码块中使用同步机制，例如互斥锁、读写锁等，来保证同一时间只有一个线程或协程能够访问共享变量，这样就可以避免数据竞争。这些使用同步机制的代码块就是临界区。

在编写并发程序时，需要仔细设计共享变量和临界区，并使用适当的同步机制来保证协程安全。

```go
package main

import (
	"fmt"
	"sync"
)

var num = 0
var mu sync.Mutex

func main() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		for i := 0; i < 10000; i++ {
			mu.Lock()
			num++
			mu.Unlock()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			mu.Lock()
			num--
			mu.Unlock()
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("num = ", num)
}
```

使用 sync.Mutex 来保护 x 的访问，保证了 x 的最终值为 0。在协程访问共享变量的代码块中，需要使用互斥锁来保证同一时间只有一个协程可以访问共享变量，这就是临界区。这样就避免了数据竞争，保证了程序的正确性。

# 6.golang协程优先级（无优先级说法

在 Go 语言中，协程（goroutine）没有单独的优先级设置，而是由 Go 运行时（runtime）来动态管理。在运行时，每个协程都被分配一个逻辑处理器（logical processor），逻辑处理器负责在物理处理器上执行协程。Go 运行时会自动在逻辑处理器之间分配协程，确保尽可能地利用所有可用的处理资源，从而提高并发性能。

由于 Go 语言的运行时会自动进行协程调度和资源分配，因此开发者无需关心协程的优先级问题。同时，Go 还提供了一些调试工具，可以帮助开发者诊断协程相关的性能问题，例如可以通过 `go tool trace` 命令查看协程的调度情况和执行时间，从而优化协程的运行效率。

# 7.golang协程安全数据类型

在 Go 语言中，有些类型是协程安全的，可以在多个协程之间安全地共享使用，而有些类型则不是。以下是一些常见的协程安全类型：

1. 基本数据类型：例如 `int`、`float64`、`bool` 等基本数据类型是协程安全的，可以在多个协程之间安全地进行读写操作。
2. 字符串（`string`）：字符串是协程安全的，可以在多个协程之间共享使用。
3. 切片（`slice`）：切片是协程安全的，但需要注意在多个协程之间对同一切片进行读写操作时，可能会出现数据竞争问题，需要使用互斥锁或通道等机制进行同步。
4. Map（`map`）：`map` 是非协程安全的，多个协程同时读写同一个 map 可能会出现数据竞争问题。需要使用互斥锁或通道等机制进行同步。
5. Channel（`channel`）：`channel` 是协程安全的，可以在多个协程之间进行安全的通信。

总的来说，在使用协程时，需要注意哪些类型是协程安全的，哪些类型不是协程安全的，避免因数据竞争问题导致程序出现异常行为。同时，在对非协程安全类型进行读写操作时，需要采用同步机制进行保护，例如使用互斥锁、读写锁、原子操作或通道等。

# 8.golang如何解决协程安全问题

在 Go 语言中，协程是由关键字 `go` 创建的，可以理解为一种轻量级的线程，Go 语言内置了一些协程相关的包和类，主要包括以下几个：

1. `go` 关键字：用于创建并发执行的协程。
2. `sync` 包：提供了互斥锁（Mutex）、读写锁（RWMutex）和条件变量（Cond）等同步原语，以及一些原子操作函数，用于协程之间的同步和互斥访问。
3. `channel` 类型：用于协程之间的通信，可以通过 `make(chan T)` 创建一个类型为 T 的通道，通过通道进行协程之间的数据传输和同步。
4. Go语言的`sync/atomic`包提供了一些原子操作的函数和类型，可以用于处理原子变量。
5. `context` 包：提供了跨协程的上下文传递和取消功能，用于在多个协程之间传递请求、参数和状态信息。
6. `timer` 包：提供了定时器功能，可以在指定时间后执行一个函数或发送一个事件到通道。
7. `select` 语句：用于监听多个通道的数据流动，可以阻塞等待任意一个通道有数据可读或有数据可写，从而实现协程之间的同步。
8. `runtime` 包：提供了与 Go 运行时相关的函数和变量，例如 `GOMAXPROCS` 变量可以设置并发执行的 CPU 核心数，`goexit` 函数可以让当前协程退出，`Gosched` 函数可以让当前协程主动让出 CPU 等待调度。

# 9.golang通道

在 Go 中，通道（channel）是一种特殊的类型，用于在不同的协程（goroutine）之间同步和传递数据。通道可以看作是一种类型安全的管道，通过它们可以安全地传递数据。

通道是 Go 语言中非常重要的并发原语，它是 Go 语言中多个协程之间同步数据交换的主要方式之一，被广泛应用于并发编程中。

通道具有两个主要的操作，即发送和接收。通过通道发送数据时，发送者会将数据传入通道，并等待接收者接收这个数据。接收者从通道中读取数据时，如果通道中没有数据，它会等待数据的到来，直到有数据可用为止。

在 Go 语言中，通道可以使用内置函数 `make()` 来创建。通道的类型指定了通道中能够传输的数据的类型。例如，一个传输字符串的通道可以使用 `make(chan string)` 来创建。

通道还可以通过使用关键字 `chan` 来定义。例如，一个传输整数的通道可以这样定义：`var ch chan int`。

```go
package main

import "fmt"

func main() {

	ch := make(chan int)

	go func() {
		for i := 1; i <= 30; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for {
			num, ok := <-ch
			if ok {
				fmt.Println("receive", num)
			} else {
				fmt.Println("channel was closed")
				break
			}
		}

	}()

	fmt.Scanln()
	fmt.Println("main goroutine was done")
}
```

首先创建了一个通道 `ch`，存储 `int` 类型的数据。然后创建了两个协程，一个向通道发送数据，另一个从通道接收数据。

向通道发送数据的协程中，我们使用 `for` 循环将 `1` 到 `5` 的整数依次发送到通道中，然后在发送完数据后调用 `close` 函数关闭通道。

从通道接收数据的协程中，我们使用 `for` 循环不断地从通道中接收数据，直到通道被关闭。在接收数据时，我们使用了特殊的语法：`num, ok := <-ch`，其中 `num` 表示接收到的数据，`ok` 表示通道是否还打开。当 `ok` 为 `false` 时，表示通道已经被关闭，我们就可以退出循环。

最后我们在主函数中调用 `fmt.Scanln()` 函数，等待协程执行完毕。当协程执行完毕后，主函数会输出 `Main goroutine is done`。

这个例子中，我们通过通道实现了两个协程之间的通信和同步，确保了数据的正确性和同步性。

# 10.golang通道缓冲

在 Golang 中，通道可以被缓冲，这意味着通道可以在未读取之前拥有多个值。

通道缓冲提供了一种机制，使发送方可以在接收方准备好接收数据之前发送多个值，而不必等待接收方。在缓冲区填满之前，发送方将阻止，并在缓冲区被读取之前，接收方将阻止。

缓冲区大小是在创建通道时指定的，如下所示：

```go
ch := make(chan int, 3) // 创建一个缓冲大小为 3 的通道
```

在这个示例中，`ch` 是一个具有 3 个缓冲区的通道。

当缓冲区已满时，发送方将被阻止，直到缓冲区有可用空间为止。同样，当缓冲区为空时，接收方将被阻止，直到有一个值可用为止。通道缓冲使通信更有效，因为它减少了 goroutine 阻塞等待的数量，从而提高了程序的性能。

在下面的示例中，我们使用一个缓冲区大小为 2 的通道来模拟生产者和消费者：

```go
package main

import "fmt"

func Producer(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
		fmt.Println("send:", i)
	}
	close(ch)
}

func Consumer(ch chan int) {
	for {
		val, ok := <-ch
		if !ok {
			break
		}
		fmt.Println("consume:", val)
	}
}

func main() {
	ch := make(chan int, 2)
	go Producer(ch)
	Consumer(ch)
}
```

## 实例

假设有两个协程，一个协程需要向另一个协程发送数据，可以使用通道来实现。

```go
package main

import (
	"fmt"
	"time"
)

func send(ch chan int) {
	fmt.Println("send the first msg")
	ch <- 1
	fmt.Println("send the second msg")
	ch <- 2
}
func main() {
	ch := make(chan int,1)
	go send(ch)
	time.Sleep(time.Second * 2)
	fmt.Println("receive 1")
	fmt.Println(<-ch)
	fmt.Println("receive 2")
	fmt.Println(<-ch)
}
```

在这个示例中，使用 `make` 创建了一个缓冲大小为 1 的通道。在 `send` 协程中，先向通道中发送了 1，然后再发送 2。在 `main` 函数中，先接收了从 `send` 协程中发送的 1，然后再接收 2。由于通道缓冲大小为 1，第一个消息可以被缓存，所以 `main` 函数不必阻塞等待 `send` 协程发送第一个消息，而可以立即接收。

如果将通道缓冲大小设置为 0，则创建的是无缓冲通道，示例代码如下：

```go
package main

import (
    "fmt"
)

func send(ch chan int) {
    fmt.Println("Sending 1st message")
    ch <- 1
    fmt.Println("Sending 2nd message")
    ch <- 2
}

func main() {
    ch := make(chan int)
    go send(ch)
    fmt.Println("Receiving 1st message")
    fmt.Println(<-ch)
    fmt.Println("Receiving 2nd message")
    fmt.Println(<-ch)
}
```

在这个示例中，创建的是一个无缓冲通道，所以在 `send` 协程中发送第一个消息后，`send` 协程会一直阻塞等待 `main` 函数接收该消息。当 `main` 函数接收了第一个消息后，`send` 协程才会被解除阻塞并发送第二个消息。可以看出，无缓冲通道保证了消息的同步传输，即发送方发送消息后会一直阻塞等待接收方接收消息。

# 11.golang通道同步

在 Go 中，通道是一种同步原语，可以用来在不同的 goroutine 之间传递消息并进行同步。

通道同步指的是：当一个 goroutine 向通道发送数据时，如果没有其他 goroutine 在接收这个数据，发送操作会被阻塞，直到有其他 goroutine 接收了这个数据为止。同样的道理，当一个 goroutine 从通道接收数据时，如果没有其他 goroutine 向这个通道发送数据，接收操作也会被阻塞，直到有其他 goroutine 向这个通道发送数据为止。

这种同步机制可以帮助我们避免 race condition（竞态条件）的发生，保证多个 goroutine 之间的数据访问安全。

## 实例

假设有两个协程 A 和 B，它们分别执行不同的任务，并且协程 A 的任务需要先执行完后，协程 B 才能继续执行，这时候可以使用通道来实现它们之间的同步。

具体实现方法是，让协程 A 在执行完任务后往一个通道中发送一个消息，然后在协程 B 中等待从该通道中接收到消息后再执行任务。这样就可以保证协程 A 先执行完任务，协程 B 再开始执行任务。

```go
package main

import (
	"fmt"
	"time"
)

func chanSync1(ch chan string) {
	fmt.Println("goroutine 1")
	ch <- "flag"
}

func chanSync2(ch chan string) {
	msg := <-ch
	fmt.Println("goroutine 2")
	fmt.Println("msg from chanSync 1", msg)
}

func main() {
	ch := make(chan string)

	go chanSync2(ch)
	go chanSync1(ch)

	time.Sleep(time.Second)
	fmt.Println("done")

}
```

在上面的代码中，`worker1` 协程执行完任务后，往通道 `ch` 中发送了一条消息，然后 `worker2` 协程从该通道中接收到消息后才开始执行任务。在 `main` 函数中使用 `Scanln` 等待用户输入，以保证协程能够执行完毕并输出结果。

协程 `worker1` 先执行完任务，并把消息发送到了通道中，然后协程 `worker2` 接收到消息后才开始执行任务。这样就保证了协程之间的同步。

# 12 .golang通道方向

在golang中，可以使用通道的方向来限制通道的发送和接收操作。通道的方向可以是只发送、只接收或双向。通过限制通道的方向，可以提高程序的安全性和可读性。

在声明通道时，可以使用<-运算符来指定通道的方向。例如，要创建一个只发送int的通道，可以使用以下声明：

```go
var sendCh chan<- int
```

这样就创建了一个sendCh通道，只能用于发送int类型的值。类似地，如果要创建一个只接收int的通道，可以使用以下声明：

```go
var recvCh <-chan int
```

这样就创建了一个recvCh通道，只能用于接收int类型的值。如果要创建一个双向的通道，可以使用以下声明：

```go
var ch chan int
```

这样就创建了一个ch通道，可以用于发送和接收int类型的值。

需要注意的是，如果试图在通道的方向不匹配的情况下进行通道操作，将会在编译时产生错误。例如，如果试图在只发送int的通道中进行接收操作，或者在只接收int的通道中进行发送操作，都会导致编译错误。

## 实例

通道方向指的是通道的发送和接收操作所允许的方向，即通道是单向的还是双向的。

在golang中，可以通过在通道类型中添加箭头符号来指定通道的方向，其中`<-`用于指定发送方向，`->`用于指定接收方向，而不加箭头符号则表示双向通道。

举个例子，假设我们有一个需要从主协程向子协程发送消息的场景，可以定义一个只允许发送的单向通道，示例如下：

```go
func main() {
    msgCh := make(chan string)

    go func(ch chan<- string) {
        ch <- "hello from child goroutine"
    }(msgCh)

    msg := <-msgCh
    fmt.Println(msg)
}
```

在上面的代码中，我们通过使用`chan<-`指定了`msgCh`通道只能用于发送，因此在子协程中我们只能往通道中发送消息，而不能从通道中接收消息，从而保证了通道方向的一致性和通道安全性。

# 13.golang通道选择器

在 Go 语言中，通道选择器（Channel Selector）是一种通过 select 语句同时等待多个通道操作的机制。通道选择器可以让程序同时等待多个通道，一旦其中任意一个通道可以进行读写操作时，程序就可以立即响应该通道的操作，而不是在其他通道等待的时间里被阻塞。

通道选择器的语法如下：

```go
select {
case <- channel1:
    // 执行 channel1 的操作
case data := <- channel2:
    // 执行 channel2 的操作
case channel3 <- data:
    // 执行 channel3 的操作
default:
    // 如果上述通道都没有操作，则执行该语句块
}
```

其中，`<-` 符号表示从通道中接收数据，`channel <- data` 表示将数据发送到通道中。

举一个简单的例子，比如我们有两个通道 `c1` 和 `c2`，我们要等待这两个通道中任意一个通道有数据，然后进行操作，我们可以使用通道选择器：

```go
select {
case <- c1:
    fmt.Println("c1 received data")
case <- c2:
    fmt.Println("c2 received data")
}
```

当其中一个通道有数据时，就会立即执行相应的操作。

## 实例

通道选择器是一种让你可以同时等待多个通道操作的机制。在某些场景下，同时等待多个通道操作可以帮助你将代码进行优化。

下面是一个使用通道选择器的例子，代码中定义了两个通道`ch1`和`ch2`，分别用来传递字符串信息。在`select`代码块中，使用`case`关键字分别监听了`ch1`和`ch2`的读取操作，如果其中有一个通道有可读取的信息，则会执行相应的分支。

```go
package main

import (
	"fmt"
	"time"
)

func main() {

	chan1 := make(chan string)
	chan2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 2)
		chan1 <- "chan1 msg"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		chan2 <- "chan2 msg"
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-chan1:
			fmt.Println("chan1 msg out")
		case <-chan2:
			fmt.Println("chan2 msg out")
		}

	}
}
```

在上面的例子中，我们开启了两个协程，分别向`ch1`和`ch2`中传递了字符串信息。在主协程中，使用`select`关键字监听了`ch1`和`ch2`的读取操作。由于我们在协程中使用了`time.Sleep`函数，因此可以模拟在不同时间收到消息的情况。

```go
chan2 msg out
chan1 msg out
```

在两个协程的消息中，`ch2`的消息先被收到，而`ch1`的消息稍晚一些。通过使用`select`关键字，我们可以在一个协程中同时监听多个通道，等待可读取的消息并执行相应的操作。

# 14.golang超时处理

在 Go 中，可以使用 `select` 语句和 `time.After()` 函数来实现超时处理。

`select` 语句用于监听多个通道的操作，当其中一个通道操作成功时，就会执行对应的操作。`time.After()` 函数可以创建一个定时器，它会在指定的时间后向通道发送一个值。将 `time.After()` 函数和 `select` 语句结合使用，可以在一定时间后取消某个操作或任务。

以下是一个简单的例子，它使用 `select` 语句和 `time.After()` 函数来实现超时处理：

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)

	// 启动一个协程，等待 3 秒后向通道发送一个值
	go func() {
		time.Sleep(3 * time.Second)
		c <- 1
	}()

	// 使用 select 语句监听 c 和一个定时器
	select {
	case result := <-c:
		fmt.Println("操作成功，结果为：", result)
	case <-time.After(2 * time.Second):
		fmt.Println("操作超时，取消任务")
	}
}
```

在上面的例子中，我们使用 `make()` 函数创建了一个整型通道 `c`，然后启动了一个协程，等待 3 秒后向通道发送一个值。在主协程中，使用 `select` 语句监听通道 `c` 和一个定时器，定时器的时间为 2 秒。当通道 `c` 接收到值时，就会执行第一个 `case` 分支；如果在 2 秒内没有接收到值，就会执行第二个 `case` 分支，输出“操作超时，取消任务”。

这样就可以通过 `select` 语句和 `time.After()` 函数来实现超时处理。



# 12.golang非阻塞通道

在 golang 中，非阻塞通道可以让发送和接收操作不阻塞当前协程，即使通道已满或为空。使用非阻塞通道可以避免在通道操作时产生死锁或阻塞的情况，提高程序的并发性能。

golang 提供了两种方式实现非阻塞通道：

1、使用 select 语句和 default 分支：可以在 select 语句中使用 default 分支，实现非阻塞的发送和接收操作。如果通道已满或为空，则 select 会立即执行 default 分支中的操作。

示例代码：

```go
ch := make(chan int, 1)

// 非阻塞发送操作
select {
case ch <- 1:
    fmt.Println("发送成功！")
default:
    fmt.Println("通道已满，发送失败！")
}

// 非阻塞接收操作
select {
case x := <-ch:
    fmt.Println("接收成功！", x)
default:
    fmt.Println("通道已空，接收失败！")
}
```

2、使用 `len` 函数进行判断

你可以使用 `len` 函数来判断通道中当前的元素数量，从而实现非阻塞的通道操作。

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 1)

	// 非阻塞发送操作
	if len(ch) < cap(ch) {
		ch <- 1
		fmt.Println("发送成功！")
	} else {
		fmt.Println("通道已满，发送失败！")
	}

	// 非阻塞接收操作
	if len(ch) > 0 {
		x := <-ch
		fmt.Println("接收成功！", x)
	} else {
		fmt.Println("通道已空，接收失败！")
	}

	// 延迟一段时间以观察结果
	time.Sleep(1 * time.Second)
}
```

这两种方法都可以在通道操作时避免阻塞当前协程，实现非阻塞的通道操作。你可以根据实际情况选择适合你需求的方法。

# 13.通道的遍历和关闭

遍历通道是指从通道中逐个读取元素，通道遍历可以通过 `range` 关键字来实现。

在golang中，通道可以通过close()函数进行关闭，以告诉接收方通道的数据已经全部发送完毕，不再有新的数据需要发送。

关闭通道的好处是可以避免接收方在等待数据时陷入无限阻塞状态，因为关闭通道后，接收方会立即收到一个零值，而不再等待新的数据。

在关闭通道后，通道仍然可以用于接收数据，但不能再发送数据。因此，如果有多个发送方和接收方，应该在使用通道前确保它们的正确关闭和同步。此外，对已关闭的通道执行发送操作会引发panic。

```go
package main

import "fmt"

func main() {
    ch := make(chan int)

    go func() {
        for i := 1; i <= 5; i++ {
            ch <- i
        }
        close(ch)
    }()

    for n := range ch {
        fmt.Println(n)
    }
}
```

# 14.golang Timer

在Go中，可以使用time包来创建定时器(timer)，定时器是一种机制，用于在未来的某个时间点触发操作。time包提供了两种类型的定时器：一种是单次定时器，一种是重复定时器。

## 14.1单次定时器：

单次定时器会在指定的时间后触发一次操作。在time包中，提供了两种创建单次定时器的方法：time.After()和time.NewTimer()。这两种方法的区别在于：time.After()返回一个时间通道，当定时器到期时，通道会被关闭；time.NewTimer()返回一个定时器，可以通过调用它的Stop()方法来停止定时器。

下面是一个使用time.After()方法创建单次定时器

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("start")
    <-time.After(2 * time.Second)
    fmt.Println("2 seconds later")
}
```

输出结果为：

```go
start
2 seconds later
```

### 14.1.1 timer原理

在 Go 语言中，time.Timer 是由一个 heap.Timer 结构体实现的，该结构体维护一个用于存储定时器的小根堆。每个定时器都关联了一个时间点，当时间点到达后，定时器会发送一个事件通知。在调用 time.NewTimer() 函数时，会创建一个定时器，并把它插入到堆中。当堆顶定时器的时间点到达时，它就会被移除并向通道中发送一个事件通知，然后继续等待下一个定时器的事件通知。由于堆的结构是动态的，所以可以轻松添加和删除定时器，因此 Timer 可以被用来实现一些高级的定时器逻辑，如周期性任务的重复执行。

在 Go 语言中，还有一个 time.Ticker 类型，它也可以用于定时器。不同于 time.Timer，time.Ticker 可以周期性地发送事件通知，而不是只发送一次。在创建一个 time.Ticker 实例时，会创建一个定时器并关联一个通道。每次定时器时间到达时，定时器会向通道发送一个事件通知，通道会被取走。这样，在使用 time.Ticker 的时候，可以通过循环读取通道的方式来周期性地执行一些任务。



### 14.1.2Timer 结构体

`heap.Timer` 是 Go 语言中标准库 `time` 包内定义的一个结构体类型，表示一个计时器。

`heap.Timer` 结构体定义如下：

```go
type Timer struct {
    C <-chan Time
    r runtimeTimer
}
```

其中 `C` 字段表示该计时器的通道，用于在计时结束时发送一个时间信号；`r` 字段则表示该计时器的运行时定时器，是一个私有结构体类型 `runtimeTimer` 的实例。

`heap.Timer` 结构体的方法有以下几个：

- `Reset(d Duration) bool`：重置计时器并返回计时器是否成功重置。如果计时器处于运行中，则会停止计时器并将其重新设置为给定的时长 `d`，并返回 true；如果计时器已经过期或已被停止，则返回 false。
- `Stop() bool`：停止计时器并返回计时器是否成功停止。如果计时器处于运行中，则会停止计时器并返回 true；如果计时器已经过期或已被停止，则返回 false。
- `ResetFunc(f func()) bool`：重置计时器并将计时器的超时行为更改为调用给定的回调函数 `f`。该方法返回计时器是否成功重置。如果计时器处于运行中，则会停止计时器并将其重新设置为零时长，并将计时器的超时行为更改为调用回调函数 `f`，并返回 true；如果计时器已经过期或已被停止，则返回 false。
- `StopFunc() bool`：停止计时器并返回计时器是否成功停止。如果计时器处于运行中，则会停止计时器的超时行为，并返回 true；如果计时器已经过期或已被停止，则返回 false。

`heap.Timer` 结构体实现了 `Reset`、`Stop`、`ResetFunc` 和 `StopFunc` 这四个方法，可以通过它们来控制计时器的行为和状态。

## 14.2重复定时器：

重复定时器会在指定的时间间隔内重复触发操作。在time包中，提供了两种创建重复定时器的方法：time.Tick()和time.NewTicker()。这两种方法的区别在于：time.Tick()返回一个时间通道，当定时器到期时，通道会被发送一个时间值；time.NewTicker()返回一个定时器，可以通过调用它的Stop()方法来停止定时器。

下面是一个使用time.Tick()方法创建重复定时器的例子：

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("start")
    ticker := time.Tick(1 * time.Second)
    for now := range ticker {
        fmt.Printf("%v\n", now)
    }
}
```

输出结果为：

```
start
2023-11-03 11:37:39.675084 +0800 CST m=+1.001348876
2023-11-03 11:37:40.675121 +0800 CST m=+2.001375626
2023-11-03 11:37:41.675136 +0800 CST m=+3.001379418
...
```

这个例子创建了一个每秒触发一次的定时器，并在每次触发时输出当前时间。

### 14.2.1  Ticker

在 Golang 中，Ticker 是一个定期触发的定时器，可以用于周期性地执行任务或重复的操作。Ticker 的实现基于 time.Timer，它会按照指定的时间间隔不断触发。

Ticker 与 Timer 的区别在于，Ticker 会定期地触发，而 Timer 只会触发一次。Ticker 提供了一个 C 属性，可以在需要时从通道中读取时间事件，以便执行相关的操作。当需要停止 Ticker 时，可以调用其 Stop 方法。

以下是 Ticker 的基本用法：

```go
func main() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Println("tick")
        }
    }
}
```

上述代码创建了一个每秒钟触发一次的 Ticker，当从 C 属性读取到时间事件时，会输出 "tick"。在函数结束时，使用 defer 关键字停止 Ticker。

需要注意的是，如果 Ticker 触发的操作耗时过长，可能会影响下一次的触发。因此，如果需要执行较耗时的操作，建议使用 time.Timer，而不是 Ticker。



### 14.2.2 Ticker源码分析

Golang中的`Ticker`类型是一个定时触发的计时器，它与`Timer`类型类似，都是基于时间的调度器。不同之处在于，`Ticker`会重复触发，而`Timer`只会触发一次。下面是`Ticker`的源码分析。

`Ticker`的定义：

```go
type Ticker struct {
    C <-chan time.Time // 计时器的触发通道，每次触发都会往通道里写入当前时间
    r runtimeTimer    // 用于计时器调度的runtimeTimer对象
    i int32           // 表示计时器是否处于运行中
}
```

从定义可以看出，`Ticker`内部维护了一个`runtimeTimer`对象，用于计时器调度。而计时器的触发事件，会往`C`通道中写入一个`time.Time`类型的值。

`Ticker`的创建方法：

```go
func NewTicker(d Duration) *Ticker {
    if d <= 0 {
        panic(errors.New("non-positive interval for NewTicker"))
    }
    c := make(chan time.Time, 1)
    t := &Ticker{
        C: c,
        r: runtimeTimer{
            when:   when(d),
            period: int64(d),
            f:      sendTime,
            arg:    c,
        },
    }
    startTimer(&t.r)
    return t
}
```

在创建`Ticker`时，会先创建一个长度为1的通道`c`，并且把它作为`runtimeTimer`的`arg`参数传入，以便于在计时器触发时往通道中写入当前时间。然后使用`startTimer`函数启动计时器，并返回`Ticker`对象。

`Ticker`的启动方法：

```go
func (t *Ticker) Reset(d Duration) {
    if d <= 0 {
        panic(errors.New("non-positive interval for Ticker.Reset"))
    }
    stopTimer(&t.r)
    t.r.when = when(d)
    t.r.period = int64(d)
    t.i = 0
    startTimer(&t.r)
}
```

`Reset`方法会重置计时器的间隔，并重新启动计时器。在重置前，会先调用`stopTimer`函数停止计时器，然后重新设置计时器的参数，最后调用`startTimer`函数启动计时器。

`Ticker`的停止方法：

```go
func (t *Ticker) Stop() {
    if t.i == 0 {
        return
    }
    stopTimer(&t.r)
    t.i = 0
    close(t.C)
}
```

`Stop`方法用于停止计时器的运行。它会先检查计时器是否已经停止，如果已经停止则直接返回，否则调用`stopTimer`函数停止计时器，并关闭`C`通道，表示计时器已经停止。在关闭通道前，会把通道里的值全部读完，以确保所有的触发事件都已经被消费。

# 15.golang协程同步

## 15.1什么是golang协程同步

在 Go 中，协程同步是指多个协程之间按照一定的顺序执行，从而避免了数据竞争和死锁等问题。

## 15.2协程同步方法

协程同步可以通过协程通信来实现，主要包括以下几种方式：

1. 互斥锁：使用 `sync.Mutex` 或 `sync.RWMutex` 等同步原语来控制并发访问共享资源的次序。
2. 条件变量：使用 `sync.Cond` 等同步原语来等待或唤醒某个特定条件的发生。
3. 信号量：使用 `chan struct{}` 类型的无缓冲通道或 `sync.WaitGroup` 等同步原语来等待一组操作完成。
4. 原子操作：使用 `sync/atomic` 包中提供的原子操作函数来执行一些针对共享资源的原子性操作，如 `AddInt32`、`SwapInt64` 等。

这些协程同步方式可以组合使用，根据具体的场景选择不同的同步方法，以实现协程之间的协作和协调。

# 16.golang协程通信

## 16.1什么是协程通信

Golang协程通信是指Golang中协程之间传递数据的过程。在Golang中，通常使用通道（Channel）进行协程之间的通信。通道是一种数据结构，可以在协程之间传递数据，协程通过向通道发送数据和从通道接收数据来实现通信。通道具有同步的特性，可以在多个协程之间传递数据，并且保证数据的顺序和完整性。通过使用通道进行协程之间的通信，可以避免数据竞争和死锁等问题。

## 16.2协程通信方法

Golang 协程之间的通信可以通过以下几种方式实现：

1. Channel 通道：Golang 中提供了 Channel 这种数据类型，可以用于协程之间的数据传输，从而实现协程的通信。
2. 共享内存：协程之间可以通过共享内存的方式进行通信，但需要注意同步的问题，否则会产生数据竞争等问题。
3. WaitGroup：WaitGroup 可以用于协程之间的同步，可以在一个协程中等待其他协程完成特定的任务后再执行后续操作。
4. Mutex 互斥锁：通过 Mutex 可以实现对共享资源的互斥访问，避免数据竞争问题。
5. Atomic 原子操作：可以用于对共享资源进行原子操作，保证并发时的数据正确性。

这些方式可以根据实际情况选择使用，对于不同的场景选择不同的通信方式可以提高程序的性能和效率。

# 17.golang对协程支持的包和类

Golang 提供了一些包和类来支持协程，这些包和类包括：

1. `go` 关键字：用于启动协程。
2. `sync` 包：提供了一些同步原语，如 Mutex、RWMutex、Cond 等，用于协程之间的同步和通信。
3. `channel`：提供了一种协程之间通信的机制，可以用于发送和接收数据。
4. `select` 语句：用于多路复用协程的通信，可以监听多个通道的数据。
5. `context` 包：用于控制协程的生命周期和传递上下文信息。
6. `runtime` 包：提供了一些和协程和调度器相关的接口，如 GOMAXPROCS、Goexit、Gosched 等。

这些包和类都是 Golang 中支持协程的重要组成部分，可以有效地实现协程之间的同步和通信。



# 18.golang等待和通知

在Go语言中，可以使用channel和sync包中的WaitGroup和Cond类型来实现等待和通知。

## 18.1 channel实现等待和通知

使用channel可以实现等待和通知的机制。在协程中，我们可以使用channel来等待某个事件的发生，并且在另一个协程中通过向channel发送消息来通知等待的协程。

例如，下面的例子中，我们创建了一个channel，然后在一个协程中等待消息的到来，而另一个协程通过向channel发送消息来通知等待的协程：

```go
package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int)

	go func() {
		fmt.Println("in goroutine,wait channel")
		<-ch
		fmt.Println("channel wait success")
	}()

	time.Sleep(time.Second * 2)

	fmt.Println("main goroutine，set num to channel")
	ch <- 1
	time.Sleep(time.Second)
}
```

在上面的例子中，创建了一个无缓冲的channel ch，然后在一个协程中等待消息的到来。另一个协程通过向channel发送消息来通知等待的协程。在主协程中，我们通过sleep函数暂停了2秒钟，然后发送消息，等待1秒钟后退出程序。



## 18.2 WaitGroup实现等待和通知

使用WaitGroup可以实现等待一组协程执行完毕的机制。我们可以使用WaitGroup的Add方法来添加需要等待的协程数量，然后在每个协程执行完毕后，通过Done方法减少需要等待的协程数量。最后，我们可以通过Wait方法来等待所有协程执行完毕。

例如，下面的例子中，我们创建了两个协程，然后使用WaitGroup来等待它们执行完毕：

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        fmt.Println("goroutine 1")
        time.Sleep(1 * time.Second)
    }()

    go func() {
        defer wg.Done()
        fmt.Println("goroutine 2")
        time.Sleep(2 * time.Second)
    }()

    wg.Wait()
    fmt.Println("all goroutines done")
}
```

在上面的例子中，我们定义了一个WaitGroup变量wg，并使用Add方法添加了需要等待的协程数量。在每个协程执行完毕后，我们使用Done方法减少需要等待的协程数量。最后，我们使用Wait方法等待所有协程执行完毕，并在主协程中打印出提示信息。

## 18.3 cond实现等待和通知

使用Cond可以实现更加复杂的等待和通知的机制。Cond可以让一个或多个协程等待某个事件的发生，并在另一个协程中通过广播或者单个通知来通知等待的协程。

例如，下面的例子中，我们创建了一个Cond变量cv，并使用它来实现等待和通知的机制：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	done := false

	go func() {
		time.Sleep(time.Second * 2)
		mu.Lock()
		fmt.Println("子协程处理")
		done = true
		cond.Signal()
		mu.Unlock()
	}()

	mu.Lock()
	if !done {
		cond.Wait()
	}
	fmt.Println("协程等待结束")
	mu.Unlock()

}
```

在上面的例子中，我们首先创建了一个Mutex变量mu，并使用它来保护共享的状态变量done。然后，我们通过NewCond方法创建了一个Cond变量cv，并将它绑定到mu上。

在另一个协程中，我们通过sleep函数暂停了2秒钟，然后获取了mu的锁，并将done设置为true，最后通过Signal方法通知等待的协程。

在主协程中，我们获取了mu的锁，并在一个循环中等待done变量的值变为true。在等待过程中，我们调用cv的Wait方法来释放mu的锁，并进入等待状态，等待Signal方法的通知。

最后，在等待结束后，我们打印出提示信息，并释放mu的锁。

需要注意的是，Cond的Wait方法会自动释放锁，并在收到通知后重新获取锁。因此，在使用Cond时，需要确保共享的状态变量是受到互斥锁保护的，以避免并发访问的问题。

# 19.生产者消费者模型

在 Go 中，可以使用 goroutine 和 channel 实现生产者消费者模型，代码如下：

```go
package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	for i := 0; i < 100; i++ {
		ch <- i
		fmt.Print(i)
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for i := range ch {
		fmt.Println("message from channel, no", i)
	}
}

func main() {
	ch := make(chan int, 10)
	go producer(ch)
	consumer(ch)
	time.Sleep(time.Second)

}
```

在上面的代码中，我们定义了两个函数，分别代表生产者和消费者。生产者函数 producer 接收一个只能写入的 channel out，并通过 for 循环不断向 channel 中写入数据。消费者函数 consumer 接收一个只能读取的 channel in，并通过 for-range 循环不断从 channel 中读取数据。

在主函数中，我们创建了一个 channel ch，并通过 go 关键字启动生产者协程。同时，在主函数中执行消费者函数，不断从 channel 中读取数据。最后，为了让程序正常结束，我们在主函数中增加了一个睡眠时间。

需要注意的是，当生产者写入完数据后，需要通过 close 关闭 channel，以通知消费者数据已经全部写入。同时，当消费者从 channel 中读取完数据后，会自动退出循环，因为 channel 已经被关闭了。

此外，在实际应用中，为了避免 channel 阻塞和协程泄露的问题，可能需要在写入和读取 channel 时使用 select 语句和超时机制。

# 20.多协程窗口卖票

在 Go 中，可以使用 goroutine 和 channel 实现多协程卖票程序，代码如下：

```go
package main

import (
	"fmt"
	"sync"
)

func sellTickets(wg *sync.WaitGroup, ch chan int, id int) {
	defer wg.Done()
	for {
		ticket, ok := <-ch
		if !ok {
			fmt.Printf("goroutine %d: channel is closed\n", id)
			return
		}
		fmt.Printf("goroutine %d:ticket sold successfully, %d\n", id, ticket)
	}
}

func main() {
	const ticketsNum = 100
	const numSellers = 4

	var wg sync.WaitGroup
	wg.Add(numSellers)

	ch := make(chan int, ticketsNum)
	for i := 0; i < ticketsNum; i++ {
		ch <- i
	}
	close(ch)

	for i := 1; i <= numSellers; i++ {
		go sellTickets(&wg, ch, i)
	}
	wg.Wait()
}
```

首先定义了一个 sellTickets 函数，该函数接收一个 sync.WaitGroup 指针 wg、一个只能读取的 channel ch，以及一个卖票程序的编号 id。在函数内部，我们使用 for 循环不断从 channel 中读取数据，并输出卖票信息。当 channel 被关闭时，我们通过判断 ok 的值来判断是否需要退出循环。

在主函数中，我们定义了两个常量 numTickets 和 numSellers，分别表示票的总数和卖票程序的数量。然后，我们创建了一个有缓冲的 channel ch，并在其中写入 numTickets 个整数。接着，我们通过 for 循环启动 numSellers 个协程，并调用 sellTickets 函数。最后，我们调用 sync.WaitGroup 的 Wait 方法，等待所有协程退出。

需要注意的是，为了避免 channel 阻塞和协程泄露的问题，我们需要在 channel 中设置缓冲区，并在写入 channel 时保证写入操作的原子性。同时，在协程退出时，需要调用 sync.WaitGroup 的 Done 方法，以通知主函数该协程已经退出。



# 21.golang协程同步共享内存的方法

在 Go 中，有几种方法可以实现协程同步和共享内存，包括：

1. Mutex 和 Condition Variable

使用 sync.Mutex 和 sync.Cond 可以实现协程之间的同步和共享内存。sync.Mutex 提供了锁机制，可以保证只有一个协程能够访问共享变量，从而避免竞态条件的发生。而 sync.Cond 可以用于协程之间的通信，以等待某个条件的发生。

1. Channel

Channel 是 Go 中最基本的同步机制之一，可以用于协程之间的通信和同步。Channel 有两种类型：无缓冲和有缓冲。无缓冲的 Channel 可以用于协程之间的同步，有缓冲的 Channel 可以用于协程之间的通信和同步。

1. WaitGroup

WaitGroup 是 Go 中用于等待一组协程完成的同步机制。使用 sync.WaitGroup 可以在主协程中等待所有子协程执行完成后再继续执行下面的代码。

1. Atomic

Go 还提供了 atomic 包，可以使用原子操作实现对共享变量的同步访问。atomic 包中提供了一些原子操作函数，如 AddInt32、AddUint32、CompareAndSwapInt32 等，这些函数可以保证对共享变量的操作是原子性的。

1. Select

在协程之间进行通信时，可以使用 select 语句来等待多个通信操作。select 语句会等待其中任意一个通信操作完成，然后执行对应的代码块。这种方式可以避免阻塞协程，提高程序的并发性能。

需要根据具体情况选择适合的同步机制，避免竞态条件和死锁等问题的发生。

# 21.锁

在 Go 中，常用的锁有以下几种：

1. Mutex

Mutex 是 Go 中最基本的锁类型，可以通过 Lock 和 Unlock 方法实现对共享变量的访问控制。当一个协程持有 Mutex 时，其他协程将被阻塞，直到 Mutex 释放。

1. RWMutex

RWMutex 是一种读写锁，可以支持多个协程同时读取共享变量，但只能有一个协程写入共享变量。当一个协程持有 RWMutex 的写锁时，其他协程将被阻塞，直到写锁释放。当一个协程持有 RWMutex 的读锁时，其他协程可以持有读锁，但不能持有写锁。

1. Cond

Cond 是 Go 中用于协程之间通信的锁类型。Cond 可以等待某个条件的发生，当条件满足时，通知其他协程继续执行。

1. Once

Once 是 Go 中用于执行一次性操作的锁类型。Once 可以保证一个操作只会执行一次，即使多个协程同时访问。

除了以上四种锁类型，Go 还提供了一些特殊的锁类型，如 atomic 包中的原子操作锁，以及 sync.Map 中的并发安全 Map。需要根据具体的需求选择适合的锁类型，避免竞态条件和死锁等问题的发生。

# 22.Mutex互斥锁

在 Go 中，Mutex 是一种基本的锁类型，用于实现对共享变量的访问控制。Mutex 的名称来源于“Mutual Exclusion”，表示互斥锁。

使用 Mutex 可以避免多个协程同时访问共享变量时出现竞态条件。通过 Lock 和 Unlock 方法，可以控制对共享变量的访问。当一个协程持有 Mutex 时，其他协程将被阻塞，直到 Mutex 释放。

Mutex 的定义如下：

```go
type Mutex struct {
    // 互斥锁状态
    state int32
    // 互斥锁等待队列
    sema  uint32
}
```

Mutex 包含一个 state 变量和一个 sema 变量。state 变量用于表示 Mutex 的状态，如果 state 的值为 0，则表示 Mutex 是未锁定的状态，如果 state 的值为 1，则表示 Mutex 是锁定的状态。sema 变量用于表示等待 Mutex 的协程数量。

Mutex 支持以下方法：

1. Lock

Lock 方法用于获取 Mutex 的锁定状态，如果 Mutex 已经被锁定，则会阻塞当前协程，直到 Mutex 可用。

```go
func (m *Mutex) Lock()
```

1. Unlock

Unlock 方法用于释放 Mutex 的锁定状态，如果 Mutex 没有被锁定，则会引发 panic。

```go
func (m *Mutex) Unlock()
```

Mutex 的使用示例：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var count = 0
var mu sync.Mutex

func addNum() {
	mu.Lock()
	count++
	mu.Unlock()
}

func main() {
	for i := 1; i <= 10; i++ {
		go addNum()
	}
	time.Sleep(time.Second)
	fmt.Println(count)
}
```

在上面的示例中，我们定义了一个 count 变量和一个 mutex 变量，用于实现对 count 变量的访问控制。在 increment 函数中，我们先使用 mutex.Lock() 获取 mutex 的锁定状态，然后对 count 变量进行操作，最后使用 mutex.Unlock() 释放 mutex 的锁定状态。在 main 函数中，我们启动了 10 个协程执行 increment 函数，使用 mutex 控制对 count 变量的访问，从而避免了竞态条件的发生。

# 23.读写锁RWMutex

在 Go 中，RWMutex 是一种读写锁，用于实现对共享变量的访问控制。RWMutex 支持多个协程同时读取共享变量，但只能有一个协程写入共享变量。RWMutex 的名称来源于“Read-Write Mutex”，表示读写锁。

使用 RWMutex 可以避免读操作和写操作之间的竞态条件。当一个协程持有 RWMutex 的写锁时，其他协程将被阻塞，直到写锁释放。当一个协程持有 RWMutex 的读锁时，其他协程可以持有读锁，但不能持有写锁。

RWMutex 的定义如下：

```go
type RWMutex struct {
    // 互斥锁状态
    w           Mutex
    // 等待写锁的协程数量
    writerSem   uint32
    // 等待读锁的协程数量
    readerSem   uint32
    // 写锁等待队列
    readerCount int32
    // 读锁等待队列
    readerWait  int32
}
```

RWMutex 包含一个 w 变量和一些状态变量。w 变量是一个普通的 Mutex，用于实现对 RWMutex 的访问控制。writerSem 和 readerSem 分别表示等待写锁和读锁的协程数量。readerCount 表示持有读锁的协程数量，readerWait 表示等待读锁的协程数量。

RWMutex 支持以下方法：

1. Lock

Lock 方法用于获取 RWMutex 的写锁，如果 RWMutex 已经被锁定，则会阻塞当前协程，直到 RWMutex 可用。

```go
func (rw *RWMutex) Lock()
```

1. Unlock

Unlock 方法用于释放 RWMutex 的写锁，如果 RWMutex 没有被锁定，则会引发 panic。

```
func (rw *RWMutex) Unlock()
```

1. RLock

RLock 方法用于获取 RWMutex 的读锁，如果 RWMutex 已经被锁定并且持有写锁，则会阻塞当前协程，直到 RWMutex 可用。

```go
func (rw *RWMutex) RLock()
```

1. RUnlock

RUnlock 方法用于释放 RWMutex 的读锁，如果 RWMutex 没有被锁定或者没有持有读锁，则会引发 panic。

```go
func (rw *RWMutex) RUnlock()
```

RWMutex 的使用示例：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var count = 0
var rwLock sync.RWMutex
var wg sync.WaitGroup

func read(i int) {
	defer wg.Done()
	rwLock.RLock()
	defer rwLock.RUnlock()
	fmt.Printf("read:%v goroutine:%d\n", i, count)
}

func write() {
	defer wg.Done()
	rwLock.Lock()
	defer rwLock.Unlock()
	count++
	fmt.Println("write:", count)
}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * 1000)
		go read(i)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go write()
	}
	wg.Wait()
	fmt.Println("final")
}
```

上面的示例中，count 是一个共享变量，用于记录读取和写入的次数。read 和 write 分别代表读操作和写操作的协程函数。在 main 函数中，启动 10 个读协程和 3 个写协程，并使用 WaitGroup 语句使程序保持同步。

在 read 函数中，首先使用 RLock 方法获取读锁，然后在 defer 语句中使用 RUnlock 方法释放读锁。在 write 函数中，同样使用 Lock 方法获取写锁，使用 Unlock 方法释放写锁。

使用 RWMutex 可以保证多个协程可以同时读取共享变量，但只有一个协程可以写入共享变量。这可以有效地避免读写操作之间的竞态条件，从而提高程序的并发性能。

# 24.WaitGroup

WaitGroup 是一个同步原语，用于等待一组协程的执行完成。WaitGroup 维护了一个计数器，用于记录还有多少个协程未完成，当计数器为 0 时，表示所有协程已经执行完毕。主协程可以使用 Wait 方法等待所有协程的执行完成。

WaitGroup 的常用方法包括：

- Add(delta int)：增加计数器的值，通常在启动协程之前调用，表示有多少个协程需要等待。
- Done()：减少计数器的值，通常在协程的 defer 语句中调用，表示该协程已经执行完毕。
- Wait()：阻塞主协程，直到计数器的值为 0，表示所有协程已经执行完毕。

下面是一个简单的示例，演示了如何使用 WaitGroup 等待一组协程的执行完成：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func GoroutineWG() {

	wg := sync.WaitGroup{}
	printOdd := func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
	}

	printEven := func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
	}
	wg.Add(2)
	go printOdd()
	go printEven()
	wg.Wait()
	fmt.Println("after group")
}

func main() {
	GoroutineWG()
}
```

上面的示例中，worker 函数表示要执行的协程函数，其中的 defer wg.Done() 语句表示该协程已经执行完毕。在 main 函数中，循环启动了 5 个协程，并在每个协程启动之前调用了 wg.Add(1) 方法，表示有 5 个协程需要等待。最后，调用 wg.Wait() 方法等待所有协程执行完毕。当所有协程执行完毕后，输出 "All workers done" 表示程序执行完成。

# 25.条件变量Cond

Cond 是 Go 语言中的条件变量，用于协程之间的同步。Cond 依赖于 Mutex 实现协程之间的同步，通常与 Mutex 一起使用。

Cond 可以用于实现协程之间的等待和通知。等待可以通过 Wait 方法实现，它会阻塞当前协程，并释放锁，等待其他协程通知该协程继续执行。通知可以通过 Signal 和 Broadcast 方法实现，它们都用于通知等待的协程继续执行，Signal 只通知一个协程，而 Broadcast 会通知所有等待的协程。

Cond 的常用方法包括：

- Wait()：等待条件变量，阻塞当前协程，并释放锁。
- Signal()：通知一个等待的协程继续执行，通常与 Lock 和 Unlock 方法一起使用。
- Broadcast()：通知所有等待的协程继续执行，通常与 Lock 和 Unlock 方法一起使用。

下面是一个示例，演示了如何使用 Cond 实现协程之间的等待和通知：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	ready := false

	for i := 0; i < 5; i++ {
		go func(i int) {
			mu.Lock()
			if !ready {
				cond.Wait()
			}
			fmt.Printf("worked %d is working\n", i)
			mu.Unlock()
		}(i)
	}

	time.Sleep(time.Second)
	mu.Lock()
	ready = true
	cond.Broadcast()
	mu.Unlock()

	time.Sleep(time.Second)
}
```

上面的示例中，使用 sync.Mutex 创建了一个互斥锁 mutex，并使用 sync.NewCond 方法创建了一个条件变量 cond。ready 变量用于记录是否可以继续执行。在循环中，启动了 5 个协程，并使用 cond.Wait() 方法等待条件变量。在主协程中，等待 1 秒后，设置 ready 为 true，并使用 cond.Broadcast() 方法通知所有等待的协程继续执行。最后，等待 1 秒后程序执行完成。

在协程中，当 ready 为 false 时，使用 cond.Wait() 方法等待条件变量，并阻塞当前协程。当 ready 为 true 时，协程将继续执行。注意，cond.Wait() 方法必须在互斥锁的保护下进行，否则可能出现竞态条件。在主协程中，设置 ready 为 true 和使用 cond.Broadcast() 方法也必须在互斥锁的保护下进行。

# 26.原子变量Atomic

在 Go 语言中，原子操作是一种特殊的操作，可以保证在并发情况下对共享变量的操作不会出现竞态条件。在并发编程中，使用原子操作可以提高程序的性能和安全性。

Go 语言中提供了 sync/atomic 包，用于实现原子操作。其中常用的函数有：

- AddInt32(addr *int32, delta int32) int32：原子加操作，将 delta 值加到 *addr 中，并返回新值。
- AddInt64(addr *int64, delta int64) int64：原子加操作，将 delta 值加到 *addr 中，并返回新值。
- CompareAndSwapInt32(addr *int32, old, new int32) bool：比较并交换操作，如果 *addr 的值等于 old，就将 *addr 的值设置为 new，返回 true，否则返回 false。
- CompareAndSwapInt64(addr *int64, old, new int64) bool：比较并交换操作，如果 *addr 的值等于 old，就将 *addr 的值设置为 new，返回 true，否则返回 false。
- SwapInt32(addr *int32, new int32) int32：交换操作，将 *addr 的值设置为 new，并返回旧值。
- SwapInt64(addr *int64, new int64) int64：交换操作，将 *addr 的值设置为 new，并返回旧值。

下面是一个使用原子操作实现的计数器示例：

```go
package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var count int32

	for i := 0; i < 10000; i++ {
		go func() {
			atomic.AddInt32(&count, 1)
		}()

	}
	time.Sleep(time.Second * 3)
	fmt.Println(count)
}
```

在上面的示例中，创建了一个 int32 类型的计数器 count，使用 atomic.AddInt32() 方法对其进行原子加操作。在循环中，创建了 10 个协程，每个协程会对 count 进行 1000000 次原子加操作。最后，程序输出计数器 count 的值。由于原子操作的保证，count 的值将正确地等于 10000000。

总之，使用原子操作可以确保在并发情况下对共享变量的操作不会出现竞态条件，提高程序的性能和安全性。但是需要注意，过度使用原子操作可能会导致程序变得复杂难以维护，应该根据具体情况进行选择。

# 27.死锁

在 Go 语言中，死锁是一种常见的并发问题，指的是在程序中存在循环等待资源的情况，导致程序无法继续执行下去。如果不及时解决，死锁可能会导致程序崩溃或出现不可预测的错误。

常见的死锁情况有：

1. 互斥锁（Mutex）的嵌套使用

如果一个协程获取了一个互斥锁，而在未释放该锁的情况下又去获取另一个互斥锁，那么就会形成死锁。这种情况下，两个协程都在等待对方释放锁，导致程序无法继续执行下去。

1. 预分配资源不足

如果在程序启动时预分配的资源（如通道、缓冲区等）不足，而程序中需要使用更多的资源，就可能会出现死锁。这种情况下，协程会阻塞等待资源的释放，而另一个协程却无法释放资源，导致程序无法继续执行下去。

1. 通道读写顺序不当

如果两个协程在通道上进行读写操作的顺序不当，就可能会出现死锁。例如，一个协程在等待通道的写操作，而另一个协程在等待通道的读操作，导致两个协程都在等待对方的操作，程序无法继续执行下去。

为了避免死锁问题，我们可以采用以下几种方法：

1. 避免嵌套使用互斥锁，尽量使用更高级别的同步原语，如读写锁（RWMutex）、条件变量（Cond）等。
2. 避免使用全局变量，尽量使用局部变量或通过通道传递数据。
3. 使用带缓冲区的通道，可以避免协程阻塞等待通道的读写操作。
4. 使用超时机制，在获取资源的操作中设置超时时间，如果超时就放弃获取资源，避免程序一直阻塞等待资源的释放。
5. 使用可重入函数，即在函数内部不会调用另一个可能会阻塞的函数。

总之，死锁是一种常见的并发问题，在编写并发程序时需要格外注意。采用合适的同步原语和编写合理的程序逻辑，可以有效避免死锁问题的发生。

## 27.1互斥锁嵌套产生的死锁

互斥锁（Mutex）的嵌套使用是常见的死锁情况之一。当一个协程已经持有了一个互斥锁，并且在未释放该锁的情况下又去获取另一个互斥锁时，就会形成死锁。

下面是一个示例，说明互斥锁的嵌套使用可能会导致死锁的情况：

```go
func deadFunc() {
	var mu1, mu2 sync.Mutex
	wg := sync.WaitGroup{}
	wg.Add(2)
	// 协程1 先获取mu1再获取mu2
	go func() {
		mu1.Lock()
		fmt.Println("goroutine 1 m1 Lock")
		mu2.Lock()
		fmt.Println("goroutine 1 m2 Lock")
		mu2.Unlock()
		fmt.Println("goroutine 1 m2 unlock")
		mu1.Unlock()
		fmt.Println("goroutine 1 m1 unlock")
		wg.Done()
	}()

	// 协程2 先获取mu2再获取mu1
	go func() {
		mu2.Lock()
		fmt.Println("goroutine 2 m2 lock")
		mu1.Lock()
		fmt.Println("goroutine 2 m1 lock")
		mu1.Unlock()
		fmt.Println("goroutine 2 m1 unlock")
		mu2.Unlock()
		fmt.Println("goroutine 2 m2 unlock")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("al goroutines have finished")
}
```

在上面的示例中，我们创建了两个互斥锁 `mu1` 和 `mu2`，然后启动了两个协程，分别尝试获取这两个互斥锁。其中，协程1先获取`mu1`再获取`mu2`，而协程2先获取`mu2`再获取`mu1`。

当协程1获取了`mu1`之后，它试图获取`mu2`，但是此时`mu2`已经被协程2获取，所以协程1会阻塞等待`mu2`的释放。同时，协程2也已经获取了`mu2`，并试图获取`mu1`，但是此时`mu1`已经被协程1获取，所以协程2也会阻塞等待`mu1`的释放。这样，两个协程就互相等待对方释放锁，形成了死锁。

为了避免互斥锁的嵌套使用，我们可以考虑使用更高级别的同步原语，如读写锁（RWMutex）、条件变量（Cond）等。同时，我们还可以采用其他方法来避免死锁问题，例如使用带缓冲区的通道、使用超时机制等。



使用RWMutex锁解决上面死锁问题

使用 RWMutex 锁可以解决上面互斥锁的嵌套使用导致的死锁问题。RWMutex 是一种读写锁，它可以同时支持多个协程对共享资源的读操作，但只能允许一个协程进行写操作，从而避免了读写之间的竞争问题。

下面是使用 RWMutex 锁改进后的代码：

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var rwLock sync.RWMutex
	var data int

	wg := sync.WaitGroup{}
	wg.Add(5)

	go func() {
		rwLock.Lock()
		defer rwLock.Unlock()

		data = 23
		fmt.Println("write operation completed")
		wg.Done()
	}()

	for i := 0; i < 4; i++ {
		go func() {
			rwLock.RLock()
			defer rwLock.RUnlock()

			fmt.Println("Read operation: ", data)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("all operations have finished:", data)
}
```



## 27.2 预分配资源不足产生的死锁

预分配资源不足也可能会造成死锁问题。这种情况通常发生在需要动态分配资源的场景中，例如协程池或者对象池等。如果预分配的资源数量不足以满足需求，就可能会发生死锁问题。

下面是一个简单的例子，演示了预分配资源不足时可能发生的死锁问题：

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// 预分配1个资源
	pool := make(chan int, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	// 协程1 从资源池中获取资源
	go func() {
		n := <-pool
		fmt.Println("goroutine 1 acquired resource", n)
		wg.Done()

		// 等待一段时间，模拟协程1在持有资源期间不释放
		for {
		}
	}()

	// 协程2 从资源池中获取资源
	go func() {
		n := <-pool
		fmt.Println("goroutine 2 acquired resource", n)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("All goroutines have finished")
}
```

在上面的代码中，我们创建了一个大小为2的资源池，向其中放入了一个资源，然后开启两个协程从资源池中获取资源并执行任务。由于资源池中只有一个资源，因此只有一个协程可以获取到资源，而另一个协程则会一直等待资源被释放，从而导致了死锁问题。

要解决这个问题，我们需要增加资源池的容量或者减少任务的数量，从而保证资源的供需平衡。同时，也可以使用带有超时或者取消机制的上下文来避免死锁问题的发生。

**解决死锁问题：**

要解决预分配资源不足造成的死锁问题，我们可以增加资源池的容量或者减少任务的数量。另外，使用带有超时或者取消机制的上下文也可以避免死锁问题的发生。

以下是修改后的代码：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	pool := make(chan int, 3)
	for i := 0; i < 3; i++ {
		pool <- i
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		n := <-pool
		fmt.Println("goroutine 1 acquired resource", n)
		// 模拟执行任务
		fmt.Println("goroutine 1 is working...")
		time.Sleep(time.Second)
		wg.Done()
		// 将资源归还到资源池中
		pool <- n
		fmt.Println("goroutine 1 released resource", n)
	}()
	// 协程2 从资源池中获取资源并执行任务
	go func() {
		n := <-pool
		fmt.Println("goroutine 2 acquired resource", n)
		// 模拟执行任务
		fmt.Println("goroutine 2 is working...")
		time.Sleep(time.Second)
		wg.Done()
		// 将资源归还到资源池中
		pool <- n
		fmt.Println("goroutine 2 released resource", n)
	}()

	// 协程3 从资源池中获取资源并执行任务
	go func() {
		n := <-pool
		fmt.Println("goroutine 3 acquired resource", n)
		// 模拟执行任务
		fmt.Println("goroutine 3 is working...")
		time.Sleep(time.Second)
		wg.Done()
		// 将资源归还到资源池中
		pool <- n
		fmt.Println("goroutine 3 released resource", n)
	}()

	// 等待所有协程执行完毕
	wg.Wait()
	fmt.Println("All goroutines have finished")
}
```

在上面的代码中，我们将资源池的容量增加到了3，保证了资源的供需平衡。同时，我们还将任务数量增加到了3，确保了所有资源都可以被利用。最终，我们使用`sync.WaitGroup`等待所有协程执行完毕，从而避免了死锁问题的发生。

## 27.3通道读写顺序不当产生的死锁

通道读写顺序不当也可能会导致死锁问题的发生。当一个协程在等待通道读取时，如果另一个协程在等待通道写入，就会产生死锁。

以下是一个例子：

```go
package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		// 等待接收数据
		<-ch
	}()

	// 向通道发送数据
	ch <- 1
	fmt.Println("Data sent to channel")
}
```



在上面的代码中，我们首先创建了一个通道`ch`，然后启动了一个协程，它等待从通道中接收数据。接着，我们向通道中发送了一个数据1，然后输出了一条消息。

在这个程序中，主协程向通道发送数据的操作和子协程等待从通道中接收数据的操作之间没有明确的同步机制。因此，当主协程向通道发送数据时，子协程可能还没有准备好从通道中接收数据，这时就会产生死锁问题。

要解决这个问题，我们需要保证通道的读写顺序是正确的，通常有以下几种方式：

1. 使用带缓冲的通道，将数据先缓存到通道中，等待子协程从通道中读取数据时再取出。

   ```go
   package main
   
   import "fmt"
   
   func main() {
   	ch := make(chan int, 1)
   	go func() {
   		// 等待接收数据
   		<-ch
   	}()
   
   	// 向通道发送数据
   	ch <- 1
   	fmt.Println("Data sent to channel")
   }
   ```

在上面的代码中，我们将通道改为带有缓冲的通道，并将容量设置为1，这样主协程向通道发送数据时不会阻塞。当子协程从通道中读取数据时，数据会被取出并处理。

2.使用同步原语，例如`sync.Mutex`、`sync.WaitGroup`等，保证通道的读写顺序正确。

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    ch := make(chan int)
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()  // 在读取通道数据之后调用wg.Done()方法
        // 等待接收数据
        <-ch
    }()

    // 向通道发送数据
    ch <- 1
    wg.Wait()
    fmt.Println("Data sent to channel")
}
```

在上面的代码中，我们使用了`sync.WaitGroup`来保证通道的读写顺序正确。首先，我们创建了一个`sync.WaitGroup`对象，并将计数器设置为1。然后，在子协程中，我们在读取通道数据之前调用了`wg.Done()`方法，这样主协程就可以向通道中发送数据了。最后，我们调用`wg.Wait()`方法，等待子协程处理完数据，从而保证了通道读写顺序的正确性。

# 28.协程可见性问题

在多协程程序中，可见性是一个常见的问题。可见性指的是当一个协程修改了某个共享变量的值后，其他协程能否立即看到这个修改。

在 golang 中，如果不进行特殊处理，共享变量的修改很可能对其他协程不可见。这是因为每个协程都有自己的本地缓存，修改操作都是在本地缓存中进行的，直到显式的同步操作才会将修改同步到主存中，其他协程才能看到这个修改。

例如，考虑以下代码：

```go
import (
    "fmt"
    "sync"
)

var count int

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            count++
            wg.Done()
        }()
    }
    wg.Wait()
    fmt.Println(count)
}
```

在这个程序中，我们启动了 100 个协程并发地对变量 count 进行递增操作，最终输出 count 的值。如果我们运行这个程序，可能会得到一个小于 100 的值，这是因为对 count 变量的修改操作没有被同步到主存中，造成其他协程看不到这个修改。

为了解决这个问题，我们可以使用 golang 中提供的原子操作和锁来保证可见性。例如，我们可以使用 Mutex 锁来保护 count 变量的读写操作，确保每次只有一个协程访问 count 变量：

```go
package main

import (
	"fmt"
	"sync"
)

func main() {

	count := 0
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			count++
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(count)
}
```

在这个程序中，我们使用 Mutex 锁保护了 count 变量的读写操作，确保每次只有一个协程访问 count 变量，从而避免了可见性问题。

# 29.不可变对象

在 golang 中，不可变对象指的是一旦创建后就不能被修改的对象。不可变对象在并发编程中有很大的优势，因为不可变对象可以避免一些并发编程中的问题，例如竞态条件、死锁等问题。

在 golang 中，字符串、数组、结构体等数据类型是不可变的。一旦这些数据类型被创建后，就不能被修改。例如，对于字符串类型，我们不能对字符串的某个字符进行修改：

```go
s := "hello"
s[0] = 'H' // 编译错误，字符串是不可变的
```

类似地，对于数组类型，我们也不能修改数组的元素：

```go
a := [3]int{1, 2, 3}
a[0] = 4 // 可以修改
a[0]++   // 可以修改
b := [3]int{4, 5, 6}
a = b // 可以修改，因为是复制整个数组
a[0] = 7 // 不能修改 b[0]，因为 a 和 b 是不同的数组
```

对于结构体类型，我们可以通过定义成员变量为只读来实现不可变对象：

```go
type Point struct {
    x, y int
}

func (p Point) X() int { return p.x }
func (p Point) Y() int { return p.y }

p := Point{1, 2}
p.x = 3 // 编译错误，结构体成员变量是只读的
```

在 golang 中，使用不可变对象可以避免很多并发编程中的问题，因为不可变对象可以确保数据的一致性，避免多个协程对同一数据进行修改而导致的问题。但需要注意的是，不可变对象也不是万能的，有时候我们需要修改数据，这时候就需要使用 golang 中提供的同步原语来保证数据的正确性。

# 30.享元模式

在 golang 中，享元模式是一种常用的设计模式，它的主要作用是减少内存的使用，提高程序的性能。享元模式将对象的状态分为内部状态和外部状态，内部状态可以被多个对象共享，而外部状态则由客户端传递。

在 golang 中，可以使用 sync.Pool 来实现享元模式。sync.Pool 是一个线程安全的对象池，可以用于存储和重用对象。对象池可以减少对象的创建和销毁，从而提高程序的性能。

下面是一个简单的例子，演示如何使用 sync.Pool 实现享元模式：

```go
package main

import (
    "fmt"
    "sync"
)

type User struct {
    ID   int
    Name string
}

func main() {
    pool := sync.Pool{
        New: func() interface{} {
            return new(User)
        },
    }

    // 从对象池中获取一个对象
    user := pool.Get().(*User)
    user.ID = 1
    user.Name = "Alice"
    fmt.Println(user)

    // 将对象放回对象池中
    pool.Put(user)

    // 从对象池中获取一个对象
    user = pool.Get().(*User)
    fmt.Println(user)
}
```

在上面的例子中，我们定义了一个名为 User 的结构体，它包含两个属性：ID 和 Name。我们使用 sync.Pool 创建了一个对象池，其中 New 属性用于创建新的对象。在主函数中，我们从对象池中获取一个对象，并设置它的 ID 和 Name 属性，然后将对象放回对象池中。最后，我们再次从对象池中获取一个对象，并打印它的值。

由于对象池中已经存在一个 User 对象，所以第二次获取对象时，就不需要创建新的对象了，而是直接返回对象池中的对象。这样可以减少对象的创建和销毁，从而提高程序的性能。

需要注意的是，由于 sync.Pool 是基于可复用的对象进行设计的，因此不能保证对象的状态一定是空的，也不能保证对象的状态不会被改变。因此，在使用对象池时，需要保证对象的状态是可复用的，或者在从对象池中获取对象后，先将对象的状态重置为初始状态。

# 31.sync.Pool

在 Go 中，`sync.Pool` 是一个可以存储和重用临时对象的对象池。它可以用来提高内存分配和回收的效率。

`sync.Pool` 的使用非常简单。我们只需要定义一个 `sync.Pool` 对象，并通过 `New` 方法设置一个新的对象生成函数。然后，我们就可以使用 `Get` 方法从对象池中获取一个对象。如果对象池中没有可用的对象，则会调用对象生成函数生成一个新的对象。使用完毕后，我们需要将对象通过 `Put` 方法放回对象池中，以便下一次使用。

下面是一个简单的例子，演示如何使用 `sync.Pool`：

```go
package main

import (
    "fmt"
    "sync"
)

type Object struct {
    data string
}

func NewObject(data string) *Object {
    return &Object{data}
}

func main() {
    pool := sync.Pool{
        New: func() interface{} {
            return NewObject("default")
        },
    }

    obj := pool.Get().(*Object)
    fmt.Println(obj.data)

    obj.data = "foo"
    pool.Put(obj)

    obj = pool.Get().(*Object)
    fmt.Println(obj.data)
}
```

在上面的例子中，我们首先创建了一个 `sync.Pool` 对象，并通过 `New` 方法设置一个新的对象生成函数。在 `main` 函数中，我们使用 `Get` 方法从对象池中获取一个对象，并打印出其属性值。然后，我们修改对象的属性值，并使用 `Put` 方法将对象放回对象池中。最后，我们再次使用 `Get` 方法从对象池中获取一个对象，并打印出其属性值。

在上面的例子中，我们通过 `New` 方法设置了一个新的对象生成函数，该函数返回一个默认值为 "default" 的 `Object` 对象。这意味着，如果对象池中没有可用的对象，则会使用该对象生成函数生成一个新的对象。因此，第一次使用 `Get` 方法从对象池中获取对象时，我们得到的是一个默认值为 "default" 的 `Object` 对象。

需要注意的是，`sync.Pool` 并不保证对象池中的对象数量、存活时间或对象存储位置。因此，我们应该尽量避免依赖这些特性，以免引发不必要的错误。此外，由于对象池中的对象可能会被随时回收，因此我们应该尽量避免在对象池中存储有状态的对象。如果需要存储有状态的对象，则应该使用 `sync.Mutex` 等同步机制来确保线程安全。

# 32.sync.Once

在Golang的`sync`包中，`Once`是一种同步原语，用于确保某个操作只被执行一次。

`Once`类型包含一个布尔值和一个互斥锁，布尔值表示该操作是否已经被执行。当调用`Do`方法时，如果布尔值为`false`，则会执行传入的操作函数并将布尔值设置为`true`，否则`Do`方法直接返回。

下面是一个使用`Once`的示例代码：

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			once.Do(func() {
				fmt.Println("only one")
			})
		}(i)

		fmt.Println("goroutine ", i)
		wg.Done()
	}
	wg.Wait()
	fmt.Println("finish")
}
```

在上面的代码中，我们创建了一个`Once`实例`once`，并在`Do`方法中传入了一个匿名函数，该函数只会被执行一次。然后我们启动了10个协程，每个协程都会等待`once.Do`方法的返回，并输出一条消息。

由于`once.Do`方法只会被执行一次，因此第一次运行时会输出"Only once"，后续运行时不会再次输出该消息。但是所有协程都会输出"Hello from goroutine"，因为`Do`方法只会阻塞第一个调用，后续调用会直接返回。

通过使用`Once`，我们可以确保某个操作只会被执行一次，避免重复执行造成的问题。在编写单例模式等需要确保只有一个实例的场景时，`Once`也是一个非常有用的工具。

# 33.Mutex 的状态

`Mutex`在Golang中有三种状态：

1. 未加锁状态：此时`Mutex`的状态为0，表示未被任何协程锁定，此时可以被任何协程加锁。
2. 已加锁状态：此时`Mutex`的状态为1，表示已被某个协程锁定，其他协程尝试加锁时将会被阻塞。
3. 已唤醒状态：此时`Mutex`的状态为-1，表示此前某个协程曾经持有锁并解锁了它，但当前没有协程加锁。此时，如果有协程尝试加锁，它将立即获取锁并将`Mutex`状态重置为0。

需要注意的是，`Mutex`状态的变化只由加锁和解锁操作触发，其他操作不会改变`Mutex`的状态。此外，在正常情况下，我们应该避免直接操作`Mutex`状态，而是应该使用`Mutex`的加锁和解锁方法来控制其状态的变化。

假设有一个共享资源`counter`，我们使用`Mutex`来保护它，下面是一个简单的例子：

```go
var (
    mutex   sync.Mutex
    counter int
)

func increment() {
    mutex.Lock()
    defer mutex.Unlock()
    counter++
}

func decrement() {
    mutex.Lock()
    defer mutex.Unlock()
    counter--
}
```

上述代码中，我们定义了一个`Mutex`变量`mutex`和一个`counter`变量，然后定义了两个函数`increment`和`decrement`，它们分别对`counter`变量进行加1和减1的操作。在这两个函数中，我们使用了`mutex.Lock()`方法和`mutex.Unlock()`方法来控制对`counter`变量的访问，这样就保证了`counter`变量的并发安全。

在初始状态下，`mutex`的状态为0，即未被任何协程锁定。当某个协程调用`increment()`或`decrement()`方法时，它会首先执行`mutex.Lock()`方法，将`mutex`状态从0改变为1，表示当前协程已经锁定了`mutex`。如果此时另一个协程尝试调用`increment()`或`decrement()`方法，它将会被阻塞，直到当前协程执行完`defer mutex.Unlock()`方法，将`mutex`状态改变为0，释放锁资源。

在`increment()`和`decrement()`方法执行过程中，由于它们都使用了`defer`关键字，所以无论函数是否执行成功，都会执行`mutex.Unlock()`方法，将`mutex`状态改变为0，释放锁资源。这样可以保证即使在函数执行过程中发生了异常，也不会导致`mutex`状态一直处于锁定状态，从而导致死锁的发生。

需要注意的是，在使用`Mutex`进行并发控制时，我们应该尽可能减小锁的粒度，避免对整个程序加锁，否则会导致并发性能下降。另外，对于一些耗时的操作，应该尽量将其放在锁的外部，避免长时间占用锁资源，导致其他协程长时间等待。

# 34.Cond的Broadcast和Signal区别

在Golang并发编程中，条件变量（`Cond`）的`Broadcast()`和`Signal()`方法都用于唤醒等待在条件变量上的 goroutine，但它们之间有一些区别。

1. `Broadcast()`: `Broadcast()`方法用于唤醒所有等待在条件变量上的 goroutine。当某个条件满足时，调用`Broadcast()`方法会同时唤醒所有等待的 goroutine，让它们有机会竞争获取资源或执行相应的操作。
2. `Signal()`: `Signal()`方法用于唤醒一个等待在条件变量上的 goroutine。当某个条件满足时，调用`Signal()`方法会选择其中一个等待的 goroutine唤醒，然后让它获取资源或执行相应的操作。需要注意的是，`Signal()`方法并不保证唤醒的是哪一个 goroutine，所以在使用时需要注意是否满足特定的调度顺序要求。

在实际应用中，应根据具体的需求来选择使用`Broadcast()`还是`Signal()`方法。如果多个 goroutine之间没有特定的执行顺序要求，且满足条件时需要同时唤醒所有等待的 goroutine，那么可以使用`Broadcast()`方法。而如果希望有一定的调度顺序，只唤醒一个 goroutine执行相应操作，那么可以使用`Signal()`方法。

> 需要注意的是，在使用条件变量时，通常需要结合互斥锁（`Mutex`）来保护共享资源的访问，以防止竞态条件的发生。

实例

下面是一个简单的示例，演示了在 Golang 中使用条件变量（`sync.Cond`）的`Broadcast()`和`Signal()`方法的区别：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)
	done := false

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			for !done {
				fmt.Printf("goroutine %d is waitting\n", i)
				cond.Wait()
				fmt.Printf("goroutine %d was awaked\n", i)
			}
			mu.Unlock()
		}(i)
	}

	time.Sleep(time.Second * 2)
	mu.Lock()
	done = true
	fmt.Println("broadcast")
	cond.Broadcast()
	mu.Unlock()

	wg.Wait()
	fmt.Println("finish")
}
```

在上述示例中，我们创建了5个 goroutine，每个 goroutine都会等待条件变量的信号。首先，我们让主 goroutine 等待2秒，模拟一些耗时操作。然后，主 goroutine 获取互斥锁，设置 `done` 标志为 `true`，并调用 `Broadcast()` 方法来唤醒所有等待的 goroutine。每个等待的 goroutine 在被唤醒后打印一条消息。最后，我们使用 `Wait()` 方法等待所有 goroutine 完成。

运行上述示例，你将看到每个等待的 goroutine 被唤醒并打印消息，然后退出。这是因为 `Broadcast()` 方法同时唤醒了所有等待的 goroutine。

如果我们将 `Broadcast()` 方法替换为 `Signal()` 方法，你会发现只有一个等待的 goroutine 被唤醒并打印消息，其他 goroutine 仍然处于等待状态。这是因为 `Signal()` 方法只唤醒一个等待的 goroutine。

# 35.GMP

GMP是Golang中调度器的三个重要组成部分之一，它是Goroutine的执行器。GMP是一个由操作系统内核线程（M）和用户级线程（G）构成的调度器。

- G：Goroutine（协程）的缩写，代表一个轻量级线程，它由调度器（Scheduler）管理。
- M：Machine的缩写，代表调度器线程（或者说是操作系统内核线程），每个M都与一个系统线程相关联，M负责在G和系统线程之间进行调度。
- P：Processor的缩写，代表逻辑处理器，P是调度器（Scheduler）的一部分，它与M一一对应，负责管理一组Goroutine的运行。

Goroutine是基于协作式调度的，调度器会在Goroutine进行IO、函数调用等阻塞操作时，主动调度其他处于就绪状态的Goroutine执行。在Goroutine被调度到某个M上运行时，该M只会负责这个Goroutine的执行，而不是像普通线程那样运行整个进程的代码。因此，在Golang中，多个Goroutine可以共享一个M，从而减少了线程切换的开销，提高了并发性能。

GMP模型是Golang调度器的核心之一，它为Golang提供了高效、轻量级的并发模型，使Golang在处理高并发任务时表现优异。

**GMP调度**

GMP是Go语言实现协程的底层机制，调度器（scheduler）是GMP系统的核心。GMP使用的是基于m:n调度器的协程实现方式，即m个协程在n个系统线程上运行。调度器在运行时，负责协程的创建、销毁和调度等工作。

GMP调度器的主要任务是将协程绑定到线程上，然后调度线程执行协程。为了实现高效的调度，GMP调度器使用了很多优化手段。其中比较重要的一项是抢占式调度。

在抢占式调度模式下，调度器会在当前协程执行完毕或者遇到I/O阻塞等情况时，强制剥夺其运行权并调度其他协程执行。这种方式可以避免某些长时间运行的协程阻塞整个调度系统，提高系统的稳定性和响应速度。

除了抢占式调度，GMP调度器还使用了许多其他优化手段，如自适应线程管理、局部性优化等，这些都可以让GMP调度器更加高效地管理协程和线程的调度。

# 36.sync.Map

`sync.Map` 是 Go 语言中用于并发读写 Map 的线程安全的容器，是在 `map` 的基础上实现的。相比于 `map`，它提供了一些额外的功能，如自动扩容和并发安全，而且使用起来也比较方便。

`sync.Map` 的特点：

- 线程安全：`sync.Map` 中的操作是并发安全的，多个 goroutine 可以同时读写 `sync.Map` 中的数据。
- 自动扩容：`sync.Map` 中的数据会自动扩容，不需要手动处理扩容问题。
- 零值：`sync.Map` 的零值是一个空的 Map，可以直接使用。

`sync.Map` 的常用方法：

- `func (m *Map) Load(key interface{}) (value interface{}, ok bool)`：获取指定 key 对应的 value，如果不存在，则返回 false。
- `func (m *Map) Store(key, value interface{})`：存储指定的 key 和 value。
- `func (m *Map) Delete(key interface{})`：删除指定的 key 和对应的 value。
- `func (m *Map) Range(f func(key, value interface{}) bool)`：遍历 Map，对每个 key-value 对执行指定的操作，如果 f 返回 false，则停止遍历。

`sync.Map` 的实现原理：

`sync.Map` 内部使用了一种叫做 "Sharded Map" 的技术来实现并发安全，它将 Map 分成多个 shard（分片），每个 shard 用一个独立的锁来保证并发安全。当访问 `sync.Map` 中的数据时，会根据 key 的哈希值将其映射到一个 shard 上，然后在该 shard 上进行操作。由于每个 shard 有自己独立的锁，因此多个 goroutine 可以同时访问不同的 shard，从而实现了并发安全。

总之，`sync.Map` 是一个线程安全、高效的 Map 实现，特别适用于多 goroutine 环境下的读写操作。但需要注意的是，`sync.Map` 并不保证操作的顺序，如果需要保证顺序，需要加锁。
