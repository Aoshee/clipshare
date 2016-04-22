package buffer

/*
import(
    "os"
)
*/

type Buffer struct {
    content bytes.Buffer
}

func (*buf Buffer) set(val string){
     buf.content.Write([]byte(val))
}

func (*buf Buffer) get() string{
     return string(buf.content)
}
