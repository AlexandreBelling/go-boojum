# Go-boojum

Boojum is an WIP protocol that aims at massively reducing the cost of zk verification on a blockchain through ZK-Snark proof aggregation.
A demo implementation of ZK-SNARK aggregation.

## Security

**The code in this repository is not suitable for production. It is PoC grade and has not undergone extensive testing. It additionally leaks memory.**

## Overview

This [article](https://eprint.iacr.org/2014/595.pdf) introduced recursive snark aggregation using cycle of elliptic curve, and created the concept of PCD and it has been seen as potential solution to instantly (30ms) verify the state of the chain through a SNARK verification. This however raised the problem of data availability as it could potentially creates situations where the state becomes inaccessible. Another problem comes from the fact that a single instruction takes up to 10sec in order to be added inside the proof which makes STARKs more adapted for that.

We propose here a circuit-agnostic solution to combine multiple ppzksnarks proofs (forged for various circuits). This is a variation of PCD comes with two circuits described in the diagram below:

![Aggregation Circuits](./docs/aggregation_circuits.png)

Both circuits are in facts almost identicals, they only differs in the sense that they are not defined on the same EC. Any proof generated on one of theses can be verified recursively inside a proof on the other one. This is a necessary conditions for constructing a practical recursive SNARKs.

One of the main difference with PCD is that boojum accepts a verification key as a public parameter. This is because the aggregation prover cannot make assumption on the two inputs he is going to verify. They could come from any circuits that has one single primary input (any circuit can be converted to a single input circuit using a hash instead of the actual inputs).

Additionnally, PCD recursive aggregation works in a sequential way: assignments are added one after the others in the proof while we describe a protocol aggregating proofs in a hierarchical fashion.

![tree](./docs/tree_of_proof.png)

The leaf nodes (ie: the batch of proofs to be aggregated together) are inputed as:

* A verification key
* A proof
* The primary inputs of the circuits.

## On-chain verification

Each parent node (ie: aggregated proof) takes the hash of the previous proofs as primary inputs. Therefore during verification of the root proof, we need to first reconstruct its input by recursively hashing the intermediary nodes.

![verification](./docs/verification.png)

### Gas Cost estimation

Although no proper benchmark has been run yet. We can estimate that currently each aggregated proofs weights 355 bytes in average (373B for MNT6 and 337B for MNT4). And each verification key (on MNT4 only) weights 717B. This adds up to (355 + 337 + 717 = 1409B) for each proof. This represents an extra cost of 88641 Gas for each proof assuming we can neglect the zero-bytes.

This estimation also does not takes into account the cost of re-hashing the merkle tree. The current implementation makes use of subset-sum hash which is natively implemented in libsnark but which is broken today.

Some of the considered options are [WIP]:

* Pedersen Hash (We could re-use zcash implementation)
* MiMc
* David-Meyers

## Improving the size of the payload

In this aggregation protocol is that we don't care this much about the intermediary proofs. The only thing that matter is that *theses proofs exist and have been successfully verified* the same applies for the leafs proofs (ie: the proofs that are submitted to the process of aggregation).

In the end what an end-user wants to prove is only that they have a valid assignment for a given public input and a given circuit. Therefore, instead of publishing the proofs on-chain we could simply publish a hash of them. The proof would have to be communicated off-chain to the aggregator pool though.

Those improvement the circuit can be represented as below:

![aggregation_circuit_improved](./docs/aggregation_circuit_improved.png)

## Off-chain aggregation

Each aggregation steps takes about 20sec, that means it would takes over 5.5 hours to aggregate 1024 proofs. However, the tree structure makes it easy to possible to distribute across a pool of worker.

The pool would load balance the process of aggregation and each worker would be rewarded for its work. For this purpose we can add an address in the verification in order to protect the worker from impersonation. Adding rewards could however also introduce the issue of byzantine behaviour : a powerfull prover would be incentivized to steal other's job. This part is still a WIP.

The aggregation protocol must also ensure that no attacker can effectively prevent nor slow down the aggregation process. The protocol [handle](https://docs.google.com/presentation/d/1fL0mBF5At4ojW0HhbvBQ2yJHA3_q8q8kiioC6WvY9g4/edit#slide=id.p) described here is currently being considered as the aggregation protocol. It is tolerant to byzantine failures and scales wells with large networks.

## [WIP] Adapted design using Handel

The aggregation mechanism as it is described above is not directly compatible with Handel. Three issues that have to be addressed in order to make the mechanism compatible with Handel.

* On handel each node manage a unique private key in order to sign an aggregate. Therefore, before the aggregation the signer already knows what job he is going to perform. On the other side, with boojum each worker can possibly have several jobs and the pool needs a consensus on who is going to aggregate which proofs. A mechanism to decentralize this should be carefully designed.

* On handel, when a worker is waiting a proof from a faulty worker (timeout or bad proof), it has the possibility to send its job to the next level and that helps guaranteeing the BFTolerance of the protocol. It is not possible to do it with the current design because we are alternating with differents EC.

* Handel needs an aggregation function that is commutative and associative. Boojum's current design uses a "kind of Merkle Tree" based on snark friendly hash functions and it is neither commutative nor associative. A few adaptations should be made in boojum in order to work. Following a discussion with N. Liochon and O. Begassat It might not be an imperative as long as anyone can know "in which order" the tree was built so it can be properly verified in the end.

## Prerequisite

In order to build the source we need the following dependencies

* Docker

## Runing the demo

In this demo, a single worker go-routine aggregates a total of 8 proofs, the aggregation process is controlled by a scheduler go-routine who then verify the aggregated proof.

### With docker

    docker build . -t demo-boojum
    docker run demo-boojum

## Related work

This works makes use of

* [Succinct Non-Interactive Zero Knowledge for a von Neumann Architecture](https://eprint.iacr.org/2013/879.pdf)
* [Incrementally Verifiable Computation or Proofs of Knowledge Imply Time/Space Efficiency](https://link.springer.com/content/pdf/10.1007%2F978-3-540-78524-8_1.pdf)
* [Scalable Zero Knowledge via Cycles of Elliptic Curves](https://eprint.iacr.org/2014/595.pdf)
* [Aggregation protocol for large scale Byzantine committee](https://docs.google.com/presentation/d/1fL0mBF5At4ojW0HhbvBQ2yJHA3_q8q8kiioC6WvY9g4/edit#slide=id.p)