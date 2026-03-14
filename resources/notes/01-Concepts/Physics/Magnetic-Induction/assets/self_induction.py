from manim import *

class SelfInduction(Scene):
    def construct(self):
        title = Text("Self-Induction: Current Growth/Decay", font_size=16).to_edge(UP)
        
        axes = Axes(
            x_range=[0, 5, 1],
            y_range=[0, 1.2, 0.5],
            x_length=7,
            y_length=3.5,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("t")
        y_label = axes.get_y_axis_label("I(t)")
        
        # Growth (increasing exponential)
        growth = axes.plot(lambda t: 1 - np.exp(-t), x_range=[0, 5], color=GREEN)
        
        # Decay (decreasing exponential)  
        decay = axes.plot(lambda t: np.exp(-t), x_range=[0, 5], color=RED)
        
        growth_label = Text("Growth", color=GREEN, font_size=14).to_corner(UL)
        decay_label = Text("Decay", color=RED, font_size=14).next_to(growth_label, DOWN)
        
        eq = MathTex(r"I(t) = I_0(1 - e^{-t/\tau})", font_size=20).to_corner(DR)
        
        self.add(title, axes, x_label, y_label)
        self.play(Create(growth), run_time=1.5)
        self.play(Create(decay), run_time=1.5)
        self.add(growth_label, decay_label, eq)
        self.wait(3)
