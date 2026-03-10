#!/usr/bin/env python3
"""
Auto-Generated Manim Animations for All Exercises
Based on actual exercise content from Study Vault
"""
from manim import *
import re


class HormoneRegulation(Scene):
    """Exercise: Hormone regulation - FSH/LH levels"""
    def construct(self):
        title = Text("Régulation Hormonale", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Axes for hormone levels
        axes = Axes(
            x_range=[0, 30], y_range=[0, 100],
            axis_config={"include_tip": True},
            x_length=10, y_length=5
        )
        labels = axes.get_axis_labels(x_label="Jours", y_label="Concentration (UI/L)")
        
        self.play(Create(axes), Write(labels))
        
        # Normal FSH curve
        fsh_normal = axes.plot(lambda x: 32 + 10*np.sin(x/5), color=BLUE, x_range=[0, 28])
        fsh_label = Text("FSH normale", font_size=20, color=BLUE).to_edge(LEFT)
        
        # Abnormal FSH (Woman M - high)
        fsh_high = axes.plot(lambda x: 92, color=RED, x_range=[0, 28])
        fsh_high_label = Text("FSH élevée (Femme M)", font_size=20, color=RED).next_to(fsh_label, DOWN)
        
        # Abnormal FSH (Woman N - low)
        fsh_low = axes.plot(lambda x: 4, color=GREEN, x_range=[0, 28])
        fsh_low_label = Text("FSH basse (Femme N)", font_size=20, color=GREEN).next_to(fsh_high_label, DOWN)
        
        self.play(Create(fsh_normal), Write(fsh_label))
        self.wait()
        self.play(Create(fsh_high), Write(fsh_high_label))
        self.wait()
        self.play(Create(fsh_low), Write(fsh_low_label))
        self.wait(2)

class NeuronSynapse(Scene):
    """Exercise: Synaptic transmission with neurotoxins"""
    def construct(self):
        title = Text("Transmission Synaptique", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Presynaptic neuron
        pre = Rectangle(width=2, height=1, color=BLUE, fill_opacity=0.3)
        pre.shift(LEFT * 3)
        pre_label = Text("Neurone présynaptique", font_size=16).move_to(pre)
        
        # Synaptic cleft
        cleft = Rectangle(width=0.5, height=1.5, color=YELLOW, fill_opacity=0.2)
        cleft.next_to(pre, RIGHT, buff=0)
        
        # Postsynaptic neuron
        post = Rectangle(width=2, height=1, color=RED, fill_opacity=0.3)
        post.next_to(cleft, RIGHT, buff=0)
        post_label = Text("Neurone postsynaptique", font_size=16).move_to(post)
        
        self.play(Create(pre), Create(cleft), Create(post))
        self.play(Write(pre_label), Write(post_label))
        
        # Neurotransmitter release
        vesicles = VGroup(*[Dot(color=GREEN, radius=0.05).move_to(pre.get_right() + UP*0.3*i) for i in range(-2, 3)])
        self.play(Create(vesicles))
        
        # Release animation
        self.play(vesicles.animate.shift(RIGHT * 2.5), run_time=2)
        self.wait(2)

class GeneticCrossing(Scene):
    """Exercise: Drosophila genetic crossing"""
    def construct(self):
        title = Text("Croisement Génétique - Drosophile", font_size=36)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Parent genotypes
        parent1 = MathTex("vg^+bw^+", color=BLUE, font_size=48)
        parent1.shift(LEFT * 3 + UP)
        p1_label = Text("Femelle", font_size=20).next_to(parent1, DOWN)
        
        parent2 = MathTex("vg\\,bw", color=RED, font_size=48)
        parent2.shift(RIGHT * 3 + UP)
        p2_label = Text("Mâle", font_size=20).next_to(parent2, DOWN)
        
        self.play(Write(parent1), Write(p1_label))
        self.play(Write(parent2), Write(p2_label))
        
        # Cross symbol
        cross = MathTex("\\times", font_size=60).move_to(ORIGIN + UP)
        self.play(Write(cross))
        self.wait()
        
        # Offspring percentages
        offspring = VGroup(
            Text("32% ailes longues, yeux rouges", font_size=20),
            Text("32% ailes vestigiales, yeux bruns", font_size=20),
            Text("18% ailes longues, yeux bruns", font_size=20),
            Text("18% ailes vestigiales, yeux rouges", font_size=20)
        ).arrange(DOWN, aligned_edge=LEFT).shift(DOWN * 1.5)
        
        self.play(Write(offspring), run_time=3)
        self.wait(2)

class ImmuneResponse(Scene):
    """Exercise: Immune cells X and Y"""
    def construct(self):
        title = Text("Réponse Immunitaire", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Cell X (B lymphocyte)
        cellX = Circle(radius=0.8, color=BLUE, fill_opacity=0.3)
        cellX.shift(LEFT * 3)
        x_label = Text("Lymphocyte B", font_size=20).next_to(cellX, DOWN)
        
        # Cell Y (Plasma cell)
        cellY = Ellipse(width=1.5, height=2, color=RED, fill_opacity=0.3)
        cellY.shift(RIGHT * 3)
        y_label = Text("Plasmocyte", font_size=20).next_to(cellY, DOWN)
        
        self.play(Create(cellX), Write(x_label))
        self.wait()
        
        # Differentiation arrow
        arrow = Arrow(cellX.get_right(), cellY.get_left(), color=YELLOW)
        diff_label = Text("Différenciation", font_size=18).next_to(arrow, UP)
        
        self.play(Create(arrow), Write(diff_label))
        self.play(Create(cellY), Write(y_label))
        
        # Antibodies
        antibodies = VGroup(*[
            Text("Y", font_size=16, color=GREEN).move_to(cellY.get_center() + 1.5*RIGHT + 0.3*UP*i)
            for i in range(-3, 4)
        ])
        
        self.play(Create(antibodies))
        self.wait(2)

# ============= PHYSICS ANIMATIONS =============

class ProjectileMotionDetailed(Scene):
    """Exercise: Projectile motion with velocity components"""
    def construct(self):
        title = Text("Mouvement de Projectile", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        axes = Axes(
            x_range=[0, 10], y_range=[0, 6],
            axis_config={"include_tip": True}
        )
        labels = axes.get_axis_labels(x_label="x (m)", y_label="y (m)")
        
        self.play(Create(axes), Write(labels))
        
        # Initial velocity vector
        v0 = 10  # m/s
        angle = 45  # degrees
        v0x = v0 * np.cos(np.radians(angle))
        v0y = v0 * np.sin(np.radians(angle))
        
        # Trajectory
        g = 9.8
        path = axes.plot(
            lambda x: x * np.tan(np.radians(angle)) - (g * x**2) / (2 * v0x**2),
            x_range=[0, 2*v0x*v0y/g],
            color=BLUE
        )
        
        self.play(Create(path))
        
        # Ball
        ball = Dot(color=RED, radius=0.1)
        ball.move_to(axes.c2p(0, 0))
        
        # Velocity components
        v_arrow = always_redraw(lambda: Arrow(
            ball.get_center(),
            ball.get_center() + RIGHT * 0.5 + UP * 0.3,
            color=YELLOW, buff=0
        ))
        
        self.play(Create(ball), Create(v_arrow))
        self.play(MoveAlongPath(ball, path), run_time=4)
        self.wait(2)

class SpringOscillationDetailed(Scene):
    """Exercise: Spring oscillation with energy"""
    def construct(self):
        title = Text("Oscillation d'un Ressort", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Fixed point
        fixed = Dot(LEFT * 3, color=GRAY)
        self.play(Create(fixed))
        
        # Spring
        spring = Line(fixed.get_center(), fixed.get_center() + RIGHT * 2, color=BLUE)
        
        # Mass
        mass = Square(side_length=0.5, color=RED, fill_opacity=1)
        mass.next_to(spring, RIGHT, buff=0)
        
        self.play(Create(spring), Create(mass))
        
        # Energy labels
        pe_label = MathTex("E_p = \\frac{1}{2}kx^2", color=BLUE).to_edge(DOWN).shift(LEFT * 3)
        ke_label = MathTex("E_c = \\frac{1}{2}mv^2", color=RED).to_edge(DOWN).shift(RIGHT * 3)
        
        self.play(Write(pe_label), Write(ke_label))
        
        # Oscillation
        for _ in range(2):
            self.play(
                mass.animate.shift(RIGHT * 1.5),
                spring.animate.stretch(1.5, 0),
                pe_label.animate.set_color(YELLOW),
                run_time=1
            )
            self.play(
                mass.animate.shift(LEFT * 3),
                spring.animate.stretch(0.5, 0),
                ke_label.animate.set_color(YELLOW),
                run_time=2
            )
            self.play(
                mass.animate.shift(RIGHT * 1.5),
                spring.animate.stretch(1, 0),
                pe_label.animate.set_color(BLUE),
                ke_label.animate.set_color(RED),
                run_time=1
            )
        
        self.wait(2)

class ElectricFieldLines(Scene):
    """Exercise: Electric field visualization"""
    def construct(self):
        title = Text("Champ Électrique", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Positive charge
        charge_pos = Circle(radius=0.3, color=RED, fill_opacity=1)
        charge_pos.shift(LEFT * 2)
        label_pos = MathTex("+Q", color=WHITE).move_to(charge_pos)
        
        # Negative charge
        charge_neg = Circle(radius=0.3, color=BLUE, fill_opacity=1)
        charge_neg.shift(RIGHT * 2)
        label_neg = MathTex("-Q", color=WHITE).move_to(charge_neg)
        
        self.play(Create(charge_pos), Create(charge_neg))
        self.play(Write(label_pos), Write(label_neg))
        
        # Field lines
        lines = VGroup()
        for i in range(8):
            angle = i * PI / 4
            start_pos = charge_pos.get_center() + 0.3 * np.array([np.cos(angle), np.sin(angle), 0])
            
            # Curved line to negative charge
            curve = CubicBezier(
                start_pos,
                start_pos + 1.5 * np.array([np.cos(angle), np.sin(angle), 0]),
                charge_neg.get_center() - 1.5 * np.array([np.cos(angle), np.sin(angle), 0]),
                charge_neg.get_center() - 0.3 * np.array([np.cos(angle), np.sin(angle), 0]),
                color=YELLOW
            )
            lines.add(curve)
        
        self.play(Create(lines), run_time=2)
        self.wait(2)

class DynamicsForces(Scene):
    """Exercise: Forces and acceleration"""
    def construct(self):
        title = Text("Dynamique: F = ma", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Object
        box = Square(side_length=1, color=BLUE, fill_opacity=0.5)
        box.shift(LEFT * 3)
        
        # Force arrow
        force = Arrow(box.get_right(), box.get_right() + RIGHT * 2, color=RED, buff=0, stroke_width=8)
        f_label = MathTex("\\vec{F}", color=RED, font_size=40).next_to(force, UP)
        
        self.play(Create(box))
        self.play(Create(force), Write(f_label))
        
        # Acceleration
        self.play(box.animate.shift(RIGHT * 5), run_time=2, rate_func=lambda t: t**2)
        
        # Formula
        formula = MathTex("a = \\frac{F}{m}", font_size=60).to_edge(DOWN)
        self.play(Write(formula))
        self.wait(2)

# ============= MATH ANIMATIONS =============

class ComplexNumbersPlane(Scene):
    """Exercise: Complex numbers on complex plane"""
    def construct(self):
        title = Text("Nombres Complexes", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        plane = ComplexPlane(
            x_range=[-5, 5], y_range=[-5, 5],
            background_line_style={"stroke_opacity": 0.4}
        )
        labels = plane.get_axis_labels(x_label="Re", y_label="Im")
        
        self.play(Create(plane), Write(labels))
        
        # Multiple complex numbers
        z1 = 3 + 2j
        z2 = -2 + 3j
        z3 = 2 - 2j
        
        dots = VGroup(
            Dot(plane.n2p(z1), color=YELLOW),
            Dot(plane.n2p(z2), color=GREEN),
            Dot(plane.n2p(z3), color=RED)
        )
        
        labels_z = VGroup(
            MathTex("z_1 = 3+2i", color=YELLOW, font_size=24).next_to(dots[0], UR),
            MathTex("z_2 = -2+3i", color=GREEN, font_size=24).next_to(dots[1], UL),
            MathTex("z_3 = 2-2i", color=RED, font_size=24).next_to(dots[2], DR)
        )
        
        self.play(Create(dots))
        self.play(Write(labels_z))
        
        # Modulus
        for i, z in enumerate([z1, z2, z3]):
            line = Line(plane.n2p(0), plane.n2p(z), color=dots[i].get_color())
            self.play(Create(line))
        
        self.wait(2)

class MatrixOperations(Scene):
    """Exercise: Matrix multiplication"""
    def construct(self):
        title = Text("Multiplication de Matrices", font_size=36)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Matrices
        A = Matrix([[1, 2], [3, 4]], left_bracket="(", right_bracket=")")
        B = Matrix([[5, 6], [7, 8]], left_bracket="(", right_bracket=")")
        equals = MathTex("=")
        C = Matrix([[19, 22], [43, 50]], left_bracket="(", right_bracket=")")
        
        A.shift(LEFT * 4)
        B.next_to(A, RIGHT, buff=0.5)
        equals.next_to(B, RIGHT, buff=0.5)
        C.next_to(equals, RIGHT, buff=0.5)
        
        self.play(Write(A))
        self.play(Write(B))
        self.wait()
        
        # Highlight calculation
        self.play(
            A.get_rows()[0].animate.set_color(YELLOW),
            B.get_columns()[0].animate.set_color(YELLOW)
        )
        self.wait()
        
        self.play(Write(equals), Write(C))
        self.wait(2)

class IntegralCalculation(Scene):
    """Exercise: Definite integral"""
    def construct(self):
        title = Text("Calcul d'Intégrale", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        axes = Axes(
            x_range=[0, 4], y_range=[0, 5],
            axis_config={"include_tip": True}
        )
        labels = axes.get_axis_labels(x_label="x", y_label="f(x)")
        
        self.play(Create(axes), Write(labels))
        
        # Function
        curve = axes.plot(lambda x: x**2, color=BLUE, x_range=[0, 3])
        func_label = MathTex("f(x) = x^2", color=BLUE).next_to(curve, UR)
        
        self.play(Create(curve), Write(func_label))
        
        # Riemann rectangles
        rects = axes.get_riemann_rectangles(
            curve, x_range=[0, 2], dx=0.25,
            color=YELLOW, fill_opacity=0.5
        )
        
        self.play(Create(rects))
        
        # Integral formula
        integral = MathTex(
            "\\int_0^2 x^2 dx = \\left[\\frac{x^3}{3}\\right]_0^2 = \\frac{8}{3}",
            font_size=36
        ).to_edge(DOWN)
        
        self.play(Write(integral))
        self.wait(2)

# ============= CHEMISTRY ANIMATIONS =============

class ChemicalReactionMechanism(Scene):
    """Exercise: Chemical reaction"""
    def construct(self):
        title = Text("Réaction Chimique", font_size=40)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Reaction
        reaction = MathTex(
            "2H_2", "+", "O_2", "\\rightarrow", "2H_2O",
            font_size=60
        )
        
        self.play(Write(reaction))
        self.wait()
        
        # Highlight reactants
        self.play(
            reaction[0].animate.set_color(BLUE),
            reaction[2].animate.set_color(RED)
        )
        self.wait()
        
        # Highlight product
        self.play(reaction[4].animate.set_color(GREEN))
        
        # Energy diagram
        axes = Axes(
            x_range=[0, 3], y_range=[0, 5],
            x_length=6, y_length=3
        ).shift(DOWN * 1.5)
        
        energy_curve = axes.plot(
            lambda x: 2 + 2*np.exp(-(x-1.5)**2/0.2) if x < 1.5 else 1,
            color=YELLOW
        )
        
        self.play(Create(axes), Create(energy_curve))
        self.wait(2)

class OrganicChemistry(Scene):
    """Exercise: Organic molecules - Alcohols"""
    def construct(self):
        title = Text("Chimie Organique: Alcools", font_size=36)
        title.to_edge(UP)
        self.play(Write(title))
        
        # Ethanol
        ethanol = MathTex("CH_3-CH_2-OH", font_size=60)
        self.play(Write(ethanol))
        
        # Highlight functional group
        oh_box = SurroundingRectangle(ethanol[0][-2:], color=RED)
        oh_label = Text("Groupe hydroxyle", font_size=24, color=RED).next_to(oh_box, DOWN)
        
        self.play(Create(oh_box), Write(oh_label))
        self.wait()
        
        # Properties
        props = VGroup(
            Text("• Soluble dans l'eau", font_size=20),
            Text("• Point d'ébullition élevé", font_size=20),
            Text("• Réactions: oxydation, estérification", font_size=20)
        ).arrange(DOWN, aligned_edge=LEFT).to_edge(DOWN)
        
        self.play(Write(props), run_time=2)
        self.wait(2)

