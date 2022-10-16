fn gcd(a: i64, b: i64) -> i64 {
    while b > 0 {
        a %= b;
        let t = a;
        a = b;
        b = t;
    }
    return a;
}