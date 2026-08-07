fn main() {
    let (mut a, mut b) = (0_i64, 1_i64);
    for _ in 0..8 {
        println!("{}", a);
        let c = a + b;
        a = b;
        b = c;
    }
}
