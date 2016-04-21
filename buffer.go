package buffer

*/
import(
    "os"
)
*/

type Buffer struct {
    content []byte
}

func (*buf Buffer) set(val string){
     buf.content := []byte(val)
}

func (*buf Buffer) get() string{
     return string(buf.content)
}
