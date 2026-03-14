from manim import *

class LenzLaw(Scene):
    def construct(self):
        title = Text("Lenz's Law", font_size=20).to_edge(UP)
        
        # Magnetic field into page (X symbols)
        for x in [-1, 0, 1]:
            cross = VGroup(
                Line(UP*0.2 + RIGHT*x, DOWN*0.2 + RIGHT*x, color=BLUE),
                Line(LEFT*0.2 + UP*x, RIGHT*0.2 + DOWN*x, color=BLUE)
            )
        
        # Coil
        coil = Ellipse(width=2.5, height=1.2, color=YELLOW, stroke_width=3)
        
        b_label = Text("B (into page)", color=BLUE, font_size=14).to_corner(UL)
        
        self.add(title, coil, b_label)
        self.wait(3)
