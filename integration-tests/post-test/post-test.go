package main
 
import (
        "flag"
        "fmt"
        "io"
        "os"
)
 
func main() {
        var filename string
        flag.StringVar(&filename, "name", "/tmp/postdat2", "Filename given by the caller")
        var data string
        flag.StringVar(&data, "data", "hello=world", "Data given by the caller")
        
        flag.Parse()
        f, err := os.Create(filename)
        if err != nil {
        fmt.Println(err)
        }
        n, err := io.WriteString(f, data)
        if err != nil {
        fmt.Println(n, err)
        }
        f.Close()
}
