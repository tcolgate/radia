#View Of Network Quality

**This is very early days, nothing usable hear yet**

A service for peer to peer network quality evaluations. Establishing
RTT and packet loss statistics for a mesh of shosts to build optimal
spanning trees.

- Built on gRPC
- Aiming to provide support for
  - Leader election
  - Leader transition
  - Message broadcast
  - Multiple network topologies overly on the initial mesh
  - simulation of network updates from data sources (like RRDs)
