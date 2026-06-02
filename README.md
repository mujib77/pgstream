![License](https://img.shields.io/badge/license-MIT-blue)
![CI](https://github.com/mujib77/pgstream/actions/workflows/ci.yml/badge.svg)
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
get deleted until pgstream has read them.

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

## Architecture

```mermaid
flowchart TD
    A[Postgres Database] --> B[WAL - Write Ahead Log]
    B --> C[Replication Slot + pgoutput]
    C --> D[Connector - connects to WAL stream]
    D --> E[Decoder - decodes raw bytes into events]
    E --> F[Handler - outputs clean JSON events]
```

## Project Structure

```
pgstream/
├── main.go
├── config/
│   └── config.go
├── internal/
│   ├── connector/
│   │   └── connector.go
│   ├── decoder/
│   │   └── decoder.go
│   ├── handler/
│   │   └── handler.go
│   └── models/
│       └── event.go
```

## Setup

Prerequisites:
- PostgreSQL with logical replication enabled
- Go 1.21 or higher

**1. Enable logical replication**

In `postgresql.conf`:

```
wal_level = logical
```

Restart Postgres after changing this.

**2. Clone the repo**

```bash
git clone https://github.com/mujib77/pgstream
cd pgstream
```

**3. Create a .env file**

```
DATABASE_URL=postgres://user:password@localhost:5432/dbname?replication=database
SLOT_NAME=pgstream_slot
PUBLICATION_NAME=pgstream_pub
```

**4. Create publication in your database**

```sql
CREATE PUBLICATION pgstream_pub FOR ALL TABLES;
```

**5. Install dependencies**

```bash
go mod tidy
```

**6. Run**

```bash
go run main.go
```

## Example Output

```
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
```

## Tech Stack

- Go
- pglogrepl - Postgres logical replication protocol
- pgx - Postgres driver for Go
- PostgreSQL logical decoding with pgoutput plugin

# test 
- surprise to see the test section
