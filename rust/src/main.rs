mod model;
mod service;
use clap::Parser;
use model::cli::CliCmd;
use service::config::parse_config;
use service::config::read_config;

use crate::service::cli::ls;
use crate::service::cli::pull;
use crate::service::cli::push;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
pub struct Args {
    #[command(subcommand)]
    pub cmd: CliCmd,
    #[arg(short, long, action)]
    pub link: bool,
    #[arg(short, long, action)]
    pub force: bool,
}

fn main() {
    let config = read_config().expect("Error reading config");
    let (ssh_servers, folders, links, work_folder) = parse_config(config);
    let args = Args::parse();
    let is_link = args.link;
    let is_force = args.force;

    match args.cmd {
        CliCmd::Ls(cmd_args) => ls(cmd_args, is_link, folders, links, ssh_servers),
        CliCmd::Pull(cmd_args) => pull(
            cmd_args,
            is_link,
            is_force,
            work_folder,
            folders,
            links,
            ssh_servers,
        ),
        CliCmd::Push(cmd_args) => push(
            cmd_args,
            is_link,
            is_force,
            work_folder,
            folders,
            links,
            ssh_servers,
        ),
    }
}
