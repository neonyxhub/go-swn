## Integration tests

`./run_e2e.sh` spawns 3 Docker containers:
- cwn1 - sends Event via swn1 to swn2
- swn1 - as cwn's peer
- swn2 - as remote target peer where Event is sent

By default, cwn1 also acts as consumer of swn2's incoming Event,
e.g. pwn in usual architecture.

If there is an external pwn ready to consume swn2's incoming Event, then
the external pwn should read `e2e/testdata/debug.yml` file,
where swn2's gRPC server address is written upon swn2 container boots up.

`./run_e2e.sh with_pwn` runs 3 Docker containers but without swn2's incoming Event
consuming, it will wait for 10 seconds (arbitrary timeout) to let external pwn
to consume swn2's Event and after timeout cwn1 tries to consume itself to verify
that there is no buffered Event in swn2.
