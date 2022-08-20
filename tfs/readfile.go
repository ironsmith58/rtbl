package tfs

import (
    "io"
    "os"
    "bufio"
)


func ReadFile (filepath string) ([]string,error){

    var result []string

    file, err := os.Open(filepath)
    if err != nil{
        return nil, err
    }
    defer file.Close()


    // Reading from a file using reader.
    reader := bufio.NewReader(file)
    var line string
    for {
        line, err = reader.ReadString('\n')
        if ( line == "" || ( err != nil && err != io.EOF ) ) {
            break
        }
        // As the line contains newline "\n" character at the end, we could remove it.
        line = line[:len(line)-1]

        result = append(result, line)
    }
    return result,nil
}
