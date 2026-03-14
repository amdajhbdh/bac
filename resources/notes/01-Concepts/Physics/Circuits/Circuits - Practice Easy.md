---
tags: [practice, circuits, easy]
---

# Circuits Practice - Easy Difficulty

## Exercices d'application sur les circuits électriques de base

### Exercice 1: Loi d'Ohm simple
Un résistor de résistance $R = 10 \Omega$ est soumis à une tension $U = 5 \text{ V}$.

1. Quel est le courant qui traverse le résistor ?
2. Si la tension est doublée, que devient le courant ?
3. Quelle est la puissance dissipée par le résistor ?

### Exercice 2: Résistances en série
Trois résistances $R_1 = 2 \Omega$, $R_2 = 3 \Omega$, et $R_3 = 5 \Omega$ sont connectées en série à une batterie de $U = 10 \text{ V}$.

1. Quelle est la résistance équivalente du montage ?
2. Quel est le courant qui traverse le circuit ?
3. Quelle est la tension aux bornes de chaque résistance ?

### Exercice 3: Résistances en parallèle
Deux résistances $R_1 = 6 \Omega$ et $R_2 = 3 \Omega$ sont connectées en parallèle à une batterie de $U = 12 \text{ V}$.

1. Quelle est la résistance équivalente du montage ?
2. Quel est le courant total fourni par la batterie ?
3. Quel est le courant qui traverse chaque branche ?

## Corrigés

### Exercice 1
1. $I = \frac{U}{R} = \frac{5 \text{ V}}{10 \Omega} = 0.5 \text{ A}$
2. Si $U' = 2U = 10 \text{ V}$, alors $I' = \frac{U'}{R} = \frac{10}{10} = 1 \text{ A}$ (le courant double aussi)
3. $P = UI = 5 \times 0.5 = 2.5 \text{ W}$ ou $P = \frac{U^2}{R} = \frac{25}{10} = 2.5 \text{ W}$ ou $P = I^2R = 0.25 \times 10 = 2.5 \text{ W}$

### Exercice 2
1. $R_{eq} = R_1 + R_2 + R_3 = 2 + 3 + 5 = 10 \Omega$
2. $I = \frac{U}{R_{eq}} = \frac{10}{10} = 1 \text{ A}$
3. $U_1 = IR_1 = 1 \times 2 = 2 \text{ V}$
   $U_2 = IR_2 = 1 \times 3 = 3 \text{ V}$
   $U_3 = IR_3 = 1 \times 5 = 5 \text{ V}$
   Vérification: $U_1 + U_2 + U_3 = 2 + 3 + 5 = 10 \text{ V} = U$

### Exercice 3
1. $\frac{1}{R_{eq}} = \frac{1}{R_1} + \frac{1}{R_2} = \frac{1}{6} + \frac{1}{3} = \frac{1}{6} + \frac{2}{6} = \frac{3}{6} = \frac{1}{2}$
   Donc $R_{eq} = 2 \Omega$
2. $I_{total} = \frac{U}{R_{eq}} = \frac{12}{2} = 6 \text{ A}$
3. $I_1 = \frac{U}{R_1} = \frac{12}{6} = 2 \text{ A}$
   $I_2 = \frac{U}{R_2} = \frac{12}{3} = 4 \text{ A}$
   Vérification: $I_1 + I_2 = 2 + 4 = 6 \text{ A} = I_{total}$

---
*Source: Adapted from physics problem sets*