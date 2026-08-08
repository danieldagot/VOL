fn main() {
    println!("{}", "  vol  ".trim());
    println!("{}", "a,b,c".split(',').collect::<Vec<_>>().join("-"));
    println!("{}", "vocabulary".contains("cab"));
    println!("{}", "a-a-a".replace('-', "+"));
}
