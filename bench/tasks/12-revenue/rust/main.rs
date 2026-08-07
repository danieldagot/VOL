fn main() {
    let revenue = [240, 175, 289, 150, 225, 199, 180, 178];
    let total: i32 = revenue.iter().sum();
    let high_value: Vec<i32> = revenue.iter().copied().filter(|&r| r >= 200).collect();
    let high_sum: i32 = high_value.iter().sum();
    let budget_count = revenue.iter().filter(|&&r| r < 200).count();
    println!("Total revenue: {}", total);
    println!("Premium orders (200+): {}", high_value.len());
    println!("Premium revenue: {}", high_sum);
    println!("Budget orders: {}", budget_count);
}
