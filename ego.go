package main
import (
"fmt"
"io"
"github.com/dustin/go-humanize"
)
//line directory.ego:1
 func DirectoryPage(w io.Writer, dir *MailDirectory) error  {
//line directory.ego:2
_, _ = fmt.Fprintf(w, "\n\n")
//line directory.ego:4
_, _ = fmt.Fprintf(w, "\n\n<html>\n<head>\n<link href=\"/assets/stampede.css\" rel=\"stylesheet\">\n</head>\n<body>\n<nav class=\"sidebar\">\n")
//line directory.ego:11
 FolderTree(w, Root) 
//line directory.ego:12
_, _ = fmt.Fprintf(w, "\n</nav>\n")
//line directory.ego:13
 s := dir.Stats() 
//line directory.ego:14
_, _ = fmt.Fprintf(w, "\n")
//line directory.ego:14
 t := s.Total() 
//line directory.ego:15
_, _ = fmt.Fprintf(w, "\n<h1>")
//line directory.ego:15
_, _ = fmt.Fprintf(w, "%v",  dir.Name )
//line directory.ego:15
_, _ = fmt.Fprintf(w, " (")
//line directory.ego:15
_, _ = fmt.Fprintf(w, "%v",  len(dir.Folders) )
//line directory.ego:15
_, _ = fmt.Fprintf(w, ")</h1>\n<div class=\"content\">\n<table>\n\t<th><td>Count</td><td>Size</td>\n\t</th><tr><td>Unread</td>\n\t\t<td>")
//line directory.ego:20
_, _ = fmt.Fprintf(w, "%v",  humanize.Comma(s.Unread.Count) )
//line directory.ego:20
_, _ = fmt.Fprintf(w, " [")
//line directory.ego:20
_, _ = fmt.Fprintf(w, "%v",  s.Unread.Count * 100 / t.Count )
//line directory.ego:20
_, _ = fmt.Fprintf(w, "%%]</td>\n\t\t<td>")
//line directory.ego:21
_, _ = fmt.Fprintf(w, "%v",  humanize.Bytes(uint64(s.Unread.Size)) )
//line directory.ego:21
_, _ = fmt.Fprintf(w, " [")
//line directory.ego:21
_, _ = fmt.Fprintf(w, "%v",  s.Unread.Size * 100 / t.Size )
//line directory.ego:21
_, _ = fmt.Fprintf(w, "%%]</dd>\n\t</tr><tr><td>Deleted</td>\n\t\t<td>")
//line directory.ego:23
_, _ = fmt.Fprintf(w, "%v",  humanize.Comma(s.Deleted.Count) )
//line directory.ego:23
_, _ = fmt.Fprintf(w, " [")
//line directory.ego:23
_, _ = fmt.Fprintf(w, "%v",  s.Deleted.Count * 100 / t.Count )
//line directory.ego:23
_, _ = fmt.Fprintf(w, "%%]</td>\n\t\t<td>")
//line directory.ego:24
_, _ = fmt.Fprintf(w, "%v",  humanize.Bytes(uint64(s.Deleted.Size)) )
//line directory.ego:24
_, _ = fmt.Fprintf(w, " [")
//line directory.ego:24
_, _ = fmt.Fprintf(w, "%v",  s.Deleted.Size * 100 / t.Size )
//line directory.ego:24
_, _ = fmt.Fprintf(w, "%%]</dd>\n\t</tr><tr>\n\t<td>Messages</td>\n\t\t<td>")
//line directory.ego:27
_, _ = fmt.Fprintf(w, "%v",  humanize.Comma(t.Count) )
//line directory.ego:27
_, _ = fmt.Fprintf(w, "</td>\n\t\t<td>")
//line directory.ego:28
_, _ = fmt.Fprintf(w, "%v",  humanize.Bytes(uint64(t.Size)) )
//line directory.ego:28
_, _ = fmt.Fprintf(w, "</dd>\n\t</tr>\n</table>\n</div>\n<div class=\"content\">\n<table>\n\t<th><td>Folder</td><td>Count</td><td>Size</td></th>\n\t")
//line directory.ego:35
 for n, f := range dir.Folders { t := f.Stats().Total() 
//line directory.ego:36
_, _ = fmt.Fprintf(w, "\n\t<tr><td>")
//line directory.ego:36
_, _ = fmt.Fprintf(w, "%v",  n )
//line directory.ego:36
_, _ = fmt.Fprintf(w, "</td>\n\t\t<td>")
//line directory.ego:37
_, _ = fmt.Fprintf(w, "%v",  humanize.Comma(t.Count) )
//line directory.ego:37
_, _ = fmt.Fprintf(w, "</td>\n\t\t<td>")
//line directory.ego:38
_, _ = fmt.Fprintf(w, "%v",  humanize.Bytes(uint64(t.Size)) )
//line directory.ego:38
_, _ = fmt.Fprintf(w, "</dd>\n\t</tr>\n\t")
//line directory.ego:40
 } 
//line directory.ego:41
_, _ = fmt.Fprintf(w, "\n</table>\n</div>\n</body>\n</html>\n")
return nil
}
//line folder.ego:1
 func FolderPage(w io.Writer, folder *MailFolder) error  {
//line folder.ego:2
_, _ = fmt.Fprintf(w, "\n\n<html>\n<head>\n<link href=\"/assets/stampede.css\" rel=\"stylesheet\">\n</head>\n<body>\n<nav class=\"sidebar\">\n")
//line folder.ego:9
 FolderTree(w, Root) 
//line folder.ego:10
_, _ = fmt.Fprintf(w, "\n</nav>\n<div class=\"content\">\n<ul class=\"msgs \">\n\t")
//line folder.ego:13
 for _, m := range folder.Messages { 
//line folder.ego:14
_, _ = fmt.Fprintf(w, "\n\t\t<li class=\"")
//line folder.ego:14
_, _ = fmt.Fprintf(w, "%v",  m.cClass() )
//line folder.ego:14
_, _ = fmt.Fprintf(w, "\"><a href=\"")
//line folder.ego:14
_, _ = fmt.Fprintf(w, "%v",  m.UrlPath() )
//line folder.ego:14
_, _ = fmt.Fprintf(w, "\">\n\t\t\t<span>")
//line folder.ego:15
_, _ = fmt.Fprintf(w, "%v",  m.hSent() )
//line folder.ego:15
_, _ = fmt.Fprintf(w, "</span>\n\t\t\t<span>")
//line folder.ego:16
_, _ = fmt.Fprintf(w, "%v",  m.hSubject() )
//line folder.ego:16
_, _ = fmt.Fprintf(w, "</span>\n\t\t\t<span>")
//line folder.ego:17
_, _ = fmt.Fprintf(w, "%v",  m.hSender() )
//line folder.ego:17
_, _ = fmt.Fprintf(w, "</span>\n\t\t</a></li>\n    ")
//line folder.ego:19
 } 
//line folder.ego:20
_, _ = fmt.Fprintf(w, "\n</ul>\n</div>\n</body>\n</html>\n")
return nil
}
//line folders.ego:1
 func FolderTree(w io.Writer, dir *MailDirectory) error  {
//line folders.ego:2
_, _ = fmt.Fprintf(w, "\n<ul>\n\t")
//line folders.ego:3
 for _, d := range dir.DirectoryList() { 
//line folders.ego:4
_, _ = fmt.Fprintf(w, "\n\t\t<li><a href=\"")
//line folders.ego:4
_, _ = fmt.Fprintf(w, "%v",  d.UrlPath() )
//line folders.ego:4
_, _ = fmt.Fprintf(w, "\">")
//line folders.ego:4
_, _ = fmt.Fprintf(w, "%v",  d.Label() )
//line folders.ego:4
_, _ = fmt.Fprintf(w, "</a>")
//line folders.ego:4
 FolderTree(w,d) 
//line folders.ego:4
_, _ = fmt.Fprintf(w, "</li>\n    ")
//line folders.ego:5
 } 
//line folders.ego:6
_, _ = fmt.Fprintf(w, "\n\t")
//line folders.ego:6
 for _, f := range dir.Folders { 
//line folders.ego:7
_, _ = fmt.Fprintf(w, "\n\t\t<li class=\"")
//line folders.ego:7
_, _ = fmt.Fprintf(w, "%v",  f.cClass() )
//line folders.ego:7
_, _ = fmt.Fprintf(w, "\"><a href=\"")
//line folders.ego:7
_, _ = fmt.Fprintf(w, "%v",  f.UrlPath() )
//line folders.ego:7
_, _ = fmt.Fprintf(w, "\">")
//line folders.ego:7
_, _ = fmt.Fprintf(w, "%v",  f.Label() )
//line folders.ego:7
_, _ = fmt.Fprintf(w, "</a></li>\n    ")
//line folders.ego:8
 } 
//line folders.ego:9
_, _ = fmt.Fprintf(w, "\n</ul>\n")
return nil
}
//line message.ego:1
 func MessagePage(w io.Writer, msg *MailMessage, body io.Reader) error  {
//line message.ego:2
_, _ = fmt.Fprintf(w, "\n<html>\n<head>\n<link href=\"/assets/stampede.css\" rel=\"stylesheet\">\n</head>\n<body>\n<nav class=\"sidebar\">\n")
//line message.ego:8
 FolderTree(w, Root) 
//line message.ego:9
_, _ = fmt.Fprintf(w, "\n</nav>\n<div class=\"content\">\n<div>\n\t<span>")
//line message.ego:12
_, _ = fmt.Fprintf(w, "%v",  msg.hSent() )
//line message.ego:12
_, _ = fmt.Fprintf(w, "</span>\n\t<span>")
//line message.ego:13
_, _ = fmt.Fprintf(w, "%v",  msg.hSubject() )
//line message.ego:13
_, _ = fmt.Fprintf(w, "</span>\n\t<span>")
//line message.ego:14
_, _ = fmt.Fprintf(w, "%v",  msg.hSender() )
//line message.ego:14
_, _ = fmt.Fprintf(w, "</span>\n</div>\n<div><pre>")
//line message.ego:16
 io.Copy(EscapeContent(w),body) 
//line message.ego:16
_, _ = fmt.Fprintf(w, "</pre></div>\n</div>\n</body>\n</html>\n")
return nil
}
