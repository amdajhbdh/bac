---
tags: [practice, electromagnetism, hard]
---

# Electromagnetism Practice - Hard Difficulty

## Exercices avancés sur l'électromagnétisme

### Exercice 1: Champ électrique d'un disque uniformément chargé
Un disque de rayon $R$ porte une charge uniforme de densité surfacique $\sigma$.

1. Calculer le champ électrique $\vec{E}$ en un point situé sur l'axe du disque à une distance $z$ du centre.
2. Quelle est la limite de ce champ lorsque $z \gg R$ (point lointain) ?
3. Quelle est la limite lorsque $z \ll R$ et $z > 0$ (près de la surface) ?
4. Tracer la variation de $E_z$ en fonction de $z$ pour $z \geq 0$.

### Exercice 2: Condensateur cylindrique
Un condensateur cylindrique est constitué de deux cylindres coaxiaux de longueurs $L$, de rayons $a$ (cylindre interne) et $b$ (cylindre externe), avec $a < b$. Les cylindres portent des charges égales et opposées $\pm Q$.

1. Calculer la capacité de ce condensateur.
2. Quelle est l'énergie stockée dans le condensateur lorsque la différence de potentiel entre les armatures est $V$ ?
3. Si le vide entre les cylindres est remplacé par un diélectrique de permitivité relative $\epsilon_r$, quelle devient la capacité ?

### Exercice 3: Induction électromagnétique dans une bobine tournante
Une bobine circulaire de rayon $R$ comportant $N$ spires tourne avec une vitesse angulaire constante $\omega$ autour d'un de ses diamètres placé dans un champ magnétique uniforme $\vec{B} = B\hat{k}$. L'axe de rotation est l'axe $x$.

1. Quel est le flux magnétique $\Phi_B$ à travers la bobine en fonction du temps ?
2. Quelle est la force électromotrice (fém) induite dans la bobine ?
3. Si la bobine est connectée à une résistance $R$, quel est le courant induit ?
4. Quel est le travail nécessaire pour maintenir la rotation à vitesse constante contre les forces électromagnétiques ?

## Corrigés

### Exercice 1
1. Considérons un anneau de rayon $r$ et de largeur $dr$ sur le disque. Sa charge est $dq = \sigma (2\pi r dr)$.
   La contribution à $E_z$ en un point sur l'axe à distance $z$ du centre est :
   $dE_z = \frac{k dq}{r^2 + z^2} \cos\theta = \frac{k \sigma (2\pi r dr)}{r^2 + z^2} \cdot \frac{z}{\sqrt{r^2 + z^2}}$
   Donc : $E_z = \int_0^R \frac{2\pi k \sigma z r}{(r^2 + z^2)^{3/2}} dr$
   
   En posant $u = r^2 + z^2$, $du = 2r dr$ :
   $E_z = 2\pi k \sigma z \int_{z^2}^{R^2+z^2} \frac{du}{2 u^{3/2}} = \pi k \sigma z \left[ -2 u^{-1/2} \right]_{z^2}^{R^2+z^2}$
   $E_z = 2\pi k \sigma z \left( \frac{1}{z} - \frac{1}{\sqrt{R^2+z^2}} \right)$
   $E_z = 2\pi k \sigma \left( 1 - \frac{z}{\sqrt{R^2+z^2}} \right)$
   
   En utilisant $k = \frac{1}{4\pi\epsilon_0}$ :
   $E_z = \frac{\sigma}{2\epsilon_0} \left( 1 - \frac{z}{\sqrt{R^2+z^2}} \right)$
   
   Pour $z > 0$ (au-dessus du disque), et par symétrie, pour $z < 0$ :
   $E_z = \frac{\sigma}{2\epsilon_0} \left( 1 - \frac{z}{\sqrt{R^2+z^2}} \right) \text{sign}(z)$
   
   Où $\text{sign}(z)$ est le signe de $z$.

2. Lorsque $z \gg R$ : $\sqrt{R^2+z^2} \approx z\sqrt{1+\frac{R^2}{z^2}} \approx z(1+\frac{R^2}{2z^2}) = z + \frac{R^2}{2z}$
   Donc : $\frac{z}{\sqrt{R^2+z^2}} \approx \frac{z}{z + \frac{R^2}{2z}} \approx 1 - \frac{R^2}{2z^2}$
   $E_z \approx \frac{\sigma}{2\epsilon_0} \left( 1 - (1 - \frac{R^2}{2z^2}) \right) = \frac{\sigma}{2\epsilon_0} \cdot \frac{R^2}{2z^2} = \frac{\sigma \pi R^2}{4\pi\epsilon_0 z^2}$
   
   Mais $\sigma \pi R^2 = Q$ (charge totale du disque), donc :
   $E_z \approx \frac{Q}{4\pi\epsilon_0 z^2}$
   
   C'est exactement le champ d'une charge ponctuelle $Q$ située à l'origine, comme attendu pour un observateur lointain.

3. Lorsque $z \ll R$ et $z > 0$ : $\sqrt{R^2+z^2} \approx R\sqrt{1+\frac{z^2}{R^2}} \approx R(1+\frac{z^2}{2R^2}) = R + \frac{z^2}{2R}$
   Donc : $\frac{z}{\sqrt{R^2+z^2}} \approx \frac{z}{R + \frac{z^2}{2R}} \approx \frac{z}{R}(1 - \frac{z^2}{2R^2})$
   $E_z \approx \frac{\sigma}{2\epsilon_0} \left( 1 - \frac{z}{R}(1 - \frac{z^2}{2R^2}) \right) \approx \frac{\sigma}{2\epsilon_0} \left(1 - \frac{z}{R}\right)$
   
   Pour $z \to 0^+$ : $E_z \to \frac{\sigma}{2\epsilon_0}$
   
   En fait, très près d'un conducteur chargé, le champ devrait être $\frac{\sigma}{\epsilon_0}$. La différence provient du fait que nous calculons le champ d'un disque fini, pas d'un plan infini.
   Pour un plan infini de charge, $E = \frac{\sigma}{2\epsilon_0}$ de chaque côté, donc en sortant du conducteur, on aurait $\frac{\sigma}{\epsilon_0}$.
   Mais ici, nous avons un disque fini, donc près du disque mais pas trop près des bords, le champ approche $\frac{\sigma}{2\epsilon_0}$ de chaque côté.

   En réalité, pour $z \ll R$, on peut considérer le disque comme un plan infini localement, donc :
   $E_z \approx \frac{\sigma}{2\epsilon_0}$ pour $z > 0$ (vers le haut)
   $E_z \approx -\frac{\sigma}{2\epsilon_0}$ pour $z < 0$ (vers le bas)

   Ce résultat peut aussi être obtenu directement de l'expression exacte :
   Lorsque $z \to 0^+$ : $E_z = \frac{\sigma}{2\epsilon_0} \left( 1 - \frac{0}{\sqrt{R^2+0}} \right) = \frac{\sigma}{2\epsilon_0}$
   Lorsque $z \to 0^-$ : $E_z = \frac{\sigma}{2\epsilon_0} \left( 1 - \frac{0}{\sqrt{R^2+0}} \right) \times (-1) = -\frac{\sigma}{2\epsilon_0}$

4. La fonction $E_z(z)$ est :
   - Positive pour $z > 0$, négative pour $z < 0$
   - Nulle à $z = 0$ (par symétrie)
   - Tend vers $0$ lorsque $z \to \pm\infty$
   - A un maximum absolu pour $z > 0$ et un minimum absolu pour $z < 0$
   
   Pour trouver l'extrémum pour $z > 0$, on différencie $E_z(z)$ par rapport à $z$ :
   $\frac{dE_z}{dz} = \frac{\sigma}{2\epsilon_0} \frac{d}{dz} \left( 1 - \frac{z}{\sqrt{R^2+z^2}} \right) = -\frac{\sigma}{2\epsilon_0} \frac{d}{dz} \left( \frac{z}{\sqrt{R^2+z^2}} \right)$
   
   $\frac{d}{dz} \left( \frac{z}{\sqrt{R^2+z^2}} \right) = \frac{\sqrt{R^2+z^2} - z \cdot \frac{z}{\sqrt{R^2+z^2}}}{R^2+z^2} = \frac{R^2+z^2 - z^2}{(R^2+z^2)^{3/2}} = \frac{R^2}{(R^2+z^2)^{3/2}}$
   
   Donc : $\frac{dE_z}{dz} = -\frac{\sigma}{2\epsilon_0} \cdot \frac{R^2}{(R^2+z^2)^{3/2}} < 0$ pour tout $z$
   
   Cela indique que $E_z(z)$ est strictement décroissante. Mais nous savons que $E_z(0) = 0$ et que $E_z(z) \to \frac{\sigma}{2\epsilon_0}$ lorsque $z \to 0^+$...
   
   En fait, j'ai fait une erreur dans la signe lors de:= 2\pi k \sigma z \left( \frac{1}{z} - \frac{1}{\sqrt{R^2+z^2}} \right)$
   
   En réexaminant :
   $E_z = 2\pi k \sigma z \left( \frac{1}{z} - \frac{1}{\sqrt{R^2+z^2}} \right) = 2\pi k \sigma \left( 1 - \frac{z}{\sqrt{R^2+z^2}} \right)$
   
   Pour $z > 0$, le terme entre parenthèses est positif car $\sqrt{R^2+z^2} > z$.
   Lorsque $z \to 0^+$ : $E_z \to 2\pi k \sigma (1 - 0) = 2\pi k \sigma = \frac{\sigma}{2\epsilon_0}$
   Lorsque $z \to +\infty$ : $E_z \to 2\pi k \sigma (1 - 1) = 0$
   
   Donc $E_z$ diminue de $\frac{\sigma}{2\epsilon_0}$ à $0$ lorsque $z$ augmente de $0$ à $+\infty$.
   
   Le maximum est donc en $z=0^+$ avec une valeur de $\frac{\sigma}{2\epsilon_0}$.

### Exercice 2
1. En utilisant le théorème de Gauss avec une surface cylindrique de rayon $r$ et de longueur $L$ coaxialement aux armatures :
   - Pour $r < a$ : $Q_{enc} = 0$ donc $E = 0$
   - Pour $a < r < b$ : $Q_{enc} = +Q$ (charge du cylindre interne)
     $\oint \vec{E} \cdot d\vec{A} = E(2\pi r L) = \frac{Q_{enc}}{\epsilon_0} = \frac{Q}{\epsilon_0}$
     Donc $E = \frac{Q}{2\pi\epsilon_0 r L}$ dirigé radialement vers l'extérieur
   - Pour $r > b$ : $Q_{enc} = Q + (-Q) = 0$ donc $E = 0$
   
   La différence de potentiel entre les armatures est :
   $V = -\int_b^a \vec{E} \cdot d\vec{r} = \int_a^b E dr = \int_a^b \frac{Q}{2\pi\epsilon_0 L} \frac{dr}{r} = \frac{Q}{2\pi\epsilon_0 L} \ln\left(\frac{b}{a}\right)$
   
   La capacité est définie par $C = \frac{Q}{V}$ :
   $C = \frac{Q}{\frac{Q}{2\pi\epsilon_0 L} \ln\left(\frac{b}{a}\right)} = \frac{2\pi\epsilon_0 L}{\ln\left(\frac{b}{a}\right)}$

2. L'énergie stockée dans un condensateur est :
   $U = \frac{1}{2} C V^2 = \frac{1}{2} Q V = \frac{Q^2}{2C}$
   
   En utilisant l'expression de $C$ :
   $U = \frac{1}{2} \left( \frac{2\pi\epsilon_0 L}{\ln\left(\frac{b}{a}\right)} \right) V^2 = \frac{\pi\epsilon_0 L V^2}{\ln\left(\frac{b}{a}\right)}$
   
   Ou encore :
   $U = \frac{Q^2}{2} \cdot \frac{\ln\left(\frac{b}{a}\right)}{2\pi\epsilon_0 L} = \frac{Q^2 \ln\left(\frac{b}{a}\right)}{4\pi\epsilon_0 L}$

3. Avec un diélectrique de permitivité relative $\epsilon_r$ entre les armatures :
   La permitivité devient $\epsilon = \epsilon_r \epsilon_0$
   Donc la capacité est multipliée par $\epsilon_r$ :
   $C' = \epsilon_r C = \epsilon_r \frac{2\pi\epsilon_0 L}{\ln\left(\frac{b}{a}\right)} = \frac{2\pi\epsilon L}{\ln\left(\frac{b}{a}\right)}$
   
   Où $\epsilon = \epsilon_r \epsilon_0$ est la permitivité absolue du diélectrique.

### Exercice 3
1. Le flux magnétique à travers la bobine est :
   $\Phi_B = \int \vec{B} \cdot d\vec{A} = B A \cos\theta$
   
   Où $\theta$ est l'angle entre $\vec{B}$ et la normale à la surface de la bobine.
   
   Si la bobine tourne autour de l'axe $x$ avec une vitesse angulaire $\omega$, et que $\vec{B} = B\hat{k}$,
   alors l'angle entre la normale à la bobine et $\vec{B}$ varie comme $\theta = \omega t$ (en supposant $\theta=0$ quand la bobine est dans le plan $xz$).
   
   La surface de la bobine est $A = \pi R^2$, donc :
   $\Phi_B(t) = B (\pi R^2) \cos(\omega t) N$
   
   Où $N$ est le nombre de spires (le flux à travers chaque spire est le même, donc on multiplie par $N$).

2. La loi de Faraday donne la fém induite :
   $\mathcal{E} = -\frac{d\Phi_B}{dt} = -\frac{d}{dt} \left[ NB\pi R^2 \cos(\omega t) \right] = NB\pi R^2 \omega \sin(\omega t)$

3. Selon la loi d'Ohm, le courant dans la bobine est :
   $I(t) = \frac{\mathcal{E}}{R} = \frac{NB\pi R^2 \omega}{R} \sin(\omega t)$
   
   Où $R$ désigne ici la résistance électrique (pas à confondre avec le rayon de la bobine).

4. Le travail nécessaire pour maintenir la rotation vient du couple électromagnétique qui s'oppose au mouvement.
   
   Le moment magnétique de la bobine est :
   $\vec{\mu} = I \vec{A}$
   
   Où $\vec{A}$ est le vecteur surface, de norme $A = \pi R^2$ et dirigé selon la normale à la bobine.
   
   Si l'angle entre $\vec{A}$ et $\vec{B}$ est $\theta = \omega t$, alors :
   $\mu = I A = I \pi R^2$
   et l'angle entre $\vec{\mu}$ et $\vec{B}$ est aussi $\theta = \omega t$.
   
   Le couple exercé sur le moment magnétique dans un champ magnétique est :
   $\vec{\tau} = \vec{\mu} \times \vec{B}$
   
   De norme : $\tau = \mu B \sin\theta = (I \pi R^2) B \sin(\omega t)$
   
   En remplaçant $I$ par son expression :
   $\tau(t) = \left( \frac{NB\pi R^2 \omega}{R} \sin(\omega t) \right) \pi R^2 B \sin(\omega t) = \frac{N B^2 \pi^2 R^4 \omega}{R} \sin^2(\omega t)$
   
   La puissance nécessaire pour lutter contre ce couple est :
   $P = \tau \omega = \frac{N B^2 \pi^2 R^4 \omega^2}{R} \sin^2(\omega t)$
   
   Le travail sur une période $T = \frac{2\pi}{\omega}$ est :
   $W = \int_0^T P dt = \frac{N B^2 \pi^2 R^4 \omega^2}{R} \int_0^{2\pi/\omega} \sin^2(\omega t) dt$
   
   En posant $u = \omega t$, $du = \omega dt$ :
   $W = \frac{N B^2 \pi^2 R^4 \omega^2}{R} \cdot \frac{1}{\omega} \int_0^{2\pi} \sin^2(u) du = \frac{N B^2 \pi^2 R^4 \omega}{R} \int_0^{2\pi} \sin^2(u) du$
   
   Puisque $\int_0^{2\pi} \sin^2(u) du = \pi$ :
   $W = \frac{N B^2 \pi^2 R^4 \omega}{R} \cdot \pi = \frac{N B^2 \pi^3 R^4 \omega}{R}$
   
   Ce travail est transformé en énergie thermique dans la résistance par effet Joule.

---
*Source: Adapted from advanced physics problem sets*