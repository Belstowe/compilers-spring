fn min(slice: &[i64]) -> Some(i64) {
    if slice.len() == 0 {
        return None;
    }

    let mut minimal: i64 = slice[0];
    for elem in slice.iter() {
        if elem < minimal {
            minimal = elem;
        }
    }
    return Some(minimal);
}