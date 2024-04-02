mod model;
mod service;
use clap::Parser;
use service::config::parse_config;
use service::config::read_config;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    #[arg(short, long)]
    cmd: String,

    #[arg(short, long)]
    folder: String,
}

fn main() {
    let config = read_config().expect("Error reading config");
    let (ssh_servers, folders) = parse_config(config);
    let args = Args::parse();
    match args.cmd.as_str() {
        "ls" => {
            let folder = service::folder::get(args.folder.clone(), folders);
            if let Some(folder) = folder {
                service::ssh::ls(folder, ssh_servers);
            } else {
                println!("Error locating folder: {}", args.folder)
            }
        }
        _ => println!("{:#?}\n{:#?}\n{:#?}\n", args, ssh_servers, folders),
    }
}
