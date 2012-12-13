mychunker
=======

Parallel chunker and dumper for MySQL tables. 
This is just a proof of concept. 

TODO
=======

If this is to continue, I would like to: 
- use transaction if Innodb
- rewrite the part that removes the last comma, there must be a way to
  use strings.Join() to do that
- improve the chunking, at least improve getChunkData so we can use another index if the table has no PK
