---
tags: [practice, dynamics, hard]
---

# Dynamics Practice - Hard Difficulty

## Exercices avancés sur les lois de Newton

### Exercice 1: Bloc avec poulie et ressort
Un bloc de masse $m_1 = 4 \text{ kg}$ repose sur une surface horizontale avec un coefficient de frottement statique $\mu_s = 0.3$ et cinétique $\mu_k = 0.2$. Il est relié par une corde passant par une poulie à un deuxième bloc de masse $m_2 = 2 \text{ kg}$ qui pend verticalement. Le système est initialement au repos, puis une force horizontale $F = 15 \text{ N}$ est appliquée sur $m_1$ vers la droite.

1. Déterminez si le système bouge lorsque la force $F$ est appliquée.
2. Si oui, quelle est l'accélération du système?
3. Quelle est la tension dans la corde?
4. Que se passe-t-il si la force $F$ est réduite à $5 \text{ N}$?

### Exercice 2: Mouvement circulaire vertical
Une petite boule de masse $m = 0.1 \text{ kg}$ est attachée à une corde de longueur $L = 0.5 \text{ m}$ et fait un mouvement circulaire vertical. À son point le plus bas, la tension dans la corde est $T_b = 2.0 \text{ N}$.

1. Quelle est la vitesse de la boule au point le plus bas?
2. Quelle est la tension dans la corde au point le plus haut de la trajectoire?
3. Quelle est la vitesse minimale que la boule doit avoir au point le plus bas pour compléter le cercle sans que la corde ne devienne slack?

### Exercice 3: Système de trois blocs
Trois blocs de masses $m_1 = 2 \text{ kg}$, $m_2 = 3 \text{ kg}$, et $m_3 = 5 \text{ kg}$ sont alignés sur une surface horizontale sans frottement. Une force $F = 20 \text{ N}$ est appliquée sur $m_1$ vers la droite.

1. Quelle est l'accélération du système?
2. Quelle est la force de contact entre $m_1$ et $m_2$?
3. Quelle est la force de contact entre $m_2$ et $m_3$?
4. Si le coefficient de frottement cinétique entre chaque bloc et la surface est $\mu_k = 0.1$, quelle est maintenant l'accélération du système?

## Corrigés

### Exercice 1
1. Force maximale de frottement statique sur $m_1$: $f_{s,max} = \mu_s m_1 g = 0.3 \times 4 \times 9.8 = 11.76 \text{ N}$
   Poids de $m_2$: $m_2 g = 2 \times 9.8 = 19.6 \text{ N}$
   Sans la force $F$, $m_2$ ferait tomber $m_1$ puisque $19.6 > 11.76$.
   Avec $F = 15 \text{ N}$ vers la droite sur $m_1$, la force nette sur $m_1$ dans la direction opposée au mouvement potentiel de $m_2$ est $F = 15 \text{ N}$.
   Force nécessaire pour empêcher le mouvement: $19.6 - 11.76 = 7.84 \text{ N}$
   Puisque $F = 15 \text{ N} > 7.84 \text{ N}$, le système bougera vers la gauche ($m_2$ descendant).
2. Lorsque le système bouge:
   Pour $m_1$ (vers la gauche est positif): $T - f_k + F = m_1 a$ où $f_k = \mu_k m_1 g = 0.2 \times 4 \times 9.8 = 7.84 \text{ N}$
   Pour $m_2$ (vers le bas est positif): $m_2 g - T = m_2 a$
   En ajoutant: $m_2 g - f_k + F = (m_1 + m_2)a$
   $a = \frac{m_2 g - \mu_k m_1 g + F}{m_1 + m_2} = \frac{19.6 - 7.84 + 15}{4 + 2} = \frac{26.76}{6} = 4.46 \text{ m/s}^2$
3. $T = m_2 g - m_2 a = 19.6 - 2 \times 4.46 = 19.6 - 8.92 = 10.68 \text{ N}$
   (Ou en utilisant l'équation pour $m_1$: $T = m_1 a + f_k - F = 4 \times 4.46 + 7.84 - 15 = 17.84 + 7.84 - 15 = 10.68 \text{ N}$)
4. Avec $F = 5 \text{ N}$:
   $a = \frac{m_2 g - \mu_k m_1 g + F}{m_1 + m_2} = \frac{19.6 - 7.84 + 5}{6} = \frac{16.76}{6} = 2.79 \text{ m/s}^2$
   Le système bouge toujours puisque $a > 0$.

### Exercice 2
1. Au point le plus bas: $T_b - mg = \frac{mv_b^2}{L}$
   Donc: $v_b^2 = \frac{(T_b - mg)L}{m} = \frac{(2.0 - 0.1 \times 9.8) \times 0.5}{0.1} = \frac{(2.0 - 0.98) \times 0.5}{0.1} = \frac{1.02 \times 0.5}{0.1} = \frac{0.51}{0.1} = 5.1$
   $v_b = \sqrt{5.1} = 2.26 \text{ m/s}$
2. Au point le plus haut: $T_t + mg = \frac{mv_t^2}{L}$
   En utilisant la conservation de l'énergie entre bas et haut:
   $\frac{1}{2}mv_b^2 = \frac{1}{2}mv_t^2 + mg(2L)$
   $v_t^2 = v_b^2 - 4gL = 5.1 - 4 \times 9.8 \times 0.5 = 5.1 - 19.6 = -14.5$
   Puisque $v_t^2$ est négatif, cela signifie que la boule n'atteint pas le point le plus haut avec cette vitesse initiale.
   Cependant, la question demande la tension en supposant qu'elle atteigne le point le plus haut, donc nous devons trouver la vitesse minimale au bas pour atteindre le haut.
3. Pour compléter le cercle, la tension au point le plus haut doit être au moins nulle:
   $T_t + mg = \frac{mv_{t,min}^2}{L}$ avec $T_t = 0$
   Donc: $mg = \frac{mv_{t,min}^2}{L}$ → $v_{t,min}^2 = gL$
   En utilisant la conservation de l'énergie: $\frac{1}{2}mv_{b,min}^2 = \frac{1}{2}mv_{t,min}^2 + mg(2L)$
   $v_{b,min}^2 = v_{t,min}^2 + 4gL = gL + 4gL = 5gL$
   $v_{b,min} = \sqrt{5gL} = \sqrt{5 \times 9.8 \times 0.5} = \sqrt{24.5} = 4.95 \text{ m/s}$
   Maintenant, la tension au point le plus haut avec cette vitesse minimale:
   $T_t + mg = \frac{mv_{t,min}^2}{L} = \frac{m(gL)}{L} = mg$
   Donc: $T_t = mg - mg = 0 \text{ N}$ (comme attendu pour le cas minimal)

### Exercice 3
1. Sans frottement: $F = (m_1 + m_2 + m_3)a$
   $a = \frac{F}{m_1 + m_2 + m_3} = \frac{20}{2 + 3 + 5} = \frac{20}{10} = 2.0 \text{ m/s}^2$
2. Force entre $m_1$ et $m_2$: Cette force doit accélérer $m_2$ et $m_3$
   $F_{12} = (m_2 + m_3)a = (3 + 5) \times 2.0 = 8 \times 2.0 = 16.0 \text{ N}$
3. Force entre $m_2$ et $m_3$: Cette force doit accélérer seulement $m_3$
   $F_{23} = m_3 a = 5 \times 2.0 = 10.0 \text{ N}$
4. Avec frottement $\mu_k = 0.1$:
   Force totale de frottement: $f_k = \mu_k (m_1 + m_2 + m_3)g = 0.1 \times 10 \times 9.8 = 9.8 \text{ N}$
   Force nette: $F_{net} = F - f_k = 20 - 9.8 = 10.2 \text{ N}$
   $a = \frac{F_{net}}{m_1 + m_2 + m_3} = \frac{10.2}{10} = 1.02 \text{ m/s}^2$

---
*Source: Adapted from advanced physics problem sets*