# Byzantine fault tolerant aggregation

## Overview

Hypercube-mesh is an aggregation protocol designed to distribute an arbitrary large SNARKs aggregation scheme across a large pool of workers. The protocol includes a reward mechanism to incentivize workers to participate. The protocol aims at operating in a non-trusted context (ie: peers do not trust each others) and enable a massively reduced SNARK verification cost.

In this paper, we gather research made in order to specify what such a protocol should be.

## Preliminary

### Definitions

- k-ary tree : *to be done*
- k-ary full tree : *to be done*
- k-ary rooted tree : *to be done*
- partition of a set : *to be done*
- refinement of a partition : *to be done*

### Common notation

The current section provide usefull notations and definition that will be necessary to build the distributed computation model.

$W = \{ w_0, w_1, ..., w_{n-1} \}$ a set of n identical processes that we will also call the workers in this paper.

$I = \{i_0, i_1, ... , i_{m-1}\}$ a set of input values taken from some set $\Gamma$. The elements if $I$ are assumed to be disctinct : $\forall j,k \in [0, m-1], (i_j = i_k) \Rightarrow (j = k)$. No assumption is made on $\Gamma$ aside the fact that $card(\Gamma) \gt m $ or $\Gamma$ is not finite.
  
$T(I)$ is the set of k-ary full trees such that $\forall t \in T(I),I = leaves(t)$. Note that $T(I)$ may contain non-rooted trees. In particular, the tree $t_0 \in T(I)$ contains only the elements of $I$ and nothing else (no edge). We will sometime refer to it as the *input tree*.
  
$O(I) = \{t \in T(I)\ : t$ *is rooted*$\}$ containing only rooted trees (ie: trees with a single root). We will sometime refer to it as the set of *output trees*.

$\Phi(t) = \mathcal{F}[parents(t), W]$, (for some $t \in T(I)$),  is called the set of retributions of $t$. It is simply a labelling of its parents nodes *(ie: nodes that are not leaves)* by *workers* from $W$. Informally, this object will be the basis of the reward mechanism. Since $parents(t_0)= \emptyset$ we add the convention $\Phi(t_0) = \emptyset$.

We can then define the notion of retributed trees $R(I)$:

$R(I) = \{(\phi ,t)$ : $t \in T(I), \phi \in \Phi(t) \}$ is refered to as the set of the retributed trees. We give $R(I)$ a partial ordering relationship $\prec$ that we define by $\forall x=(t_x, \phi _x), \forall y=(t_y, \phi_y) \in R(I), x \prec y \Leftrightarrow t_x \subset t_y$ and $\phi _x$ is the restriction of $\phi_y$ to $parents(t_y)$.

### Modelization

We are using a discrete time modelization. The processes $w$ taking parts in the protocol will have their state updated one each at a time in any order. 

$s_n(w) \subset R(I)$ denotes the state of process $w \in W$ at time $n \in \mathbb{N}^{+}$. Processes can maintain several retributed tree, this is due to the fact that trees can be overlapped. We provide the power-set of $R(I)$, $S(I)=\mathfrak{P}(R(I))$ a partial ordering relationship : $\prec$. We define it as the following, Let $\sigma$, $s \in S(I)$, then  $\sigma \prec s \Leftrightarrow \forall \rho \in \sigma, \exists r \in s$, such that $\sigma \prec s$. This notation will be usefull later in order to define byzantine behaviour. Moreover, as the state of a worker can evolve during a distributed protocol, we will sometime denote $s_n(w_i)$.

In our modelization, the state of a worker can only change in the following ways :

- **Message passing** : If a worker $w_i$ can write a substate $\sigma \prec s_{i, n}$ in the state of another worker $w_j$. Its state is updated in the following way : $s_{j, n+1} = s_{j, n} \cup \sigma$. Crash failures can be depicted as absence of communication, this makes this model naturally encompass failure and unreliable communication.

- **Aggregation** : A worker $w \in W$ can spontenaously aggregate his state. The new state $s_{n+1}(w)$ is such that:

  - $\sigma = s_{n}(w) \restriction s_{n+1}(w)$ contains a single element $\rho \in R(I)$
  - $\sigma ' = s_{n+1}(w) \restriction s_n(w)$ contains a single element $r \in R(I)$
  - $\sigma ' \prec \sigma$ and $\sigma '$ contains only an additional vertices

## Specification

In order to be applicable the protocol must address the following issues:

* Being tolerant to crash fault
* Being tolerant to byzantine fault
* 

### Crash resilience

