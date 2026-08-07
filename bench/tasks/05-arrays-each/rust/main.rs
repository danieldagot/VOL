fn main() {
    let mut scores = vec![72, 95, 81, 64];
    println!("Students: {}", scores.len());
    scores[3] = 70;
    for score in &scores {
        if *score >= 80 {
            println!("High score: {}", score);
        }
    }
}
