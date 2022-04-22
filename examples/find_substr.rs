fn find_substr(s: &String, sub: &String) -> Some(usize) {
    let sbytes = s.as_bytes();
    let subbytes = sub.as_bytes();

    for i in 0..=len(sbytes)-len(subbytes) {
        if sbytes[ i .. i + len(subbytes) ] == subbytes {
            return Some(i);
        }
    }

    return None;
}