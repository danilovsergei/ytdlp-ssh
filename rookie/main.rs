use rookie::chrome;

fn main() {
    let domains = vec!["youtube.com".to_string()];
    let cookies = chrome(Some(domains)).unwrap();

    println!("# Netscape HTTP Cookie File");

    for cookie in cookies {
        // Format the cookie data according to the Netscape format
        println!(
            "{}\t{}\t{}\t{}\t{}\t{}\t{}",
            cookie.domain,
            if cookie.domain.starts_with('.') { "TRUE" } else { "FALSE" }, // Include subdomains?
            cookie.path,
            if cookie.secure { "TRUE" } else { "FALSE" },
            cookie.expires.map_or(0, |t| t as i64), // Convert u64 to i64
            cookie.name,
            cookie.value
        );
    }
}
