mod model;
mod service;
use clap::Parser;
use service::config::parse_config;
use service::config::read_config;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    #[command(subcommand)]
    cmd: CliCmd,
}
#[derive(Parser, Debug)]
struct CmdArgs {
    target: String,
}
#[derive(Parser, Debug)]
enum CliCmd {
    Ls(CmdArgs),
}

fn main() {
    let config = read_config().expect("Error reading config");
    let (ssh_servers, folders) = parse_config(config);
    let args = Args::parse();
    match args.cmd {
        CliCmd::Ls(cmd_args) => {
            let args_folder = cmd_args.target;
            let folder = service::folder::get(args_folder.clone(), folders);
            if let Some(folder) = folder {
                service::ssh::ls(folder, ssh_servers);
            } else {
                println!("Error locating folder: {}", args_folder)
            }
        } // _ => println!("{:#?}\n{:#?}\n{:#?}\n", args, ssh_servers, folders),
    }
}
