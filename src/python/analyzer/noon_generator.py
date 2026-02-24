"""
BAC Unified - Robust Noon Code Generator
Generates Noon animation code following the correct API patterns
"""

import re
from typing import Dict, List, Any


class NoonCodeGenerator:
    """
    Generates Noon (Rust) animation code for BAC problems.
    Uses the correct Noon API: .text().with_text("string") and u32 font sizes
    """

    def __init__(self):
        pass

    def generate_from_solution(self, problem: str, solution: Dict) -> str:
        """Main entry point - generate complete Noon code."""

        viz_type = solution.get("visualization", {}).get("type", "general")

        if "integral" in viz_type.lower() or "∫" in problem:
            return self.generate_integral(problem, solution)
        elif "complex" in viz_type.lower() or "polynome" in problem.lower():
            return self.generate_complex(problem, solution)
        elif "projectile" in viz_type.lower() or "projectile" in problem.lower():
            return self.generate_projectile(problem, solution)
        elif "circuit" in viz_type.lower() or "circuit" in problem.lower():
            return self.generate_circuit(problem, solution)
        elif "molecule" in viz_type.lower() or "molecule" in problem.lower():
            return self.generate_molecule(problem, solution)
        else:
            return self.generate_generic(problem, solution)

    def generate_integral(self, problem: str, solution: Dict) -> str:
        """Generate code for integral/calculus problems."""

        steps = solution.get("steps", [])
        answer = solution.get("answer", "?")

        code = """// Auto-generated: Integral Visualization
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text()
        .with_text("Integral Calculus")
        .with_position(-4.0, 5.0)
        .with_color(Color::YELLOW)
        .with_font_size(36)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // Problem
    let problem = scene.text()
        .with_text("Problem: f(x) = x^2")
        .with_position(-5.0, 4.0)
        .with_color(Color::WHITE)
        .with_font_size(20)
        .make();
    scene.play(problem.show_creation()).run_time(1.0);
    
    // Axes
    let x_axis = scene.line()
        .from(-7.0, 0.0)
        .to(7.0, 0.0)
        .with_color(Color::WHITE)
        .make();
    let y_axis = scene.line()
        .from(0.0, -4.0)
        .to(0.0, 6.0)
        .with_color(Color::WHITE)
        .make();
    scene.play(x_axis.show_creation()).run_time(0.8);
    scene.play(y_axis.show_creation()).run_time(0.8);
    
    // Function curve
    let mut prev: Option<(f32, f32)> = None;
    for i in 0..100 {
        let x = -6.0 + (i as f32 / 100.0 * 12.0);
        let y = x * x;
        
        let _dot = scene.circle()
            .with_radius(0.05)
            .with_position(x, y)
            .with_color(Color::CYAN)
            .make();
        
        if let Some((px, py)) = prev {
            let line = scene.line()
                .from(px, py)
                .to(x, y)
                .with_color(Color::CYAN)
                .make();
            scene.play(line.show_creation()).run_time(0.01);
        }
        prev = Some((x, y));
    }
    
    // Area under curve
    for i in 0..20 {
        let x = 0.0 + (i as f32 / 20.0) * 2.0;
        let y = x * x;
        
        let _shade = scene.line()
            .from(x, 0.0)
            .to(x, y)
            .with_color(Color::CYAN)
            .make();
    }
    scene.wait();
    
    // Solution steps
"""

        step_y = -1.5
        for i, step in enumerate(steps[:3]):
            code += f'''    let step{i} = scene.text()
        .with_text("{step}")
        .with_position(-5.0, {step_y - i * 0.6})
        .with_color(Color::GREEN)
        .with_font_size(16)
        .make();
    scene.play(step{i}.show_creation()).run_time(0.8);
'''

        code += f"""    // Final answer
    let ans = scene.text()
        .with_text("Answer: {answer}")
        .with_position(-5.0, -4.5)
        .with_color(Color::RED)
        .with_font_size(28)
        .make();
    scene.play(ans.show_creation()).run_time(1.0);
    
    scene
}}

fn main() {{
    noon::app(|_| scene(app.window_rect())).run();
}}
"""
        return code

    def generate_complex(self, problem: str, solution: Dict) -> str:
        """Generate code for complex numbers problems."""

        points = solution.get("visualization", {}).get("points", [])

        if not points:
            points = [
                {"label": "A", "z": "2", "x": 2.0, "y": 0.0, "color": "GREEN"},
                {"label": "B", "z": "4i", "x": 0.0, "y": 4.0, "color": "RED"},
                {"label": "C", "z": "3+3i", "x": 3.0, "y": 3.0, "color": "BLUE"},
            ]

        code = """// Auto-generated: Complex Numbers Visualization
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text()
        .with_text("Complex Plane")
        .with_position(-3.0, 5.5)
        .with_color(Color::YELLOW)
        .with_font_size(36)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // Axes
    let real_axis = scene.line()
        .from(-6.0, 0.0)
        .to(6.0, 0.0)
        .with_color(Color::WHITE)
        .make();
    let imag_axis = scene.line()
        .from(0.0, -4.0)
        .to(0.0, 6.0)
        .with_color(Color::WHITE)
        .make();
    scene.play(real_axis.show_creation()).run_time(0.8);
    scene.play(imag_axis.show_creation()).run_time(0.8);
    
    // Labels
    let re_label = scene.text()
        .with_text("Re(z)")
        .with_position(5.0, -0.5)
        .with_color(Color::WHITE)
        .with_font_size(16)
        .make();
    let im_label = scene.text()
        .with_text("Im(z)")
        .with_position(0.3, 5.5)
        .with_color(Color::WHITE)
        .with_font_size(16)
        .make();
    scene.play(re_label.show_creation()).run_time(0.3);
    scene.play(im_label.show_creation()).run_time(0.3);
"""

        for pt in points:
            color = pt.get("color", "GREEN")
            code += f'''    // Point {pt["label"]}: {pt["z"]}
    let point_{pt["label"].lower()} = scene.circle()
        .with_radius(0.25)
        .with_position({pt["x"]}, {pt["y"]})
        .with_color(Color::{color})
        .make();
    let label_{pt["label"].lower()} = scene.text()
        .with_text("{pt["label"]}({pt["z"]})")
        .with_position({pt["x"] + 0.3}, {pt["y"] + 0.3})
        .with_color(Color::{color})
        .with_font_size(18)
        .make();
    scene.play(point_{pt["label"].lower()}.show_creation()).run_time(0.5);
    scene.play(label_{pt["label"].lower()}.show_creation()).run_time(0.3);
'''

        if len(points) >= 3:
            code += """    // Triangle connections
    let line_ab = scene.line()
        .from(2.0, 0.0)
        .to(0.0, 4.0)
        .with_color(Color::WHITE)
        .make();
    let line_bc = scene.line()
        .from(0.0, 4.0)
        .to(3.0, 3.0)
        .with_color(Color::WHITE)
        .make();
    let line_ca = scene.line()
        .from(3.0, 3.0)
        .to(2.0, 0.0)
        .with_color(Color::WHITE)
        .make();
    scene.play(line_ab.show_creation()).run_time(0.5);
    scene.play(line_bc.show_creation()).run_time(0.5);
    scene.play(line_ca.show_creation()).run_time(0.5);
"""

        code += """    // Solutions panel
    let sol_title = scene.text()
        .with_text("Solutions:")
        .with_position(-6.0, -1.5)
        .with_color(Color::YELLOW)
        .with_font_size(24)
        .make();
    scene.play(sol_title.show_creation()).run_time(0.5);
    
    let sol = scene.text()
        .with_text("z1=2, z2=4i, z3=3+3i")
        .with_position(-6.0, -2.5)
        .with_color(Color::WHITE)
        .with_font_size(18)
        .make();
    scene.play(sol.show_creation()).run_time(0.5);
    
    let triangle_type = scene.text()
        .with_text("Triangle: Right-angled isosceles")
        .with_position(-6.0, -3.5)
        .with_color(Color::ORANGE)
        .with_font_size(18)
        .make();
    scene.play(triangle_type.show_creation()).run_time(0.5);
    
    scene
}}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
"""
        return code

    def generate_projectile(self, problem: str, solution: Dict) -> str:
        """Generate code for projectile motion."""

        code = """// Auto-generated: Projectile Motion
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text()
        .with_text("Projectile Motion")
        .with_position(-4.0, 5.5)
        .with_color(Color::YELLOW)
        .with_font_size(36)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // Ground
    let ground = scene.rectangle()
        .with_size(16.0, 0.5)
        .with_position(0.0, -4.5)
        .with_color(Color::GREEN)
        .make();
    scene.play(ground.fade_in()).run_time(0.5);
    
    // Parameters
    let v0 = 25.0;
    let angle = 45.0;
    let g = 9.8;
    let rad = angle * 3.14159 / 180.0;
    
    let params = scene.text()
        .with_text("v0 = 25 m/s, angle = 45 degrees")
        .with_position(-5.0, 4.5)
        .with_color(Color::WHITE)
        .with_font_size(18)
        .make();
    scene.play(params.show_creation()).run_time(0.8);
    
    // Ball
    let ball = scene.circle()
        .with_radius(0.3)
        .with_position(-6.0, -4.0)
        .with_color(Color::RED)
        .make();
    scene.play(ball.show_creation()).run_time(0.5);
    
    // Animate trajectory
    let mut t = 0.0;
    let dt = 0.05;
    while t < 3.0 {
        let x = -6.0 + v0 * rad.cos() * t;
        let y = -4.0 + v0 * rad.sin() * t - 0.5 * g * t * t;
        
        if y < -4.25 { break; }
        
        scene.play(ball.move_to(x, y)).run_time(dt);
        t += dt;
    }
    
    // Physics equations
    let eq1 = scene.text()
        .with_text("x(t) = v0*cos(theta)*t")
        .with_position(2.0, 4.0)
        .with_color(Color::GREEN)
        .with_font_size(16)
        .make();
    let eq2 = scene.text()
        .with_text("y(t) = v0*sin(theta)*t - gt^2/2")
        .with_position(2.0, 3.3)
        .with_color(Color::GREEN)
        .with_font_size(16)
        .make();
    scene.play(eq1.show_creation()).run_time(0.5);
    scene.play(eq2.show_creation()).run_time(0.5);
    
    // Range
    let range = v0 * v0 * 2.0 * rad.sin() * rad.cos() / g;
    let ans = scene.text()
        .with_text(&format!("Range: {:.1} m", range))
        .with_position(2.0, 2.0)
        .with_color(Color::RED)
        .with_font_size(24)
        .make();
    scene.play(ans.show_creation()).run_time(1.0);
    
    scene
}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
"""
        return code

    def generate_circuit(self, problem: str, solution: Dict) -> str:
        """Generate code for circuit diagrams."""

        code = """// Auto-generated: Electric Circuit
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text()
        .with_text("RC Circuit")
        .with_position(-2.0, 5.5)
        .with_color(Color::YELLOW)
        .with_font_size(36)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // Battery
    let battery = scene.rectangle()
        .with_size(0.4, 1.5)
        .with_position(-5.0, 0.0)
        .with_color(Color::GRAY)
        .make();
    scene.play(battery.show_creation()).run_time(0.5);
    
    // Wire top
    let wire1 = scene.line()
        .from(-4.8, 0.5)
        .to(-4.8, 2.5)
        .with_color(Color::YELLOW)
        .make();
    let wire2 = scene.line()
        .from(-4.8, 2.5)
        .to(-2.0, 2.5)
        .with_color(Color::YELLOW)
        .make();
    scene.play(wire1.show_creation()).run_time(0.3);
    scene.play(wire2.show_creation()).run_time(0.3);
    
    // Resistor
    let resistor = scene.rectangle()
        .with_size(1.5, 0.4)
        .with_position(-1.0, 2.5)
        .with_color(Color::BROWN)
        .make();
    let res_label = scene.text()
        .with_text("R = 10 Ohm")
        .with_position(-1.8, 3.2)
        .with_color(Color::WHITE)
        .with_font_size(14)
        .make();
    scene.play(resistor.show_creation()).run_time(0.5);
    scene.play(res_label.show_creation()).run_time(0.3);
    
    // Wire after resistor
    let wire3 = scene.line()
        .from(0.0, 2.5)
        .to(0.0, 1.5)
        .with_color(Color::YELLOW)
        .make();
    scene.play(wire3.show_creation()).run_time(0.3);
    
    // Capacitor
    let cap1 = scene.line()
        .from(-0.15, 1.5)
        .to(-0.15, 0.0)
        .with_color(Color::BLUE)
        .make();
    let cap2 = scene.line()
        .from(0.15, 1.5)
        .to(0.15, 0.0)
        .with_color(Color::BLUE)
        .make();
    let cap_label = scene.text()
        .with_text("C = 100 uF")
        .with_position(0.5, 0.8)
        .with_color(Color::WHITE)
        .with_font_size(14)
        .make();
    scene.play(cap1.show_creation()).run_time(0.3);
    scene.play(cap2.show_creation()).run_time(0.3);
    scene.play(cap_label.show_creation()).run_time(0.3);
    
    // Wire return
    let wire4 = scene.line()
        .from(0.0, 0.0)
        .to(0.0, -1.5)
        .with_color(Color::YELLOW)
        .make();
    let wire5 = scene.line()
        .from(0.0, -1.5)
        .to(-5.0, -1.5)
        .with_color(Color::YELLOW)
        .make();
    let wire6 = scene.line()
        .from(-5.0, -1.5)
        .to(-5.0, -0.5)
        .with_color(Color::YELLOW)
        .make();
    scene.play(wire4.show_creation()).run_time(0.3);
    scene.play(wire5.show_creation()).run_time(0.3);
    scene.play(wire6.show_creation()).run_time(0.3);
    
    // Equations
    let eq1 = scene.text()
        .with_text("V = IR")
        .with_position(2.0, 4.0)
        .with_color(Color::GREEN)
        .with_font_size(24)
        .make();
    let eq2 = scene.text()
        .with_text("tau = RC")
        .with_position(2.0, 3.3)
        .with_color(Color::GREEN)
        .with_font_size(24)
        .make();
    scene.play(eq1.show_creation()).run_time(0.5);
    scene.play(eq2.show_creation()).run_time(0.5);
    
    scene
}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
"""
        return code

    def generate_molecule(self, problem: str, solution: Dict) -> str:
        """Generate code for molecule visualization."""

        code = """// Auto-generated: Molecule Structure
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text()
        .with_text("Ethanol: C2H5OH")
        .with_position(-3.0, 5.5)
        .with_color(Color::YELLOW)
        .with_font_size(36)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // Carbon 1
    let c1 = scene.circle()
        .with_radius(0.35)
        .with_position(-2.5, 0.0)
        .with_color(Color::GRAY)
        .make();
    let c1_label = scene.text()
        .with_text("C")
        .with_position(-2.65, -0.15)
        .with_color(Color::WHITE)
        .with_font_size(18)
        .make();
    scene.play(c1.show_creation()).run_time(0.5);
    scene.play(c1_label.show_creation()).run_time(0.3);
    
    // Carbon 2
    let c2 = scene.circle()
        .with_radius(0.35)
        .with_position(-0.5, 0.0)
        .with_color(Color::GRAY)
        .make();
    let c2_label = scene.text()
        .with_text("C")
        .with_position(-0.65, -0.15)
        .with_color(Color::WHITE)
        .with_font_size(18)
        .make();
    scene.play(c2.show_creation()).run_time(0.5);
    scene.play(c2_label.show_creation()).run_time(0.3);
    
    // C-C bond
    let cc_bond = scene.line()
        .from(-2.15, 0.0)
        .to(-0.85, 0.0)
        .with_color(Color::WHITE)
        .make();
    scene.play(cc_bond.show_creation()).run_time(0.5);
    
    // Oxygen
    let o = scene.circle()
        .with_radius(0.3)
        .with_position(1.2, 0.0)
        .with_color(Color::RED)
        .make();
    let o_label = scene.text()
        .with_text("O")
        .with_position(1.05, -0.15)
        .with_color(Color::WHITE)
        .with_font_size(18)
        .make();
    scene.play(o.show_creation()).run_time(0.5);
    scene.play(o_label.show_creation()).run_time(0.3);
    
    // C-O bond
    let co_bond = scene.line()
        .from(-0.15, 0.0)
        .to(0.9, 0.0)
        .with_color(Color::WHITE)
        .make();
    scene.play(co_bond.show_creation()).run_time(0.5);
    
    // O-H bond
    let oh_bond = scene.line()
        .from(1.5, 0.0)
        .to(1.8, -0.5)
        .with_color(Color::WHITE)
        .make();
    
    // H on OH
    let h_oh = scene.circle()
        .with_radius(0.2)
        .with_position(2.0, -0.7)
        .with_color(Color::WHITE)
        .make();
    scene.play(oh_bond.show_creation()).run_time(0.3);
    scene.play(h_oh.show_creation()).run_time(0.3);
    
    // Multiple H atoms on carbons
    let h1a = scene.circle().with_radius(0.2).with_position(-2.8, 0.8).with_color(Color::WHITE).make();
    let h1b = scene.circle().with_radius(0.2).with_position(-2.8, -0.8).with_color(Color::WHITE).make();
    let h1c = scene.circle().with_radius(0.2).with_position(-3.2, 0.0).with_color(Color::WHITE).make();
    scene.play(h1a.show_creation()).run_time(0.2);
    scene.play(h1b.show_creation()).run_time(0.2);
    scene.play(h1c.show_creation()).run_time(0.2);
    
    let h2a = scene.circle().with_radius(0.2).with_position(-0.8, 0.8).with_color(Color::WHITE).make();
    let h2b = scene.circle().with_radius(0.2).with_position(-0.8, -0.8).with_color(Color::WHITE).make();
    scene.play(h2a.show_creation()).run_time(0.2);
    scene.play(h2b.show_creation()).run_time(0.2);
    
    // Info panel
    let info = scene.text()
        .with_text("Molar Mass: 46.07 g/mol")
        .with_position(-5.0, -2.0)
        .with_color(Color::CYAN)
        .with_font_size(20)
        .make();
    let boil = scene.text()
        .with_text("Boiling: 78.4 C")
        .with_position(-5.0, -2.8)
        .with_color(Color::WHITE)
        .with_font_size(20)
        .make();
    scene.play(info.show_creation()).run_time(0.5);
    scene.play(boil.show_creation()).run_time(0.5);
    
    // Functional group highlight
    let oh_group = scene.text()
        .with_text("-OH (Hydroxyl)")
        .with_position(1.5, -2.8)
        .with_color(Color::YELLOW)
        .with_font_size(18)
        .make();
    scene.play(oh_group.show_creation()).run_time(0.5);
    
    scene
}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
"""
        return code

    def generate_generic(self, problem: str, solution: Dict) -> str:
        """Generate generic visualization for any problem."""

        answer = solution.get("answer", "?")

        code = f"""// Auto-generated: Generic Problem Visualization
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {{
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text()
        .with_text("Problem Solution")
        .with_position(-3.0, 5.0)
        .with_color(Color::YELLOW)
        .with_font_size(36)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // Problem text
    let problem = scene.text()
        .with_text("Problem")
        .with_position(-1.5, 3.5)
        .with_color(Color::WHITE)
        .with_font_size(24)
        .make();
    scene.play(problem.show_creation()).run_time(0.8);
    
    // Solution
    let sol_title = scene.text()
        .with_text("Solution:")
        .with_position(-1.5, 1.5)
        .with_color(Color::GREEN)
        .with_font_size(24)
        .make();
    scene.play(sol_title.show_creation()).run_time(0.8);
    
    // Answer
    let ans = scene.text()
        .with_text("Answer: {answer}")
        .with_position(-2.0, -0.5)
        .with_color(Color::RED)
        .with_font_size(32)
        .make();
    scene.play(ans.show_creation()).run_time(1.0);
    
    scene
}}

fn main() {{
    noon::app(|_| scene(app.window_rect())).run();
}}
"""
        return code


if __name__ == "__main__":
    gen = NoonCodeGenerator()

    sol = {
        "steps": ["Expand function", "Find antiderivative", "Apply bounds"],
        "answer": "8/3",
        "visualization": {"type": "integral", "function": "x^2"},
    }
    code = gen.generate_from_solution("∫x^2 dx from 0 to 2", sol)
    print("Integral code generated!")
    print(code[:500])
