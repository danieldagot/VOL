fn main() {
    let temps = [22, 18, 25, 31, 29, 17, 24, 28, 20, 26];
    let total: i32 = temps.iter().sum();
    let avg = total / temps.len() as i32;
    let hot = temps.iter().filter(|&&t| t >= 28).count();
    let mild = temps.iter().filter(|&&t| t >= 20 && t < 28).count();
    let cold = temps.iter().filter(|&&t| t < 20).count();
    println!("Days measured: {}", temps.len());
    println!("Average: {}", avg);
    println!("Hot days (28+): {}", hot);
    println!("Mild days: {}", mild);
    println!("Cold days (<20): {}", cold);
}
