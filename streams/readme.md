# Redis Streams: Full Command Reference with Use Case

## Use Case: Order Processing System

We are building a distributed **Order Processing System** using Redis Streams. Services like **payment**, **inventory**, and **shipping** consume new orders asynchronously using Redis Streams and **Consumer Groups**.

---

## 1. `XADD` – Add Entry to a Stream

Adds an entry to the stream.

```bash
XADD orders * customer_id 123 item "book" quantity 2
```

### Options:

- `*` – Auto-generate ID
- `<ID>` – Custom ID
- `MAXLEN [~] <count>` – Trim stream length
- `MINID [~] <ID>` – Remove entries older than given ID

---

## 2. `XRANGE` – Read Range of Entries

Read entries from a stream between ID ranges.

```bash
XRANGE orders - +
XRANGE orders 0 1695119950000-0 COUNT 5
```

---

## 3. `XREVRANGE` – Read in Reverse

Read latest entries first.

```bash
XREVRANGE orders + - COUNT 3
```

---

## 4. `XREAD` – Read New Entries

Useful for simple consumers (not using groups).

```bash
XREAD BLOCK 0 STREAMS orders $
```

---

## 5. `XGROUP CREATE` – Create Consumer Group

Create a group for collaborative consumption.

```bash
XGROUP CREATE orders order_group $
```

### Options:

- `$` – Read new entries only
- `0` – Read from beginning
- `MKSTREAM` – Create stream if not exist

---

## 6. `XREADGROUP` – Read as Consumer in Group

Read entries using consumer group context.

```bash
XREADGROUP GROUP order_group consumer1 BLOCK 5000 COUNT 10 STREAMS orders >
```

- `>` – New entries (not previously delivered)

---

## 7. `XACK` – Acknowledge Entry

Mark a message as successfully processed.

```bash
XACK orders order_group 1695119945821-0
```

---

## 8. `XPENDING` – View Unacknowledged Entries

Check pending messages in a group.

```bash
XPENDING orders order_group
XPENDING orders order_group - + 10 consumer1
```

---

## 9. `XCLAIM` – Claim Stuck Messages

Claim a message that another consumer failed to acknowledge.

```bash
XCLAIM orders order_group consumer2 60000 1695119945821-0
```

### Options:

- `<min-idle-time>` – Time (in ms) a message must be idle before being claimable
- `JUSTID` – Return only IDs
- `IDLE <ms>` – Reset idle time
- `TIME <unix-time>` – Override current time
- `RETRYCOUNT <n>` – Set retry attempts

---

## 10. `XAUTOCLAIM` – Auto-Claim Stuck Messages

Automatically claim entries that have been idle.

```bash
XAUTOCLAIM orders order_group consumer2 60000 0 COUNT 10
```

### Returns:

- Next start ID
- List of claimed entries

### Options:

- `<min-idle-time>` – Required idle duration
- `<start>` – Starting ID for scanning
- `COUNT <n>` – Max number of entries to return
- `JUSTID` – Only return IDs

---

## 11. `XDEL` – Delete an Entry

Remove specific message from the stream.

```bash
XDEL orders 1695119945821-0
```

---

## 12. `XTRIM` – Trim Stream Size

Limit stream size or age.

```bash
XTRIM orders MAXLEN 1000
XTRIM orders MINID 1695119900000-0
```

---

## 13. `XINFO` – Get Stream Metadata

```bash
XINFO STREAM orders
XINFO GROUPS orders
XINFO CONSUMERS orders order_group
```

---

## Summary

| Command      | Purpose                            |
| ------------ | ---------------------------------- |
| `XADD`       | Add new order to stream            |
| `XRANGE`     | Get order history                  |
| `XREAD`      | Basic consumption                  |
| `XGROUP`     | Create group of order handlers     |
| `XREADGROUP` | Group-based reading of orders      |
| `XACK`       | Acknowledge processed order        |
| `XPENDING`   | Monitor unprocessed orders         |
| `XCLAIM`     | Take over stuck orders             |
| `XAUTOCLAIM` | Auto-retry failed order processing |
| `XDEL`       | Remove an entry                    |
| `XTRIM`      | Cap stream size                    |
| `XINFO`      | Introspect stream and consumers    |

Redis Streams provide a powerful, fault-tolerant pattern for event-driven systems like our Order Processing System. Using **Consumer Groups**, `XACK`, `XCLAIM`, and `XAUTOCLAIM`, we ensure **reliable message delivery** even when failures occur.

---

✅ Ideal for microservices, real-time processing, and asynchronous workflows.
