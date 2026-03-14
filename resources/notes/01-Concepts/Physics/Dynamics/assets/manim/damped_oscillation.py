from manim import *

class DampedOscillation(Scene):
    def construct(self):
        title = Text("Damped Oscillation", font_size=24).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 6, 1],
            y_range=[-1.5, 1.5, 1],
            x_length=8,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("t")
        y_label = axes.get_y_axis_label("x(t)")
        
        def damped(t):
            return np.exp(-0.5 * t) * np.cos(2 * np.pi * t)
        
        graph = axes.plot(damped, x_range=[0, 6], color=BLUE)
        
        def envelope_pos(t):
            return np.exp(-0.5 * t)
        
        envelope_up = axes.plot(envelope_pos, x_range=[0, 6], color=GRAY)
        envelope_down = axes.plot(lambda t: -np.exp(-0.5 * t), x_range=[0, 6], color=GRAY)
        
        eq = MathTex(r"x(t) = Ae^{-\gamma t}\cos(\omega't)", font_size=28).to_edge(DOWN)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(envelope_up), run_time=1)
        self.play(Create(envelope_down), run_time=1)
        self.play(Create(graph), run_time=3)
        self.play(Write(eq))
        self.wait(3)
