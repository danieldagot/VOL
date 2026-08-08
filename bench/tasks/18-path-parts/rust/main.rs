use std::path::{Path, PathBuf};

fn main() {
    let joined: PathBuf = ["reports", "2026", "q1.json"].iter().collect();
    println!("{}", joined.display());
    let path = Path::new("reports/2026/q1.json");
    println!("{}", path.file_name().unwrap().to_string_lossy());
    println!("{}", path.parent().unwrap().display());
    println!(".{}", path.extension().unwrap().to_string_lossy());
}
