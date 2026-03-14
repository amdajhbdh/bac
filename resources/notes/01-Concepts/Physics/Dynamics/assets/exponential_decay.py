from manim import *

class ExponentialDecay(Scene):
    def construct(self):
        title = Text("Radioactive Decay: N(t) = N0 * e^(-λt)", font_size=18).to_edge(UP)
        
        axes = Axes(
            x_range=[0, 10, 2],
            y_range=[0, 1.2, 0.5],
            x_length=8,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("Time")
        y_label = axes.get_y_axis_label("N(t)/N0")
        
        decay = axes.plot(lambda t: np.exp(-0.5*t), x_range=[0, 10], color=BLUE)
        half_life = DashedLine(axes.coords_to_point(1.386, 0.5), axes.coords_to_point(1.386, 0), color=RED)
        
        self.add(title, axes, x_label, y_label)
        self.play(Create(decay), run_time=2)
        self.play(Create(half_life))
        self.wait(3)
