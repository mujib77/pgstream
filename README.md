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