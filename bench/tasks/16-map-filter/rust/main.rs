fn main() {
    let xs = [1, 2, 3, 4, 5, 6, 7, 8];
    let mapped: Vec<i32> = xs.iter().map(|x| x * 3).collect();
    let mid_count = mapped.iter().filter(|&&n| n > 10 && n < 20).count();
    let high_sum: i32 = mapped.iter().copied().filter(|&n| n > 10).sum();
    println!("{}", mid_count);
    println!("{}", high_sum);
}
