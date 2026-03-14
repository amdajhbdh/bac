from manim import *

class ActionPotential(Scene):
    def construct(self):
        title = Text("Neuron Action Potential", font_size=20).to_edge(UP)
        self.add(title)
        
        axes = Axes(
            x_range=[0, 5, 1],
            y_range=[-90, 40, 20],
            x_length=8,
            y_length=4,
            axis_config={"include_tip": True}
        )
        
        x_label = axes.get_x_axis_label("Time (ms)")
        y_label = axes.get_y_axis_label("Potential (mV)")
        
        # Resting, depolarization, repolarization, hyperpolarization
        def potential(t):
            if t < 0.5:
                return -70
            elif t < 1.5:
                return -70 + 100*(t-0.5)
            elif t < 2.5:
                return 30 - 60*(t-1.5)
            elif t < 4:
                return -60 + 10*(t-2.5)
            else:
                return -70 + 10*(t-4)
        
        curve = axes.plot(potential, x_range=[0, 5], color=BLUE)
        
        labels = VGroup(
            Text("Resting", font_size=12).next_to(axes.coords_to_point(0.2, -70), DOWN),
            Text("Depolarization", font_size=12).next_to(axes.coords_to_point(1, 0), UP),
            Text("Repolarization", font_size=12).next_to(axes.coords_to_point(2, -20), DOWN),
            Text("Hyperpolarization", font_size=12).next_to(axes.coords_to_point(3.5, -60), DOWN),
        )
        
        self.add(axes, x_label, y_label, title)
        self.play(Create(curve), run_time=3)
        self.play(Write(labels))
        self.wait(3)
