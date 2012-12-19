mychunker
=======

Parallel chunker and dumper for MySQL tables. 
This is just a proof of concept. 
It dumps a specific table (-schema . -table) by chunks, chunking by PK only for now. 
It can also run N parallel threads to do the chunking (defaults to 4) though this will only be of benefit if the I/O subsystem has more than 1 disk. 


TODO
=======

If this is to continue, I would like to: 
- use transaction if Innodb (not possible with -threads > 1)
- rewrite the part that removes the last comma, there must be a way to
  use strings.Join() to do that
- improve the chunking, at least improve getChunkData so we can use another index if the table has no PK
