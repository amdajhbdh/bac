---
tags: [practice, electromagnetism, medium]
---

# Electromagnetism Practice - Medium Difficulty

## Exercices d'application sur les champs électriques et magnétiques

### Exercice 1: Dipôle électrique
Un dipôle électrique est constitué de deux charges égales et opposées $+q$ et $-q$ séparées par une distance $d$. Le moment électrique du dipôle est défini comme $\vec{p} = q\vec{d}$, où $\vec{d}$ pointe de la charge négative vers la charge positive.

1. Calculer le potentiel électrique $V$ en un point situé sur l'axe du dipôle à une distance $r$ du centre ($r \gg d$).
2. Calculer le champ électrique $\vec{E}$ en ce même point.
3. Quel est l'énergie potentielle $U$ du dipôle placé dans un champ électrique uniforme $\vec{E}_0$ ?

### Exercice 2: Loi de Gauss et symétrie cylindrique
Un câble coaxial est constitué d'un conducteur interne de rayon $a$ et d'un conducteur externe de rayon interne $b$ et externe $c$. Une charge linéaire $\lambda$ est répartie sur le conducteur interne.

1. En utilisant le théorème de Gauss, déterminer le champ électrique $\vec{E}$ pour :
   - $r < a$ (à l'intérieur du conducteur interne)
   - $a < r < b$ (entre les conducteurs)
   - $b < r < c$ (dans le conducteur externe)
   - $r > c$ (à l'extérieur du câble)
2. Si le conducteur externe est mis à la terre, quelle est la différence de potentiel entre les conducteurs ?

### Exercice 3: Mouvement d'une particule dans des champs croisés
Une particule de charge $q > 0$ et de masse $m$ est lancée avec une vitesse $\vec{v}_0 = v_0\hat{i}$ dans une région où il existe un champ électrique uniforme $\vec{E} = E\hat{j}$ et un champ magnétique uniforme $\vec{B} = B\hat{k}$.

1. Quelle est la force initiale agissant sur la particule ?
2. Déterminer la trajectoire de la particule.
3. Sous quelle condition la particule continuera-t-elle à se déplacer en ligne droite sans déviation ?

## Corrigés

### Exercice 1
1. Le potentiel d'une charge ponctuelle est $V = k\frac{q}{r}$. Pour le dipôle :
   $V = k\left[\frac{q}{r - \frac{d}{2}\cos\theta} - \frac{q}{r + \frac{d}{2}\cos\theta}\right]$
   Pour $r \gg d$, en utilisant l'approximation $\frac{1}{1 \pm x} \approx 1 \mp x$ pour $x \ll 1$ :
   $V \approx kq\left[\frac{1}{r}\left(1 + \frac{d}{2r}\cos\theta\right) - \frac{1}{r}\left(1 - \frac{d}{2r}\cos\theta\right)\right]$
   $V \approx kq\left[\frac{d}{r^2}\cos\theta\right] = k\frac{p\cos\theta}{r^2}$ où $p = qd$
   
   Sur l'axe, $\theta = 0$ ou $\pi$, donc :
   $V = \pm k\frac{p}{r^2}$ (le signe dépend de l'orientation)

2. $\vec{E} = -\nabla V$
   En coordonnées sphériques, pour un potentiel qui ne dépend que de $r$ et $\theta$ :
   $E_r = -\frac{\partial V}{\partial r} = -\frac{\partial}{\partial r}\left(k\frac{p\cos\theta}{r^2}\right) = 2k\frac{p\cos\theta}{r^3}$
   $E_\theta = -\frac{1}{r}\frac{\partial V}{\partial \theta} = -\frac{1}{r}\frac{\partial}{\partial \theta}\left(k\frac{p\cos\theta}{r^2}\right) = k\frac{p\sin\theta}{r^3}$
   $E_\phi = 0$
   
   Sur l'axe ($\theta = 0$) : $\vec{E} = 2k\frac{p}{r^3}\hat{i}$
   Sur l'axe ($\theta = \pi$) : $\vec{E} = -2k\frac{p}{r^3}\hat{i}$

3. $U = -\vec{p} \cdot \vec{E}_0 = -pE_0\cos\theta$
   Où $\theta$ est l'angle entre $\vec{p}$ et $\vec{E}_0$

### Exercice 2
1. En utilisant une surface gaussienne cylindrique de rayon $r$ et de longueur $L$ :
   - Pour $r < a$ : $Q_{enc} = 0$ donc $E = 0$
   - Pour $a < r < b$ : $Q_{enc} = \lambda L$
     $\oint \vec{E} \cdot d\vec{A} = E(2\pi r L) = \frac{Q_{enc}}{\epsilon_0} = \frac{\lambda L}{\epsilon_0}$
     Donc $E = \frac{\lambda}{2\pi\epsilon_0 r}$ dirigé radialement vers l'extérieur si $\lambda > 0$
   - Pour $b < r < c$ : Dans un conducteur à l'équilibre électrostatique, $E = 0$
   - Pour $r > c$ : $Q_{enc} = \lambda L$ (la charge sur le conducteur externe est sur sa surface interne)
     Donc $E = \frac{\lambda}{2\pi\epsilon_0 r}$ (même expression que pour $a < r < b$)
   
   Cependant, si le conducteur externe est mis à la terre, sa potentiel est fixé à 0, ce qui induit une charge $-\lambda L$ sur sa surface interne et $+\lambda L$ sur sa surface externe si celui-ci était initialement neutre. Mais s'il est mis à la terre, il peut échanger librement des charges avec la terre, donc pour maintenir $V=0$, la charge totale sur le conducteur externe sera telle que le potentiel produit par $\lambda L$ sur le conducteur interne soit annulé.
   
   En réalité, pour un câble coaxial avec le conducteur externe mis à la terre :
   - $a < r < b$ : $E = \frac{\lambda}{2\pi\epsilon_0 r}$
   - $r > b$ : $E = 0$ (puisque le conducteur externe blindé empêche le champ de s'échapper vers l'extérieur)

2. La différence de potentiel $V(a) - V(b)$ est :
   $V(a) - V(b) = -\int_b^a \vec{E} \cdot d\vec{r} = -\int_b^a \frac{\lambda}{2\pi\epsilon_0 r} dr = \frac{\lambda}{2\pi\epsilon_0} \ln\left(\frac{b}{a}\right)$

### Exercice 3
1. Force initiale : $\vec{F} = q(\vec{E} + \vec{v}_0 \times \vec{B})$
   $\vec{v}_0 \times \vec{B} = v_0\hat{i} \times B\hat{k} = v_0B (\hat{i} \times \hat{k}) = v_0B (-\hat{j}) = -v_0B\hat{j}$
   Donc $\vec{F} = q(E\hat{j} - v_0B\hat{j}) = q(E - v_0B)\hat{j}$

2. Les équations du mouvement sont :
   $m\frac{dv_x}{dt} = q(v_y B)$
   $m\frac{dv_y}{dt} = q(E - v_x B)$
   
   En posant $\omega_c = \frac{qB}{m}$ (fréquence de cyclotron) et $v_E = \frac{E}{B}$ (vitesse de dérive) :
   $\frac{dv_x}{dt} = \omega_c v_y$
   $\frac{dv_y}{dt} = \omega_c(v_E - v_x)$
   
   La solution générale est une combinaison d'un mouvement circulaire de fréquence $\omega_c$ et d'une translation uniforme à la vitesse $\vec{v}_E = \frac{\vec{E} \times \vec{B}}{B^2}$.
   
   Avec les conditions initiales $v_x(0) = v_0$, $v_y(0) = 0$ :
   $v_x(t) = v_E + (v_0 - v_E)\cos(\omega_c t)$
   $v_y(t) = -(v_0 - v_E)\sin(\omega_c t)$
   
   En intégrant :
   $x(t) = v_E t + \frac{v_0 - v_E}{\omega_c}\sin(\omega_c t)$
   $y(t) = \frac{v_0 - v_E}{\omega_c}(\cos(\omega_c t) - 1)$
   
   La trajectoire est une trochoïde (cycloïde générale).

3. Pour qu'il n'y ait pas de déviation, la force nette doit être nulle :
   $\vec{F} = q(\vec{E} + \vec{v} \times \vec{B}) = 0$
   Donc $\vec{E} + \vec{v} \times \vec{B} = 0$
   $\vec{v} \times \vec{B} = -\vec{E}$
   
   En effectuant un produit vectoriel par $\vec{B}$ des deux côtés :
   $(\vec{v} \times \vec{B}) \times \vec{B} = -\vec{E} \times \vec{B}$
   En utilisant l'identité du triple produit vectoriel :
   $\vec{v}(\vec{B}\cdot\vec{B}) - \vec{B}(\vec{v}\cdot\vec{B}) = -\vec{E} \times \vec{B}$
   
   Si $\vec{v} \perp \vec{B}$ (comme dans notre cas avec $\vec{v} = v\hat{i}$ et $\vec{B} = B\hat{k}$), alors $\vec{v}\cdot\vec{B} = 0$ :
   $vB^2 = -(\vec{E} \times \vec{B})\cdot\hat{v}$ (après quelques manipulations)
   
   Plus simplement, la condition pour aucune déviation est que la vitesse soit égale à la vitesse de dérive :
   $\vec{v} = \frac{\vec{E} \times \vec{B}}{B^2}$
   
   Dans notre cas spécifique : $\vec{E} = E\hat{j}$, $\vec{B} = B\hat{k}$
   $\vec{E} \times \vec{B} = EB(\hat{j} \times \hat{k}) = EB\hat{i}$
   Donc $\vec{v} = \frac{EB}{B^2}\hat{i} = \frac{E}{B}\hat{i}$
   
   Donc la particule se déplacera en ligne droite si $v_0 = \frac{E}{B}$.

---
*Source: Adapted from physics problem sets*