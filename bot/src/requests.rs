use hypertext::{html_elements, maud, Renderable};
use log::info;
use reqwest::Client;
use serde::{Deserialize, Serialize};
use teloxide::types::ChatId;

use crate::{templates::ShowWithTemplate, HandlerResult};

pub enum PromoVariants<T: ShowWithTemplate> {
    Ok(Promo),
    Empty(T),
}

type EmptyListString = String;
const BASE_URL: &str = "http://prod-team-29-4254c2ee.REDACTED/api/v1/promos";

#[derive(Clone, Debug, Deserialize)]
pub struct Promo {
    pub promo_id: String,
    pub partner_id: String,
    title: String,
    description: String,
}

#[derive(Debug, Serialize)]
pub struct Verdict {
    partner_id: String,
    promo_id: String,
    verdict: bool,
}

impl Verdict {
    pub fn new(partner_id: String, promo_id: String, verdict: bool) -> Self {
        Verdict {
            partner_id,
            promo_id,
            verdict,
        }
    }
    pub async fn send(&self) -> HandlerResult<()> {
        let client = Client::new();
        let body = serde_json::to_string(self)?;
        info!("{}", body);
        client
            .post(format!("{BASE_URL}/moderation"))
            .header("Content-Type", "application/json")
            .body(body)
            .send()
            .await?
            .error_for_status()?;
        Ok(())
    }
}

impl Promo {
    // WARN: for tests
    pub fn new(id: ChatId) -> Self {
        Promo {
            promo_id: id.to_string(),
            partner_id: "1".to_string(),
            title: "1".to_string(),
            description: "1".to_string(),
        }
    }
    pub async fn generate() -> HandlerResult<PromoVariants<EmptyListString>> {
        let promo = reqwest::get(format!("{BASE_URL}/moderation"))
            .await?
            .error_for_status()?;
        if promo.status() == 204 {
            return Ok(PromoVariants::Empty(
                "List for moderation is empty".to_string(),
            ));
        };
        let text = promo.text().await?;
        info!("{}", text);
        let promo: Promo = serde_json::from_str(&text)?;
        Ok(PromoVariants::Ok(promo))
    }
}

impl ShowWithTemplate for Promo {
    fn render(&self) -> hypertext::Rendered<String> {
        maud!("Name: " b { (&self.title) } "\n" pre {(&self.description)}).render()
    }
}

impl ShowWithTemplate for EmptyListString {
    fn render(&self) -> hypertext::Rendered<String> {
        maud!( b { (self) }).render()
    }
}

impl ShowWithTemplate for PromoVariants<EmptyListString> {
    fn render(&self) -> hypertext::Rendered<String> {
        match self {
            PromoVariants::Ok(s) => s.render(),
            PromoVariants::Empty(s) => s.clone().render(),
        }
    }
}
