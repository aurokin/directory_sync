use std::process::Command;

fn main() {
    let mut custom_print = Command::new("echo");
    custom_print.arg("meow");
    let output = custom_print.output().expect("Error running echo");
    let output_str = String::from_utf8(output.stdout).unwrap();
    // custom_print("Twiggle says MEOW");
    println!("{output_str}");
}
