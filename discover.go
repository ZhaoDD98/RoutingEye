package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "time"
    "flag"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
    "io/ioutil"
)

func init() {
    file := "./" +"mesg"+ ".txt" // Give the logger file path
    logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
    if err != nil {
        panic(err)
    }
    log.SetOutput(logFile)
    log.SetPrefix("[ICMP_Trace]") // Logger header
    log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
    return
}

/*func trace(mid, seqnum, ip) {
    var dst net.IPAddr
    dst.IP = ip
    if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
        log.Fatal(err)
    }
    wm := icmp.Message{
        Type: ipv4.ICMPTypeEcho, Code: 0,
        Body: &icmp.Echo{
            ID:   mid & 0xffff,
            Data: []byte("HELLO-R-U-THERE"),
        },
    }
    wm.Body.(*icmp.Echo).Seq = seqnum
    wb, err := wm.Marshal(nil)
    if err != nil {
        log.Fatal(err)
    }
    if err := p.SetTTL(i); err != nil {
        log.Fatal(err)
    }

    // In the real world usually there are several
    // multiple traffic-engineered paths for each hop.
    // You may need to probe a few times to each hop.
    begin := time.Now()
    if _, err := p.WriteTo(wb, nil, &dst); err != nil {
        log.Fatal(err)
    }
    if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
        log.Fatal(err)
    }
    n, cm, peer, err := p.ReadFrom(rb)
    if err != nil {
        if err, ok := err.(net.Error); ok && err.Timeout() {
            fmt.Printf("%v\t*\n", i)
            continue
        }
        log.Fatal(err)
    }
    rm, err := icmp.ParseMessage(1, rb[:n])
    if err != nil {
        log.Fatal(err)
    }

}*/

func main() {
    // Get input fields
    var host string
    flag.StringVar(&host, "host", "", "The host which IP packet route to")
    var ipGiven string
    flag.StringVar(&ipGiven, "ip", "", "String type to give the destination IP")
    var filePath string
    flag.StringVar(&filePath, "file", "", "File type to give the destination IPs")
    flag.Parse()

    // If don't know the usage
    Unargs := flag.Args()
    for _, value := range Unargs {
        if( value == "?"){ // Use ? to show the usage of 
            flag.Usage()
        }
        log.Fatal("The above is the parameter configuration instructions")
    }
    var dst net.IPAddr

    // If host field is not null.
    if host != "" {
        ips, err := net.LookupIP(host) // Domain2IP
        if err != nil {
            log.Fatal(err)
        }
        for _, ip:= range ips{
            if ip.To4() != nil { // Find one ip can be obtained
                dst.IP = ip
                break
            }
        }
        if dst.IP == nil {
            log.Fatal("no A record found")
        }
    }

    if ipGiven != "" {
        dst.IP = net.ParseIP(string(ipGiven)) // Transfer string type to Addr type
    }
    if filePath != "" {
        data, err := ioutil.ReadFile(filePath)
        if err != nil {
            log.Fatal("File reading error", err)
        }
        dst.IP = net.ParseIP(string(data)) // Transfer string type to Addr type
    }
    c, err := net.ListenPacket("ip4:1", "0.0.0.0") // ICMP for IPv4
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
    p1 := ipv4.NewPacketConn(c) // Add new connection p1
    p2 := ipv4.NewPacketConn(c) // Add new connection p2

    if err := p1.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
        log.Fatal(err)
    }
    if err := p2.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
        log.Fatal(err)
    }
    wm := icmp.Message{
        Type: ipv4.ICMPTypeEcho, Code: 0,
        Body: &icmp.Echo{
            ID:   os.Getpid() & 0xffff, // Use process id to be the ID of package 
            Data: []byte("HELLO-R-U-THERE"),
        },
    }
    rb := make([]byte, 1500) // Response buffer zone
    for i := 1; i <= 20; i++ { // Up to 64 hops
        wm.Body.(*icmp.Echo).Seq = i // Set Seq number.
        wb, err := wm.Marshal(nil) // Read Message from write byte stream
        if err != nil {
            log.Fatal(err)
        }
        begin := time.Now() // Time send package
        var n int
        var cm *ipv4.ControlMessage
        var peer net.Addr
        var errs error
        if i % 2 == 0 {
            if err := p1.SetTTL(i); err != nil { //Set TTL(Time To Live) number. Tips : this is used to control the hops to jump
                log.Fatal(err)
            }
            if _, err := p1.WriteTo(wb, nil, &dst); err != nil { // Send
                log.Fatal(err)
            }
            if err := p1.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
                log.Fatal(err)
            }
            n, cm, peer, errs = p1.ReadFrom(rb) // Read response message from response byte stream
        } else {
            if err := p2.SetTTL(i); err != nil { //Set TTL(Time To Live) number. Tips : this is used to control the hops to jump
                log.Fatal(err)
            }
            if _, err := p2.WriteTo(wb, nil, &dst); err != nil { // Send
                log.Fatal(err)
            }
            if err := p2.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil { // Read response message from response byte stream
                log.Fatal(err)
            }
            n, cm, peer, errs = p2.ReadFrom(rb) // Read response message from response byte stream
        }
        if errs != nil {
            if err, ok := errs.(net.Error); ok && err.Timeout() { // Package time out
                continue
            }
            log.Fatal(err)
        }
        rm, err := icmp.ParseMessage(1, rb[:n]) // Response controll Message
        if err != nil {
            log.Fatal(err, cm)
        }
        rtt := time.Since(begin) // Rtt
        pkid := rb[32:34] // Package id, also the Process id
        pkseq := rb[34:36] // Package Seq number
        switch rm.Type {
        case ipv4.ICMPTypeTimeExceeded: // If message type is timeexceeded 
            log.Printf("Got TimeExceeded Mes From %-15v,  Id = %-4v , Seq = %-4v, time = %-4v \n", peer, pkid, pkseq, rtt)
        case ipv4.ICMPTypeEchoReply: // Package arrive the desination
            pkid = rb[4:6] // Revalue the pkid, because the information location is different
            pkseq = rb[6:8] // Revalue the Seq number, because the information location is different
            log.Printf("Got EchoReply    Mes From %-15v,  Id = %-4v , Seq = %-4v, time = %-4v \n", peer, pkid, pkseq, rtt)
            return
        default:
            log.Printf("unknown ICMP message: %15-v \n", rm)
        }
    }
    fmt.Println("Tadk Finish")
}

