use std::mem;

fn gcd(mut a: i64, mut b: i64) -> i64 {
    while b > 0 {
        a %= b;
        mem::swap(&mut a, &mut b);
    }
    return a;
}