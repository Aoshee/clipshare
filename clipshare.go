package main

import (
	"fmt"
        "net"
        "runtime"
	"strings"	
	"encoding/gob"
        "os"
        "os/exec"
	"time"
        "github.com/facebookgo/pidfile"
)

var (
	pidfile_path = "/var/run/clipshare"
	pidfile_name = "clipshare.pid"
)

func handleConnection(conn net.Conn, clip chan string) {
	// receive the message
	var text string
	err := gob.NewDecoder(conn).Decode(&text)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Received", text)

		// ClipboardUpdate here 
		clip <- text
		close(clip)
		fmt.Println("Updated Clipboard")
	}
	conn.Close()
}

func clipshare_server(clip chan string) {
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
	    go handleConnection(conn, clip)
        }
}

func clipshare_client(hosts []string) {
	// Connect to host
	var host string
	for key := range(hosts) {
	    host = hosts[key] + ":8002"
	    fmt.Printf("Connecting to host %s\n", host)
	    c, err := net.Dial("tcp", host)
	    if err != nil {
		fmt.Println(err)
		return
	    }
	    // send the clip_text
	    text := get_clip_text()
	    fmt.Println("Sending", text)
	    fmt.Println("Calling host")
	    err = gob.NewEncoder(c).Encode(text)
	    if err != nil {
	       	fmt.Println(err)
  	    }
	    // close connection
	    c.Close()
	}
}

func set_clip_text(text string) {
	// Set clip text here
	cmd := exec.Command("xsel", "-b", "-i")
	cmd_stdin, err := cmd.StdinPipe()
	if err!= nil {
	   	panic(err)
	}
	_, err = cmd_stdin.Write([]byte(text))
	if err!= nil {
	   	panic(err)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Still output error")
	}
}

func get_clip_text() string{
	// get clipboard data from xsel
	cmd := exec.Command("xsel", "-b", "-o")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Still output error")
	}
	fmt.Printf("in get_clip_text %s", out)
	return string(out)
}

func handleReq(conn *net.UnixConn, clip chan string) {
	var buf [1024]byte
        n, err := conn.Read(buf[:])
	if err != nil {
		panic(err)
	}
	in := string(buf[:n])
	// Check if input is to set or get
	input := string(in)
	input = strings.Trim(input, " ")
	if(strings.Contains(input, "get")) {
	    fmt.Println("Get clip text")
	    text := <-clip
	    fmt.Printf("Got text %s", text)
	} else if (strings.Contains(input, "set")) {
	    // todo: accept host here
	    fmt.Println("Set clip text")
	    hosts := strings.Fields(input)
	    fmt.Printf("sending text to hosts %v", hosts[1:])
	    clipshare_client(hosts[1:])
	}
}

func clipshare_local (clip chan string) {
	// set up unix socket here
	l, err := net.ListenUnix("unix",  &net.UnixAddr{"/tmp/clipshare_local", "unix"})
	if err != nil {
		panic(err)
	}   
	defer os.Remove("/tmp/clipshare_local")
	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			panic(err)
		}
		go handleReq(conn, clip)
        }
}

func connect_local_sock(args []string) {
	raddr := net.UnixAddr{"/tmp/clipshare_local", "unix"}
	conn, err := net.DialUnix("unix", nil, &raddr)
	if err != nil {
		panic(err)
	}
	var string_args string 
	for key := range args {
		string_args = string_args + " " + args[key]
	}
	_, err = conn.Write([]byte(string_args))

	if err != nil {
		panic(err)
	}   
	conn.Close()
}

func clipshare_init(){
	fmt.Printf("Starting Clipshare...")

	// Create directory for clipshare and set path
	err := os.Mkdir(pidfile_path, 777)
	if (err != nil) {
		panic(err)
	}
	pidfile.SetPidfilePath(pidfile_path+ "/" + pidfile_name)
	err = pidfile.Write()
	if( err != nil ){
		panic(err)
	}

	// Listen for external connections parallely
	runtime.GOMAXPROCS(2)

	clip := make(chan string)
	
	// Start listening to receive data from other peers
	go clipshare_server(clip)

	// Listen for local messages
	clipshare_local(clip)
}


func process_running() bool {
	// check if pidfile exists
	name := pidfile_path+ "/" +pidfile_name
	_, err := os.Stat(name)
	if (os.IsNotExist(err)) {
		fmt.Println("PIDFILE does not exist")
		return false
	}else {
		fmt.Println("PIDFILE does exists")
		return true
	}
}

func main() {
        // check if process is already running
	if ( !process_running()) {
		time.Sleep (5)
		clipshare_init()
	}else {
		// if args passed, send to open socket
		if len(os.Args) > 1 {
			args := os.Args[1:]
			connect_local_sock(args)
		}else {
		      fmt.Printf("Usage: clipshare get| set host(s)")
		}
	}
}
