use bot::dispatch;
use teloxide::Bot;

// So --
// 1. Sign In -> 1 Command, 1 Dialogue
// 2. CRUD Campaigns -> 4 Commands, 4 Dialogues
// 3. Statistics -> 4 Commands
// --->
// Plots?
// Architecture?
// Testing?

#[tokio::main]
async fn main() {
    pretty_env_logger::init();
    log::info!("Starting bot...");

    let bot = Bot::from_env();
    dispatch(bot).await.unwrap();
}
