// BAC Unified - Complex Numbers Visualization
// Problem: P(z) = z³ - (5+7i)z² + (-6+26i)z + 24 - 24i
// Solutions: z1 = 2, z2 = 4i, z3 = 3 + 3i

use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);

    // === Title ===
    let title = scene.text("Complex Numbers: Polynomial Solutions")
        .with_position(-6.0, 6.5)
        .with_color(Color::YELLOW)
        .with_font_size(36.0)
        .make();
    scene.play(title.show_creation()).run_time(1.5);

    // === Draw Complex Plane ===
    let real_axis = scene.line()
        .from(-8.0, 0.0)
        .to(8.0, 0.0)
        .with_color(Color::WHITE)
        .with_thickness(2.0)
        .make();
    
    let imag_axis = scene.line()
        .from(0.0, -6.0)
        .to(0.0, 8.0)
        .with_color(Color::WHITE)
        .with_thickness(2.0)
        .make();
    
    let real_label = scene.text("Re(z)")
        .with_position(7.0, -0.5)
        .with_color(Color::WHITE)
        .with_font_size(18.0)
        .make();
    
    let imag_label = scene.text("Im(z)")
        .with_position(0.3, 7.0)
        .with_color(Color::WHITE)
        .with_font_size(18.0)
        .make();
    
    scene.play(real_axis.show_creation()).run_time(1.0);
    scene.play(imag_axis.show_creation()).run_time(1.0);
    scene.play(real_label.show_creation()).run_time(0.5);
    scene.play(imag_label.show_creation()).run_time(0.5);

    // Grid
    for i in -7..=7 {
        if i != 0 {
            let _v_line = scene.line()
                .from(i as f32, -6.0)
                .to(i as f32, 8.0)
                .with_color(Color::GRAY.with_alpha(0.3))
                .make();
        }
    }
    for i in -5..=7 {
        if i != 0 {
            let _h_line = scene.line()
                .from(-8.0, i as f32)
                .to(8.0, i as f32)
                .with_color(Color::GRAY.with_alpha(0.3))
                .make();
        }
    }

    // === Point A: z = 2 ===
    let point_a = scene.circle()
        .with_radius(0.3)
        .with_position(1.6, 0.0)
        .with_color(Color::GREEN)
        .make();
    
    let label_a = scene.text("A(2)")
        .with_position(1.9, 0.3)
        .with_color(Color::GREEN)
        .with_font_size(20.0)
        .make();
    
    let a_vertical = scene.line()
        .from(1.6, 0.0)
        .to(1.6, 0.0)
        .with_color(Color::GREEN.with_alpha(0.5))
        .make();
    
    scene.play(point_a.show_creation()).run_time(0.8);
    scene.play(label_a.show_creation()).run_time(0.5);

    // === Point B: z = 4i ===
    let point_b = scene.circle()
        .with_radius(0.3)
        .with_position(0.0, 3.2)
        .with_color(Color::RED)
        .make();
    
    let label_b = scene.text("B(4i)")
        .with_position(0.3, 3.5)
        .with_color(Color::RED)
        .with_font_size(20.0)
        .make();
    
    let b_horizontal = scene.line()
        .from(0.0, 3.2)
        .to(0.0, 3.2)
        .with_color(Color::RED.with_alpha(0.5))
        .make();
    
    scene.play(point_b.show_creation()).run_time(0.8);
    scene.play(label_b.show_creation()).run_time(0.5);

    // === Point C: z = 3 + 3i ===
    let point_c = scene.circle()
        .with_radius(0.3)
        .with_position(2.4, 2.4)
        .with_color(Color::BLUE)
        .make();
    
    let label_c = scene.text("C(3+3i)")
        .with_position(2.7, 2.7)
        .with_color(Color::BLUE)
        .with_font_size(20.0)
        .make();
    
    let c_vertical = scene.line()
        .from(2.4, 0.0)
        .to(2.4, 2.4)
        .with_color(Color::BLUE.with_alpha(0.5))
        .make();
    
    let c_horizontal = scene.line()
        .from(0.0, 2.4)
        .to(2.4, 2.4)
        .with_color(Color::BLUE.with_alpha(0.5))
        .make();
    
    scene.play(point_c.show_creation()).run_time(0.8);
    scene.play(label_c.show_creation()).run_time(0.5);
    scene.play(c_vertical.show_creation()).run_time(0.5);
    scene.play(c_horizontal.show_creation()).run_time(0.5);

    // === Triangle ABC ===
    let triangle_ab = scene.line()
        .from(1.6, 0.0)
        .to(0.0, 3.2)
        .with_color(Color::WHITE)
        .with_thickness(2.0)
        .make();
    
    let triangle_bc = scene.line()
        .from(0.0, 3.2)
        .to(2.4, 2.4)
        .with_color(Color::WHITE)
        .with_thickness(2.0)
        .make();
    
    let triangle_ca = scene.line()
        .from(2.4, 2.4)
        .to(1.6, 0.0)
        .with_color(Color::WHITE)
        .with_thickness(2.0)
        .make();
    
    scene.play(triangle_ab.show_creation()).run_time(1.0);
    scene.play(triangle_bc.show_creation()).run_time(1.0);
    scene.play(triangle_ca.show_creation()).run_time(1.0);

    // === Locus Circle (diameter AB) ===
    // Center: (0.8, 1.6), Radius: 1.6
    let center_x = 0.8;
    let center_y = 1.6;
    let radius = 1.6;
    
    for i in 0..=50 {
        let angle = (i as f32 / 50.0) * 6.28318;
        let x = center_x + radius * angle.cos();
        let y = center_y + radius * angle.sin();
        
        let point = scene.dot()
            .with_position(x, y)
            .with_color(Color::YELLOW.with_alpha(0.6))
            .with_radius(0.08)
            .make();
    }

    // === Solutions Panel ===
    let solutions_title = scene.text("Solutions of P(z) = 0:")
        .with_position(-8.0, -2.5)
        .with_color(Color::GREEN)
        .with_font_size(24.0)
        .make();
    
    let sol_a = scene.text("z1 = 2")
        .with_position(-8.0, -3.5)
        .with_color(Color::GREEN)
        .with_font_size(22.0)
        .make();
    
    let sol_b = scene.text("z2 = 4i")
        .with_position(-8.0, -4.2)
        .with_color(Color::RED)
        .with_font_size(22.0)
        .make();
    
    let sol_c = scene.text("z3 = 3 + 3i")
        .with_position(-8.0, -4.9)
        .with_color(Color::BLUE)
        .with_font_size(22.0)
        .make();
    
    scene.play(solutions_title.show_creation()).run_time(0.8);
    scene.play(sol_a.show_creation()).run_time(0.5);
    scene.play(sol_b.show_creation()).run_time(0.5);
    scene.play(sol_c.show_creation()).run_time(0.5);

    // === Triangle Analysis ===
    let analysis_title = scene.text("Triangle Analysis:")
        .with_position(3.0, -2.0)
        .with_color(Color::ORANGE)
        .with_font_size(24.0)
        .make();
    
    let ratio = scene.text("(zC-zA)/(zC-zB) = i")
        .with_position(3.0, -3.0)
        .with_color(Color::WHITE)
        .with_font_size(20.0)
        .make();
    
    let nature = scene.text("Triangle ABC is:")
        .with_position(3.0, -4.0)
        .with_color(Color::WHITE)
        .with_font_size(20.0)
        .make();
    
    let right_iso = scene.text("RIGHT-ANGLED ISOSCELES")
        .with_position(3.0, -4.8)
        .with_color(Color::YELLOW)
        .with_font_size(24.0)
        .make();
    
    let at_c = scene.text("at point C")
        .with_position(3.0, -5.5)
        .with_color(Color::YELLOW)
        .with_font_size(20.0)
        .make();
    
    scene.play(analysis_title.show_creation()).run_time(0.8);
    scene.play(ratio.show_creation()).run_time(0.5);
    scene.play(nature.show_creation()).run_time(0.5);
    scene.play(right_iso.show_creation()).run_time(0.5);
    scene.play(at_c.show_creation()).run_time(0.5);

    // === Locus Explanation ===
    let locus_title = scene.text("Locus Gamma:")
        .with_position(3.0, -7.0)
        .with_color(Color::YELLOW)
        .with_font_size(22.0)
        .make();
    
    let locus_desc = scene.text("Circle with diameter [AB]")
        .with_position(3.0, -7.7)
        .with_color(Color::CYAN)
        .with_font_size(20.0)
        .make();
    
    scene.play(locus_title.show_creation()).run_time(0.5);
    scene.play(locus_desc.show_creation()).run_time(0.5);

    scene
}

fn main() {
    noon::app(model).update(update).view(view).run();
}

fn model<'a>(app: &App) -> Scene {
    app.new_window()
        .size(1920, 1080)
        .title("BAC Unified - Complex Numbers")
        .view(view)
        .build()
        .unwrap();
    scene(app.window_rect())
}

fn update(app: &App, scene: &mut Scene, _update: Update) {
    scene.update(app.time, app.window_rect());
}

fn view(app: &App, scene: &Scene, frame: Frame) {
    let draw = app.draw();
    draw.background().color(Color::rgb(0.08, 0.1, 0.15));
    scene.draw(draw.clone());
    draw.to_frame(app, &frame).unwrap();
}
