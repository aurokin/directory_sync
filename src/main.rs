mod model;
mod service;
use clap::Parser;
use model::cli::CliCmd;
use service::config::parse_config;
use service::config::read_config;

use crate::service::cli::ls;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
pub struct Args {
    #[command(subcommand)]
    pub cmd: CliCmd,
    #[arg(short, long, action)]
    pub link: bool,
}

fn main() {
    let config = read_config().expect("Error reading config");
    let (ssh_servers, folders, links) = parse_config(config);
    let args = Args::parse();
    let is_link = args.link;
    match args.cmd {
        CliCmd::Ls(cmd_args) => ls(cmd_args, folders, ssh_servers, links, is_link),
        CliCmd::Pull(pull_args) => {
            println!("{:?}", pull_args);
        }
        CliCmd::Push(push_args) => {
            println!("{:?}", push_args);
        }
    }
}
