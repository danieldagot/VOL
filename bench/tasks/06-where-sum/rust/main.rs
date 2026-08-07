fn main() {
    let numbers = [4, 7, 2, 9, 12];
    let total: i32 = numbers.iter().copied().filter(|&n| n > 5).sum();
    println!("Sum: {}", total);
    assert_eq!(total, 28, "Unexpected collection total.");
}
