---
tags: [practice, dynamics, medium]
---

# Dynamics Practice - Medium Difficulty

## Exercices d'application sur les lois de Newton

### Exercice 1: Bloc sur plan incliné
Un bloc de masse $m = 5 \text{ kg}$ est placé sur un plan incliné faisant un angle $\theta = 30^\circ$ avec l'horizontale. Le coefficient de frottement cinétique entre le bloc et le plan est $\mu_k = 0.2$.

1. Dessinez le diagramme des forces agissant sur le bloc.
2. Calculez l'accélération du bloc le long du plan.
3. Quelle est la force normale exercée par le plan sur le bloc?

### Exercice 2: Système de poulies
Deux masses $m_1 = 3 \text{ kg}$ et $m_2 = 5 \text{ kg}$ sont connectées par une corde légère passant par une poulie fixe sans frottement. La masse $m_1$ est sur une surface horizontale sans frottement tandis que $m_2$ pend verticalement.

1. Quelle est l'accélération du système?
2. Quelle est la tension dans la corde?
3. Si le système part du repos, quelle distance parcourt $m_2$ en 2 secondes?

### Exercice 3: Mouvement circulaire
Une voiture de masse $m = 1200 \text{ kg}$ tourne sur une route horizontale de rayon $r = 80 \text{ m}$. Le coefficient de frottement statique entre les pneus et la route est $\mu_s = 0.7$.

1. Quelle est la vitesse maximale que la voiture peut avoir sans déraper?
2. Quelle est l'accélération centripète à cette vitesse maximale?
3. Si la route était bancée à un angle de $15^\circ$, quelle serait la vitesse maximale sans frottement nécessaire?

## Corrigés

### Exercice 1
1. Forces: poids $mg$ vers le bas, normale $N$ perpendiculaire au plan, frottement $f_k$ opposé au mouvement.
2. Décomposer le poids: $mg\sin\theta$ le long du plan (vers le bas), $mg\cos\theta$ perpendiculaire au plan.
   Équation selon l'axe parallèle au plan: $mg\sin\theta - f_k = ma$
   Avec $f_k = \mu_k N = \mu_k mg\cos\theta$
   Donc: $mg\sin\theta - \mu_k mg\cos\theta = ma$
   $a = g(\sin\theta - \mu_k\cos\theta) = 9.8(\sin 30^\circ - 0.2\cos 30^\circ) = 9.8(0.5 - 0.2\times 0.866) = 9.8(0.5 - 0.173) = 9.8\times 0.327 = 3.2 \text{ m/s}^2$
3. $N = mg\cos\theta = 5\times 9.8\times \cos 30^\circ = 49\times 0.866 = 42.4 \text{ N}$

### Exercice 2
1. Pour $m_1$: $T = m_1a$
   Pour $m_2$: $m_2g - T = m_2a$
   En ajoutant: $m_2g = (m_1 + m_2)a$
   Donc: $a = \frac{m_2g}{m_1 + m_2} = \frac{5\times 9.8}{3 + 5} = \frac{49}{8} = 6.125 \text{ m/s}^2$
2. $T = m_1a = 3\times 6.125 = 18.375 \text{ N}$
3. $d = \frac{1}{2}at^2 = \frac{1}{2}\times 6.125\times 2^2 = \frac{1}{2}\times 6.125\times 4 = 12.25 \text{ m}$

### Exercice 3
1. Force de frottement maximale: $f_{s,max} = \mu_s N = \mu_s mg$
   Cette force fournit l'accélération centripète: $f_{s,max} = \frac{mv^2}{r}$
   Donc: $\mu_s mg = \frac{mv^2}{r}$
   $v^2 = \mu_s gr$
   $v = \sqrt{\mu_s gr} = \sqrt{0.7\times 9.8\times 80} = \sqrt{548.8} = 23.4 \text{ m/s}$
2. $a_c = \frac{v^2}{r} = \frac{548.8}{80} = 6.86 \text{ m/s}^2$
3. Sur route bancée sans frottement: $N\sin\theta = \frac{mv^2}{r}$ et $N\cos\theta = mg$
   En divisant: $\tan\theta = \frac{v^2}{rg}$
   Donc: $v = \sqrt{rg\tan\theta} = \sqrt{80\times 9.8\times \tan 15^\circ} = \sqrt{784\times 0.268} = \sqrt{210.1} = 14.5 \text{ m/s}$

---
*Source: Adapted from various physics problem sets*