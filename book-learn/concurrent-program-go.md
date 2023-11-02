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
