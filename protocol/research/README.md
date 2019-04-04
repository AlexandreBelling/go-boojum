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
  
$T(I)$ is the set of k-ary full trees such that $\forall t \in T(I), \mathit{leaves}(t) \subset I$. In particular, the subset $T_0(I) = \mathit{trivials}(I)$ of $T(I)$ contains all the trivial trees (which are those with a single vertex included in I).
  
$O(I) = \{t \in T(I)\ \vert$ $ \mathit{leaves}(t) = I\}$ containing only rooted trees (ie: trees with a single root). We will sometime refer to it as the set of *output trees*.

$\Phi_t = \mathcal{F}[\mathit{parents}(t), W]$, (for some $t \in T(I)$),  is called the set of contributions of $t$. It is simply a labelling of its parents nodes *(ie: nodes that are not leaves)* by *workers* from $W$. Informally, this object will be the basis of the reward mechanism. As $parents(t \in T_0(I))= \emptyset$ it follows $\Phi(t_0) = \emptyset \rightarrow W$, the empty mapping.

$R(I) = \{(\phi ,t)$ : $t \in T(I), \phi \in \Phi_t \}$ is refered to as the set of the contributed trees. We give $R(I)$ a partial ordering relationship $\prec$ that we define by $\forall x=(t_x, \phi _x), \forall y=(t_y, \phi_y) \in R(I), x \prec y \Leftrightarrow t_x \subset t_y$ and $\phi _x$ is the restriction of $\phi_y$ to $parents(t_y)$.  Additionally, we will often use the notation $\phi_r$ to denote the contribution function of the tree $r \in R(I)$ and $t_r$ its tree.

### Computation and communication model

The current section aims at 

We are using a discrete time modelization. The processes $w$ taking parts in the protocol will have their state updated one each at a time in any order.

$s_n(w) \subset R(I)$ denotes the state of process $w \in W$ at time $n \in \mathbb{N}^{+}$. We provide the power-set of $R(I)$, $S(I)=\mathcal{P}(R(I))$ a partial ordering relationship $\prec$. We define it as the following, Let $s_1, s_2 \in S(I)$, then  $s_1 \prec s_2 \Leftrightarrow \forall \rho_1 \in s_1, \exists r_2 \in s_2$, such that $r_1 \prec r_2$. This notation will be usefull later in order to define byzantine behaviour. Moreover, as the state of a worker can evolve during a distributed protocol, we will sometime note it $s_n(w_i)$.

In our modelization, the state of a worker can only change in the following ways :

- *Message passing* : If a worker $w_i$ can write a substate $\sigma \prec s_{i, n}$ in the state of another worker $w_j$. Its state is updated in the following way : $s_{j, n+1} = s_{j, n} \cup \sigma$. Crash failures can be depicted as absence of communication, this lets the model naturally encompass failure and thus represents efficiently unreliable communication.

- *Aggregation* : A worker $w \in W$ can spontenaously aggregate his state. The new state $s_{n+1}(w)$ is such that:

$$(1) \text{: }card(\sigma) = k \text{ with } \sigma = s_{n}(w) \restriction s_{n+1}(w) = \{\rho_1 , \rho_2, ... \rho_k \} \subset R(I)
$$$$(2) \text{: }card(\sigma') = k \text{ with }\sigma ' = s_{n+1}(w) \restriction s_n(w) = \{r\} \text{ with } r \in R(I)
$$$$(3) \text{: }\forall i,j \lt k + 1 \text{ with } i \neq j, t_i \cap t_j = \emptyset 
$$$$(4) \text{: The children of the root of } t_r \text{ are exactly the roots of the } t_{\rho}$$

$(3)$ implies that $t_i$ and $t_j$ do not share any leaves.

### Protocol output and retribution

Given $W$, a set of identical process, $I$ a set of inputs known to every $w \in W$, a strictly positive integer $k$ that we call the arity of the aggregation.

The execution of the protocol is *successfull* at time $n$ if $\exists w \in W, \exists r \in s(w), o=t_r \in O(I)$ (ie: $t_r$ is a full rooted tree whose leaves are exactly the elements of $I$). In this case $r$ is the output of the protocol.

Let $\pi: O(I) \rightarrow (W \rightarrow \mathbb{N}^{+})$. We will sometime write $\pi_o$ instead of $\pi(o)$. $\pi$ is a *retribution rule* if:

We now define a *retribution rule* to be a mapping  where $$\forall o \in O(I)\text{, } \sum_{w \in W} \pi_o(w) = N_{\pi}$$ where $N_{\pi}$ is some strictly positive constant integer. We will sometime write $\pi_o$ instead of $\pi(o)$.

* A retribution rule is *fair* if $\forall o_1, o_2 \in O(I)$ and $\forall w \in W$ we have,
  
$$t_{o1} = t_{o2} \text{ and } \phi_{o1}^{-1}(w) \subset \phi_{o2}^{-1}(w) \implies \pi_{o1}(w) \leq \pi_{o2}(w)$$

A process $w \in W$ is $\pi\text{-rational}$ for some retribution rule $\pi$ if it behaves in a way that maximizes the ratio:

$$\lambda = \frac{\pi_o(w)}{\text{\#Aggregations performed in the protocol}}$$

It will try it appears as much