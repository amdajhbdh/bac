from manim import *

class Transformer(Scene):
    def construct(self):
        title = Text("Transformer", font_size=20).to_edge(UP)
        
        # Primary coil
        primary = VGroup()
        for i in range(3):
            rect = Rectangle(width=0.3, height=0.8, color=BLUE)
            rect.shift(LEFT*0.5 + UP*i*0.5)
            primary.add(rect)
        
        # Core
        core = Rectangle(width=2, height=2, color=GRAY, fill_opacity=0.3)
        
        # Secondary coil
        secondary = VGroup()
        for i in range(3):
            rect = Rectangle(width=0.3, height=0.8, color=RED)
            rect.shift(RIGHT*0.5 + UP*i*0.5)
            secondary.add(rect)
        
        # Labels
        p_label = Text("Primary", color=BLUE, font_size=14).next_to(primary, LEFT)
        s_label = Text("Secondary", color=RED, font_size=14).next_to(secondary, RIGHT)
        
        eq = MathTex(r"\frac{V_p}{V_s} = \frac{N_p}{N_s}", font_size=24).to_corner(DR)
        
        self.add(title, core, primary, secondary, p_label, s_label)
        self.play(Write(eq))
        self.wait(3)
