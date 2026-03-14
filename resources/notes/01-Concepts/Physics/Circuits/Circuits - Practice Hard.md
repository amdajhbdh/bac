---
tags: [practice, circuits, hard]
---

# Circuits Practice - Hard Difficulty

## Exercices avancés sur les circuits électriques

### Exercice 1: Circuit RLC en série
Un circuit RLC en série est composé d'une résistance $R = 50 \Omega$, d'une inductance $L = 0.1$ H, et d'une capacité $C = 10 \mu$F. Le circuit est alimenté par une source de tension sinusoïdale de $U(t) = U_0 \sin(\omega t)$ avec $U_0 = 100$ V.

1. Calculer la pulsation de résonance $\omega_0$ du circuit.
2. Déterminer l'impédance du circuit à la résonance.
3. Calculer le courant efficace dans le circuit à la résonance.
4. Quelle est la tension aux bornes de la bobine et du condensateur à la résonance ?
5. Si la fréquence d'alimentation est de $f = 50$ Hz, calculer le facteur de puissance du circuit.

### Exercice 2: Filtre passe-bas RC
Un filtre passe-bas du premier ordre est constitué d'une résistance $R = 1$ k$\Omega$ en série avec une capacité $C = 100$ nF, suivie d'un ampli-op en suivi de tension (pour éviter l'effet de charge).

1. Calculer la fréquence de coupure $f_c$ du filtre.
2. Quel est l'atténuation (en dB) à une fréquence de $f = 2f_c$ ?
3. Quelle est la phase de la sortie par rapport à l'entrée à la fréquence de coupure ?
4. Si on veut obtenir une fréquence de coupure de $f_c = 10$ kHz avec la même résistance, quelle devrait être la valeur de la capacité ?

### Exercice 3: Amplificateur opérationnel en inverseur
Un ampli-op idéal est monté en configuration inverseur avec une résistance d'entrée $R_1 = 10$ k$\Omega$ et une résistance de réaction $R_f = 100$ k$\Omega$. Il est alimenté par des tensions $\pm V_{sat} = \pm 12$ V.

1. Quelle est la gain en tension du montage ?
2. Quelle est la tension de sortie pour une entrée de $U_e = -0,5$ V ?
3. Quelle est la plus grande amplitude sinusoïdale que l'on peut appliquer à l'entrée sans saturer la sortie ?
4. Quelle est la tension différentielle d'entrée de l'ampli-op dans ces conditions ?
5. Si la bande passante en boucle ouverte de l'ampli-op est de $f_0 = 1$ MHz, quelle est la bande passante en boucle fermée du montage ?

## Corrigés

### Exercice 1
1. La pulsation de résonance d'un circuit RLC en série est :
   $\omega_0 = \frac{1}{\sqrt{LC}} = \frac{1}{\sqrt{0.1 \times 10 \times 10^{-6}}} = \frac{1}{\sqrt{1 \times 10^{-6}}} = \frac{1}{10^{-3}} = 1000$ rad/s
   
   En fréquence : $f_0 = \frac{\omega_0}{2\pi} = \frac{1000}{2\pi} \approx 159.2$ Hz

2. À la résonance, l'impédance d'un circuit RLC en série est purement résistive et égale à $R$ :
   $Z_0 = R = 50 \Omega$
   
   (À la résonance, les réactances du inducteur et du condensateur s'annulent : $X_L = X_C$)

3. Le courant efficace à la résonance est :
   $I_{rms} = \frac{U_{rms}}{Z_0} = \frac{U_0/\sqrt{2}}{R} = \frac{100/\sqrt{2}}{50} = \frac{100}{50\sqrt{2}} = \frac{2}{\sqrt{2}} = \sqrt{2} \approx 1.41$ A
   
   (où $U_{rms} = \frac{U_0}{\sqrt{2}}$ pour une sinusoidale)

4. À la résonance :
   - Tension aux bornes de la résistance : $U_R = RI_{rms} = 50 \times \sqrt{2} = 50\sqrt{2} \approx 70.7$ V
   - Tension aux bornes de l'inductance : $U_L = IX_L = I\omega_0L$
   - Tension aux bornes de la capacité : $U_C = IX_C = \frac{I}{\omega_0C}$
   
   Puisque $X_L = X_C = \omega_0L = \frac{1}{\omega_0C}$ à la résonance :
   $X_L = \omega_0L = 1000 \times 0.1 = 100 \Omega$
   $U_L = I_{rms} X_L = \sqrt{2} \times 100 = 100\sqrt{2} \approx 141.4$ V
   $U_C = I_{rms} X_C = \sqrt{2} \times 100 = 100\sqrt{2} \approx 141.4$ V
   
   Remarquons que $U_L$ et $U_C$ sont en opposition de phase (de 180°), donc leur somme est nulle, ce qui est cohérent avec la loi des mailles : $U(t) = U_R + U_L + U_C$.

5. Pour $f = 50$ Hz : $\omega = 2\pi f = 2\pi \times 50 = 100\pi$ rad/s
   
   Réactances :
   $X_L = \omega L = 100\pi \times 0.1 = 10\pi \approx 31.4 \Omega$
   $X_C = \frac{1}{\omega C} = \frac{1}{100\pi \times 10 \times 10^{-6}} = \frac{1}{100\pi \times 10^{-5}} = \frac{10^5}{100\pi} = \frac{1000}{\pi} \approx 318.3 \Omega$
   
   Impédance :
   $Z = \sqrt{R^2 + (X_L - X_C)^2} = \sqrt{50^2 + (31.4 - 318.3)^2} = \sqrt{2500 + (-286.9)^2} = \sqrt{2500 + 82311.6} = \sqrt{84811.6} \approx 291.2 \Omega$
   
   Facteur de puissance : $\cos\varphi = \frac{R}{Z} = \frac{50}{291.2} \approx 0.172$
   
   Comme $X_C > X_L$, le circuit est à dominante capacitive, donc le courant avance sur la tension ($\varphi < 0$).

### Exercice 2
1. La fréquence de coupure d'un filtre passe-bas RC du premier ordre est :
   $f_c = \frac{1}{2\pi RC} = \frac{1}{2\pi \times 1000 \times 100 \times 10^{-9}} = \frac{1}{2\pi \times 10^{-4}} = \frac{10^4}{2\pi} \approx \frac{10000}{6.28} \approx 1592$ Hz $\approx 1.59$ kHz
   
2. Le gain en module d'un filtre passe-bas RC est :
   $|G(f)| = \frac{1}{\sqrt{1 + (f/f_c)^2}}$
   
   À $f = 2f_c$ :
   $|G(2f_c)| = \frac{1}{\sqrt{1 + (2)^2}} = \frac{1}{\sqrt{1+4}} = \frac{1}{\sqrt{5}}$
   
   L'atténuation en dB est :
   $A_{dB} = 20 \log_{10}(|G|) = 20 \log_{10}\left(\frac{1}{\sqrt{5}}\right) = -20 \log_{10}(\sqrt{5}) = -20 \times \frac{1}{2} \log_{10}(5) = -10 \log_{10}(5)$
   $A_{dB} \approx -10 \times 0.699 = -6.99$ dB ≈ -7 dB
   
   (Le signe moins indique une atténuation, on dit souvent "une atténuation de 7 dB")

3. La phase d'un filtre passe-bas RC est :
   $\varphi(f) = -\arctan\left(\frac{f}{f_c}\right)$
   
   À la fréquence de coupure : $f = f_c$
   $\varphi(f_c) = -\arctan(1) = -45^\circ$
   
   La sortie est en retard de 45° par rapport à l'entrée.

4. Pour $f_c = 10$ kHz = $10^4$ Hz avec $R = 1$ k$\Omega$ = $10^3$ $\Omega$ :
   $f_c = \frac{1}{2\pi RC}$
   $C = \frac{1}{2\pi R f_c} = \frac{1}{2\pi \times 10^3 \times 10^4} = \frac{1}{2\pi \times 10^7} = \frac{10^{-7}}{2\pi} \approx 1.59 \times 10^{-8}$ F = 15.9 nF

### Exercice 3
1. Pour un ampli-op idéal en configuration inverseur :
   $A_v = -\frac{R_f}{R_1} = -\frac{100}{10} = -10$
   
   Le gain est de -10, ce qui signifie une amplification de 10 fois avec inversion de phase.

2. $U_s = A_v \times U_e = -10 \times (-0,5) = 5$ V
   
   (Cette valeur est inférieure à la tension de saturation de 12 V, donc l'ampli-op n'est pas saturé)

3. La sortie ne peut pas dépasser les tensions d'alimentation : $-V_{sat} \leq U_s \leq V_{sat}$
   Donc : $-12 \leq U_s \leq 12$
   
   Puisque $U_s = A_v U_e = -10 U_e$ :
   $-12 \leq -10 U_e \leq 12$
   
   En divisant par -10 (et en inversant les inégalités) :
   $1.2 \geq U_e \geq -1.2$
   
   Donc : $-1.2 \leq U_e \leq 1.2$ V
   
   La plus grande amplitude sinusoïdale centrée sur 0 que l'on peut appliquer est donc :
   $U_{e,max} = 1.2$ V (amplitude)
   
   Ce qui donne une sortie de $U_{s,max} = -10 \times 1.2 = -12$ V (ou +12 V selon la phase)

4. Pour un ampli-op idéal, la tension différentielle d'entrée est nulle lorsque l'ampli-op n'est pas saturé, grâce à la rétroaction négative qui maintient $v_+ = v_-$.
   
   Dans notre cas, avec $U_e = -0,5$ V et $U_s = 5$ V :
   La tension à l'entrée inverseur ($v_-$) est déterminée par le diviseur $R_1$-$R_f$ entre $U_e$ et $U_s$ :
   $v_- = U_e + (U_s - U_e)\frac{R_1}{R_1+R_f} = -0,5 + (5 - (-0,5))\frac{10}{10+100} = -0,5 + 5.5 \times \frac{10}{110} = -0,5 + 5.5 \times 0,0909 = -0,5 + 0,5 = 0$ V
   
   L'entrée non inverseur $v_+$ est reliée à la masse, donc $v_+ = 0$ V.
   
   Donc la tension différentielle d'entrée est : $v_d = v_+ - v_- = 0 - 0 = 0$ V
   
   (Ceci est vrai tant que l'ampli-op n'est pas saturé)

5. La bande passante en boucle fermée d'un ampli-op est liée à la bande passante en boucle ouverte par :
   $f_{BF} = \frac{f_{BO}}{|A_{BF}|}$
   
   Où $f_{BO}$ est la bande passante en boucle ouverte et $A_{BF}$ est le gain en boucle fermé.
   
   Dans notre cas : $f_{BF} = \frac{1\ \text{MHz}}{|-10|} = \frac{10^6}{10} = 100\ \text{kHz}$
   
   Alternativement, on peut utiliser le produit gain-bande passante qui est constant pour un ampli-op donné :
   $G \times f = \text{constant}$
   
   En boucle ouverte : $G_{BO} \times f_{BO} = \text{constant}$ (mais $G_{BO}$ varie avec la fréquence)
   En réalité, pour un ampli-op compensé, le produit gain en bande passante est approximativement constant :
   $|A_v| \times f_{BF} = f_0$ où $f_0$ est la fréquence unité gain (où le gain en boucle ouverte tombe à 1)
   
   Donc : $f_{BF} = \frac{f_0}{|A_v|} = \frac{1\ \text{MHz}}{10} = 100\ \text{kHz}$

---
*Source: Adapted from advanced physics problem sets*