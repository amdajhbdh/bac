from manim import *

class ProjectileMotion(Scene):
    def construct(self):
        title = Text("Projectile Motion", font_size=20).to_edge(UP)
        axes = Axes(x_range=[0,8,2], y_range=[0,5,1], x_length=8, y_length=4)
        trajectory = axes.plot(lambda x: x - x**2/16, x_range=[0,16], color=BLUE)
        self.add(title, axes, trajectory)
        self.wait(3)
