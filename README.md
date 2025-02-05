# files209

files209 is an infinite file store

ext4 has a file limit on linux. files209 aims to be an infinite files
store on top of ext4 using a new archive.

## Production Installation Instructions

1.  Launch a **Ubuntu 24.04** server and ssh into it.
1.  Install with the command **sudo snap install files209**
1.  Generate ssl keys by running **sudo files209.genssl**
1.  Make production ready by running **sudo files209.prod mpr**
1.  Restart the files209 with **sudo snap restart files209.f2store**
1.  Run **sudo files209.prod r** to get your key string. Needed in your program to connect to your files209 server.
1.  You would also need the server's IP address for your program

1.  The programs' default port is 31822.


## Operating Instructions

This project ships with a golang API "github.com/saenuma/files209" and a CLI

The CLI is bundled with the snapcraft program and can be called with **files209.cli**