---
tags: [electromagnetism, magnetic-field]
---

# Champ Magnétique

## Définition
Un champ magnétique est un champ physique qui décrit l'influence magnétique sur les charges électriques en mouvement, les courants électriques et les matériaux magnétiques. Il est représenté par le vecteur $\vec{B}$ et mesuré en teslas (T).

## Propriétés du champ magnétique
1. **Formé par des mouvements de charges** : Les courants électriques et les spins électroniques créent des champs magnétiques
2. **Lignes de champ** : Comme le champ électrique, on représente le champ magnétique par des lignes de champ qui :
   - Sont fermées (elles commencent et finissent sur le pôle nord d'un aimant)
   - Ne se croisent jamais
   - Sont plus denses là où le champ est plus intense
3. **Force de Lorentz** : Une charge $q$ se déplaçant avec un vitesse $\vec{v}$ dans un champ magnétique $\vec{B}$ subit une force :
   $$\vec{F} = q\vec{v} \times \vec{B}$$
   Cette force est toujours perpendiculaire à la fois à $\vec{v}$ et à $\vec{B}$

## Loi de Biot-Savart
La loi de Biot-Savart permet de calculer le champ magnétique créé par un élément de courant électrique :

$$d\vec{B} = \frac{\mu_0}{4\pi} \frac{I d\vec{l} \times \hat{r}}{r^2}$$

Où :
- $d\vec{B}$ est l'élément de champ magnétique créé par $Id\vec{l}$
- $\mu_0 = 4\pi \times 10^{-7} \text{ T·m/A}$ est la perméabilité du vide
- $I$ est le courant en ampères (A)
- $d\vec{l}$ est l'élément de longueur du conducteur dans la direction du courant
- $\hat{r}$ est le vecteur unité pointant de l'élément de courant vers le point d'observation
- $r$ est la distance entre l'élément de courant et le point d'observation

## Applications de la loi de Biot-Savart

### Fil conducteur droit et infini
Pour un fil droit et infini portant un courant $I$, à une distance $r$ du fil :
$$B = \frac{\mu_0 I}{2\pi r}$$

### Boucle circulaire de rayon $R$
Au centre d'une boucle circulaire de rayon $R$ portant un courant $I$ :
$$B = \frac{\mu_0 I}{2R}$$

Sur l'axe de la boucle, à une distance $x$ du centre :
$$B = \frac{\mu_0 I R^2}{2(R^2 + x^2)^{3/2}}$$

## Loi d'Ampère-Maxwell
La loi d'Ampère (forme intégrale) relie le champ magnétique le long d'un chemin fermé au courant traversant la surface délimitée par ce chemin :

$$\oint \vec{B} \cdot d\vec{l} = \mu_0 I_{enc}$$

Où $I_{enc}$ est le courant total traversant la surface délimitée par le chemin d'intégration.

### Applications de la loi d'Ampère

#### Solénoïde infini
Pour un solénoïde de $n$ spires par unité de longueur portant un courant $I$ :
$$B = \mu_0 n I$$
(À l'intérieur du solénoïde, uniforme et parallèle à l'axe)

#### Toroïde
Pour un toroïde de rayon moyen $R$ et $N$ spires totales portant un courant $I$ :
$$B = \frac{\mu_0 N I}{2\pi r}$$
où $r$ est la distance par rapport au centre du toroïde

## Forces entre conducteurs portant un courant
Deux fils parallèles infiniment longs, séparés par une distance $d$ et portant des courants $I_1$ et $I_2$, exercent une force l'un sur l'autre :

$$\frac{F}{L} = \frac{\mu_0 I_1 I_2}{2\pi d}$$

La force est attractive si les courants sont dans la même direction, répulsive s'ils sont dans des directions opposées.

## Moment magnétique
Une boucle de courant $I$ portant une surface $\vec{A}$ (perpendiculaire au plan de la boucle) possède un moment magnétique :
$$\vec{\mu} = I\vec{A}$$

Dans un champ magnétique uniforme $\vec{B}$, cette boucle subit un couple :
$$\vec{\tau} = \vec{\mu} \times \vec{B}$$

Et une énergie potentielle :
$$U = -\vec{\mu} \cdot \vec{B}$$

## Matériaux magnétiques
1. **Diamagnétiques** : Légèrement répulsés par les champs magnétiques ($\chi < 0$)
2. **Paramagnétiques** : Légèrement attirés par les champs magnétiques ($\chi > 0$)
3. **Ferromagnétiques** : Forte attirance, peuvent conserver une aimantation ($\chi >> 0$)

## Unité
L'unité légale du champ magnétique est le tesla (T), où :
$$1 \text{ T} = 1 \text{ N/(A·m)} = 1 \text{ Wb/m}^2$$

D'autres unités utilisées :
- Gauss (G) : $1 \text{ T} = 10^4 \text{ G}$
- Millitesla (mT) : $1 \text{ mT} = 10^{-3} \text{ T}$
- Microtesla (µT) : $1 \text{ µT} = 10^{-6} \text{ T}$

## Exemples

### Exemple 1: Champ autour d'un fil droit
Un fil droit porte un courant de $I = 5 \text{ A}$. Quelle est l'intensité du champ magnétique à une distance $r = 2 \text{ cm}$ du fil ?

$$B = \frac{\mu_0 I}{2\pi r} = \frac{(4\pi \times 10^{-7}) \times 5}{2\pi \times 0.02} = \frac{20\pi \times 10^{-7}}{0.04\pi} = 5 \times 10^{-5} \text{ T} = 50 \text{ µT}$$

### Exemple 2: Champ au centre d'une bobine
Une bobine circulaire de rayon $R = 0.1 \text{ m}$ comporte $N = 100$ spires et porte un courant $I = 2 \text{ A}$. Quel est le champ magnétique au centre de la bobine ?

$$B = \frac{\mu_0 N I}{2R} = \frac{(4\pi \times 10^{-7}) \times 100 \times 2}{2 \times 0.1} = \frac{800\pi \times 10^{-7}}{0.2} = 4000\pi \times 10^{-7} = 1.26 \times 10^{-3} \text{ T} = 1.26 \text{ mT}$$

---
*Source: Adapted from physics textbooks and educational resources*