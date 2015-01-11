package main
import (
"fmt"
"io"
)
//line directory.ego:1
 func DirectoryPage(w io.Writer, dir *MailDirectory) error  {
//line directory.ego:2
_, _ = fmt.Fprintf(w, "\n\n<html>\n<nav>\n\t<ul>\n\t\t")
//line directory.ego:6
 for n, d := range dir.Directories { 
//line directory.ego:7
_, _ = fmt.Fprintf(w, "\n\t\t\t<li><a href=\"")
//line directory.ego:7
_, _ = fmt.Fprintf(w, "%v",  d.UrlPath() )
//line directory.ego:7
_, _ = fmt.Fprintf(w, "\">")
//line directory.ego:7
_, _ = fmt.Fprintf(w, "%v",  n )
//line directory.ego:7
_, _ = fmt.Fprintf(w, "</a></li>\n\t    ")
//line directory.ego:8
 } 
//line directory.ego:9
_, _ = fmt.Fprintf(w, "\n\t\t")
//line directory.ego:9
 for n, f := range dir.Folders { 
//line directory.ego:10
_, _ = fmt.Fprintf(w, "\n\t\t\t<li><a href=\"")
//line directory.ego:10
_, _ = fmt.Fprintf(w, "%v",  f.UrlPath() )
//line directory.ego:10
_, _ = fmt.Fprintf(w, "\">")
//line directory.ego:10
_, _ = fmt.Fprintf(w, "%v",  n )
//line directory.ego:10
_, _ = fmt.Fprintf(w, "</a></li>\n\t    ")
//line directory.ego:11
 } 
//line directory.ego:12
_, _ = fmt.Fprintf(w, "\n\t</ul>\n</nav>\n</html>\n")
return nil
}
//line folder.ego:1
 func FolderPage(w io.Writer, folder *MailFolder) error  {
//line folder.ego:2
_, _ = fmt.Fprintf(w, "\n\n<html>\n<nav>\n\t<ul>\n\t\t")
//line folder.ego:6
 for _, m := range folder.Messages { 
//line folder.ego:7
_, _ = fmt.Fprintf(w, "\n\t\t\t<li><a href=\"")
//line folder.ego:7
_, _ = fmt.Fprintf(w, "%v",  m.UrlPath() )
//line folder.ego:7
_, _ = fmt.Fprintf(w, "\">")
//line folder.ego:7
_, _ = fmt.Fprintf(w, "%v",  m.String() )
//line folder.ego:7
_, _ = fmt.Fprintf(w, "</a></li>\n\t    ")
//line folder.ego:8
 } 
//line folder.ego:9
_, _ = fmt.Fprintf(w, "\n\t</ul>\n</nav>\n</html>\n")
return nil
}
