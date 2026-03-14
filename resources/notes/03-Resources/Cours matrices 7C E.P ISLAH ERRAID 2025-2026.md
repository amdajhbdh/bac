---
title: Cours matrices 7C E.P ISLAH ERRAID 2025-2026.pdf
source: 03-Resources/Cours matrices 7C E.P ISLAH ERRAID 2025-2026.pdf
date: 2026-03-14
tags: [ai-ocr, extracted]
---

> [!summary]
> Cette section présente les concepts fondamentaux des matrices en mathématiques : définition, vocabulaire, types particuliers (matrices carrées, diagonales, identité, nulle, transposée) et opérations élémentaires (addition, multiplication par un scalaire, produit matriciel).

---

## Définition et vocabulaire

Une matrice $n \times p$ est un tableau rectangulaire de nombres à $n$ lignes et $p$ colonnes.

Les nombres qui.compose la matrice sont appelés les **éléments de la matrice** (ou aussi les coefficients).

Une matrice à $n$ lignes et $p$ colonnes est dite matrice d'ordre $(n, p)$ ou de dimension $n \times p$.

En général on note $a_{ij}$ le coefficient qui se trouve à la $i^{ème}$ ligne et à la $j^{ème}$ colonne, où $1 \leq i \leq n$ et $1 \leq j \leq p$.

$$
A = \begin{pmatrix}
a_{11} & a_{12} & \dots & a_{1p} \\
a_{21} & a_{22} & \dots & a_{2p} \\
\vdots & \vdots & \ddots & \vdots \\
a_{n1} & a_{n2} & \dots & a_{np}
\end{pmatrix}
$$

On représente la matrice sous la forme $A = (a_{ij})_{1 \leq i \leq n}$ ou encore $A = (a_{ij})_{1 < i < n}^{1 < j < p}$.

### Cas particuliers

- Une matrice ne contenant qu'une seule ligne s'appelle **matrice ligne** ou **vecteur ligne**.

- Une matrice ne contenant qu'une seule colonne s'appelle **matrice colonne** ou **vecteur colonne**.

- Une matrice ayant autant de lignes que de colonnes est appelée **matrice carrée**.

###  et types spéciaux

Pour une matrice carrée $A$ on a :

- La **diagonale principale** de $A$ est formée par les coefficients $a_{11}, a_{22}, \dots, a_{nn}$.

> [!note]
> **Matrice diagonale** : Si tous les coefficients n'appartenant pas à la diagonale principale sont nuls, on dit que $A$ est une matrice diagonale.
> $$A \text{ est diagonale } \Leftrightarrow \forall i \neq j ; a_{ij} = 0$$

> [!note]
> **Matrice triangulaire supérieure** : Si tous les coefficients situ és au-dessus de la diagonale principale sont nuls, on dit que $A$ est une matrice triangulaire supérieure.
> $$A \text{ est triangulaire supérieure } \Leftrightarrow \forall i > j ; a_{ij} = 0$$

> [!note]
> **Matrice triangulaire inférieure** : Si tous les coefficients situ és au-dessous de la diagonale principale sont nuls, on dit que $A$ est une matrice triangulaire inférieure.
> $$A \text{ est triangulaire inférieure } \Leftrightarrow \forall i < j ; a_{ij} = 0$$

> [!note]
> **Matrice identité** : On appelle matrice identité d'ordre $n$, toute matrice carrée d'ordre $n$, dont les coefficients de la diagonale principale sont égaux à 1 et tous les autres sont égaux à 0. On la note $I_n$.

> [!note]
> **Matrice nulle** : On appelle matrice nulle de dimension $n \times p$, la matrice dont tous les coefficients sont nuls. On la note $O_{n,p}$ ou simplement $O$.

> [!note]
> **Transposée d'une matrice $A$** : La transposée d'une matrice $A$ est la matrice notée $A^T$ ou $A^t$, obtenue en échangeant les lignes et les colonnes de $A$. Si $A$ est de dimension $n \times p$, alors $A^t$ est de dimension $p \times n$.

**Exemple :**

$$
A = \begin{pmatrix}
1 & 2 \\
3 & 4
\end{pmatrix}
\quad \Rightarrow \quad
A^t = \begin{pmatrix}
1 & 3 \\
2 & 4
\end{pmatrix}
$$

**Propriétés :**

- $(A + B)^t = A^t + B^t$
- $(A \times B)^t = B^t \times A^t$

---

## Opérations sur les matrices

### Addition

Soient $A$ et $B$ deux matrices de même dimension. La matrice $A + B$ est la matrice de même dimension obtenue en additionnant les coefficients de $A$ avec leurs correspondants dans $B$.

Si $A = (a_{ij})_{1 \leq i \leq n}^{1 \leq j \leq p}$ et $B = (b_{ij})_{1 \leq i \leq n}^{1 \leq j \leq p}$ alors :
$$A + B = (a_{ij} + b_{ij})_{1 \leq i \leq n}^{1 \leq j \leq p}$$

> [!note]
> **Propriétés :**
> - **Commutativité** : $A + B = B + A$
> - **Élément neutre** : $A + O = A$
> - **Associativité** : $(A + B) + C = A + (B + C) = A + B + C$

---

### Produit par un réel

Si $A = (a_{ij})_{1 \leq i \leq n}^{1 \leq j \leq p}$ alors :
$$\lambda A = (\lambda a_{ij})_{1 \leq i \leq n}^{1 \leq j \leq p}$$

On multiplie chaque coefficient de $A$ par $\lambda$.

> [!note]
> **Propriétés :**
> - La matrice $(-1) \times A$ se note $-A$ et s'appelle la matrice opposée de $A$
> - La matrice $A - B = A + (-B)$
> - $\lambda(A + B) = \lambda A + \lambda B$
> - $(\lambda + \mu)A = \lambda A + \mu A$
> - $\lambda(\mu A) = (\lambda \mu)A$

---

### Produit de deux matrices

Soit $A$ une matrice de dimension $n \times p$ et $B$ une matrice de dimension $p \times m$. On appelle produit $A \times B$ la matrice de dimension $n \times m$ obtenue en multipliant chaque ligne de $A$ par chaque colonne de $B$. Le coefficient de la $i^{ème}$ ligne et la $j^{ème}$ colonne de $A \times B$ est obtenu en multipliant la $i^{ème}$ ligne de $A$ par la $j^{ème}$ colonne de $B$.

Si $A = (a_{ij})_{1 \leq i \leq n}^{1 \leq j \leq p}$ et $B = (b_{ij})_{1 \leq i \leq p}^{1 \leq j \leq m}$ alors :
$$A \times B = C = (c_{ij})_{1 \leq i \leq n}^{1 \leq j \leq m} \quad \text{et} \quad c_{ij} = \sum_{k=1}^{p} a_{ik} \cdot b_{kj}$$

> [!warning]
> Pour que le produit $A \times B$ existe, il faut obligatoire que le nombre de colonnes de $A$ soit égal au nombre de lignes de $B$.

> [!note]
> **Propriétés :**
> - $I_n \times A_{np} = A_{np}$ et $A_{np} \times I_m = A_{np}$
> - **Le produit des matrices n'est pas commutatif en général** : $A \times B \neq B \times A$
> - $A \times B = 0$ n'implique pas $A = 0$ ou $B = 0$
> - $A \times C = B \times C$ n'implique pas $A = B$
> - **Associativité** : $A \times (B \times C) = (A \times B) \times C$
> - **Distributivité** par rapport à la somme :
> - $A \times (B + C) = A \times B + A \times C$
> - $(A + B) \times C = A \times C + B \times C$
> - La matrice nulle est un élément absorbant : $A \times O = O \times A = O$