fn fibonacci(n: i64) -> i64 {
    if n <= 2 {
        return 1
    }
    return fibonacci(n - 2) + fibonacci(n - 1)
}

fn main() {
    ruster::writeln_i64(fibonacci(10))
}