// rustc-only: no serde_json. Fixed-schema extract for this task's JSON shape.
fn json_field<'a>(raw: &'a str, key: &str) -> &'a str {
    let needle = format!("\"{}\":", key);
    let start = raw.find(&needle).expect("key") + needle.len();
    let rest = raw[start..].trim_start();
    if let Some(stripped) = rest.strip_prefix('"') {
        let end = stripped.find('"').expect("string end");
        &stripped[..end]
    } else {
        let end = rest
            .find(|c| c == ',' || c == '}')
            .expect("value end");
        rest[..end].trim()
    }
}

fn main() {
    let raw = r#"{"n":3,"name":"vol"}"#;
    let n = json_field(raw, "n");
    println!("{}", n);
    println!("{}", json_field(raw, "name"));
    println!("{{\"n\":{}}}", n);
}
