//! MMR Runner
use crate::{
    api::ethereum::epoch,
    mmr::{hash::MergeHash, helper, H256},
    result::Error,
    ShadowShared,
};
use cmmr::MMR;
use cmmr::{MMRStore};
use primitives::rpc::RPC;
use std::{env, thread, time};

/// MMR Runner
#[derive(Clone)]
pub struct Runner(ShadowShared);

impl AsRef<ShadowShared> for Runner {
    fn as_ref(&self) -> &ShadowShared {
        &self.0
    }
}

impl AsMut<ShadowShared> for Runner {
    fn as_mut(&mut self) -> &mut ShadowShared {
        &mut self.0
    }
}

impl From<ShadowShared> for Runner {
    fn from(s: ShadowShared) -> Self {
        Self(s)
    }
}

impl Runner {
    /// Async epoch
    pub fn epoch(block: u64) {
        if !epoch(block) {
            thread::sleep(time::Duration::from_secs(10));
            Self::epoch(block);
        }
    }

    /// translate from hash to network name
    pub fn network_name(h: &str) -> &str {
        match h {
            "41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d" => "Ropsten",
            "d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3" => "Mainnet",
            "6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177" => "Rinkeby",
            "bf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a" => "Goerli",
            _ => "unknown",
        }
    }

	/// Start the runner
    pub async fn start(&mut self) {
        while let Err(err) = self.run().await {
            error!("{:?}", err);
			actix_rt::time::delay_for(time::Duration::from_secs(10)).await;
        }
    }

    async fn run(&mut self) -> Result<(), Error> {
        let mut mmr_size = helper::mmr_size_from_store(&self.as_ref().db);
        info!("last mmr size {}", mmr_size);

        if mmr_size > 0 {
            let valid_hash = H256::hex(&(&self.as_ref().store).get_elem(0).unwrap().unwrap_or([0;32]));
            let hash_from_ethereum = self.0.eth.get_header_by_number(0).await?.hash;

            if let Some(hash) = hash_from_ethereum {
                let rpc_hash = H256::hex(&hash);
                assert!(valid_hash == rpc_hash, "RPC network should be {} but {}", Runner::network_name(&valid_hash), Runner::network_name(&rpc_hash));
            } else {
                panic!("rpc request is unreachable");
            }
        }

        let last_leaf = helper::mmr_size_to_last_leaf(mmr_size as i64);
        let mut ptr = if last_leaf == 0 { 0 } else { last_leaf + 1 };

        // Using a cache rpc block number to optimize and reduce rpc call.
        let mut last_rpc_block_number = self.0.eth.block_number().await?;

        loop {
            // checking finalization
            if last_rpc_block_number < (ptr as u64 + 12) {
                trace!("Pause 10s due to finalization checking, prepare to push block {}, last block number from rpc is {}", ptr, last_rpc_block_number);
                last_rpc_block_number = self.0.eth.block_number().await?;
                actix_rt::time::delay_for(time::Duration::from_secs(10)).await;
                continue;
            }

            // Note:
            //
            // This trigger is ugly, need better solution in the future, ptr % 30000 is to compatible with existing production, can be removed later
            // if (ptr + 15000) % 30000 == 0 || ptr % 30000 == 0 {
            //     thread::spawn(move || Self::epoch((ptr + 15000) as u64))
            //         .join()
            //         .unwrap_or_default();
            // }

            match self.push(ptr, mmr_size).await {
                Err(_e) => {
                    actix_rt::time::delay_for(time::Duration::from_secs(10)).await;
                }
                Ok(mmr_size_new) => {
                    if ptr
                        % env::var("MMR_LOG")
                            .unwrap_or_else(|_| "10000".to_string())
                            .parse::<i64>()
                            .unwrap_or(10000)
                        == 0
                    {
                        trace!("Pushed mmr {} into database", ptr);
                    }

                    mmr_size = mmr_size_new;
                    ptr += 1;
                }
            }
        }
    }

    /// Push new header hash into storage
    pub async fn push(&mut self, number: i64, mmr_size: u64) -> Result<u64, Error> {
        let mut mmr = MMR::<_, MergeHash, _>::new(mmr_size, &self.as_ref().store);
        let hash_from_ethereum = self.0.eth.get_header_by_number(number as u64).await?.hash;
        if let Some(hash) = hash_from_ethereum {
            mmr.push(hash)?;
            let mmr_size_new = mmr.mmr_size();

            mmr.commit()?;
            Ok(mmr_size_new)
        } else {
            Err(Error::Primitive(format!(
                "Get Ethereum header {} from ethereum rpc failed",
                number
            )))
        }
    }
}
