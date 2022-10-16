fn factorial(n: i64) -> i64 {
    let result: i64 = 1;
    while n > 1 {
        result = result * n;
        n = n - 1;
    }
    return result;
}

fn main() {
    ruster::writeln_i64(factorial(10));
}