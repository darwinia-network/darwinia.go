//! format bridge
#![macro_use]
use scale::{Decode, Encode};
use std::fmt::{Debug, Display, Formatter, Result as FmtResult};

#[macro_export]
/// Convert bytes to hex
macro_rules! hex {
    ($bytes:expr) => {{
        let mut s = String::new();
        for i in $bytes {
            s.push_str(&format!("{:02x}", i));
        }
        s
    }};
}

#[macro_export]
/// Convert hext to `Vec<u8>` or `[u8; n]`
macro_rules! bytes {
    // Convert hex to Vec<u8>
    ($hex:expr) => {{
        let mut h = $hex;
        if h.starts_with("0x") {
            h = &h[2..];
        }

        (0..h.len())
            .step_by(2)
            .map(|i| u8::from_str_radix(&h[i..i + 2], 16))
            .collect::<Result<Vec<u8>, _>>()
            .unwrap_or_default()
    }};

    // Convert hex to [u8; $bits]
    ($hex:expr, $bits:expr) => {{
        let mut hash = [0_u8; $bits];
        hash.copy_from_slice(&bytes!($hex));
        hash
    }};
}

macro_rules! u256 {
    // Convert hex to [u64; 4]
    ($hex:expr) => {{
        let mut h = $hex;
        if h.starts_with("0x") {
            h = &h[2..];
        }

        let mut u256 = [0_u64; 4];
        let bytes = (0..h.len())
            .step_by(2)
            .map(|i| u64::from_str_radix(&h[i..i + 2], 16))
            .collect::<Result<Vec<u64>, _>>()
            .unwrap_or_default();
        u256.copy_from_slice(&bytes);
        u256
    }};
}

#[derive(Decode, Encode)]
struct Bloom(pub [u8; 256]);

impl Display for Bloom {
    fn fmt(&self, f: &mut Formatter<'_>) -> FmtResult {
        f.write_str(&hex!(self.0.as_ref()))
    }
}

impl Debug for Bloom {
    fn fmt(&self, f: &mut Formatter<'_>) -> FmtResult {
        f.debug_list().entries(self.0.iter()).finish()
    }
}

impl Default for Bloom {
    fn default() -> Bloom {
        Bloom([0; 256])
    }
}

/// Eth header
#[derive(Decode, Encode, Debug)]
pub struct EthHeader {
    parent_hash: [u8; 32],
    timestamp: u64,
    number: u64,
    author: [u8; 20],
    transactions_root: [u8; 32],
    uncles_hash: [u8; 32],
    extra_data: Vec<u8>,
    state_root: [u8; 32],
    receipts_root: [u8; 32],
    log_bloom: Bloom,
    gas_used: [u64; 4],
    gas_limit: [u64; 4],
    difficulty: [u64; 4],
    seal: Vec<Vec<u8>>,
    hash: Option<[u8; 32]>,
}

impl EthHeader {
    /// New EthHeader from string array
    pub fn from_go_ffi(
        parent_hash: &str,
        timestamp: u64,
        number: u64,
        author: &str,
        transactions_root: &str,
        uncles_hash: &str,
        extra_data: &str,
        state_root: &str,
        receipts_root: &str,
        log_bloom: &str,
        gas_used: &str,
        gas_limit: &str,
        difficulty: &str,
        mixh: &str,
        nonce: &str,
        hash: &str,
    ) -> EthHeader {
        EthHeader {
            parent_hash: bytes!(parent_hash, 32),
            timestamp,
            number,
            author: bytes!(author, 20),
            transactions_root: bytes!(transactions_root, 32),
            uncles_hash: bytes!(uncles_hash, 32),
            extra_data: bytes!(extra_data),
            state_root: bytes!(state_root, 32),
            receipts_root: bytes!(receipts_root, 32),
            log_bloom: Bloom(bytes!(log_bloom, 256)),
            gas_used: u256!(gas_used),
            gas_limit: u256!(gas_limit),
            difficulty: u256!(difficulty),
            seal: match mixh.is_empty() && nonce.is_empty() {
                true => vec![],
                false => vec![bytes!(mixh), bytes!(nonce)],
            },
            hash: match hash.is_empty() {
                true => None,
                false => Some(bytes!(hash, 32)),
            },
        }
    }
}

impl Default for EthHeader {
    fn default() -> EthHeader {
        EthHeader::from_go_ffi(
            "0x0000000000000000000000000000000000000000000000000000000000000000",
            0,
            0,
            "0x0000000000000000000000000000000000000000",
            "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
            "",
            "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "0x00000000",
            "0x00000000",
            "0x00000000",
            "",
            "",
            ""
        )
    }
}

// pub struct DoubleNodeWithMerkleProof {
//     pub dag_nodes: [[u8; 64]; 2],
//     pub proof: Vec<[u8; 16]>,
// }