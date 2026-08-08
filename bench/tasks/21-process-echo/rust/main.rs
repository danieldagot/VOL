use std::process::Command;

fn main() {
    let out = Command::new("echo")
        .arg("vol")
        .output()
        .expect("echo");
    println!("{}", out.status.code().unwrap_or(1));
    println!("{}", String::from_utf8_lossy(&out.stdout).trim());
}
