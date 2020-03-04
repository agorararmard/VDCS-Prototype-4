# VDCS-Prototype-3
This is an experiment for milestone 3 using Go
## To Do:
### Client:
- All modifications in the VDCS.go library:
- Request Cycle: Send circuit description(# of gates) to directory service
- SendServerGarble: Prepare array of messages (generate symmetric keys to encrypt the messages) then use public keys to encrypt the symmetric keys.
- SendServerEval: Send input wires only to server n of the cycle encrypted using the same scheme.
### Directory Service:
#### Registration:
- Registration Servers: IP, Port. Public Key, # of Gates willing to work on. Fees per gate.
- Token returned.
- Registration Clients: Register once compute many.
#### Cycle Construction:
- Client sends # of gates and fees paid.
- Find available servers. Servers keep pinging you every 1 minutes.
- Filter servers based on number of gates willing to compute.
- Filter rest of servers based on who’s cheaper per gate.
- Randomly choose set (n) return info (IP, Port, Public Key) to client.
#### Failure Handling: 
- if Server i tries to send to server i+1 and finds it offline:
- Request a new server i+1 from directory of servers.
- Directory of service chooses the new server, and send its info to the client.
- Client re-encrypts the symmetric key and send it to server i.
- Server i replaces the previous index of key with the new one w pass to the new server i+1.
### Servers:
- Unit servers,
- I receive a general message. Check the type.
- If Garble: Extract circuit field and run garble. Then send to next after encrypting the garbled circuit. Array of messages with the first omitted.
- If Rerand: Do nothing. Then send to next after encrypting the garbled circuit. Array of messages with first and last omitted. 
- If Evaluate: Extract garbled circuit. Wait for input wires. Evaluate and return.
### VDCS.go Structs:
Done

### Questions
- Client’s anonymity?
- Validating computing power: Challenge + hardware hash? Statistical model?
- Rerandomization?
