package main

import (
  "net"
  "fmt"
  "flag"
  "log"
  "io"
  "time"
)

func main(){
  laddrStr := flag.String("laddr", "", "local address")
  raddrStr := flag.String("raddr", "", "remote address")
  rdelayInmsec := flag.Int("rdelay", 5000, "recive delay in msec")
  sdelayInmsec := flag.Int("sdelay", 5000, "send delay in msec")
    
  flag.Parse()
  fmt.Printf("laddr=%s raddr=%s\n", *laddrStr, *raddrStr)
  
  laddr, _ := net.ResolveTCPAddr("tcp", *laddrStr)
  raddr, _ := net.ResolveTCPAddr("tcp", *raddrStr)
    
  listener, err := net.ListenTCP("tcp", laddr)
  if err != nil {
    log.Fatal(err)  
  }

  rcon, err := net.DialTCP("tcp", nil, raddr)
  if err != nil {
    log.Fatal(err)  
  }
  
  lcon, err := listener.AcceptTCP()
  if err != nil {
    log.Fatal(err)
  }
  
  go transfer(lcon, rcon, *sdelayInmsec)
  transfer(rcon, lcon, *rdelayInmsec)      
}

func transfer(src *net.TCPConn, dst *net.TCPConn, delayInmsec int) {
  for {
    readBuffer := make([] byte, 4096)
    readByte, err := src.Read(readBuffer)
    if err != nil && err != io.EOF {
      log.Fatal(err)
    }
    if err == io.EOF {
      src.CloseRead()
      fmt.Printf("close %p", src)
    } else {
      fmt.Printf("src=[%p] dst=[%p] read len=[%d] data=[%s]\n", src, dst, readByte, readBuffer[:32])      
    }
    
    go func(){
      time.Sleep(time.Duration(delayInmsec) * time.Millisecond)
      if err == io.EOF {
        dst.CloseWrite()
        fmt.Printf("close %p", src)        
      } else {
        _, err := dst.Write(readBuffer[:readByte])
        if err != nil {
          log.Fatal(err)
        }      
        fmt.Printf("src=[%p] dst=[%p] write len=[%d] data=[%s]\n", src, dst, readByte, readBuffer[:32]) 
      }
    }()
    
    if err == io.EOF {
      break
    }    
  }   
}