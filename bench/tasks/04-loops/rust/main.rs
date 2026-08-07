fn main() {
    println!("Countdown");
    let mut remaining = 3;
    while remaining > 0 {
        println!("{}", remaining);
        remaining -= 1;
    }
    for _ in 0..2 {
        println!("Go!");
    }
}
