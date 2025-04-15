use log::{debug, info, warn};
use teloxide::{dispatching::dialogue::GetChatId, prelude::Requester, types::CallbackQuery, Bot};

use crate::{
    buttons::{failure_kb, moderation_kb},
    memory::{PromoStore, Store},
    message::IntoMessage,
    requests::{Promo, PromoVariants, Verdict},
    HandlerResult,
};

async fn generate(bot: &Bot, c: &CallbackQuery, d: &Store) -> HandlerResult<()> {
    // ** do request to server **
    let msg_id = c.clone().message.unwrap().id();
    let chat_id = c.chat_id().unwrap();
    bot.edit_message_text(chat_id, msg_id, "Wait for it...")
        .await?;
    let promo = Promo::generate().await;
    match promo {
        Ok(v) => {
            v.edit_message(bot, msg_id, chat_id, moderation_kb())
                .await?;
            if let PromoVariants::Ok(p) = v {
                d.update(PromoStore::Last { promo: Some(p) }).await?;
            } else {
                d.update(PromoStore::Last { promo: None }).await?;
            };
        }
        Err(e) => {
            warn!("{}", e);
            "Something went wrong"
                .to_string()
                .edit_message(bot, msg_id, chat_id, failure_kb())
                .await?;
            d.update(PromoStore::Last { promo: None }).await?;
        }
    };
    Ok(())
}

pub async fn handle_approve(
    bot: Bot,
    c: CallbackQuery,
    p: PromoStore,
    d: Store,
) -> HandlerResult<()> {
    debug!("{:?}", p);
    if let Some(promo) = p.promo() {
        let verdict = Verdict::new(promo.partner_id, promo.promo_id, true);
        let status = verdict.send().await;
        warn!("{:?}", status);
    }
    generate(&bot, &c, &d).await?;
    Ok(())
}

pub async fn handle_ban(bot: Bot, c: CallbackQuery, p: PromoStore, d: Store) -> HandlerResult<()> {
    // ** do request to server for ban **
    info!("{:?}", p);
    if let Some(promo) = p.promo() {
        let verdict = Verdict::new(promo.partner_id, promo.promo_id, false);
        let _ = verdict.send().await;
    }
    generate(&bot, &c, &d).await?;
    Ok(())
}

pub async fn handle_failure(
    bot: Bot,
    c: CallbackQuery,
    promo: PromoStore,
    d: Store,
) -> HandlerResult<()> {
    info!("{:?}", promo.promo());
    generate(&bot, &c, &d).await
}
