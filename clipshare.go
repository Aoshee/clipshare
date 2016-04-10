package main

import (
       "fmt"
       "net"
       "runtime"
       "strings"
       "encoding/gob"
//       "bufio"
       "os"
       "os/exec"
//       "log"
       "github.com/facebookgo/pidfile"
	
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

func get_clip_text() {
     cmd := exec.Command("xsel", "-b", "-o")
     //cmd.Stdin = strings.NewReader("some input")
     //var out bytes.Buffer
     //cmd.Stdout = &out
     out, err := cmd.Output()
     if err != nil {
     	fmt.Printf("Still output error")
     }
     fmt.Printf("%s", out)
}

func main() {
     	//pid_exists, _ := pidfile.Read()
	_, err := os.Stat("/home/chandrika/clipshare/clipshare.pid")
	if (os.IsNotExist(err)) {
  	    fmt.Printf("Here")
		pidfile.SetPidfilePath("/home/chandrika/clipshare/clipshare.pid")
		errr := pidfile.Write()
		if( errr != nil ){
			fmt.Printf("Error encountered %v", errr)
	        }
	    fmt.Printf("Starting clipshare...\n")
	    runtime.GOMAXPROCS(2)
	    // Start listening to receive data from other peers
	    go clipshare_server()
	}
	
	var in string
	for {
	      // Accept the text to be copied
	      // reader := bufio.NewReader(os.Stdin)
   	      // text, _ := reader.ReadString('\n')*/
	      fmt.Scanf("%s", &in)
	      if ( strings.Compare(in, ">>") == 0) {
	      	 get_clip_text()
	      }
	      // Add which client to connect to or a register process
	      // Client connects to the server to send text
	      
	      //go clipshare_client(text)
	}
}
