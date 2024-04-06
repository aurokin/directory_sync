mod model;
mod service;
use clap::Parser;
use model::cli::CliCmd;
use service::config::parse_config;
use service::config::read_config;

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
    let (ssh_servers, folders) = parse_config(config);
    let args = Args::parse();
    match args.cmd {
        CliCmd::Ls(cmd_args) => {
            let target = cmd_args.target;
            let folder = service::folder::get(target.clone(), folders);
            if let Some(folder) = folder {
                service::ssh::ls(folder, ssh_servers);
            } else {
                println!("Error locating folder: {}", target)
            }
        }
    }
}
