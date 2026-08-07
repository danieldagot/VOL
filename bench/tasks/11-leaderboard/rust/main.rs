fn main() {
    let player1 = [8, 6, 9, 5, 10, 7];
    let player2 = [7, 8, 6, 9, 8, 8];
    let p1_total: i32 = player1.iter().sum();
    let p2_total: i32 = player2.iter().sum();
    let winner = if p1_total > p2_total { "Player 1" } else { "Player 2" };
    let p1_strong = player1.iter().filter(|&&s| s >= 8).count();
    let p2_strong = player2.iter().filter(|&&s| s >= 8).count();
    println!("Player 1 total: {}", p1_total);
    println!("Player 2 total: {}", p2_total);
    println!("Winner: {}", winner);
    println!("P1 rounds 8+: {}", p1_strong);
    println!("P2 rounds 8+: {}", p2_strong);
}
