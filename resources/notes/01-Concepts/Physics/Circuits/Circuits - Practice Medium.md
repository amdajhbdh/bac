---
tags: [practice, circuits, medium]
---

# Circuits Practice - Medium Difficulty

## Exercices d'application sur les circuits électriques

### Exercice 1: Diviseur de tension
Un diviseur de tension est constitué de deux résistances $R_1 = 4 \Omega$ et $R_2 = 6 \Omega$ en série alimentées par une tension $U = 20 \text{ V}$.

1. Quelle est la tension aux bornes de $R_2$ ?
2. Si on veut obtenir une tension de $U_2 = 8 \text{ V}$ aux bornes de $R_2$ avec $R_1 = 4 \Omega$, quelle devrait être la valeur de $R_2$ ?
3. Quelle est la puissance consommée par chaque résistance ?

### Exercice 2: Pont de Wheatstone
Un pont de Wheatstone est constitué de quatre résistances : $R_1 = 100 \Omega$, $R_2 = 200 \Omega$, $R_3 = 150 \Omega$, et $R_x$ (inconnue). Le pont est équilibré lorsque le galvanomètre montre zéro courant.

1. Quelle est la valeur de $R_x$ pour que le pont soit équilibré ?
2. Si $R_3$ varie entre $100 \Omega$ et $200 \Omega$, quel est l'intervalle de valeurs possibles pour $R_x$ lorsque le pont reste équilibré ?
3. Quelle est la sensibilité du montage, définie comme $\frac{\Delta R_x}{\Delta R_3}$ ?

### Exercice 3: Circuit avec source réelle
Une batterie fournit une tension à vide $U_0 = 12 \text{ V}$ et possède une résistance interne $r = 0.5 \Omega$. Elle est connectée à une charge externe $R$.

1. Exprimer la tension aux bornes de la charge $U$ en fonction de $R$.
2. Quelle valeur de $R$ maximise la puissance dissipée dans la charge ?
3. Quelle est cette puissance maximale ?
4. Quelle est l'efficacité du transfert d'énergie à cette condition ?

## Corrigés

### Exercice 1
1. Dans un diviseur de tension : $U_2 = U \frac{R_2}{R_1 + R_2} = 20 \times \frac{6}{4+6} = 20 \times \frac{6}{10} = 20 \times 0.6 = 12 \text{ V}$
   
2. $U_2 = U \frac{R_2}{R_1 + R_2}$
   $8 = 20 \times \frac{R_2}{4 + R_2}$
   $\frac{8}{20} = \frac{R_2}{4 + R_2}$
   $0.4 = \frac{R_2}{4 + R_2}$
   $0.4(4 + R_2) = R_2$
   $1.6 + 0.4R_2 = R_2$
   $1.6 = 0.6R_2$
   $R_2 = \frac{1.6}{0.6} = 2.67 \Omega$
   
3. Courant dans le circuit : $I = \frac{U}{R_1 + R_2} = \frac{20}{4+6} = 2 \text{ A}$
   Puissance dans $R_1$ : $P_1 = I^2R_1 = 4 \times 4 = 16 \text{ W}$
   Puissance dans $R_2$ : $P_2 = I^2R_2 = 4 \times 6 = 24 \text{ W}$
   Puissance totale : $P_{tot} = UI = 20 \times 2 = 40 \text{ W}$ (vérification : $16 + 24 = 40$)

### Exercice 2
1. Pour un pont de Wheatstone équilibré : $\frac{R_1}{R_2} = \frac{R_3}{R_x}$
   Donc : $R_x = R_3 \frac{R_2}{R_1} = 150 \times \frac{200}{100} = 150 \times 2 = 300 \Omega$
   
2. Lorsque le pont est équilibré : $R_x = R_3 \frac{R_2}{R_1}$
   Avec $R_1 = 100 \Omega$ et $R_2 = 200 \Omega$ fixes :
   $R_x = R_3 \times \frac{200}{100} = 2R_3$
   
   Si $R_3 \in [100, 200] \Omega$, alors $R_x \in [2 \times 100, 2 \times 200] = [200, 400] \Omega$
   
3. Sensibilité : $\frac{\Delta R_x}{\Delta R_3} = \frac{R_2}{R_1} = \frac{200}{100} = 2$
   Cela signifie qu'une variation de $1 \Omega$ de $R_3$ provoque une variation de $2 \Omega$ de $R_x$ pour maintenir l'équilibre.

### Exercice 3
1. Le modèle d'une source réelle est une tension idéale $U_0$ en série avec une résistance interne $r$.
   Le courant dans le circuit est : $I = \frac{U_0}{R + r}$
   La tension aux bornes de la charge est : $U = U_0 - Ir = U_0 - \frac{U_0 r}{R + r} = U_0\left(1 - \frac{r}{R + r}\right) = U_0 \frac{R}{R + r}$
   
   Donc : $U = U_0 \frac{R}{R + r} = 12 \frac{R}{R + 0.5}$
   
2. La puissance dissipée dans la charge est : $P = UI = \left(U_0 \frac{R}{R + r}\right) \left(\frac{U_0}{R + r}\right) = \frac{U_0^2 R}{(R + r)^2}$
   
   Pour maximiser $P$ par rapport à $R$, on dérivée et on égale à zéro :
   $\frac{dP}{dR} = U_0^2 \frac{(R + r)^2 - R \cdot 2(R + r)}{(R + r)^4} = U_0^2 \frac{(R + r)[(R + r) - 2R]}{(R + r)^4} = U_0^2 \frac{r - R}{(R + r)^3}$
   
   $\frac{dP}{dR} = 0$ lorsque $r - R = 0$, donc $R = r = 0.5 \Omega$
   
3. $P_{max} = \frac{U_0^2 r}{(r + r)^2} = \frac{U_0^2 r}{4r^2} = \frac{U_0^2}{4r} = \frac{12^2}{4 \times 0.5} = \frac{144}{2} = 72 \text{ W}$
   
4. L'efficacité est définie comme : $\eta = \frac{P_{charge}}{P_{source}} = \frac{UI}{U_0 I} = \frac{U}{U_0}$
   
   À la condition de puissance maximale : $U = U_0 \frac{R}{R + r} = U_0 \frac{r}{r + r} = U_0 \frac{1}{2} = \frac{U_0}{2}$
   
   Donc : $\eta = \frac{U_0/2}{U_0} = \frac{1}{2} = 50\%$
   
   Cela signifie que la moitié de l'énergie produite par la source est dissipée dans sa résistance interne lorsque la puissance transmise à la charge est maximale.

---
*Source: Adapted from physics problem sets*