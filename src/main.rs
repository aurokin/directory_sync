use std::path::Path;

fn main() {
    println!("Hello World");
    let path = Path::new("/Users/auro/.dirsync.toml");
    let path_exists = path.exists();
    let path_is_file = path.is_file();
    println!("{path_exists}");
    println!("{path_is_file}");
}
