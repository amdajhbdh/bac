---
title: Génétique - Catalogue Complet des Cas de Croisements
niveau: 7ème/Terminale
matière: Biology
date: 2026-03-13
tags:
  - génétique
  - croisements
  - catalogue
  - complet
  - tous-les-cas
cssclasses:
  - tables
---

# Catalogue Complet des Cas de Croisements Génétiques

> **Source:** Résumé exhaustif pour BAC Mauritanie
> **Dernière mise à jour:** 2026-03-13

---

## Vue d'Ensemble: Tous les Cas

```mermaid
graph TD
    A[Croisements Génétiques] --> B[Monohybridisme<br/>(1 gène)]
    A --> C[Dihybridisme<br/>(2 gènes)]
    A --> D[Trihybridisme<br/>(3 gènes)]
    
    B --> B1[Dominance Complète]
    B --> B2[Dominance Incomplète]
    B --> B3[Codominance]
    B --> B4[Allèles Multiples]
    B --> B5[Allèles Léthaux]
    
    C --> C1[Indépendance]
    C --> C2[Liaison Absolue]
    C --> C3[Liaison Partielle]
    C --> C4[Épistasie]
    
    D --> D1[Ratio 27:9:9:9:3:3:3:1]
```

---

# PARTIE A: MONOHYBRIDISME (1 Gène)

## A.1 Types de Dominance

| Cas       | Description              | Ratio F2 | Exemple             |
| --------- | ------------------------ | -------- | ------------------- |
| **A.1.1** | Dominance complète       | 3:1      | Pois jaune/vert     |
| **A.1.2** | Dominance incomplète     | 1:2:1    | Muflier rouge/blanc |
| **A.1.3** | Codominance              | 1:2:1    | Groupes sanguins AB |
| **A.1.4** | Surdominance (hétérosis) | 1:2:1    | Hémoglobine HbS     |
| **A.1.5** | Sous-domination          | 1:2:1    | Certaines plantes   |

## A.2 Allèles Multiples

| Cas       | Description      | Nombre d'Allèles | Exemple              |
| --------- | ---------------- | ---------------- | -------------------- |
| **A.2.1** | 2 allèles (pair) | 2                | Majorité des gènes   |
| **A.2.2** | 3 allèles        | 3                | Groupes sanguins ABO |
| **A.2.3** | 4+ allèles       | 4+               | Gènes HLA            |

## A.3 Allèles Léthaux

| Cas       | Description     | Ratio F2 | Exemple           |
| --------- | --------------- | -------- | ----------------- |
| **A.3.1** | Récessif léthal | 2:1      | Souris jaune (Ay) |
| **A.3.2** | Dominant léthal | 1:0      | Choléra           |

---

# PARTIE B: DIHYBRIDISME (2 Gènes)

## B.1 Indépendance (Loi d'Assortiment Indépendant)

### B.1.1 Croisement Standard: AaBb × AaBb

| Cas | Description | Ratio F2 | Formule |
|-----|-------------|----------|---------|
| **B.1.1.1** | Dominance complète, 2 gènes | **9:3:3:1** | (3:1)² |
| **B.1.1.2** | Incomplète + Complète | 3:6:3:4 | - |
| **B.1.1.3** | Codominance + Complète | Variables | - |

### B.1.2 Croisement Test: AaBb × aabb

| Cas | Ratio Observé | Interprétation |
|-----|--------------|----------------|
| **B.1.2.1** | 1:1:1:1 | Indépendance (gènes sur chromosomes différents) |
| **B.1.2.2** | ≠1:1:1:1 | Liaison entre gènes |

---

## B.2 Liaison Génétique (Linkage)

### B.2.1 Types de Liaison selon Distance

| Cas | Distance | % Recombinants | Ratio F1×F1 | Formule Gamètes |
|-----|----------|----------------|-------------|-----------------|
| **B.2.1.1** | Liaison absolue | 0 cM | 0% | 1:2:1 | Parentaux: 100% |
| **B.2.1.2** | Liaison très forte | 1-5 cM | 1-5% | Proche 1:2:1 | Parentaux: 95-99% |
| **B.2.1.3** | Liaison modérée | 5-25 cM | 5-25% | 9:3:3:1 modifié | Parentaux: 75-95% |
| **B.2.1.4** | Liaison faible | 25-50 cM | 25-50% | Proche 9:3:3:1 | Parentaux: 50-75% |
| **B.2.1.5** | Indépendance | >50 cM | 50% | **9:3:3:1** | Tous: 25% |

### B.2.2 Configuration CIS vs TRANS

| Cas | Configuration | Description | Croisement Test × ab/ab |
|-----|---------------|-------------|------------------------|
| **B.2.2.1** | CIS | $AB/ab$ | 50% AB + 50% ab (sans CO) |
| **B.2.2.2** | TRANS | $Ab/aB$ | 50% Ab + 50% aB (sans CO) |
| **B.2.2.3** | CIS + CO | $AB/ab$ + crossing | Apparition Ab et aB |
| **B.2.2.4** | TRANS + CO | $Ab/aB$ + crossing | Apparition AB et ab |

---

## B.3 Épistasie (Interactions Géniques)

### B.3.1 Tableau Complet des Types d'Épistasie

| Cas | Type | Ratio F2 | Mécanisme | Exemple |
|-----|------|----------|-----------|---------|
| **B.3.1.1** | Épistasie récessive | **9:3:4** | aa masque B | Poulet coloré |
| **B.3.1.2** | Épistasie dominante | **12:3:1** | A- masque bb | Souris agouti |
| **B.3.1.3** | Épistasie récessive duplicate | **9:7** | aa ou bb → même phénotype | Pois |
| **B.3.1.4** | Épistasie dominante duplicate | **15:1** | A- ou B- → même phénotype | Fruits |
| **B.3.1.5** | Épistasie complémentaire | **9:7** | A- et B- requis | Fleurs |
| **B.3.1.6** | Épistasie additive | **9:3:3:1** | Effets additifs | - |
| **B.3.1.7** | Épistasie double | **13:3** | Interaction complexe | - |
| **B.3.1.8** | Épistasie multiplicative | Variables | Effets multipliés | - |

### B.3.2 Épistasie Détaillée

#### B.3.2.1 Épistasie Réscessive (9:3:4)

```
Gène 1: E (extension du pigment) - E- = pigmenté, ee = pas pigment
Gène 2: C (présence pigment) - C- = couleur, cc = albinos

Génotypes → Phénotypes:
E- C- → Coloré (9)
E- cc → Blanc (3)
ee C- → Blanc (4)
ee cc → Blanc (4)
```

#### B.3.2.2 Épistasie Dominante (12:3:1)

```
Gène 1: A (couleur) - A- = agouti, aa = noir
Gène 2: C (pigment) - C- = pigmenté, cc = albinos

Génotypes → Phénotypes:
A- C- → Agouti (12)
aa C- → Noir (3)
A- cc → Albinos (3+1)
aa cc → Albinos
```

#### B.3.2.3 Épistasie Réscessive Duplicate (9:7)

```
Gène 1: A - A- = pigment, aa = pas pigment
Gène 2: B - B- = pigment, bb = pas pigment

Seul A- ET B- ensemble produisent du pigment!

Génotypes → Phénotypes:
A- B- → Coloré (9)
aa B- → Blanc (3)
A- bb → Blanc (3)
aa bb → Blanc (1)
Total blanc: 7
```

#### B.3.2.4 Épistasie Dominante Duplicate (15:1)

```
Gène 1: A - A- = forme allongée
Gène 2: B - B- = forme allongée

Soit A- OU B- donne la forme allongée!

Génotypes → Phénotypes:
A- B- → Allongé (9)
A- bb → Allongé (3)
aa B- → Allongé (3)
aa bb → Rond (1)
```

---

# PARTIE C: F1 × F1 - TOUS LES CAS

## C.1 Monohybridisme F1 × F1

| Cas | Croisement | Ratio Génotypique | Ratio Phénotypique |
|-----|-----------|-------------------|-------------------|
| **C.1.1** | Aa × Aa (complet) | 1:2:1 | 3:1 |
| **C.1.2** | Aa × Aa (incomplet) | 1:2:1 | 1:2:1 |
| **C.1.3** | Aa × Aa (codominance) | 1:2:1 | 1:2:1 |

## C.2 Dihybridisme F1 × F1 - Selon Type de Liaison

| Cas | Situation | Ratio Phénotypique | Condition |
|-----|-----------|-------------------|-----------|
| **C.2.1** | Indépendance | **9:3:3:1** | Chromosomes différents ou >50cM |
| **C.2.2** | Liaison partielle | Variable (1-50cM) | 1-50% recombinants |
| **C.2.3** | Liaison absolue | **1:2:1** | 0% recombinants |

### C.2.1 Dihybridisme Indépendant: AaBb × AaBb

```
9/16 = A- B- (dominant pour les deux)
3/16 = A- bb (dominant A, récessif b)
3/16 = aa B- (récessif a, dominant B)
1/16 = aa bb (récessif pour les deux)
```

### C.2.2 Liaison Absolue: (AB/ab) × (AB/ab)

```
Sans recombinaison:
Gamètes possibles: seulement AB et ab

Descendants:
1/4 = AABB → A- B-
1/2 = AaBb → A- B-
1/4 = aabb → aa bb

Ratio: 3:1 (comme monohybridisme!)
```

### C.2.3 Liaison Partielle: d = 20%

```
Pour chaque parent AB/ab:
- AB: (1-d)/2 = 0.4
- ab: (1-d)/2 = 0.4
- Ab: d/2 = 0.1
- aB: d/2 = 0.1

Grille de croisement 4×4 donne:
A- B-: 0.36 = 36%
A- bb: 0.12 = 12%
aa B-: 0.12 = 12%
aa bb: 0.16 = 16%

Ratio observé: 9:3:3:4 (proche de 9:3:3:1)
```

## C.3 Croisement Test F1 × aabb

| Cas | Situation | Ratio Observé | Interprétation |
|-----|-----------|----------------|----------------|
| **C.3.1** | Indépendance | 1:1:1:1 | Gènes indépendants |
| **C.3.2** | Liaison | ≠1:1:1:1 | Gènes liés |

---

# PARTIE D: CROSSING-OVER

## D.1 Fréquence de Recombinaison

$$p = \frac{\text{gamètes recombinés}}{\text{total gamètes}} \times 100\%$$

| Cas | Distance (cM) | Fréquence | Dénomination |
|-----|---------------|-----------|--------------|
| **D.1.1** | 0 cM | 0% | Liaison absolue |
| **D.1.2** | 1-10 cM | 1-10% | Liaison forte |
| **D.1.3** | 10-25 cM | 10-25% | Liaison modérée |
| **D.1.4** | 25-50 cM | 25-50% | Liaison faible |
| **D.1.5** | >50 cM | >50% | Indépendance |

## D.2 Calcul de Distance Génétique

### Formule

$$d = \frac{N_{\text{recombinants}}}{N_{\text{total}}} \times 100\%$$

### Exemples

| Nombre Total | Recombinants | Distance |
|--------------|--------------|----------|
| 1000 | 50 | 5 cM |
| 1000 | 100 | 10 cM |
| 1000 | 250 | 25 cM |
| 1000 | 500 | 50 cM |

---

# PARTIE E: TABLEAU RÉCAPITULATIF FINAL

## E.1 Tous les Ratios Possibles

| Ratio F2 | Cas Correspondants | Section |
|----------|---------------------|---------|
| **1:1** | Croisement test monohybrid | C.1 |
| **1:2:1** | Liaison absolue F1×F1,codominance | C.2.2, A.1.2 |
| **3:1** | Dominance complète monohybrid | A.1.1 |
| **3:1** | Liaison absolue (phénotype) | C.2.2 |
| **9:3:3:1** | Indépendance dihybrid | B.1 |
| **9:3:4** | Épistasie récessive | B.3.1.1 |
| **12:3:1** | Épistasie dominante | B.3.1.2 |
| **9:7** | Épistasie récessive duplicate | B.3.1.3 |
| **15:1** | Épistasie dominante duplicate | B.3.1.4 |
| **13:3** | Épistasie double | B.3.1.7 |
| **27:9:9:9:3:3:3:1** | Trihybridisme | D.1 |

## E.2 Arbre de Décision

```
Quel type de croisement?
│
├─ 1 seul caractère → MONOHYBRIDISME
│   │
│   ├─ Dominance complète → Ratio 3:1
│   ├─ Dominance incomplète → Ratio 1:2:1
│   └─ Codominance → Ratio 1:2:1
│
├─ 2 caractères → DIHYBRIDISME
│   │
│   ├─ Gènes SUR chromosomes DIFFÉRENTS → Ratio 9:3:3:1
│   │
│   ├─ Gènes LIÉS sur même chromosome
│   │   │
│   │   ├─ Distance = 0 cM → Liaison ABSOLUE → Ratio 1:2:1
│   │   ├─ Distance 1-50 cM → Liaison PARTIELLE → Ratio modifié
│   │   └─ Distance > 50 cM → Ratio 9:3:3:1
│   │
│   └─ Interaction (ÉPISTASIE)
│       │
│       ├─ aa masque B → Ratio 9:3:4
│       ├─ A- masque bb → Ratio 12:3:1
│       ├─ aa OU bb → Ratio 9:7
│       └─ A- OU B- → Ratio 15:1
│
└─ 3 caractères → TRIHYBRIDISME → Ratio (3:1)³ = 27:9:9:9:3:3:3:1
```

---

# PARTIE F: FORMULES UTILES

## F.1 Probabilités pour Gènes Indépendants

| Situation | Formule |
|-----------|---------|
| Monohybrid Aa×Aa | P(dominant) = 3/4, P(récessif) = 1/4 |
| Dihybrid | P(phénotype) = P(caractère1) × P(caractère2) |
| n gènes | Ratio = $(3:1)^n$ |

## F.2 Probabilités pour Gènes Liés

| Situation | Formule |
|-----------|---------|
| Gamète parental | $P = \frac{1-d}{2}$ |
| Gamète recombinant | $P = \frac{d}{2}$ |

Où $d$ = distance en fraction (ex: 0.2 pour 20%)

---

# ANNEXE: CAS SPÉCIAUX

## G.1 Gènes Pléiotropes

Un gène affecte plusieurs caractères.

## G.2 Hérédités Quantitatives (Polygeny)

Plusieurs gènes affectent un même caractère.

## G.3 Interaction Gène-Environnement

Le phénotype dépend du génotype ET de l'environnement.

## G.4 Imprinting Génétique

Expression différente selon origine parentale.

## G.5 Hérédité Mitochondriale

Transmission exclusively maternelle.

---

*Document complet - Tous les cas de croisements génétiques*
*Voir aussi: [[Genetique-Croisements]] | [[Genetique-Graph-Complet.canvas]]*
