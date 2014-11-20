package main
import (
"fmt"
"io"
)
//line directory.ego:1
 func Directory(w io.Writer, dir *MailDirectory) error  {
//line directory.ego:2
_, _ = fmt.Fprintf(w, "\n\n<html>\n<nav>\n\t")
//line directory.ego:5
 for n := range dir.Directories { 
//line directory.ego:6
_, _ = fmt.Fprintf(w, "\n        <div>")
//line directory.ego:6
_, _ = fmt.Fprintf(w, "%v",  n)
//line directory.ego:6
_, _ = fmt.Fprintf(w, "</div>\n    ")
//line directory.ego:7
 } 
//line directory.ego:8
_, _ = fmt.Fprintf(w, "\n\t")
//line directory.ego:8
 for n := range dir.Folders { 
//line directory.ego:9
_, _ = fmt.Fprintf(w, "\n        <div>")
//line directory.ego:9
_, _ = fmt.Fprintf(w, "%v",  n)
//line directory.ego:9
_, _ = fmt.Fprintf(w, "</div>\n    ")
//line directory.ego:10
 } 
//line directory.ego:11
_, _ = fmt.Fprintf(w, "\n</nav>\n</html>\n")
return nil
}
