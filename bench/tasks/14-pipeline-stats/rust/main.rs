fn main() {
    let nums = [3, 8, 1, 12, 5, 9, 4, 15, 2, 11];
    let large: Vec<i32> = nums.iter().copied().filter(|&n| n > 5).collect();
    let large_count = large.len();
    let large_sum: i32 = large.iter().sum();
    let large_double: i32 = large.iter().map(|n| n * 2).sum();
    let small_count = nums.iter().filter(|&&n| n < 5).count();
    println!("{}", large_count);
    println!("{}", large_sum);
    println!("{}", large_double);
    println!("{}", small_count);
}
