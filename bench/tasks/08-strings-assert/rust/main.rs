fn main() {
    let word = "hello";
    let length = word.chars().count();
    let doubled = word.to_owned() + word;
    println!("{}", length);
    println!("{}", doubled);
    assert_eq!(length, 5, "Expected length 5");
}
