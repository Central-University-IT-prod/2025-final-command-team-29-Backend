use teloxide::{dispatching::dialogue::InMemStorage, prelude::Dialogue};

use crate::requests::Promo;

#[derive(Clone, Default, Debug)]
pub enum PromoStore {
    #[default]
    Start,
    Last {
        promo: Option<Promo>,
    },
}

impl PromoStore {
    pub fn promo(&self) -> Option<Promo> {
        match self {
            PromoStore::Start => None,
            PromoStore::Last { promo } => promo.to_owned(),
        }
    }
}
pub type Store = Dialogue<PromoStore, InMemStorage<PromoStore>>;
