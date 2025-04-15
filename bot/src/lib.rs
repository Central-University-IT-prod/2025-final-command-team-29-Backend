use commands::BaseCommand;
use log::info;
use memory::PromoStore;
use moderation_system::{handle_approve, handle_ban, handle_failure};
use teloxide::{
    dispatching::{dialogue::InMemStorage, HandlerExt, UpdateFilterExt},
    dptree,
    prelude::Dispatcher,
    types::{CallbackQuery, Message, Update},
    Bot,
};

mod buttons;
mod commands;
mod memory;
mod message;
mod moderation_system;
mod requests;
mod templates;

type HandlerResult<T> = Result<T, Box<dyn std::error::Error + Send + Sync>>;

pub async fn dispatch(bot: Bot) -> HandlerResult<()> {
    let command_schema = Update::filter_message()
        .inspect(|m: Message| info!("msg from: {}", m.chat.id))
        .filter_command::<BaseCommand>()
        .branch(
            dptree::entry()
                .filter(|cmd: BaseCommand| matches!(cmd, BaseCommand::Start))
                .endpoint(commands::handle_start),
        )
        .branch(
            dptree::entry()
                .filter(|cmd: BaseCommand| matches!(cmd, BaseCommand::List))
                .enter_dialogue::<Message, InMemStorage<PromoStore>, PromoStore>()
                .endpoint(commands::handle_list),
        )
        .branch(
            dptree::entry()
                .filter(|cmd: BaseCommand| matches!(cmd, BaseCommand::Help))
                .endpoint(commands::handle_help),
        );
    let callback_schema = Update::filter_callback_query()
        .inspect(|c: CallbackQuery| info!("msg from: {}", c.from.id))
        .branch(
            dptree::entry()
                .filter(|c: CallbackQuery| c.data.unwrap_or_default() == *"approve")
                .enter_dialogue::<CallbackQuery, InMemStorage<PromoStore>, PromoStore>()
                .endpoint(handle_approve),
        )
        .branch(
            dptree::entry()
                .filter(|c: CallbackQuery| c.data.unwrap_or_default() == *"failure")
                .enter_dialogue::<CallbackQuery, InMemStorage<PromoStore>, PromoStore>()
                .endpoint(handle_failure),
        )
        .branch(
            dptree::entry()
                .filter(|c: CallbackQuery| c.data.unwrap_or_default() == *"ban")
                .enter_dialogue::<CallbackQuery, InMemStorage<PromoStore>, PromoStore>()
                .endpoint(handle_ban),
        );

    let schema = dptree::entry()
        .branch(callback_schema)
        .branch(command_schema);

    Dispatcher::builder(bot, schema)
        .dependencies(dptree::deps![InMemStorage::<PromoStore>::new()])
        .enable_ctrlc_handler()
        .build()
        .dispatch()
        .await;
    Ok(())
}
