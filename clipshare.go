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

var (
	pidfile_path = "/var/run/clipshare"
	pidfile_name = "clipshare.pid"
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
		set_clip_text(text)
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

func clipshare_client(hosts string) {
	// Connect to peer
	c, err := net.Dial("tcp", hosts":8002")
	if err != nil {
		fmt.Println(err)
		return
	}
	// send the clip_text
	text := get_clip_text()
	fmt.Println("Sending", text)
	err = gob.NewEncoder(c).Encode(text)
	if err != nil {
		fmt.Println(err)
  	}
	// close connection
	c.Close()
}

func set_clip_text(text string) {
	// Set clip text here
}

func get_clip_text() string{
	// get clipboard data from xsel
	cmd := exec.Command("xsel", "-b", "-o")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Still output error")
	}
	fmt.Printf("%s", out)
	return out
}

func init(){
	fmt.Printf("Starting Clipshare...")

	// Create directory for clipshare and set path
	err_dir := Mkdir(pidfile_path)
	pidfile.SetPidfilePath(pidfile_path)
	err_pid := pidfile.Write()
	if( err_pid != nil ){
		fmt.Printf("Error encountered %v", errr)
	}

	// Listen for external connections parallely
	runtime.GOMAXPROCS(2)

	// Start listening to receive data from other peers
	go clipshare_server()

	// Listen for local messages
	clipshare_local()
}

func clipshare_local () {
	// set up unix socket here

	// receive messgaes and get hosts
	msg := ">> hosts"
	hosts := get_hosts(msg)
	
	// send text to hosts
	clipshare_client(text)
}

func process_running() bool {
	// check if pidfile exists
	_, err := os.Stat(pidfile_path)
	if (os.IsNotExist(err)) {
		return true
	}else {
		return false
	}
}

func main() {
        // check if process is already running
	if ( !process_running()) {
		init()
	}
	
	// if args passed, send to open socket
}
