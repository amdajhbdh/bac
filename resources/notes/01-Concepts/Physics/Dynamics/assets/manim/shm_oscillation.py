from manim import *

class SHMOscillation(Scene):
    def construct(self):
        title = Text("Simple Harmonic Motion", font_size=24).to_edge(UP)
        self.add(title)
        
        axes_x = Axes(
            x_range=[0, 4, 1],
            y_range=[-1.5, 1.5, 1],
            x_length=8,
            y_length=1.5,
            axis_config={"include_tip": False}
        ).shift(UP * 2)
        
        position_graph = axes_x.plot(
            lambda t: np.cos(2 * np.pi * t),
            x_range=[0, 4],
            color=BLUE
        )
        
        axes_v = Axes(
            x_range=[0, 4, 1],
            y_range=[-1.5, 1.5, 1],
            x_length=8,
            y_length=1.5,
            axis_config={"include_tip": False}
        ).shift(UP * 0.3)
        
        velocity_graph = axes_v.plot(
            lambda t: -np.sin(2 * np.pi * t),
            x_range=[0, 4],
            color=GREEN
        )
        
        axes_a = Axes(
            x_range=[0, 4, 1],
            y_range=[-1.5, 1.5, 1],
            x_length=8,
            y_length=1.5,
            axis_config={"include_tip": False}
        ).shift(DOWN * 1.4)
        
        accel_graph = axes_a.plot(
            lambda t: -np.cos(2 * np.pi * t),
            x_range=[0, 4],
            color=RED
        )
        
        self.add(axes_x, axes_v, axes_a)
        self.play(Create(position_graph), run_time=2)
        self.play(Create(velocity_graph), run_time=2)
        self.play(Create(accel_graph), run_time=2)
        
        phase_note = Text("v leads x by 90 deg", font_size=14).to_edge(DOWN)
        self.add(phase_note)
        
        self.wait(3)
