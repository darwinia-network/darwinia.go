use crate::{
    mmr::{MergeHash, H256},
    ShadowShared,
};
use actix_web::{web, Responder};
use cmmr::MMR;
use primitives::{chain::ethereum::EthereumRelayHeaderParcelJson, rpc::RPC};

/// Proof target header
///
/// ```
/// use actix_web::web;
/// use darwinia_shadow::{api::ethereum, ShadowShared};
///
/// // GET `/ethereum/parcel/19`
/// ethereum::parcel(web::Path::from("19".to_string()), web::Data::new(ShadowShared::new(None)));
/// ```
#[allow(clippy::eval_order_dependence)]
pub async fn handle(block: web::Path<String>, shared: web::Data<ShadowShared>) -> impl Responder {
    let num: u64 = block.to_string().parse().unwrap_or(0);
    let root = if num == 0 {
        "0000000000000000000000000000000000000000000000000000000000000000".to_string()
    } else {
        H256::hex(
            &MMR::<_, MergeHash, _>::new(cmmr::leaf_index_to_mmr_size(num - 1), &shared.store)
                .get_root()
                .unwrap_or_default(),
        )
    };

    web::Json(EthereumRelayHeaderParcelJson {
        header: shared
            .eth
            .get_header_by_number(num)
            .await
            .unwrap_or_default()
            .into(),
        mmr_root: format!("0x{}", root),
    })
}
