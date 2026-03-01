import gleam/io
import lustre
import lustre/attribute.{class, placeholder, value}
import lustre/element.{button, div, h1, h2, h3, input, p, text}
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

pub fn init() -> Model {
  Model(count: 0, problem: "")
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
  div([class("container")], [
    h1([], [text("BAC Unifié")]),
    p([], [text("Compteur: " <> int.to_string(model.count))]),
    div([class("buttons")], [
      button([on_click(Decrement)], [text("-")]),
      button([on_click(Increment)], [text("+")]),
    ]),
    p([], [text("Problème: " <> model.problem)]),
    input([value(model.problem), on_input(UpdateProblem)]),
    button([on_click(Submit)], [text("Soumettre")]),
  ])
}

pub fn main() {
  let app = lustre.application(init, update, view)
  lustre.start(app, "#app")
}
