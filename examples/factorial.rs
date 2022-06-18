fn factorial(n: i64) -> i64 {
    if n <= 1 {
        return 1;
    }
    return factorial(n - 1);
}

fn main() {
    writeln(factorial(10));
}