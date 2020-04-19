package main
// based on https://github.com/golang/net/blob/master/icmp/example_test.go
import (
    "time"
    "os"
    "log"
    "fmt"
    "os/signal"
    "syscall"
    "net"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Invalid target to ping!")
        os.Exit(1)
    }
    // the target to ping
    target := os.Args[1]

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    var sucess float64  = 0
    var packets float64  = 0
    var loss float64  = 0.0

    ping := func(address string) {
        destination, duration, err := Ping(address)
        if err != nil {
            log.Printf("Ping %s (%s): %s\n", address, destination, err)
            packets++
            return
        }
        packets++
        sucess++
        log.Printf("Ping %s (%s): %s\n", address, destination, duration)
    }

    for {
        ping(target)
        time.Sleep(1 * time.Second)
        // goroutine to execute the following concurrently to display packet loss
        go func() {
            <-sigs
            fmt.Println()
            loss = (1 - (sucess / packets)) * 100
            fmt.Print("Packet loss: ")
            fmt.Printf("%.2f", loss)
            fmt.Println("%")
            os.Exit(1)
        }()
    }
}

func Ping(address string) (*net.IPAddr, time.Duration, error) {
    const ProtocolICMP = 1
    // Start listening for icmp replies
    listener, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
        return nil, 0, err
    }
    defer listener.Close()

    // Gets the ipv4 address of the target
    destination, err := net.ResolveIPAddr("ip4", address)
    if err != nil {
        return nil, 0, err
    }

    // create a new ICMP message
    message := icmp.Message {
        Type: ipv4.ICMPTypeEcho, Code: 0,
        Body: &icmp.Echo{
            ID: os.Getpid() & 0xffff, Seq: 1,
            Data: []byte("HELLO!"),
        },
    }
    bytes, err := message.Marshal(nil)
    if err != nil {
        return destination, 0, err
    }

    // send to target
    start := time.Now()
    response, err := listener.WriteTo(bytes, destination)
    if err != nil {
        return destination, 0, err
    }

    reply := make([]byte, 1500)
    // wait for a reply for 3 seconds
    err = listener.SetReadDeadline(time.Now().Add(3 * time.Second))
    if err != nil {
        return destination, 0, err
    }

    response, peer, err := listener.ReadFrom(reply)
    if err != nil {
        return destination, 0, err
    }
    duration := time.Since(start)

    // parse reply
    parsed, err := icmp.ParseMessage(ProtocolICMP, reply[:response])
    if err != nil {
        return destination, 0, err
    }

    // returns results if response is good
    if (parsed.Type == ipv4.ICMPTypeEchoReply) {
      return destination, duration, nil
    } else {
      return destination, 0, fmt.Errorf("got %+v from %v; want echo reply", parsed, peer)
    }
}
