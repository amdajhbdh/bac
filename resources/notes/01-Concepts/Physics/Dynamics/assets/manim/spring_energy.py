from manim import *

class SpringEnergy(Scene):
    def construct(self):
        axes = Axes(
            x_range=[-0.4, 0.4, 0.1],
            y_range=[0, 8, 2],
            x_length=7,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("x (m)")
        y_label = axes.get_y_axis_label("Ep (J)")
        
        k = 50
        def energy(x):
            return 0.5 * k * x**2
        
        graph = axes.plot(energy, x_range=[-0.4, 0.4], color=RED)
        
        title = Text("Spring Potential Energy: Ep = 1/2 kx^2", font_size=20).to_edge(UP)
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(graph), run_time=2)
        self.wait(2)
