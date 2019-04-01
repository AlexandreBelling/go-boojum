# Byzantine fault tolerant aggregation

## Overview

Hypercube-mesh is an aggregation protocol designed to distribute an arbitrary large SNARKs aggregation scheme across a large pool of workers. The protocol includes a reward mechanism to incentivize workers to participate. The protocol aims at operating in a non-trusted context (ie: peers do not trust each others) and enable a massively reduced SNARK verification cost.

In this paper, we gather research made in order to specify what such a protocol should be.

## Problem Statement

Let be :

* $W = \{ w_0, w_1, ..., w_{n-1} \}$ a set of n identical processes that we will also call the workers in this paper.

* $I = \{i_0, i_1, ... , i_{m-1}\}$ a set of input values taken from some set $E$. 
  * The elements if $I$ are assumed to be disctinct : $\forall j,k \in [0, m-1], (i_j = i_k) \implies (j = k)$
  * No assumption on $E$ is throughout the paper
  
* $T(I)$ is a set of binary such that $\forall t \in T(I),I = leaves(t)$
  * $T(I)$ contains non-rooted trees. 
  * In particular, the tree $t_0$ containing only the elements of $I$ and nothing else (no edge) is included in $T(I)$. We will sometime refer to it as the **input tree**.
  * $O(I)$ is a restriction of $T(I)$ containing only rooted trees (ie: trees with a single root). We will sometine refer to it as the set of **output trees**

## Specification

In order to be applicable the protocol must address the following issues:

* Being tolerant to crash fault
* Being tolerant to byzantine fault
* 

### Crash resilience

