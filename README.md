# bmining

Balanced mininig.  Balanced in terms of what coin to mine. The goal is maintain an even hashrate distribution across all chains being mined.

## Glossary

**bminer** - Balanced miner implemented in a distributed fashion.

## Notes

This package uses a state machine inspired by Rob Pike's discussion about lexer design in this [talk](https://www.youtube.com/watch?v=HxaD_trXwRE).

## Directories

**bminer** - State machine of a Drone that mines crypto during its idle time.

**hashrate** - Fetch network hashrate from a reliable source.  (Only Monero supported at this time.)

**mvp.Mar29** - Minimal Viable Product (MVP) to demo after a one week sprint.
