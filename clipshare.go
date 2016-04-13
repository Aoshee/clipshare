package main

import (
       "fmt"
       "net"
       "runtime"
	//"strings"	
	"encoding/gob"
	//"bufio"
       "os"
       "os/exec"
	//"log"
	"time"
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
	// Connect to peers
	// for each host
	// host = host + ":8002"
	host := "127.0.0.1:8002"
	c, err := net.Dial("tcp", host)
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
	return string(out)
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

	// Start listening to receive data from other peers
	go clipshare_server()

	// Listen for local messages
	clipshare_local()
}

func handleReq(conn *net.UnixConn) {
	var buf [1024]byte
        n, err := conn.Read(buf[:])
	if err != nil {
		panic(err)
	}
        fmt.Printf("%s\n", string(buf[:n]));
}

func clipshare_local () {
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
		go handleReq(conn)
        }
	 

	/* based on the message connect to
	 * client <- if '>>'
	 *        send hosts to client
	 *        client sends the clipboard contents to the hosts
	 * server <- if '<<'
	 *        server displays/copies the message from clipboard
	 * receive messgaes and get hosts
         */
	// msg := ">> hosts"
	//hosts := get_hosts(msg)
	
	// send text to hosts
	//clipshare_client(text)
}

func process_running() bool {
	// check if pidfile exists
	name := pidfile_path+ "/" +pidfile_name
	fmt.Printf("PIDFILE %v", name)
	_, err := os.Stat(name)
	if (os.IsNotExist(err)) {
		fmt.Printf("PIDFILE does not exist")
		return false
	}else {
		fmt.Printf("PIDFILE does exists")
		return true
	}
}

func connect_local_sock() {
	raddr := net.UnixAddr{"/tmp/clipshare_local", "unix"}
	conn, err := net.DialUnix("unix", nil, &raddr)
	if err != nil {
		panic(err)
	}   
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}   
	conn.Close()
}

func main() {
        // check if process is already running
	if ( !process_running()) {
		time.Sleep (5)
		clipshare_init()
	}
	
	// if args passed, send to open socket
	connect_local_sock()
}
