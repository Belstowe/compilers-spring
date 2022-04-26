use std::mem;

fn gcd(a: mut i64, b: mut i64) -> i64 {
    while b > 0 {
        a %= b;
        mem::swap(&mut a, &mut b);
    }
    return a;
}