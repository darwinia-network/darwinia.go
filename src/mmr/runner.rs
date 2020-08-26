//! MMR Runner
use crate::{
    chain::eth::EthHeaderRPCResp,
    hash::{MergeHash, H256},
    helper,
    pool::{ConnPool, PooledConn},
    result::Error,
    schema::mmr_store::dsl::*,
    store::Store,
};
use cmmr::MMR;
use diesel::{dsl::count, prelude::*};
use reqwest::Client;
use std::time;

/// MMR Runner
#[derive(Clone)]
pub struct Runner {
    store: Store,
    pool: ConnPool,
    client: Client,
}

impl Runner {
    /// Get Pooled connection
    pub fn conn(&self) -> Result<PooledConn, Error> {
        let cfp = self.pool.get();
        if cfp.is_err() {
            return Err(Error::Custom("Connect to database failed".into()));
        }

        Ok(cfp.unwrap())
    }

    /// Start with sqlite3 conn
    pub fn with(pool: ConnPool) -> Runner {
        let store = Store::with(pool.clone());
        Runner {
            store,
            pool,
            client: Client::new(),
        }
    }

    /// Start the runner
    pub async fn start(&mut self) -> Result<(), Error> {
        let mut ptr = {
            let last_leaf = helper::mmr_size_to_last_leaf(self.mmr_count()?);
            if last_leaf == 0 {
                0
            } else {
                last_leaf + 1
            }
        };

        loop {
            if let Err(e) = self.push(ptr).await {
                error!("Push block to mmr_store failed: {:?}", e);
                trace!("MMR service restarting after 10s...");
                async_std::task::sleep(time::Duration::from_secs(10)).await;
            } else {
                ptr += 1;
            }
        }
    }

    /// Get block hash by number
    pub async fn get_hash(&mut self, block: i64) -> Result<String, Error> {
        Ok(EthHeaderRPCResp::get(&self.client, block as u64)
            .await?
            .result
            .hash)
    }

    /// Trim mmr
    pub fn trim(&mut self, leaf: u64) -> Result<(), Error> {
        let mpos = cmmr::leaf_index_to_pos(leaf);
        let conn = self.conn()?;
        diesel::delete(mmr_store.filter(pos.ge(mpos as i64))).execute(&conn)?;
        Ok(())
    }

    /// Push new header hash into storage
    pub async fn push(&mut self, pnumber: i64) -> Result<(), Error> {
        let mut mmr = MMR::<_, MergeHash, _>::new(
            cmmr::leaf_index_to_mmr_size((pnumber - 1) as u64),
            &self.store,
        );
        if let Err(e) = mmr.push(H256::from(
            &EthHeaderRPCResp::get(&self.client, pnumber as u64)
                .await?
                .result
                .hash,
        )) {
            return Err(Error::MMR(e));
        }

        mmr.commit()?;
        Ok(())
    }

    /// Get the count of mmr store
    pub fn mmr_count(&self) -> Result<i64, Error> {
        let conn = self.conn()?;
        let res = mmr_store.select(count(elem)).first::<i64>(&conn);
        if let Err(e) = res {
            Err(Error::Diesel(e))
        } else {
            Ok(res?)
        }
    }
}