use std::env;

fn main() {
    env::set_var("VOL_DENSITY_PORT", "9090");
    println!(
        "{}",
        env::var("VOL_DENSITY_PORT").unwrap_or_else(|_| "missing".into())
    );
    println!(
        "{}",
        env::var("VOL_DENSITY_NO_SUCH_KEY").unwrap_or_else(|_| "fallback".into())
    );
}
