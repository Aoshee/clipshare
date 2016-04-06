package main

import (
       "fmt"
       "net"
       "encoding/gob"
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
     	fmt.Printf("In server\n")
	ln, err := net.Listen("tcp",":8001")
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
     c, err := net.Dial("tcp", "127.0.0.1:8001")
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
	// Start the server so that others can send data to this node
	go clipshare_server()
	var text string
	for {
	      // Accept the text to be copied
	      fmt.Scanf("%s", &text)

	      // Add which client to connect to or a register process?

	      // Client connects to the server to send text
	      go clipshare_client(text)
	}
}