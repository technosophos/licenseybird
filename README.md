# licenseybird: Add license front matter to source files.

This program adds license information to the front matter of files
specified on the command line. The license material is commented based
on the rules of the programming language (identified by extension or, in
certain cases, file name).

## Usage

Add a header to each of four files:

```
$ licenseybird foo.go bar.py Makefile Dockerfile
```
