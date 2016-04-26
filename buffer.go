package clipshare

/*
import(
    "os"
)
*/

// todo: change to queue

type Buffer struct {
    content string
}

func (buf *Buffer) set(val string){
     buf.content = val
}

func (buf *Buffer) get() string{
     return buf.content
}
