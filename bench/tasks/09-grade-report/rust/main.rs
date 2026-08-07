fn main() {
    let scores = [85, 72, 91, 60, 78, 95, 55, 68, 88, 74];
    let total: i32 = scores.iter().sum();
    let avg = total / scores.len() as i32;
    let a_grades = scores.iter().filter(|&&s| s >= 90).count();
    let b_grades = scores.iter().filter(|&&s| s >= 80 && s < 90).count();
    let passing = scores.iter().filter(|&&s| s >= 60).count();
    let failing = scores.iter().filter(|&&s| s < 60).count();
    println!("Class average: {}", avg);
    println!("A grades: {}", a_grades);
    println!("B grades: {}", b_grades);
    println!("Passing: {}", passing);
    println!("Failing: {}", failing);
}
