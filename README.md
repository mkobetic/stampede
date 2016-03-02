*** Early WIP ***

Stampede is a simple personal web-mail service that fetches email from one or more accounts and stores it in an mbox style mail archive (similar to Thunderbird, etc).

The basic model is straightforward, a mail archive is a nested structure of MailDirectories containing MailFolders which represent individual mbox files. MailFolder then contains individual MailMessages. MailMessage knows where in the mbox file it is and how to pull its headers and contents out of it.

Opening a MailFolder or MailDirectory recursively scans the underlying file system structure creating the corresponding model elements along the way. Opening a MailFolder does only a fast scan to find message boundaries and pull out basic headers that are needed right away (subject, sender, date, ...), any further message parsing will happen lazily.

TODO:
* [x] scan an existing mail archive, create a structure directories, folders and messages
* [x] present the archive contents in very basic UI
* [ ] parse message contents and improve the presentation
* [ ] add concept of Inbox and basic configuration
* [ ] fetch mail (POP)
* [ ] send mail (SMTP)
* [ ] UI for writing messages
* [ ] UI for more complex messages (attachment handling ...)
* [ ] rule based sorting
