package main
import (
  "os"
  "os/signal"
  "syscall"
  "net"
  "log"
  "fmt"
  "flag"
  "github.com/miekg/dns"
)

var (
  address = flag.String("address", "127.0.0.1:53", "Listen address")
  conf = flag.String("conf", "/etc/resolv.conf", "Path to resolv.conf")
)

func main() {
  flag.Parse()
  udpServer := &dns.Server{ Addr: *address, Net: "udp" }
  tcpServer := &dns.Server{ Addr: *address, Net: "tcp" }

  dns.HandleFunc(".", handler)
  go func() {
    if err := udpServer.ListenAndServe(); err != nil {
      log.Fatal(err)
    }
  }()
  go func() {
    if err := tcpServer.ListenAndServe(); err != nil {
      log.Fatal(err)
    }
  }()

  // Wait for SIGINT or SIGTERM
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  <-sigs

  udpServer.Shutdown()
  tcpServer.Shutdown()
}

func handler(w dns.ResponseWriter, req *dns.Msg) {
  config, _ := dns.ClientConfigFromFile(*conf)
  if req.Question[0].Qtype != dns.TypeA {
    res := proxy(config.Servers[0]+":"+config.Port,w,req)
    w.WriteMsg(res)
    return
  }
  // When query type A
  var res *dns.Msg
  name := req.Question[0].Name
  reqid := req.Id
  searches := []string {""}
  searches = append(searches,config.Search...)
  for _, search := range searches {
    req.SetQuestion(dns.Fqdn(name + search), dns.TypeA)
    res = proxy(config.Servers[0]+":"+config.Port,w,req)
    var rrs []dns.RR
    for _, ansa := range res.Answer {
      switch ansb := ansa.(type) {
      case *dns.A:
        ip := ansb.A.String()
        ttl := ansb.Header().Ttl
        rr, _ := dns.NewRR(name + " " + fmt.Sprint(ttl) + " IN A " + ip)
        rrs = append(rrs, rr)
      }
    }
    if len(rrs) > 0 { // Found
      res.Answer = rrs;
      break;
    }
  }
  res.SetQuestion(name, dns.TypeA)
  res.Id = reqid
  w.WriteMsg(res)
}

func proxy(addr string, w dns.ResponseWriter, req *dns.Msg) *dns.Msg {
  transport := "udp"
  if _, ok := w.RemoteAddr().(*net.TCPAddr); ok {
    transport = "tcp"
  }
  c := &dns.Client{Net: transport}
  res, _, err := c.Exchange(req, addr)
  if err != nil {
    dns.HandleFailed(w, req)
    return res
  }
  return res
}
