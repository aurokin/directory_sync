mod model;
mod service;
use service::config::parse_config;
use service::config::read_config;

fn main() {
    let config = read_config().expect("Error reading config");
    let (ssh_servers, folders) = parse_config(config);

    for ssh_server_map in ssh_servers {
        let ssh_server = ssh_server_map.1;
        println!("{:?}", ssh_server);
    }
    for folder_map in folders {
        let folder = folder_map.1;
        println!("{:?}", folder);
    }
}
