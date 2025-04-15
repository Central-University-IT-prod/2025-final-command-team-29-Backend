use hypertext::html_elements;
use hypertext::{self, Renderable, Rendered};

pub trait ShowWithTemplate {
    fn render(&self) -> Rendered<String>;
}

pub fn greeting(name: &str) -> Rendered<String> {
    hypertext::maud! {"Hello, " b { (name) } "!\n"
        "This is telegram bot to moderate promos of your company!\n"
    }
    .render()
}
