//! `shadow` commands
use crate::result::Error;
use structopt::{clap::AppSettings, StructOpt};

mod count;
mod import;
mod run;
mod trim;

#[derive(StructOpt)]
#[structopt(setting = AppSettings::InferSubcommands)]
enum Opt {
    /// Current block height in mmr store
    Count,
    /// Start shadow service
    Run {
        /// Http server port
        #[structopt(short, long, default_value = "3000")]
        port: u16,
        /// Verbose mode
        #[structopt(short, long)]
        verbose: bool,
    },
    /// Imports mmr from geth
    Import {
        /// Datadir of geth
        #[structopt(short, long)]
        path: String,
        /// Header limits
        #[structopt(short, long)]
        limit: i32,
    },
    /// Trim mmr from target leaf
    Trim {
        /// The target leaf
        #[structopt(short, long)]
        leaf: u64,
    },
}

/// Exec `shadow` binary
pub async fn exec() -> Result<(), Error> {
    match Opt::from_args() {
        Opt::Count => count::exec(),
        Opt::Run { port, verbose } => run::exec(port, verbose).await,
        Opt::Import { path, limit } => import::exec(path, limit),
        Opt::Trim { leaf } => trim::exec(leaf),
    }?;

    Ok::<(), Error>(())
}
