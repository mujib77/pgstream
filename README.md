# pgstream

A real-time PostgreSQL WAL reader built in Go. 
Captures every INSERT, UPDATE, and DELETE from 
your Postgres database as it happens.

## What is pgstream?

Every time you insert, update, or delete a row in 
PostgreSQL, Postgres writes that change to a file 
called the WAL (Write Ahead Log) before touching 
the actual data. This is how Postgres makes sure 
nothing gets lost if the server crashes.

pgstream taps into that log and streams every change 
out in real time. You can see exactly what changed, 
in which table, and what the old and new values were.

This is the same mechanism that powers production 
CDC (Change Data Capture) tools like Debezium and 
pglogical used by companies running Postgres at scale.

## How it works

**Replication Slot**

pgstream creates a replication slot in your Postgres 
database. Think of it as a bookmark in the WAL stream. 
Postgres tracks your position and makes sure no changes 
get deleted until pgstream has read them. This means 
you can stop and restart pgstream without missing 
any changes.

**pgoutput**

WAL records are stored in raw binary. pgoutput is a 
built-in Postgres plugin that decodes that binary into 
readable messages like INSERT, UPDATE, DELETE along 
with the actual row data.

**LSN (Log Sequence Number)**

Every record in the WAL has a unique position called 
an LSN. pgstream tracks the LSN of every event so you 
always know exactly where in the stream each change 
came from.

**Event flow**

Postgres writes a change to WAL. pgoutput decodes it. 
pgstream reads it from the replication slot and outputs 
it as a clean JSON event with the table name, operation 
type, and row data.

## Architecture

```mermaid
flowchart TD
    A[Postgres Database] --> B[WAL - Write Ahead Log]
    B --> C[Replication Slot + pgoutput]
    C --> D[Connector - connects to WAL stream]
    D --> E[Decoder - decodes raw bytes into events]
    E --> F[Handler - outputs clean JSON events]
```

### Project Structure

```
pgstream/
├── main.go                         - entry point, wires everything together
├── config/
│   └── config.go                   - loads config from environment variables
├── internal/
│   ├── connector/
│   │   └── connector.go            - connects to Postgres replication stream
│   ├── decoder/
│   │   └── decoder.go              - decodes pgoutput messages into events
│   ├── handler/
│   │   └── handler.go              - handles and prints each event
│   └── models/
│       └── event.go                - WalEvent struct definition
```

## Setup

**Prerequisites**
- PostgreSQL with logical replication enabled
- Go 1.21 or higher

**Enable logical replication in postgresql.conf**
wal_level = logical
Restart Postgres after changing this.

**Clone the repo**
git clone https://github.com/mujib77/pgstream
cd pgstream

**Create a .env file**
DATABASE_URL=postgres://user:password@localhost:5432/dbname?replication=database
SLOT_NAME=pgstream_slot
PUBLICATION_NAME=pgstream_pub

**Create publication in your database**
```sql
CREATE PUBLICATION pgstream_pub FOR ALL TABLES;
```

**Install dependencies**
go mod tidy

**Run**
go run main.go

## Example Output

connecting to postgres...
connected!
replication slot created: pgstream_slot
starting replication...
listening for changes...
[INSERT] table=users lsn=0/16C752F8
data={
"email": "wal@gmail.com",
"id": "1",
"name": "Mujib"
}
[UPDATE] table=users lsn=0/16C75410
old={
"name": "Mujib"
}
new={
"name": "Bum"
}
[DELETE] table=users lsn=0/16C754A8
data={
"email": "wal@gmail.com",
"id": "1",
"name": "Bum"
}

## Tech Stack

- Go
- pglogrepl - Postgres logical replication protocol
- pgx - Postgres driver for Go
- PostgreSQL logical decoding with pgoutput plugin