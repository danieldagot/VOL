fn greet(name: &str) -> String {
    format!("Hello, {}", name)
}

fn square(number: i32) -> i32 {
    number * number
}

fn main() {
    println!("{}", greet("friend"));
    println!("Six squared is {}", square(6));
}
