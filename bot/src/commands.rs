use crate::{
    buttons::{failure_kb, moderation_kb},
    memory::{PromoStore, Store},
    message::IntoMessage,
    requests::{Promo, PromoVariants},
    templates::{self},
    HandlerResult,
};
use log::warn;
use teloxide::{payloads::SendMessageSetters, prelude::*, utils::command::BotCommands};

#[derive(BotCommands, Clone)]
#[command(
    rename_rule = "lowercase",
    description = "These commands are supported:"
)]
pub enum BaseCommand {
    #[command(description = "You're here")]
    Help,
    #[command(description = "View list of not approved promos")]
    List,
    #[command(description = "Greeting from bot :)")]
    Start,
}

pub async fn handle_start(bot: Bot, msg: Message) -> HandlerResult<()> {
    let greet = templates::greeting(msg.chat.first_name().unwrap_or("stranger"));
    bot.send_message(msg.chat.id, greet)
        .parse_mode(teloxide::types::ParseMode::Html)
        .await?;
    Ok(())
}

pub async fn handle_list(bot: Bot, msg: Message, d: Store) -> HandlerResult<()> {
    let promo = Promo::generate().await;
    match promo {
        Ok(v) => {
            v.send(&bot, msg.chat.id, moderation_kb()).await?;
            if let PromoVariants::Ok(p) = v {
                d.update(PromoStore::Last { promo: Some(p) }).await?;
            } else {
                d.update(PromoStore::Last { promo: None }).await?;
            };
        }
        Err(e) => {
            let promo = Promo::new(msg.chat.id);
            d.update(PromoStore::Last { promo: Some(promo) }).await?;
            warn!("{}", e);
            "Something went wrong"
                .to_string()
                .send(&bot, msg.chat.id, failure_kb())
                .await?;
        }
    };
    Ok(())
}

pub async fn handle_help(bot: Bot, msg: Message) -> HandlerResult<()> {
    bot.send_message(msg.chat.id, BaseCommand::descriptions().to_string())
        .await?;
    Ok(())
}
