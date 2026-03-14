---
title: Génétique - Croisements et Hérédité
niveau: 7ème/Terminale
matière: Biology
date: 2026-03-13
tags:
  - génétique
  - croisements
  - mendel
  - hybridisme
  - monohybridisme
  - dihybridisme
  - épistasie
aliases:
  - Genetics-Crosses
  - Croisements-Génétiques
cssclasses:
  - tables
---

# Génétique - Croisements et Hérédité

## Table des Matières

```dataview
TABLE WITHOUT ID
file.link as "Chapitre",
niveau as "Niveau"
FROM "01-Concepts/Biology"
WHERE contains(tags, "génétique")
SORT file.name
```

---

> **Source:** Cours de Biology BAC Mauritanie
> **Dernière mise à jour:** 2026-03-13

---

# Partie 1: Croisement Monohybride

## 1.1 Définitions Fondamentales

### 1.1.1 Qu'est-ce qu'un Monohybride?

Un croisement **monohybride** est une expérience de génétique qui étudie l'hérédité d'un **seul caractère** (une seule paire d'allèles) à la fois.

> [!definition]
> **Caractère:** Une caractéristique héréitaire observable (couleur des yeux, forme des graines, groupe sanguin, etc.)
> 
> **Allèle:** Une des différentes versions d'un gène.
> 
> **Génotype:** L'ensemble des allèles présents chez un individu.
> 
> **Phénotype:** L'expression visible du génotype (l'aspect physique).

### 1.1.2 Notation Génétique

| Symbole | Signification |
|---------|---------------|
| $A$ | Allèle dominant (s'exprime toujours) |
| $a$ | Allèle récessif (masqué par A) |
| $AA$ | Homozygote dominant |
| $Aa$ | Hétérozygote (porteur) |
| $aa$ | Homozygote récessif |
| $A-$ | Dominant (AA ou Aa) |

### 1.1.3 Cycle de Croisement

$$P \xrightarrow{\text{croisement}} F_1 \xrightarrow{\text{croisement}} F_2$$

- **P (Parents):** Génération initiale
- **F1 (Première génération):** Descendants du croisement P × P
- **F2 (Deuxième génération):** Descendants du croisement F1 × F1

---

## 1.2 Types de Dominance

### A. Dominance Complète (Dominance Totale)

#### Définition Détaillée

> [!tip] Définition
> Dans la **dominance complète**, l'allèle dominant **masque complètement** l'expression de l'allèle récessif chez l'hétérozygote. Le phénotype de l'hétérozygote est identique à celui de l'homozygote dominant.

#### Mécanisme Moléculaire

1. L'allèle dominant produit une protéine fonctionnelle
2. Cette protéine est présente dans la cellule
3. L'allèle récessif peut être:
   - Non fonctionnel (pas de protéine)
   - Moins efficace
   - Absent

#### Exemples Connus

| Espèce | Caractère | Dominant | Récessif |
|--------|-----------|----------|----------|
| Pois | Couleur graine | Jaune (A) | Vert (a) |
| Pois | Forme graine | Lisse (L) | Ridé (l) |
| Haricot | Couleur fleur | Violet (P) | Blanc (p) |
| Souris | Couleur pelage | Noir (B) | Albinos (b) |
| Drosophile | Couleur yeux | Rouge (E) | Blanc (e) |

#### Croisement Détaillé: $Aa \times Aa$

**Parents:**
- Parent 1: Hétérozygote (Aa)
- Parent 2: Hétérozygote (Aa)

**Gamètes produits:**
- Parent 1: $A$, $a$ (un gamète reçoit A, l'autre reçoit a)
- Parent 2: $A$, $a$

**Grille de Punnett:**

| | **A** | **a** |
|---|-------|-------|
| **A** | $AA$ | $Aa$ |
| **a** | $Aa$ | $aa$ |

**Résultats F2 - Analyse:**

| Génotype | Nombre | Fréquence | Phénotype |
|----------|--------|-----------|-----------|
| $AA$ | 1 | 25% | Dominant |
| $Aa$ | 2 | 50% | Dominant |
| $aa$ | 1 | 25% | Récessif |

**Ratio:**
- Génotypique: **1:2:1** (AA:Aa:aa)
- Phénotypique: **3:1** (Dominant:Récessif)

#### Formule de Probabilité

$$P(\text{dominant}) = P(AA) + P(Aa) = \frac{1}{4} + \frac{2}{4} = \frac{3}{4}$$

$$P(\text{récessif}) = P(aa) = \frac{1}{4}$$

#### Preuve Expérimentale (Mendel)

Mendel a croisé des pois homozygotes Jaunes (AA) avec des pois homozygotes Verts (aa):

1. **F1:** 100% Jaunes (Aa)
2. **F1 × F1:** 
   - 25% AA (Jaunes)
   - 50% Aa (Jaunes)
   - 25% aa (Verts)

> [!success] Conclusion de Mendel
> Les facteurs hérédites (allèles) se séparent lors de la formation des gamètes et se recombinent au hasard lors de la fécondation.

---

### B. Dominance Incomplète (Hérédité Intermédiaire)

#### Définition Détaillée

> [!warning] Définition
> Dans la **dominance incomplète** (aussi appelée **hétérodominance**), l'hétérozygote présente un phénotype **intermédiaire** entre les deux homozygotes. Aucun allèle ne domine l'autre.

**Différence avec dominance complète:**

| Dominance Complète | Dominance Incomplète |
|-------------------|---------------------|
| Aa = phénotype A | Aa = phénotype intermédiaire |
| Pas de mélange visible | Mélange visible |
| Ratio F2: 3:1 | Ratio F2: 1:2:1 |

#### Mécanisme Moléculaire

1. Chaque allèle produit une protéine
2. Chez l'hétérozygote, les deux protéines sont présentes
3. Elles interagissent pour produire un phénotype intermédiaire
4. Exemple: pigment rouge + pigment blanc → pigment rose

#### Exemples Détaillés

##### Exemple 1: Fleur de Muflier (Antirrhinum)

| Génotype | Phénotype | Couleur |
|----------|-----------|---------|
| $RR$ | Homozygote rouge | Rouge foncé |
| $RB$ | Hétérozygote | Rose |
| $BB$ | Homozygote blanc | Blanc |

**Croisement:** $RR \times BB$

**F1:** 100% $RB$ (roses)

**Croisement F1:** $RB \times RB$

| | **R** | **B** |
|---|-------|-------|
| **R** | $RR$ (rouge) | $RB$ (rose) |
| **B** | $RB$ (rose) | $BB$ (blanc) |

**Ratio F2:** 1 rouge : 2 roses : 1 blanc

##### Exemple 2: Couleur du pelage chez le bétail

| Génotype | Phénotype |
|----------|-----------|
| $WW$ | Blanc |
| $WB$ | Gris (w Gris) |
| $BB$ | Noir |

##### Exemple 3: Couleur des cheveux humains (simplifié)

| Génotype | Phénotype |
|----------|-----------|
| $BB$ | Cheveux noirs |
| $Bb$ | Cheveux châtains |
| $bb$ | Cheveux blonds |

> [!note] Important
> La dominance incomplète n'est pas un "mélange" au niveau génétique. Les allèles restent distincts et se séparent lors de la méiose. C'est l'expression phénotypique qui semble être un mélange.

---

### C. Codominance

#### Définition Détaillée

> [!danger] Définition
> Dans la **codominance**, les deux alleles sont **complètement et simultanément exprimés** chez l'hétérozygote. Aucun ne domine l'autre, et les deux produits sont visibles.

**Différences clés:**

| Dominance Incomplète | Codominance |
|---------------------|-------------|
| Phénotype intermédiaire | Les DEUX phénotypes visibles |
| Aa = mélange | Aa = phénotype 1 + phénotype 2 |
| Ex: fleur rose | Ex: sang AB (A + B) |

#### Exemples Détaillés

##### Exemple 1: Groupes Sanguins Humains (Système ABO)

Le système ABO est l'exemple classique de codominance avec allèles multiples.

**Les trois allèles:**

| Allèle | Description | Effet |
|--------|-------------|-------|
| $I^A$ | Antigène A | Produce enzyme qui ajoute糖A aux GR |
| $I^B$ | Antigène B | Produce enzyme qui ajoute糖B aux GR |
| $i$ | Pas d'antigène | Pas d'enzyme active |

**Relations de dominance:**
- $I^A$ > $i$ (A dominant sur O)
- $I^B$ > $i$ (B dominant sur O)
- $I^A$ = $I^B$ (codominance!)

**Tableau des génotypes et phénotypes:**

| Génotype | Phénotype (Groupe) | Antigènes présents |
|----------|-------------------|-------------------|
| $I^A I^A$ ou $I^A i$ | A | A seulement |
| $I^B I^B$ ou $I^B i$ | B | B seulement |
| $I^A I^B$ | AB | A et B (CODOMINANTS!) |
| $ii$ | O | Aucun |

**Croisement AB × AB:**

| | $I^A$ | $I^B$ |
|---|-------|-------|
| **$I^A$** | $I^A I^A$ (A) | $I^A I^B$ (AB) |
| **$I^B$** | $I^A I^B$ (AB) | $I^B I^B$ (B) |

**Ratio:** 1 A : 2 AB : 1 B

> [!tip] Point Clé
> Le groupe AB montre que les deux antigènes A et B sont exprimés simultanément sur les globules rouges. C'est la **codominance pure**.

##### Exemple 2: Robe "Roan" du bétail

| Génotype | Phénotype |
|----------|-----------|
| $RR$ | Rouge |
| $RW$ | **Roan** (poils rouges ET blancs ensemble) |
| $WW$ | Blanc |

##### Exemple 3: Planches de sang chez les poules

某些 races de poules ont des plumes avec des motifs de couleurs distincts qui s'expriment simultanément.

##### Exemple 4: Gène de la drépanocytose

Le gène de l'hémoglobine:
- $H^A$: Hémoglobine normale
- $H^S$: Hémoglobine S (drépanocytose)

Les hétérozygotes $H^A H^S$ produisent les deux types d'hémoglobine → phénomène de **codominance** au niveau moléculaire (mais phénotype généralement normal).

---

### D. Allèles Multiples (Polymorphisme)

#### Définition

> [!info] Définition
> Les **allèles multiples** (ou polymorphisme) surviennent lorsqu'il existe **plus de deux allèles différents** pour un gène donné dans une population.

Un individu ne peut avoir que **deux allèles** (un de chaque parent), mais la **population** peut en avoir beaucoup plus.

#### Système ABO - Trois Allèles

$$I^A, \quad I^B, \quad i$$

**Dans la population:** 3 allèles possibles
**Chez un individu:** Maximum 2 allèles (sur les 3)

#### Fréquences Alléliques (Hardy-Weinberg)

Pour un gène avec 2 allèles (simplifié):

$$p + q = 1$$

- $p$ = fréquence de l'allèle dominant
- $q$ = fréquence de l'allèle récessif

**Fréquences génotypiques:**

$$p^2 + 2pq + q^2 = 1$$

| Génotype | Fréquence |
|----------|-----------|
| Homozygote dominant | $p^2$ |
| Hétérozygote | $2pq$ |
| Homozygote récessif | $q^2$ |

**Pour 3 allèles (ABO):**

$$p + q + r = 1$$

$$(p + q + r)^2 = p^2 + q^2 + r^2 + 2pq + 2pr + 2qr = 1$$

| Génotype | Fréquence |
|----------|-----------|
| $I^A I^A$ | $p^2$ |
| $I^B I^B$ | $q^2$ |
| $ii$ | $r^2$ |
| $I^A I^B$ | $2pq$ |
| $I^A i$ | $2pr$ |
| $I^B i$ | $2qr$ |

#### Exemples d'Allèles Multiples

##### 1. Groupe sanguin MN

| Génotype | Phénotype |
|----------|-----------|
| $L^M L^M$ | M |
| $L^M L^N$ | MN |
| $L^N L^N$ | N |

##### 2. Couleur de la robe des lapins

5 allèles contrôlent la pigmentation:
- $C$ = couleur complète
- $c^{ch}$ = chinchilla
- $c^h$ = himalayen
- $c^a$ = albinos

##### 3. Gènes HLA (immunité)

Plus de 100 allèles différents pour certains gènes HLA!

---

## 1.3 Résumé: Comparaison des Types de Monohybridisme

| Type | Génotype F2 | Phénotype F2 | Ratio G | Ratio P | Exemple |
|------|-------------|--------------|---------|---------|---------|
| Dominance complète | 1:2:1 | 3:1 | AA:Aa:aa | Dom:Rec | Pois Jaune/Vert |
| Dominance incomplète | 1:2:1 | 1:2:1 | RR:RB:BB | Rouge:Rose:Blanc | Muflier |
| Codominance | 1:2:1 | 1:2:1 | $I^A I^A$:$I^A I^B$:$I^B I^B$ | A:AB:B | Groupes sanguins |
| Allèles multiples | Variable | Variable | - | - | ABO |

---

## 1.4 Exercices Monohybridisme

### Exercice 1: Croisement Simple

**Énoncé:** Chez le pois, le caractère "graine lisse" est dominant (L) sur "graine ridée" (l).

1. Quels sont les génotypes et phénotypes possibles du croisement Ll × ll?
2. Quelle est la probabilité d'obtenir un descendant ridé?

**Solution:**

| | L | l |
|---|---|---|
| **l** | Ll (lisse) | ll (ridé) |
| **l** | Ll (lisse) | ll (ridé) |

- Lisse (Ll): 50%
- Ridé (ll): 50%

$$P(\text{ridé}) = \frac{2}{4} = 50\%$$

---

### Exercice 2: Groupe Sanguin

**Énoncé:** Un père de groupe A (homozygote) et une mère de groupe B (homozygote) ont un enfant. Quels sont tous les groupes sanguins possibles?

**Solution:**

Père: $I^A I^A$ × Mère: $I^B I^B$

| | $I^A$ | $I^A$ |
|---|---|---|
| **$I^B$** | $I^A I^B$ (AB) | $I^A I^B$ (AB) |
| **$I^B$** | $I^A I^B$ (AB) | $I^A I^B$ (AB) |

**100% Groupe AB!**

---

### Exercice 3: Dominance Incomplète

**Énoncé:** Chez le muflier, une fleur rouge (RR) croise avec une fleur blanche (BB). 
1. Donnez le phénotype F1
2. Croisez ensuite les F1 entre eux
3. Donnez le ratio F2

**Solution:**

1. F1: 100% RB = Rose (hétérozygote = intermédiaire)
2. F1 × F1: RB × RB

| | R | B |
|---|---|---|
| **R** | RR (rouge) | RB (rose) |
| **B** | RB (rose) | BB (blanc) |

3. **Ratio F2:** 
   - 1 Rouge : 2 Roses : 1 Blanc

---

# Partie 2: Croisement Dihybride

## 2.1 Définition

> [!definition]
> Un croisement **dihybride** étudie l'hérédité de **deux caractères simultanément**, chacun contrôlé par un gène différent.

### Notation

Pour deux gènes:
- Gène 1: allèles $A$ et $a$
- Gène 2: allèles $B$ et $b$

### Condition d'Indépendance

Pour que les deux gènes se transmettent indépendamment, ils doivent:
1. Être sur des **chromosomes différents**, OU
2. Être très **éloignés** sur le même chromosome

---

## 2.2 Indépendance (Loi d'Assortiment Indépendant)

### Énoncé de la Loi

> [!quote] Mendel (1865)
> "Les facteurs (allèles) pour différents caractères sont distribués indépendamment les uns des autres lors de la formation des gamètes."

### Croisement Type: $AaBb \times AaBb$

#### Gamètes Possibles

Chaque parent hétérozygote pour deux gènes produit **4 types de gamètes**:

$$AB, \ Ab, \ aB, \ ab$$

> [!tip] Règle
> Pour $n$ gènes hétérozygotes indépendants: $2^n$ gamètes différents
> - 1 gène: $2^1 = 2$ gamètes
> - 2 gènes: $2^2 = 4$ gamètes
> - 3 gènes: $2^3 = 8$ gamètes

#### Grille de Punnett 4×4

| | **AB** | **Ab** | **aB** | **ab** |
|---|--------|--------|--------|--------|
| **AB** | AABB | AABb | AaBB | AaBb |
| **Ab** | AABb | AAll | AaBb | Aall |
| **aB** | AaBB | AaBb | aaBB | aaBb |
| **ab** | AaBb | Aall | aaBb | aabb |

#### Ratio Phénotypique F2: 9:3:3:1

| Phénotype | Génotypes | Nombre | Fraction |
|-----------|-----------|--------|----------|
| A- B- | 9 types | 9 | 9/16 |
| A- bb | 3 types | 3 | 3/16 |
| aa B- | 3 types | 3 | 3/16 |
| aa bb | 1 type | 1 | 1/16 |

#### Exemple Détaillé: Pois Jaune-Lisse × Vert-Ridé

**Caractères:**
- $A$ = Jaune, $a$ = Vert
- $L$ = Lisse, $l$ = Ridé

**Parents:** $AALL$ × $aabb$

**Gamètes F1:** tous $AaLl$

**Croisement F1 × F1:** $AaLl$ × $AaLl$

**Résultats F2:**

| Phénotype | Description | Nombre | Pourcentage |
|-----------|-------------|--------|-------------|
| Jaune-Lisse | A- L- | 9 | 56.25% |
| Jaune-Ridé | A- ll | 3 | 18.75% |
| Vert-Lisse | aa L- | 3 | 18.75% |
| Vert-Ridé | aa ll | 1 | 6.25% |

---

## 2.3 Calcul des Probabilités

### Méthode multiplicative

Pour des gènes **indépendants**, on calcule chaque caractère séparément, puis on multiplie:

$$P(\text{phénotype combiné}) = P(\text{caractère 1}) \times P(\text{caractère 2})$$

### Tableau des Probabilités

Pour $AaBb \times AaBb$:

| Phénotype | Calcul | Résultat |
|-----------|--------|----------|
| A- B- | $\frac{3}{4} \times \frac{3}{4}$ | 9/16 |
| A- bb | $\frac{3}{4} \times \frac{1}{4}$ | 3/16 |
| aa B- | $\frac{1}{4} \times \frac{3}{4}$ | 3/16 |
| aa bb | $\frac{1}{4} \times \frac{1}{4}$ | 1/16 |

### Formule Générale

Pour $n$ gènes indépendants avec dominance complète:

$$\text{Ratio} = (3:1)^n$$

| $n$ | Ratio |
|-----|-------|
| 1 | 3:1 |
| 2 | 9:3:3:1 |
| 3 | 27:9:9:9:3:3:3:1 |
| 4 | 81:27:27:27:9:9:9:9:3:3:3:3:1 |

---

## 2.4 Croisement Test (Test Cross)

### Définition

> [!definition]
> Le **croisement test** consiste à croiser un individu de phénotype dominant (génotype inconnu) avec un **homozygote récessif**.

### Utilité

Il permet de déterminer si l'indominant est:
- **Homozygote dominant (AA)** → Tous descendants dominants
- **Hétérozygote (Aa)** → 50% dominants, 50% récessifs

### Exemple

**Problème:** Un pois à graines lisses (L-). Quel est son génotype?

**Croisement test:** L- × ll

| Si parent est AA | Si parent est Aa |
|------------------|------------------|
| 100% lisse | 50% lisse, 50% ridés |

**Observation:** Si des descendants ridés apparaissent, le parent était hétérozygote!

---

## 2.5 Exercices Dihybridisme

### Exercice 1: Probabilité Simple

**Énoncé:** Croisement $AaBb \times aabb$. Calculez:
1. P(A-bb)
2. P(aaB-)
3. P(aabb)

**Solution:**

| Gamète | A- | aa | B- | bb |
|--------|----|----|----|----|
| AB | ✓ | | ✓ | |
| Ab | ✓ | | | ✓ |
| aB | | ✓ | ✓ | |
| ab | | ✓ | | ✓ |

1. **P(A-bb):** $P(A-) \times P(bb) = \frac{1}{2} \times \frac{1}{2} = \frac{1}{4}$
2. **P(aaB-):** $P(aa) \times P(B-) = \frac{1}{2} \times \frac{1}{2} = \frac{1}{4}$
3. **P(aabb):** $P(aa) \times P(bb) = \frac{1}{2} \times \frac{1}{2} = \frac{1}{4}$

---

### Exercice 2: Ratio Attend

**Énoncé:** Quel est le ratio phénotypique attendu pour le croisement $AABb \times AAbb$?

**Solution:**

- Gène A: AA × aa → 100% Aa (tous A-)
- Gène B: Bb × bb → 50% B-, 50% bb

| Phénotype | Ratio |
|-----------|-------|
| A- B- | 1/2 |
| A- bb | 1/2 |

**Ratio:** 1:1 (pour le caractère B)

---

# Partie 3: Gènes Liés et Types de Liaison

## 3.1 Liaison Génétique (Linkage) - Vue d'Ensemble

### Définition

> [!definition]
> Deux gènes sont **liés** (linked) s'ils sont situés sur le **même chromosome**. Ils tendent à être hérités ensemble.

### Conséquence

Les gènes liés **ne suivent pas** la loi d'indépendance de Mendel. Leurs combinaisons parentales sont conservées plus fréquemment.

### Fréquence des Gamètes

Pour gènes liés:
- Gamètes **parentaux** (non recombinant): > 50%
- Gamètes **recombinants**: < 50%

### Notation

On représente les chromosomes avec les allèles sur la même ligne:

$$\frac{AB}{ab} = \text{couplage cis}$$
$$\frac{Ab}{aB} = \text{couplage trans}$$

---

## 3.2 Types de Liaison Génétique

### A. Liaison Absolue (Complète)

#### Définition

> [!danger] Liaison Absolue
> La **liaison absolue** (ou **complète**) se produit lorsque deux gènes sont tellement proches qu'ils ne se séparent **jamais** par crossing-over. La fréquence de recombinaison est **0%**.

#### Caractéristiques

| Propriété | Valeur |
|-----------|--------|
| Distance | 0 cM |
| Fréquence recombinaison | 0% |
| Gamètes parentaux | 100% |
| Gamètes recombinés | 0% |

#### Exemple

```
Chromosome parental:  A --- B
Chromosome homologue:  a --- b

Après méiose:
→ 50% des gamètes: AB
→ 50% des gamètes: ab

AUCUN gamète recombinant! (AB, ab seulement)
```

#### Résultat du Croisement

Pour $\frac{AB}{ab} \times \frac{ab}{ab}$ (croisement test):

| Phénotype | Génotype | Fréquence |
|-----------|----------|-----------|
| Dominant 1 + Dominant 2 | AB/ab | 50% |
| Récessif 1 + Récessif 2 | ab/ab | 50% |

**Ratio: 1:1** (comme si c'était un seul gène!)

> [!tip]
> Avec liaison absolue, le dihybridisme se comporte comme un monohybridisme.

---

### B. Liaison Partielle

#### Définition

> [!warning] Liaison Partielle
> La **liaison partielle** se produit lorsque deux gènes sont sur le même chromosome mais **séparés par une distance modérée** (1-50 cM). Le crossing-over peut les séparer, mais moins fréquemment que s'ils étaient indépendants.

#### Caractéristiques

| Propriété | Valeur |
|-----------|--------|
| Distance | 1-50 cM |
| Fréquence recombinaison | 1-50% |
| Gamètes parentaux | 50-99% |
| Gamètes recombinés | 1-50% |

#### Mécanisme

```
Avant crossing-over:
Chromosome 1:  A ------- B
Chromosome 2:  a ------- b

Après crossing-over (1 événement sur 4):
Chromosome 1:  A ----- b
Chromosome 2:  a ----- B
```

#### Fréquence des Gamètes

Pour une distance $d$ (en fraction):

| Type | Fréquence |
|------|-----------|
| Parentaux | $1 - d$ |
| Recombinés | $d$ |

#### Croisement Test avec Liaison Partielle

**Exemple:** Distance = 20% (0.2)

$$\frac{AB}{ab} \times \frac{ab}{ab}$$

| Gamète | Type | Fréquence | Descendants |
|--------|------|-----------|--------------|
| AB | Parental | 40% | A+B+ |
| ab | Parental | 40% | a+b- |
| Ab | Recombiné | 10% | A+b- |
| aB | Recombiné | 10% | a+B+ |

**Ratio observed:** 40:40:10:10 = **4:4:1:1**

> [!important]
> Contrairement à l'indépendance (1:1:1:1), on observe plus de types parentaux!

---

### C. Liaison Indépendante

#### Définition

> [!success] Liaison Indépendante
> Deux gènes sont **indépendants** (non liés) s'ils sont situés sur:
> - Des **chromosomes différents**, OU
> - Une distance > 50 cM sur le même chromosome

#### Caractéristiques

| Propriété | Valeur |
|-----------|--------|
| Distance | > 50 cM (ou chromosomes différents) |
| Fréquence recombinaison | 50% |
| Gamètes parentaux | 50% |
| Gamètes recombinés | 50% |

#### Résultat du Croisement Test

Pour $AaBb \times aabb$ (gènes indépendants):

| Gamète | Fréquence | Type |
|--------|------------|------|
| AB | 25% | Parental |
| Ab | 25% | Recombiné |
| aB | 25% | Recombiné |
| ab | 25% | Parental |

**Ratio: 1:1:1:1**

> [!note]
> Avec liaison indépendante, le ratio est identique pour tous les gamètes car il n'y a pas de différence entre gamètes parentaux et recombinés!

---

## 3.3 Tableau Comparatif: Types de Liaison

| Type | Distance | % Recombinants | % Parentaux | Croisement Test (AB/ab × ab/ab) |
|------|----------|----------------|-------------|--------------------------------|
| **Absolue** | 0 cM | 0% | 100% | 1:0:0:1 (seulement AB et ab) |
| **Partielle** | 1-50 cM | 1-50% | 50-99% | Ratio intermédiaire |
| **Indépendante** | >50 cM | 50% | 50% | 1:1:1:1 |

---

## 3.4 Croisement F1 × F1: Analyse Complète

### Introduction

Le croisement **F1 × F1** (croisement entre hybrides de première génération) est fundamental pour déterminer:
- Le type de dominance
- La liaison entre gènes
- Les proportions attendues

### Cas 1: Monohybridisme F1 × F1

#### Croisement: $Aa \times Aa$

| | A | a |
|---|---|---|
| **A** | AA (25%) | Aa (25%) |
| **a** | Aa (25%) | aa (25%) |

**Ratio génotypique:** 1:2:1
**Ratio phénotypique:** 3:1 (dominance complète)

---

### Cas 2: Dihybridisme Indépendant F1 × F1

#### Croisement: $AaBb \times AaBb$

Quand les gènes sont sur des chromosomes différents:

| Phénotype | Ratio attendu |
|-----------|---------------|
| A- B- | 9 |
| A- bb | 3 |
| aa B- | 3 |
| aa bb | 1 |

**Ratio: 9:3:3:1**

---

### Cas 3: Dihybridisme avec Liaison Partielle F1 × F1

#### Configuration: $\frac{AB}{ab} \times \frac{AB}{ab}$

Soit une distance $d$ entre les gènes.

**Calcul des gamètes:**

Pour chaque parent $\frac{AB}{ab}$:
- Gamètes parentaux: $AB$, $ab$ (fréquence: $\frac{1-d}{2}$ chacun)
- Gamètes recombinés: $Ab$, $aB$ (fréquence: $\frac{d}{2}$ chacun)

**Grille de Croisement (avec $d$ = 20%):**

| Parent 1\2 | AB (0.4) | ab (0.4) | Ab (0.1) | aB (0.1) |
|------------|----------|----------|----------|----------|
| AB (0.4) | AABB | AaBb | AABb | AaBB |
| ab (0.4) | AaBb | aabb | Aabb | aaBb |
| Ab (0.1) | AABb | Aabb | AAbb | AaBb |
| aB (0.1) | AaBB | aaBb | AaBb | aaBB |

**Fréquences phénotypiques (pour d = 20%):**

| Phénotype | Génotypes | Fréquence |
|-----------|-----------|------------|
| A- B- | AABB, AABb(2), AaBB(2), AaBb(4) | 0.4×0.4 + 2×0.4×0.1 + 2×0.4×0.1 + 4×0.1×0.1 = 0.16 + 0.08 + 0.08 + 0.04 = **0.36** |
| A- bb | AAbb, Aabb(2) | 0.4×0.1 + 2×0.4×0.1 = 0.04 + 0.08 = **0.12** |
| aa B- | aaBB, aaBb(2) | 0.1×0.4 + 2×0.1×0.4 = 0.04 + 0.08 = **0.12** |
| aa bb | aabb | 0.4×0.4 = **0.16** |

**Ratio observé:** 36:12:12:16 = **9:3:3:4** (proche de 9:3:3:1 mais modifié!)

> [!important]
> Avec liaison partielle, le ratio s'éloigne de 9:3:3:1!

---

### Cas 4: Liaison Absolue F1 × F1 ($\frac{AB}{ab} \times \frac{AB}{ab}$)

Quand $d = 0$ (pas de recombinaison):

**Gamètes:** Seulement $AB$ et $ab$ (50% chacun)

| | AB (0.5) | ab (0.5) |
|---|----------|----------|
| AB (0.5) | AABB | AaBb |
| ab (0.5) | AaBb | aabb |

**Ratio phénotypique:** 1:2:1 (comme monohybridisme!)

| Phénotype | Fréquence |
|-----------|-----------|
| A- B- | 75% |
| aa bb | 25% |

---

### Cas 5: F1 × F1 avec Gènes Indépendants (Récapitulatif)

Pour $AaBb \times AaBb$ avec gènes indépendants:

| Gène A | Gène B | Probabilité combinée |
|--------|--------|---------------------|
| AA (1/4) | BB (1/4) | 1/16 (A-B-) |
| AA (1/4) | bb (1/4) | 1/16 (A-bb) |
| aa (1/4) | BB (1/4) | 1/16 (aaB-) |
| aa (1/4) | bb (1/4) | 1/16 (aabb) |
| Aa (2/4) | BB (1/4) | 2/16 (A-B-) |
| Aa (2/4) | bb (1/4) | 2/16 (A-bb) |
| aa (1/4) | Bb (2/4) | 2/16 (aaB-) |
| Aa (2/4) | Bb (2/4) | 4/16 (A-B-) |

**Total A-B-:** 1+2+2+4 = 9/16 ✓
**Total A-bb:** 1+2 = 3/16 ✓
**Total aaB-:** 1+2 = 3/16 ✓
**Total aabb:** 1/16 ✓

---

---

## 3.6 Configuration CIS vs TRANS (Suite)

### Croisement Test: Déterminer CIS vs TRANS

---

## 3.2 Configuration CIS vs TRANS

### Configuration CIS (En Cis)

> [!tip]
> Les **deux allèles dominants** sont sur un chromosome, les **deux récessifs** sur l'autre.

$$\frac{AB}{ab}$$

- Chromosome 1: $AB$
- Chromosome 2: $ab$

**Signification:**
- Dominants ensemble
- Récessifs ensemble
- "En CIS" = du même côté

### Configuration TRANS (En Trans)

> [!warning]
> Chaque chromosome porte un **dominant** et un **récessif**.

$$\frac{Ab}{aB}$$

- Chromosome 1: $Ab$
- Chromosome 2: $aB$

**Signification:**
- Dominant d'un gène avec récessif de l'autre
- "En TRANS" = opposés

### Tableau Comparatif

| Configuration | Représentation | Description |
|--------------|----------------|-------------|
| **CIS** | $AB/ab$ | $AB$ sur un chromosome, $ab$ sur l'autre |
| **TRANS** | $Ab/aB$ | $Ab$ sur un chromosome, $aB$ sur l'autre |

### Croisement Test: Déterminer CIS vs TRANS

On croise avec un homozygote récessif:

**CIS ($AB/ab$):**
$$\frac{AB}{ab} \times \frac{ab}{ab}$$
→ Sans crossing-over: 50% AB/ab (A+B+), 50% ab/ab (a+b-)

**TRANS ($Ab/aB$):**
$$\frac{Ab}{aB} \times \frac{ab}{ab}$$
→ Sans crossing-over: 50% Ab/ab (A+b-), 50% aB/ab (a+B-)

**Avec crossing-over:** Apparition de nouvelles combinaisons!

---

## 3.3 Crossing-over (Recombinaison)

### Définition

> [!definition]
> Le **crossing-over** est l'échange de segments entre chromosomes homologues pendant la **prophase I** de la méiose.

### Processus Détaillé

```
PROPHASE I:

Chromosome 1:  A --- B --- C --- D
                    ╲         ╱
Chromosome 2:       a --- b --- c --- d

↓ (enjambement/crossing-over)

Chromosome 1:  A --- b --- c --- D
Chromosome 2:       a --- B --- C --- d
```

### Résultat

| Avant crossing-over | Après crossing-over |
|---------------------|---------------------|
| $AB$ / $ab$ | $Ab$ / $aB$ |

### Types de Gamètes

1. **Gamètes parentaux:** Conservent la combinaison originale
   - $AB$, $ab$ (pour configuration cis)

2. **Gamètes recombinés:** Nouvelle combinaison
   - $Ab$, $aB$ (pour configuration cis)

### Fréquence de Recombinaison

$$p = \frac{\text{nombre de gamètes recombinés}}{\text{nombre total de gamètes}} \times 100\%$$

> [!important]
> La fréquence de recombinaison est égale à la **distance génétique** entre les gènes.

---

## 3.4 Distance Génétique et Carte Génétique

### Unité: centiMorgan (cM)

$$1\% \text{ de recombinaison} = 1 \text{ centiMorgan (cM)} = 1 \text{ map unit}$$

### Calcul de la Distance

**Exemple:**
Sur 1000 gamètes observés, 80 sont recombinés.

$$d = \frac{80}{1000} \times 100\% = 8\% = 8 \text{ cM}$$

### Interprétation

| Distance | Interpretation |
|----------|----------------|
| < 10 cM | Liaison forte |
| 10-25 cM | Liaison modérée |
| > 25 cM | Indépendance apparente |
| 50 cM | Indépendance théorique |

> [!caution] Maximum
> La distance maximale mesurable est de 50 cM. Au-delà, les gènes se comportent comme indépendants.

### Carte Génétique (Carte de Liaison)

```
Chromosome:
[A]----5cM----[B]----3cM----[C]

A et B: 5 cM
B et C: 3 cM
A et C: 8 cM (5 + 3)
```

---

## 3.5 Exercices: Gènes Liés

### Exercice: Distance Génétique

**Énoncé:** Chez la drosophile, deux gènes liés (V et M) sont croisés. Les résultats du croisement test donnent:
- 460 descendants V+M+
- 440 descendants v+m-
- 35 descendants V+m-
- 35 descendants v+M+

**Calculez la distance entre ces gènes.**

**Solution:**

**Total descendants:** 460 + 440 + 35 + 35 = 970

**Recombinants:** 35 + 35 = 70

$$d = \frac{70}{970} \times 100\% = 7.22\%$$

**Distance:** ≈ 7.2 cM

---

# Partie 4: Épistasie

## 4.1 Définition

> [!definition]
> L'**épistasie** est l'interaction entre deux gènes différents où un gène **masque ou modifie l'expression** de l'autre gène.

### Terminologie

- **Gène épistatique:** Le gène qui masque (le "masque")
- **Gène hypostatique:** Le gène qui est masqué (le "masqué")

### Analogie

> [!example]
> L'épistasie, c'est comme si un gène avait un "interrupteur" qui peut empêcher l'autre gène de s'exprimer, même si celui-ci est présent.

---

## 4.2 Types d'Épistasie

### A. Épistasie Réscessive (Ratio 9:3:4)

**Condition:** Le génotype **homozygote récessif** (aa) d'un gène masque l'autre gène.

**Exemple:** Couleur des plumes de poulet

| Gène | Allèles | Effet |
|------|---------|-------|
| E | E-/ee | Extension du pigment |
| C | C-/cc | Présence de pigment |

**Phénotypes:**

| Génotype | Phénotype |
|----------|-----------|
| E- C- | Colour (coloré) |
| E- cc | Blanc (pas de pigment) |
| ee C- | Blanc (pas d'extension) |
| ee cc | Blanc (albinos) |

**Ratio F2:** 9:3:4 (pour 9 Colored, 3 White (récessif), 4 White (épistatique))

> [!note]
> Deux génotypes différents peuvent donner le même phénotype (blanc)!

### B. Épistasie Dominante (Ratio 12:3:1)

**Condition:** La présence d'au moins un **allèle dominant** (A-) du gène épistatique masque l'autre gène.

**Exemple:** Couleur de la souris

| Gène | Allèles | Effet |
|------|---------|-------|
| A | A-/aa | Agouti/Noir |
| C | C-/cc | Albinos |

**Phénotypes:**

| Génotype | Phénotype |
|----------|-----------|
| A- C- | Agouti |
| aa C- | Noir |
| A- cc | Albinos |
| aa cc | Albinos |

**Ratio F2:** 12:3:1

### C. Épistasie Réscessive Duplicate (Ratio 9:7)

**Condition:** Chaque gène récessif homozygote peut, à lui seul, produire le même phénotype (albinos/blanc).

**Exemple:** Couleur des fleurs de pois

| Gène | Allèles | Effet |
|------|---------|-------|
| A | A-/aa | Pigment présent/absent |
| B | B-/bb | Pigment présent/absent |

**Phénotypes:**

| Génotype | Phénotype |
|----------|-----------|
| A- B- | Coloré |
| aa B- | Blanc |
| A- bb | Blanc |
| aa bb | Blanc |

**Ratio F2:** 9:7
(9 coloré + 7 blancs = 3 + 3 + 1)

### D. Épistasie Dominante Duplicate (Ratio 15:1)

**Condition:** La présence d'au moins un allèle dominant de **l'un OU l'autre** gène produit le même phénotype.

**Exemple:** Forme du fruit

| Gène | Allèles | Effet |
|------|---------|-------|
| A | A-/aa | Allongé/Rond |
| B | B-/bb | Allongé/Rond |

**Phénotypes:**

| Génotype | Phénotype |
|----------|-----------|
| A- B- | Allongé |
| A- bb | Allongé |
| aa B- | Allongé |
| aa bb | Rond |

**Ratio F2:** 15:1

### E. Épistasie Complémentaire (Ratio 9:7)

**Condition:** Les **deux gènes doivent avoir** au moins un allèle dominant pour produire le phénotype.

**Exemple:** Synthèse de pigment chez les fleurs

**Mécanisme:**
- Gène A: Enzyme 1 (convertit précurseur en pigment intermédiaire)
- Gène B: Enzyme 2 (convertit pigment intermédiaire en pigment final)

| Génotype | Phénotype |
|----------|-----------|
| A- B- | Coloré |
| aa B- | Blanc |
| A- bb | Blanc |
| aa bb | Blanc |

**Ratio F2:** 9:7

> [!tip]
> L'épistasie complémentaire et réscessive duplicate donnent toutes deux 9:7, mais par des mécanismes différents!

### F. Épistasie à Double Effet (Ratio 13:3)

**Condition:** Un gène dominant produit un phénotype, mais homozygote récessif de l'autre gène le modifie.

---

## 4.3 Tableau Récapitulatif des Épistasies

| Type | Ratio | Gène Épistatique | Exemple |
|------|-------|------------------|---------|
| Réscessive | 9:3:4 | aa | Poulet coloré |
| Dominante | 12:3:1 | A- | Souris agouti |
| Réscessive duplicate | 9:7 | aa ou bb | Pois |
| Dominante duplicate | 15:1 | A- ou B- | Fruits |
| Complémentaire | 9:7 | A- et B- requis | Fleurs |
| Double effet | 13:3 | complexe | - |

---

## 4.4 Exercices: Épistasie

### Exercice

**Énoncé:** Chez les poules, le gène C contrôle la pigmentation (C = pigmenté, cc = albinos). Le gène E contrôle l'extension (E = couleur complète, ee = blanc spotted).

Pour le croisement CcEe × CcEe:
1. Donnez tous les phénotypes possibles
2. Calculez le ratio

**Solution:**

| Génotype | Phénotype |
|----------|-----------|
| C- E- | Coloré |
| C- ee | Blanc spotted |
| cc E- | Albinos |
| cc ee | Albinos |

**Ratio:** 9 Coloré : 3 Blanc spotted : 4 Albinos = **9:3:4**

---

# Partie 5: Calculs de Probabilités (p)

## 5.1 Règles Fondamentales

### A. Règle de Multiplication (ET)

Pour deux événements **indépendants**:

$$P(A \cap B) = P(A) \times P(B)$$

> [!example]
> P(Aa ET bb) = P(Aa) × P(bb) = (1/2) × (1/4) = 1/8

### B. Règle d'Addition (OU)

Pour deux événements **mutuellement exclusifs**:

$$P(A \cup B) = P(A) + P(B)$$

> [!example]
> P(AA ou aa) = P(AA) + P(aa) = 1/4 + 1/4 = 1/2

---

## 5.2 Applications aux Croisements

### Monohybride $Aa \times Aa$

| Événement | Probabilité |
|-----------|-------------|
| AA | 1/4 |
| Aa | 2/4 = 1/2 |
| aa | 1/4 |
| Phénotype dominant (A-) | 3/4 |
| Phénotype récessif (a) | 1/4 |

### Dihybride $AaBb \times AaBb$

| Phénotype | Calcul | Probabilité |
|-----------|--------|-------------|
| A- B- | 3/4 × 3/4 | 9/16 |
| A- bb | 3/4 × 1/4 | 3/16 |
| aa B- | 1/4 × 3/4 | 3/16 |
| aa bb | 1/4 × 1/4 | 1/16 |

### Gènes Liés

Sans crossing-over: gamètes parentaux 50% chacun
Avec distance $d$:
$$P(\text{recombiné}) = \frac{d}{2}$$
$$P(\text{parental}) = \frac{1-d}{2}$$

---

## 5.3 Test du Chi-carré ($\chi^2$)

### Définition

Le test $\chi^2$ permet de vérifier si les résultats expérimentaux correspondent aux attentes théoriques.

$$\chi^2 = \sum \frac{(O - E)^2}{E}$$

Où:
- $O$ = valeur observée
- $E$ = valeur attendue

### Degrés de Liberté

$$ddl = n - 1$$

Où $n$ = nombre de catégories

### Interprétation

| $\chi^2$ calculé | $\chi^2$ таблица | Décision |
|-----------------|-----------------|----------|
| < | > | **Accepter** hypothèse |
| > | < | **Rejeter** hypothèse |

> [!warning]
> Si $\chi^2_{calculé} > \chi^2_{table}$, on rejette l'hypothèse d'indépendance!

---

# Partie 6: Fiches Résumées

## Fiche 1: Ratio des Croisements

| Type | Croisement | Ratio F2 |
|------|------------|----------|
| Monohybride complet | Aa × Aa | 3:1 |
| Monohybride incomplet | idem | 1:2:1 |
| Codominance | idem | 1:2:1 |
| Dihybride indépendant | AaBb × AaBb | 9:3:3:1 |
| Épistasie réscessive | - | 9:3:4 |
| Épistasie dominante | - | 12:3:1 |
| Épistasie duplicate | - | 9:7 ou 15:1 |

## Fiche 2: Formule Distance Génétique

$$d = \frac{N_{\text{recombinants}}}{N_{\text{total}}} \times 100\%$$

## Fiche 3: Formule Général pour n Gènes

$$\text{Ratio} = (3:1)^n$$

---

# Annexes

## Vidéos et Animations

Voir [[Genetique-Graph-Complet.canvas]] pour les animations visuelles.

## Exercices Complémentaires

Voir [[02-Practice/Biology/Genetics-Exercises]] pour plus de problèmes.

---

*Notes créées avec Obsidian*
*Liens liés: [[Genetics-Basics]] | [[02-Practice/Biology/Genetics-Exercises]] | [[00-Meta/MOCs/Natural Sciences MOC]]*
