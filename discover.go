package main

import (
	"log"
	//"fmt"
    "net"
	"os"
	"sync"
    "time"
    "flag"
    "golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"./funcs"
)

var wg sync.WaitGroup

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



func trace(pkid int, seqnum int, ip string, c net.PacketConn) {
    var dst net.IPAddr
	dst.IP = net.ParseIP(ip)
	p := ipv4.NewPacketConn(c)
    if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
        log.Fatal(err)
    }
    wm := icmp.Message{
        Type: ipv4.ICMPTypeEcho, Code: 0,
        Body: &icmp.Echo{
            ID:   pkid & 0xffff,
            Data: []byte("HELLO-R-U-THERE"),
        },
    }
    rb := make([]byte, 1000) // Response buffer zone
    for i := 1; i <= 35; i++ { // Up to 64 hops
        wm.Body.(*icmp.Echo).Seq = i + seqnum // Set Seq number.
        wb, err := wm.Marshal(nil) // Read Message from write byte stream
        if err != nil {
            log.Fatal(err)
        }
        begin := time.Now() // Time send package
		if err := p.SetTTL(i); err != nil { //Set TTL(Time To Live) number. Tips : this is used to control the hops to jump
			log.Fatal(err)
		}
		if _, err := p.WriteTo(wb, nil, &dst); err != nil { // Send
			log.Fatal(err)
        }
        
		//log.Printf("Write one Message %v, %v", i + seqnum, pkid)
		if err := p.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			log.Fatal(err)
        }
        time.Sleep(time.Duration(10) * time.Microsecond)
		n, cm, peer, errs := p.ReadFrom(rb) // Read response message from response byte stream
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
			//ch <- "Got TimeExceeded Mes From " + peer.String() + ",  Id = " + string(pkid) + " , Seq =" + string(pkseq)
        case ipv4.ICMPTypeEchoReply: // Package arrive the desination
            pkid = rb[4:6] // Revalue the pkid, because the information location is different
            pkseq = rb[6:8] // Revalue the Seq number, because the information location is different
			log.Printf("Got EchoReply    Mes From %-15v,  Id = %-4v , Seq = %-4v, time = %-4v \n", peer, pkid, pkseq, rtt)
			//ch <- "Got EchoReply    Mes From " + peer.String() + ",  Id = " + string(pkid) + " , Seq =" + string(pkseq) 
            return
        default:
            log.Printf("unknown ICMP message: %v \n", rm)
        }
        
    }

}

func main() {
    // Get input fields
    var host string
    flag.StringVar(&host, "host", "", "The host which IP packet route to")
    var ipGiven string
    flag.StringVar(&ipGiven, "ip", "", "String type to give the destination IP")
	var ipsGiven string
	flag.StringVar(&ipsGiven, "ips", "", "String type to give the destination IPs")

    flag.Parse()

    // If don't know the usage
    Unargs := flag.Args()
    for _, value := range Unargs {
        if( value == "?"){ // Use ? to show the usage of 
            flag.Usage()
        }
        log.Fatal("The above is the parameter configuration instructions")
    }
    var dst string
    
    // If host field is not null.
    if host != "" {
        ips, err := net.LookupIP(host) // Domain2IP
        if err != nil {
            log.Fatal(err)
        }
        for _, ip:= range ips{
            if ip.To4() != nil { // Find one ip can be obtained
                dst = string(ip)
                break
            }
        }
        if dst == "" {
            log.Fatal("no A record found")
        }
	}
	c, err := net.ListenPacket("ip4:1", "0.0.0.0") // ICMP for IPv4
    if err != nil {
        log.Fatal(err)
    }
	defer c.Close()
    if ipGiven != "" {
		dst = ipGiven
		trace(1, 1, dst, c)
		//println(<-ch1)
    }
	if ipsGiven != "" {
		dst = ipsGiven
		Seg2Ip := funcs.Format(dst)
		wg.Add(len(Seg2Ip))
		//ch := make(chan string)
		for i, ip:= range Seg2Ip{
			if(ip != ""){
				go trace(i , i, ip, c)
			}			
		}
		wg.Wait()
    }
	
}

