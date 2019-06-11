# preallocate
[![GoDoc](https://godoc.org/git.sr.ht/~tslocum/preallocate?status.svg)](https://godoc.org/git.sr.ht/~tslocum/preallocate)
[![builds.sr.ht status](https://builds.sr.ht/~tslocum/preallocate.svg)](https://builds.sr.ht/~tslocum/preallocate?)

File preallocation library

## Features

- Allocates files efficiently (via syscall) on the following platforms:
  - [Linux](http://man7.org/linux/man-pages/man2/fallocate.2.html)
  - [Windows](https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/content/ntifs/nf-ntifs-ntsetinformationfile)
- Falls back to writing null bytes

## Documentation

Docs are hosted on [godoc.org](https://godoc.org/git.sr.ht/~tslocum/preallocate).

## Support

Issues and suggestions are hosted on [todo.sr.ht](https://todo.sr.ht/~tslocum/preallocate).
