from manim import *

class MatrixTransform(Scene):
    def construct(self):
        title = Text("Matrix as Linear Transformation", font_size=18).to_edge(UP)
        
        # Grid
        grid = NumberPlane(x_range=[-3, 3], y_range=[-3, 3])
        
        # Original unit square
        square = Polygon([-1,-1,0], [1,-1,0], [1,1,0], [-1,1,0], color=BLUE, fill_opacity=0.3)
        
        # Transformed square (shear)
        def transform(point):
            x, y, z = point
            return [x + 0.5*y, y, 0]
        
        transformed = square.apply_function(transform)
        transformed.set_color(RED)
        
        labels = VGroup(
            Text("Original", color=BLUE, font_size=14),
            Text("Transformed", color=RED, font_size=14)
        ).arrange(DOWN).to_corner(DR)
        
        self.add(title, grid, square, transformed, labels)
        self.wait(3)
