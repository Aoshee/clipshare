package main

import (
       "fmt"
       "net"
       "runtime"
       "encoding/gob"
       "bufio"
       "os"
)


func handleConnection(conn net.Conn) {
  // receive the message
  var text string
  err := gob.NewDecoder(conn).Decode(&text)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Received", text)
    // ClipboarUpdate here 
  }
  conn.Close()
}

func clipshare_server() {
     	fmt.Printf("Listening\n")
	ln, err := net.Listen("tcp",":8002")
	if err != nil {
	// handle error
	}
	for {
	    conn, err := ln.Accept()
	    if err != nil {
	    	// handle error
	    }
	    go handleConnection(conn)
        }
}

func clipshare_client(text string) {
     // change this to accept peer to send text to
     c, err := net.Dial("tcp", "127.0.0.1:8002")
     if err != nil {
     	fmt.Println(err)
     return
     }
     // send the message
     fmt.Println("Sending", text)
     err = gob.NewEncoder(c).Encode(text)
     if err != nil {
     	fmt.Println(err)
  	}
     c.Close()
}

func main() {
     	fmt.Printf("Starting clipshare...\n")
	runtime.GOMAXPROCS(2)
	// Start listening to receive data from other peers
	go clipshare_server()
	for {
	      // Accept the text to be copied
	       reader := bufio.NewReader(os.Stdin)
   	       text, _ := reader.ReadString('\n')
	      // Add which client to connect to or a register process?

	      // Client connects to the server to send text
	      go clipshare_client(text)
	}
}
