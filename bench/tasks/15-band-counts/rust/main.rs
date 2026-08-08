fn main() {
    let vals = [22, 18, 25, 31, 29, 17, 24, 28, 20, 26];
    let mut total = 0;
    let mut hot = 0;
    let mut mild = 0;
    let mut cold = 0;
    for &v in &vals {
        total += v;
        if v >= 28 {
            hot += 1;
        } else if v >= 20 {
            mild += 1;
        } else {
            cold += 1;
        }
    }
    println!("{}", vals.len());
    println!("{}", total / vals.len() as i32);
    println!("{}", hot);
    println!("{}", mild);
    println!("{}", cold);
}
