use std::fs;
use std::path::Path;

// struct Folder {
//     name: String,
//     path: String,
// }

fn main() {
    println!("Hello World");
    let config = read_config().expect("Error Reading File");
    println!("{config}");
}

fn read_config() -> Option<String> {
    let path = Path::new("/Users/auro/.dirsync.toml");
    let path_exists = path.exists();
    let path_is_file = path.is_file();

    if path_exists && path_is_file {
        let file = fs::read_to_string(path).expect("Error reading file");
        return Some(file);
    }

    return None;
}
