# Graphiques Visuels - Génétique

## Diagramme 1: Arbre des Cas de Croisements

```mermaid
flowchart TD
    A[Croisements<br/>Génétiques] --> B[Monohybridisme<br/>1 Gène]
    A --> C[Dihybridisme<br/>2 Gènes]
    A --> D[Trihybridisme<br/>3 Gènes]
    
    B --> B1[Dominance Complète<br/>3:1]
    B --> B2[Dom. Incomplète<br/>1:2:1]
    B --> B3[Codominance<br/>1:2:1]
    B --> B4[Allèles Multiples<br/>ABO]
    B --> B5[Allèles Léthaux<br/>2:1]
    
    C --> C1[Indépendance<br/>9:3:3:1]
    C --> C2[Liaison Génétique]
    C --> C3[Épistasie]
    
    C2 --> C2a[Liaison Absolue<br/>0 cM → 1:2:1]
    C2 --> C2b[Liaison Partielle<br/>1-50 cM → modifié]
    C2 --> C2c[Indépendance<br/>>50 cM → 9:3:3:1]
    
    C3 --> E1[Récessive<br/>9:3:4]
    C3 --> E2[Dominante<br/>12:3:1]
    C3 --> E3[Réss. Duplicate<br/>9:7]
    C3 --> E4[Dom. Duplicate<br/>15:1]
    C3 --> E5[Complémentaire<br/>9:7]
    C3 --> E6[Double<br/>13:3]
    
    D --> D1[(3:1)³<br/>27:9:9:9:3:3:3:1]
    
    style A fill:#e0f2fe,stroke:#0284c7
    style B fill:#fce7f3,stroke:#db2777
    style C fill:#dcfce7,stroke:#16a34a
    style C3 fill:#fef3c7,stroke:#d97706
```

---

## Diagramme 2: Types de Liaison

```mermaid
flowchart LR
    subgraph LIAISON["LIAISON GÉNÉTIQUE"]
        direction TB
        LA[Liaison Absolue<br/>d = 0 cM<br/>0% recombinants] --> LTF[Liaison Très Forte<br/>d = 1-5 cM<br/>1-5% recombinants]
        LTF --> LM[Liaison Modérée<br/>d = 5-25 cM<br/>5-25% recombinants]
        LM --> LF[Liaison Faible<br/>d = 25-50 cM<br/>25-50% recombinants]
        LF --> LI[Indépendance<br/>d > 50 cM<br/>50% recombinants]
    end
    
    LA -->|"Ratio F1×F1"| R1["1:2:1<br/>(comme mono)"]
    LTF -->|"Ratio F1×F1"| R2["Proche 1:2:1"]
    LM -->|"Ratio F1×F1"| R3["9:3:3:1<br/>modifié"]
    LF -->|"Ratio F1×F1"| R4["Proche 9:3:3:1"]
    LI -->|"Ratio F1×F1"| R5["9:3:3:1"]
    
    style LIAISON fill:#f3e8ff,stroke:#9333ea
    style R1 fill:#fee2e2,stroke:#dc2626
    style R5 fill:#d1fae5,stroke:#059669
```

---

## Diagramme 3: Configuration CIS vs TRANS

```mermaid
flowchart LR
    subgraph CROSS["CROSSING-OVER"]
        direction TB
        
        subgraph CIS["Configuration CIS"]
            direction LR
            C1[(Chromosome 1)] -->|"AB"| C2[(Chromosome 2)]
            C2 -->|"ab"| C3[(Chromosome 1)]
        end
        
        subgraph TRANS["Configuration TRANS"]
            direction LR
            T1[(Chromosome 1)] -->|"Ab"| T2[(Chromosome 2)]
            T2 -->|"aB"| T3[(Chromosome 1)]
        end
    end
    
    subgraph RESULTATS["Résultats du Croisement Test"]
        direction TB
        R1["CIS × ab/ab<br/>→ Sans CO: AB, ab<br/>→ Avec CO: Ab, aB"]
        R2["TRANS × ab/ab<br/>→ Sans CO: Ab, aB<br/>→ Avec CO: AB, ab"]
    end
    
    CROSS -->|"Crossing-over"| RESULTATS
    
    style CIS fill:#dbeafe,stroke:#2563eb
    style TRANS fill:#ffedd5,stroke:#ea580c
    style RESULTATS fill:#f0fdf4,stroke:#16a34a
```

---

## Diagramme 4: Épistasie

```mermaid
flowchart TD
    E[Épistasie<br/>Interaction Génique] --> E1[Résessive<br/>aa masque B<br/>Ratio: 9:3:4]
    E --> E2[Dominante<br/>A- masque b<br/>Ratio: 12:3:1]
    E --> E3[Réscessive Duplicate<br/>aa OU bb → même phén.<br/>Ratio: 9:7]
    E --> E4[Dominante Duplicate<br/>A- OU B- → même phén.<br/>Ratio: 15:1]
    E --> E5[Complémentaire<br/>A- ET B- requis<br/>Ratio: 9:7]
    E --> E6[Double<br/>Interaction complexe<br/>Ratio: 13:3]
    
    E1 -->|"Exemple"| Ex1[Poulet:<br/>Couleur + Plumes]
    E2 -->|"Exemple"| Ex2[Souris:<br/>Agouti + Albinos]
    E3 -->|"Exemple"| Ex3[Pois:<br/>Fleurs colorées]
    E4 -->|"Exemple"| Ex4[Fruits:<br/>Forme]
    E5 -->|"Exemple"| Ex5[Fleurs:<br/>Pigment]
    
    style E fill:#fef3c7,stroke:#d97706
    style E1 fill:#fee2e2,stroke:#dc2626
    style E2 fill:#fee2e2,stroke:#dc2626
    style E3 fill:#fee2e2,stroke:#dc2626
    style E4 fill:#fee2e2,stroke:#dc2626
    style E5 fill:#fee2e2,stroke:#dc2626
    style E6 fill:#fee2e2,stroke:#dc2626
```

---

## Diagramme 5: F1 × F1 selon Configuration

```mermaid
flowchart TD
    START[Croisement F1 × F1] --> Q1{1 ou 2<br/>caractères?}
    
    Q1 -->|1 caractère| M1{Type de<br/>dominance?}
    M1 -->|Complète| MR1[3:1]
    M1 -->|Incomplète| MR2[1:2:1]
    M1 -->|Codominance| MR3[1:2:1]
    
    Q1 -->|2 caractères| Q2{Gènes liés<br/>ou indépendants?}
    
    Q2 -->|Indépendants| DR1[9:3:3:1]
    Q2 -->|Liés| Q3{Distance?}
    
    Q3 -->|d = 0 cM| LR1[1:2:1<br/>Liaison absolue]
    Q3 -->|1-50 cM| LR2[Ratio modifié<br/>Liaison partielle]
    Q3 -->|>50 cM| LR3[9:3:3:1<br/>Indépendance]
    
    style START fill:#e0f2fe,stroke:#0284c7
    style MR1 fill:#dcfce7,stroke:#16a34a
    style DR1 fill:#dcfce7,stroke:#16a34a
    style LR1 fill:#fee2e2,stroke:#dc2626
```

---

## Diagramme 6: Calculs de Probabilités

```mermaid
flowchart TD
    P[Probabilités] --> P1[Monohybridisme<br/>Aa × Aa]
    P --> P2[Dihybridisme<br/>AaBb × AaBb]
    P --> P3[Gènes Liés]
    
    P1 --> P1a[P(dominant) = 3/4]
    P1 --> P1b[P(récessif) = 1/4]
    P1 --> P1c[Ratio: 3:1]
    
    P2 --> P2a[P(A- B-) = 3/4 × 3/4 = 9/16]
    P2 --> P2b[P(A- bb) = 3/4 × 1/4 = 3/16]
    P2 --> P2c[P(aa B-) = 1/4 × 3/4 = 3/16]
    P2 --> P2d[P(aabb) = 1/4 × 1/4 = 1/16]
    P2 --> P2e[Ratio: 9:3:3:1]
    
    P3 --> P3a[P(parental) = (1-d)/2]
    P3 --> P3b[P(recombinant) = d/2]
    P3 --> P3c[d = distance en fraction]
    
    style P fill:#e0f2fe,stroke:#0284c7
    style P1 fill:#fce7f3,stroke:#db2777
    style P2 fill:#dcfce7,stroke:#16a34a
    style P3 fill:#fef3c7,stroke:#d97706
```

---

## Tableau Récapitulatif des Ratios

| Section | Cas | Ratio F2 | Formule |
|---------|-----|----------|---------|
| **A** | Monohybridisme complet | 3:1 | - |
| **A** | Monohybridisme incomplet | 1:2:1 | - |
| **A** | Codominance | 1:2:1 | - |
| **A** | Allèles léthaux | 2:1 | - |
| **B** | Dihybridisme indépendant | 9:3:3:1 | (3:1)² |
| **B** | Liaison absolue | 1:2:1 | - |
| **B** | Liaison partielle | Variable | (d)/2 |
| **B** | Épistasie récessive | 9:3:4 | - |
| **B** | Épistasie dominante | 12:3:1 | - |
| **B** | Épistasie réss. duplicate | 9:7 | - |
| **B** | Épistasie dom. duplicate | 15:1 | - |
| **B** | Épistasie complémentaire | 9:7 | - |
| **B** | Épistasie double | 13:3 | - |
| **D** | Trihybridisme | 27:9:9:9:3:3:3:1 | (3:1)³ |

---

## Formule Générale

Pour **n gènes indépendants** avec dominance complète:

$$\text{Ratio} = (3:1)^n$$

| n | Ratio |
|---|-------|
| 1 | 3:1 |
| 2 | 9:3:3:1 |
| 3 | 27:9:9:9:3:3:3:1 |
| 4 | 81:27:27:27:9:9:9:9:3:3:3:3:1 |

---

*Voir aussi: [[Genetique-Catalogue-Complet]] | [[Genetique-Graph-Complet.canvas]]*
