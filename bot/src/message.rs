use teloxide::{
    payloads::{EditMessageTextSetters, SendMessageSetters},
    prelude::Requester,
    types::{ChatId, InlineKeyboardMarkup, Message, MessageId},
    Bot,
};

use crate::{
    buttons::failure_kb, requests::PromoVariants, templates::ShowWithTemplate, HandlerResult,
};

pub trait IntoMessage {
    async fn send(
        &self,
        bot: &Bot,
        chat_id: ChatId,
        kb: InlineKeyboardMarkup,
    ) -> HandlerResult<Message>
    where
        Self: ShowWithTemplate,
    {
        let msg = bot
            .send_message(chat_id, self.render())
            .reply_markup(kb)
            .parse_mode(teloxide::types::ParseMode::Html)
            .await?;
        Ok(msg)
    }
    async fn edit_message(
        &self,
        bot: &Bot,
        message_id: MessageId,
        chat_id: ChatId,
        kb: InlineKeyboardMarkup,
    ) -> HandlerResult<Message>
    where
        Self: ShowWithTemplate,
    {
        let msg = bot
            .edit_message_text(chat_id, message_id, self.render())
            .reply_markup(kb)
            .parse_mode(teloxide::types::ParseMode::Html)
            .await?;
        Ok(msg)
    }
}

impl IntoMessage for String {}

impl IntoMessage for PromoVariants<String> {
    async fn send(
        &self,
        bot: &Bot,
        chat_id: ChatId,
        kb: InlineKeyboardMarkup,
    ) -> HandlerResult<Message>
    where
        Self: ShowWithTemplate,
    {
        let msg = bot
            .send_message(chat_id, self.render())
            .parse_mode(teloxide::types::ParseMode::Html);
        let msg = match self {
            PromoVariants::Ok(_) => msg.reply_markup(kb),
            PromoVariants::Empty(_) => msg.reply_markup(failure_kb()),
        }
        .await?;
        Ok(msg)
    }

    async fn edit_message(
        &self,
        bot: &Bot,
        message_id: MessageId,
        chat_id: ChatId,
        kb: InlineKeyboardMarkup,
    ) -> HandlerResult<Message>
    where
        Self: ShowWithTemplate,
    {
        let msg = bot
            .edit_message_text(chat_id, message_id, self.render())
            .parse_mode(teloxide::types::ParseMode::Html);
        let msg = match self {
            PromoVariants::Ok(_) => msg.reply_markup(kb),
            PromoVariants::Empty(_) => msg.reply_markup(failure_kb()),
        }
        .await?;
        Ok(msg)
    }
}
