---
tags: [genetics, practice, mendel]
difficulty: easy-medium
date: 2026-03-08
---

# Mendel's First Law - Practice Problems

## Easy Problems

### Problem 1: Basic Monohybrid Cross
In pea plants, tall (T) is dominant over short (t). Cross two heterozygous tall plants.

**Question**: What are the genotype and phenotype ratios?

<details>
<summary>Solution</summary>

Tt * Tt

```
     T    t
T   TT   Tt
t   Tt   tt
```

- **Genotype ratio**: 1 TT : 2 Tt : 1 tt
- **Phenotype ratio**: 3 tall : 1 short (75% tall, 25% short)
</details>

---

### Problem 2: Test Cross
A purple flower (unknown genotype) is crossed with a white flower (pp). All offspring are purple.

**Question**: What is the genotype of the purple parent?

<details>
<summary>Solution</summary>

Since all offspring are purple, the purple parent must be **PP** (homozygous dominant).

If it were Pp:
- Pp * pp would give 50% Pp (purple) and 50% pp (white)
- But we got 100% purple, so parent = PP
</details>

---

### Problem 3: Probability
Two heterozygous brown-eyed parents (Bb) have children.

**Question**: What's the probability their first child has blue eyes (bb)?

<details>
<summary>Solution</summary>

Bb * Bb

```
     B    b
B   BB   Bb
b   Bb   bb
```

Probability of bb = **1/4 or 25%**
</details>

---

## Medium Problems

### Problem 4: Multiple Offspring
In guinea pigs, black fur (B) is dominant over white (b). Two heterozygous black guinea pigs mate.

**Question**: If they have 4 offspring, what's the probability that exactly 3 are black?

<details>
<summary>Solution</summary>

P(black) = 3/4, P(white) = 1/4

Using binomial probability: C(4,3) * (3/4)**3 * (1/4)**1

= 4 * (27/64) * (1/4)
= 4 * 27/256
= **108/256 ≈ 42.2%**
</details>

---

### Problem 5: Genotype Determination
In snapdragons, red flowers (RR) and white flowers (rr) show incomplete dominance (Rr = pink).

A pink flower is crossed with a white flower.

**Question**: What are the expected offspring ratios?

<details>
<summary>Solution</summary>

Rr * rr

```
     R    r
r   Rr   rr
r   Rr   rr
```

- **Genotype ratio**: 1 Rr : 1 rr
- **Phenotype ratio**: 1 pink : 1 white (50% each)
</details>

---

### Problem 6: Real-World Application
In humans, the ability to taste PTC is dominant (T) over non-tasting (t). A taster woman whose father was a non-taster marries a non-taster man.

**Question**: What's the probability their children will be tasters?

<details>
<summary>Solution</summary>

Woman's genotype: Must be **Tt** (taster with non-taster father)
Man's genotype: **tt** (non-taster)

Tt * tt

```
     T    t
t   Tt   tt
t   Tt   tt
```

Probability of taster children (Tt) = **50%**
</details>

---

## Challenge Problem

### Problem 7: Multiple Generations
A homozygous dominant tall plant (TT) is crossed with a short plant (tt). One F1 offspring is then crossed with a short plant.

**Question**: What fraction of the F2 generation will be tall?

<details>
<summary>Solution</summary>

**P generation**: TT * tt  ->  All F1 are Tt

**F1 cross**: Tt * tt

```
     T    t
t   Tt   tt
t   Tt   tt
```

**F2 result**: 50% Tt (tall), 50% tt (short)

Answer: **1/2 or 50% tall**
</details>
