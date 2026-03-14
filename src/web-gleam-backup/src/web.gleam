import gleam/int
import gleam/io
import lustre
import lustre/attribute.{class, placeholder, value}
import lustre/effect.{type Effect}
import lustre/element.{text}
import lustre/element/html
import lustre/event.{on_click, on_input}

pub type Model {
  Model(count: Int, problem: String)
}

pub type Msg {
  Increment
  Decrement
  UpdateProblem(String)
  Submit
}

pub fn init(
  model: Model,
  effect: effect.Effect(Msg),
) -> #(Model, effect.Effect(Msg)) {
  #(Model(count: 0, problem: ""), effect.none())
}

pub fn update(model: Model, msg: Msg) -> #(Model, Nil) {
  case msg {
    Increment -> #(Model(..model, count: model.count + 1), Nil)
    Decrement -> #(Model(..model, count: model.count - 1), Nil)
    UpdateProblem(p) -> #(Model(..model, problem: p), Nil)
    Submit -> #(Model(..model, count: model.count + 1), Nil)
  }
}

pub fn view(model: Model) -> element.Element(Msg) {
  html.div([class("container")], [
    html.h1([], [text("BAC Unifié")]),
    html.span([], [text("Compteur: " <> int.to_string(model.count))]),
    html.div([class("buttons")], [
      html.button([on_click(Decrement)], [text("-")]),
      html.button([on_click(Increment)], [text("+")]),
    ]),
    html.span([], [text("Problème: " <> model.problem)]),
    html.input([value(model.problem), on_input(UpdateProblem)]),
    html.button([on_click(Submit)], [text("Soumettre")]),
  ])
}

pub fn main() {
  let app = lustre.application(init, update, view)
  lustre.start(app, "#app")
}
