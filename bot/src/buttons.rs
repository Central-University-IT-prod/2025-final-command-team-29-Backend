use teloxide::types::{InlineKeyboardButton, InlineKeyboardMarkup};

macro_rules! inline_kb {
    ($($row:expr),*) => {{
        let mut kb = InlineKeyboardMarkup::default();
        $(
            let buttons = $row
                .into_iter()
                .map(|(val, callback)| button(val, callback))
                .collect::<Vec<InlineKeyboardButton>>();
            kb = kb.append_row(buttons);
        )*
        kb
    }};
}

fn button<T: Into<String>>(val: T, callback: T) -> InlineKeyboardButton {
    InlineKeyboardButton::callback(val, callback)
}

pub fn moderation_kb() -> InlineKeyboardMarkup {
    inline_kb! { [("✅", "approve"), ("❌", "ban") ] }
}

pub fn failure_kb() -> InlineKeyboardMarkup {
    inline_kb! { [ ("Retry", "failure") ] }
}
