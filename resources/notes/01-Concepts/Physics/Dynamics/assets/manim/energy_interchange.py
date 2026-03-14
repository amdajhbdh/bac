from manim import *

class EnergyInterchange(Scene):
    def construct(self):
        title = Text("Energy Interchange in SHM", font_size=24).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 4, 1],
            y_range=[0, 2, 0.5],
            x_length=8,
            y_length=4,
            axis_config={"include_tip": False}
        )
        
        x_label = axes.get_x_axis_label("t")
        y_label = axes.get_y_axis_label("Energy")
        
        total_energy = axes.plot(
            lambda t: 1 + 0*t,
            x_range=[0, 4],
            color=PURPLE,
            stroke_width=3
        )
        
        kinetic = axes.plot(
            lambda t: np.sin(2 * np.pi * t)**2,
            x_range=[0, 4],
            color=GREEN
        )
        
        potential = axes.plot(
            lambda t: np.cos(2 * np.pi * t)**2,
            x_range=[0, 4],
            color=RED
        )
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(total_energy), run_time=1)
        self.play(Create(kinetic), run_time=2)
        self.play(Create(potential), run_time=2)
        
        self.wait(3)
